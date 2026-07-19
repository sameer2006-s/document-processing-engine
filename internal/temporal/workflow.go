package temporal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type DocumentProcessingWorkflowResult struct {
	Success   bool
	Message   string
	OCRResult string
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

	return DocumentProcessingWorkflowResult{
		Success:   true,
		Message:   "Document processing workflow completed successfully",
		OCRResult: ocrResult,
	}, nil
}
