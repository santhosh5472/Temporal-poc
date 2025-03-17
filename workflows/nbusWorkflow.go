package workflows

import (
	"temporalPoc/activities"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowInput struct {
	AccountID     string  `json:"accountId"`
	CustomerID    string  `json:"customerId"`
	TotalAmount   float64 `json:"totalAmount"`
	PaymentMethod string  `json:"paymentMethod"`
	BankAccount   string  `json:"bankAccount"`
	AccountType   string  `json:"accountType"`
}

func NBUSPostSetupWorkflow(ctx workflow.Context, input WorkflowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("NBUS Post-setup Workflow started", "AccountID", input.AccountID)

	signalChan := workflow.GetSignalChannel(ctx, "update_workflow_input")
	workflow.Go(ctx, func(ctx workflow.Context) {
        for {
            var newInput WorkflowInput
            if !signalChan.Receive(ctx, &newInput) {
                return 
            }
            logger.Info("Received updated input via signal", 
                "AccountID", newInput.AccountID,
                "TotalAmount", newInput.TotalAmount,
			)
            input = newInput
        }
    })

	syncRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 1,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute * 1,
		MaximumAttempts:    5,
	}

	asyncRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 5,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute * 10,
		MaximumAttempts:    10,
	}

	ctx1 := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy:         syncRetryPolicy,
	})
	err := workflow.ExecuteActivity(ctx1, activities.CreateInstallmentScheduleActivity, input.AccountID, input.TotalAmount).Get(ctx1, nil)
	if err != nil {
		logger.Error("CreateInstallmentSchedule failed", "Error", err)
		return err
	}
	logger.Info("CreateInstallmentSchedule completed")


	downpaymentAmount := input.TotalAmount * 0.1
	err = workflow.ExecuteActivity(ctx1, activities.AllocateDownpaymentActivity, input.AccountID, downpaymentAmount, input.PaymentMethod).Get(ctx1, nil)
	if err != nil {
		logger.Error("AllocateDownpayment failed", "Error", err)
		return err
	}
	logger.Info("AllocateDownpayment completed")

	ctx2 := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
		RetryPolicy:         asyncRetryPolicy,
	})

	autopayFuture := workflow.ExecuteActivity(ctx2, activities.HandleAutopayEnrollmentActivity,
		input.AccountID, input.PaymentMethod, input.BankAccount)

	equityReviewFuture := workflow.ExecuteActivity(ctx2, activities.EquityReviewActivity, input.AccountID)

	newAccountCreatedFuture := workflow.ExecuteActivity(ctx2, activities.NewAccountCreatedActivity,
		input.CustomerID, input.AccountType)

	err = autopayFuture.Get(ctx2, nil)
	if err != nil {
		logger.Error("HandleAutopayEnrollment failed", "Error", err)
	} else {
		logger.Info("HandleAutopayEnrollment completed successfully")
	}

	err = equityReviewFuture.Get(ctx2, nil)
	if err != nil {
		logger.Error("EquityReview failed", "Error", err)
	} else {
		logger.Info("EquityReview completed successfully")
	}

	err = newAccountCreatedFuture.Get(ctx2, nil)
	if err != nil {
		logger.Error("NewAccountCreated failed", "Error", err)
	} else {
		logger.Info("NewAccountCreated completed successfully")
	}

	logger.Info("NBUS Post-setup Workflow completed")
	return nil
}
