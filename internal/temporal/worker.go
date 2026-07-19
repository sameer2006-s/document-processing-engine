package temporal

import (
	"log"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type TemporalWorker struct {
	worker worker.Worker
	DocumentProcessingWorkflow func(ctx workflow.Context, documentID string) (DocumentProcessingWorkflowResult,error)
}

func NewTemporalWorker(temporalClient *TemporalClient) *TemporalWorker {
	return &TemporalWorker{
		worker: worker.New(temporalClient.Client, "document-processing-task-queue", worker.Options{}),
		DocumentProcessingWorkflow: DocumentProcessingWorkflow,
	}
}

func (w *TemporalWorker) RegisterWorkflows() error {
	w.worker.RegisterWorkflow(w.DocumentProcessingWorkflow)
	return nil
}

func (w *TemporalWorker) RegisterActivities() error {
	w.worker.RegisterActivity(NewOCRActivity().RunOCRActivity)
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