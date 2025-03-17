package main

import (
	"temporalPoc/starter"
	"temporalPoc/worker"
	"temporalPoc/workflows"
	"time"
)

func main() {
	go starter.NbusStarter()
	go worker.NbusWorker()
	workflows.UpdateRunningWorkflow(workflows.WorkflowInput{
		AccountID:     "test",
		CustomerID:    "testcust",
		BankAccount:   "test",
		TotalAmount:   5000,
		PaymentMethod: "test",
		AccountType:   "test",
	})
	time.Sleep(time.Hour * 1)
}
