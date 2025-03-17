package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InstallmentSchedule struct {
	AccountID       string        `json:"accountId"`
	TotalAmount     float64       `json:"totalAmount"`
	DownPayment     float64       `json:"downPayment"`
	InstallmentDays []int         `json:"installmentDays"`
	StartDate       time.Time     `json:"startDate"`
	EndDate         time.Time     `json:"endDate"`
	Installments    []Installment `json:"installments"`
}

type Installment struct {
	Number  int       `json:"number"`
	Amount  float64   `json:"amount"`
	DueDate time.Time `json:"dueDate"`
	Status  string    `json:"status"`
}

type DownpaymentAllocation struct {
	AccountID         string    `json:"accountId"`
	DownpaymentAmount float64   `json:"downpaymentAmount"`
	PaymentMethod     string    `json:"paymentMethod"`
	TransactionID     string    `json:"transactionId"`
	AllocationDate    time.Time `json:"allocationDate"`
	Status            string    `json:"status"`
}

type AutopayEnrollment struct {
	AccountID      string    `json:"accountId"`
	EnrollmentDate time.Time `json:"enrollmentDate"`
	PaymentMethod  string    `json:"paymentMethod"`
	BankAccount    string    `json:"bankAccount"`
	Status         string    `json:"status"`
}

type EquityReview struct {
	AccountID        string    `json:"accountId"`
	EquityAmount     float64   `json:"equityAmount"`
	EquityDate       time.Time `json:"equityDate"`
	DueAmount        float64   `json:"dueAmount"`
	DueDate          time.Time `json:"dueDate"`
	ReviewStatus     string    `json:"reviewStatus"`
	InvoiceScheduled bool      `json:"invoiceScheduled"`
}

type NewAccount struct {
	AccountID    string    `json:"accountId"`
	CustomerID   string    `json:"customerId"`
	CreationDate time.Time `json:"creationDate"`
	AccountType  string    `json:"accountType"`
	Status       string    `json:"status"`
}

func main() {
	router := gin.Default()
	// router.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(http.StatusOK)
	// 		return
	// 	}
	// 	c.Next()
	// })
	nbus := router.Group("/api/nbus")
	{
		nbus.POST("/installment-schedule", createInstallmentSchedule)
		nbus.POST("/downpayment-allocation", allocateDownpayment)
		nbus.POST("/autopay-enrollment", handleAutopayEnrollment)
		nbus.POST("/equity-review", equityReview)
		nbus.POST("/new-account", newAccountCreated)
	}
	router.Run(":8080")
}

func createInstallmentSchedule(c *gin.Context) {
	var req struct {
		AccountID   string  `json:"accountId"`
		TotalAmount float64 `json:"totalAmount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	installments := make([]Installment, 6)
	for i := 0; i < 6; i++ {
		installments[i] = Installment{
			Number:  i + 1,
			Amount:  req.TotalAmount / 6,
			DueDate: now.AddDate(0, i, 0),
			Status:  "Scheduled",
		}
	}
	response := InstallmentSchedule{
		AccountID:       req.AccountID,
		TotalAmount:     req.TotalAmount,
		DownPayment:     req.TotalAmount * 0.1,
		InstallmentDays: []int{1, 1, 1, 1, 1, 1},
		StartDate:       now,
		EndDate:         now.AddDate(0, 6, 0),
		Installments:    installments,
	}
	c.JSON(http.StatusOK, response)
}

func allocateDownpayment(c *gin.Context) {
	var req struct {
		AccountID     string  `json:"accountId"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"paymentMethod"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := DownpaymentAllocation{
		AccountID:         req.AccountID,
		DownpaymentAmount: req.Amount,
		PaymentMethod:     req.PaymentMethod,
		TransactionID:     "TRX-" + time.Now().Format("20060102150405"),
		AllocationDate:    time.Now(),
		Status:            "Allocated",
	}
	c.JSON(http.StatusOK, response)
}

func handleAutopayEnrollment(c *gin.Context) {
	var req struct {
		AccountID     string `json:"accountId"`
		PaymentMethod string `json:"paymentMethod"`
		BankAccount   string `json:"bankAccount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := AutopayEnrollment{
		AccountID:      req.AccountID,
		EnrollmentDate: time.Now(),
		PaymentMethod:  req.PaymentMethod,
		BankAccount:    req.BankAccount,
		Status:         "Enrolled",
	}

	c.JSON(http.StatusOK, response)
}

func equityReview(c *gin.Context) {
	var req struct {
		AccountID string `json:"accountId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create sample response
	response := EquityReview{
		AccountID:        req.AccountID,
		EquityAmount:     15000.00,
		EquityDate:       time.Now(),
		DueAmount:        2500.00,
		DueDate:          time.Now().AddDate(0, 1, 0),
		ReviewStatus:     "Completed",
		InvoiceScheduled: true,
	}

	c.JSON(http.StatusOK, response)
}

func newAccountCreated(c *gin.Context) {
	var req struct {
		CustomerID  string `json:"customerId"`
		AccountType string `json:"accountType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := NewAccount{
		AccountID:    "ACCT-" + time.Now().Format("20060102150405"),
		CustomerID:   req.CustomerID,
		CreationDate: time.Now(),
		AccountType:  req.AccountType,
		Status:       "Active",
	}

	c.JSON(http.StatusOK, response)
}
