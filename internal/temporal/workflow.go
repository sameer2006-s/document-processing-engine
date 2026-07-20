package temporal

import (
	"time"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"go.temporal.io/sdk/workflow"
)

type DocumentProcessingWorkflowResult struct {
	Success   bool
	Message   string
	OCRResult string
	Thumbnail string
}

func DocumentProcessingWorkflow(ctx workflow.Context, documentID string) (DocumentProcessingWorkflowResult, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting document processing workflow", "documentID", documentID)

	var activity *OCRActivity
	var ocrResult string
	err := workflow.ExecuteActivity(activityCtx, activity.RunOCRActivity, documentID).Get(activityCtx, &ocrResult)
	if err != nil {
		logger.Error("Failed to run OCR activity", "error", err)
		return DocumentProcessingWorkflowResult{
			Success:   false,
			Message:   err.Error(),
			OCRResult: "failed",
		}, err
	}

	var thumbnail string
	err = workflow.ExecuteActivity(activityCtx, activity.RunThumbnailActivity, documentID).Get(activityCtx, &thumbnail)
	if err != nil {
		logger.Error("Failed to run thumbnail activity", "error", err)
		return DocumentProcessingWorkflowResult{
			Success:   false,
			Message:   err.Error(),
			Thumbnail: thumbnail,
		}, err
	}

	err = workflow.ExecuteActivity(activityCtx, activity.UpdateDocumentStatusActivity, documentID, document.DocumentStatusDone).Get(activityCtx, nil)
	if err != nil {
		logger.Error("Failed to update document status", "error", err)
		return DocumentProcessingWorkflowResult{
			Success:   false,
			Message:   err.Error(),
			OCRResult: ocrResult,
			Thumbnail: thumbnail,
		}, err
	}

	return DocumentProcessingWorkflowResult{
		Success:   true,
		Message:   "Document processing workflow completed successfully",
		OCRResult: ocrResult,
		Thumbnail: thumbnail,
	}, nil
}
