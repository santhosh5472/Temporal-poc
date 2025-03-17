package activities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBaseURL = "http://localhost:8080/api/nbus"

func CreateInstallmentScheduleActivity(accountID string, totalAmount float64) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"accountId":   accountID,
		"totalAmount": totalAmount,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		apiBaseURL+"/installment-schedule",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("error calling installment schedule API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error parsing API response: %w", err)
	}

	fmt.Printf("Installment Schedule created successfully: %v\n", response)
	return nil
}

func AllocateDownpaymentActivity(accountID string, amount float64, paymentMethod string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"accountId":     accountID,
		"amount":        amount,
		"paymentMethod": paymentMethod,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		apiBaseURL+"/downpayment-allocation",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("error calling downpayment allocation API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Println("Downpayment allocated successfully")
	return nil
}

func HandleAutopayEnrollmentActivity(accountID string, paymentMethod string, bankAccount string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"accountId":     accountID,
		"paymentMethod": paymentMethod,
		"bankAccount":   bankAccount,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}
	resp, err := http.Post(
		apiBaseURL+"/autopay-enrollment",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("error calling autopay enrollment API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Println("Autopay enrollment handled successfully")
	return nil
}

func EquityReviewActivity(accountID string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"accountId": accountID,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		apiBaseURL+"/equity-review",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("error calling equity review API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Println("Equity review completed successfully")
	return nil
}

func NewAccountCreatedActivity(customerID string, accountType string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"customerId":  customerID,
		"accountType": accountType,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}
	resp, err := http.Post(
		apiBaseURL+"/new-account",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return fmt.Errorf("error calling new account API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Println("New account created successfully")
	return nil
}
