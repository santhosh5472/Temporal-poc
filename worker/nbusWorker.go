package worker

import (
	"log"
	"temporalPoc/activities"
	"temporalPoc/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func NbusWorker() {
	// Create the client object
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	// Create a worker
	w := worker.New(c, "nbus-task-queue", worker.Options{})

	// Register workflow and activities
	w.RegisterWorkflow(workflows.NBUSPostSetupWorkflow)

	// Register all activities
	w.RegisterActivity(activities.CreateInstallmentScheduleActivity)
	w.RegisterActivity(activities.AllocateDownpaymentActivity)
	w.RegisterActivity(activities.HandleAutopayEnrollmentActivity)
	w.RegisterActivity(activities.EquityReviewActivity)
	w.RegisterActivity(activities.NewAccountCreatedActivity)

	// Start the worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker:", err)
	}
}
