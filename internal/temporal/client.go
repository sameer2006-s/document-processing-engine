package temporal

import (
	"context"

	"go.temporal.io/sdk/client"
)

const DocumentProcessingTaskQueue = "document-processing-task-queue"

type TemporalClient struct {
	Client client.Client
}

func NewTemporalClient() (*TemporalClient, error) {
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		return nil, err
	}
	return &TemporalClient{Client: c}, nil
}

func (c *TemporalClient) Close() {
	c.Client.Close()
}

// StartDocumentProcessing starts the document workflow (fire-and-forget).
func (c *TemporalClient) StartDocumentProcessing(documentID string) error {
	_, err := c.Client.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        "document-processing-workflow-" + documentID,
			TaskQueue: DocumentProcessingTaskQueue,
		},
		DocumentProcessingWorkflow,
		documentID,
	)
	return err
}
