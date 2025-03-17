package workflows

import (
	"context"

	"go.temporal.io/sdk/client"
)

// UpdateRunningWorkflow sends updated input to a running workflow
func UpdateRunningWorkflow(newInput WorkflowInput) error {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	// Send signal to the running workflow
	err = c.SignalWorkflow(context.Background(),
		"nbus-post-setup-workflow-1", // The workflow ID
		"",                           // Empty means latest run
		"update_workflow_input",      // Signal name
		newInput)                     // Signal payload

	if err != nil {
		return err
	}

	return nil
}
