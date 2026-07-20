package temporal

import (
	"time"
	"github.com/sameer2006-s/document-processing-engine/internal/document"
	"go.temporal.io/sdk/workflow"
	"go.temporal.io/sdk/temporal"
)

type DocumentProcessingWorkflowResult struct {
	Success   bool
	Message   string
	OCRResult string
	Thumbnail string
	Tags string
}

func DocumentProcessingWorkflow(ctx workflow.Context, documentID string) (DocumentProcessingWorkflowResult, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
			InitialInterval: 1 * time.Second,
			MaximumInterval: 10 * time.Second,
			BackoffCoefficient: 2,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting document processing workflow", "documentID", documentID)

	var activity *DocumentActivity
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

	var tags string
	err = workflow.ExecuteActivity(activityCtx, activity.RunTagActivity, documentID).Get(activityCtx, &tags)
	if err != nil {
		logger.Error("Failed to run tag activity", "error", err)
		return DocumentProcessingWorkflowResult{
			Success:   false,
			Message:   err.Error(),
			Tags: tags,
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
			Tags: tags,
		}, err
	}

	return DocumentProcessingWorkflowResult{
		Success:   true,
		Message:   "Document processing workflow completed successfully",
		OCRResult: ocrResult,
		Thumbnail: thumbnail,
		Tags: tags,
	}, nil
}
