package starter

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"temporalPoc/workflows" // Import your workflow package
)

func NbusStarter() {
	// Create the client object
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	// Start the workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "nbus-post-setup-workflow-1", // Consider making this dynamic
		TaskQueue: "nbus-task-queue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.NBUSPostSetupWorkflow)
	if err != nil {
		log.Fatalln("Unable to start workflow:", err)
	}

	log.Println("Started workflow", "WorkflowID:", we.GetID(), "RunID:", we.GetRunID())

	// Wait for workflow completion (optional)
	var result interface{}
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Workflow failed:", err)
	}
	log.Println("Workflow completed successfully")
}
