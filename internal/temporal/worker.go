package temporal

import (
	"log"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type TemporalWorker struct {
	worker                     worker.Worker
	ocrActivity                *OCRActivity
	DocumentProcessingWorkflow func(ctx workflow.Context, documentID string) (DocumentProcessingWorkflowResult, error)
}

func NewTemporalWorker(temporalClient *TemporalClient, ocrActivity *OCRActivity) *TemporalWorker {
	return &TemporalWorker{
		worker:                     worker.New(temporalClient.Client, DocumentProcessingTaskQueue, worker.Options{}),
		ocrActivity:                ocrActivity,
		DocumentProcessingWorkflow: DocumentProcessingWorkflow,
	}
}

func (w *TemporalWorker) RegisterWorkflows() error {
	w.worker.RegisterWorkflow(w.DocumentProcessingWorkflow)
	return nil
}

func (w *TemporalWorker) RegisterActivities() error {
	w.worker.RegisterActivity(w.ocrActivity)
	return nil
}

func (w *TemporalWorker) Run() error {
	log.Println("Starting worker...")
	err := w.worker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal("Failed to run worker: ", err)
	}
	log.Println("Worker started successfully.")
	return nil
}
