package temporal

import (
	"go.temporal.io/sdk/client"
)

type TemporalClient struct {
	Client client.Client
	StartWorkflowOptions client.StartWorkflowOptions
}

func NewTemporalClient(id string, taskQueue string) (*TemporalClient, error) {
	startWorkflowOptions := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: taskQueue,
	}
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		return nil, err
	}

	return &TemporalClient{Client: c, StartWorkflowOptions: startWorkflowOptions}, nil
}
