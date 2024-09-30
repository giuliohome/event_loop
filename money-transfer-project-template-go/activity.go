package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// @@@SNIPSTART money-transfer-project-template-go-activity-withdraw
func Withdraw(ctx context.Context, jsonData PaymentDetails) (string, error) {
	log.Printf("Withdrawing $%d from account %s running on process PID: %dn.\n\n",
		jsonData.Amount,
		jsonData.SourceAccount,
		os.Getpid(),
	)

	// Serialize the data to JSON
	data, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Serialization failure")
		return "error", fmt.Errorf("serialization error: %v", err)
	}

	ex, err := os.Executable()
	if err != nil {
		return "error", fmt.Errorf("executable error: %v", err)
	}

	// Write the JSON data to a temporary file
	tempFileIn, err := os.CreateTemp("", "data.json")
	if err != nil {
		return "error", fmt.Errorf("Error creating temp file: %v", err)

	}
	defer os.Remove(tempFileIn.Name()) // Clean up

	if err := os.WriteFile(tempFileIn.Name(), data, 0644); err != nil {
		return "error", fmt.Errorf("Error writing temp file: %v", err)
	}

	tempFileOut, err := os.CreateTemp("", "data.json")
	if err != nil {
		return "error", fmt.Errorf("Error creating temp file: %v", err)

	}
	defer os.Remove(tempFileOut.Name()) // Clean up

	cmd := exec.Command(ex, "WithdrawProcess", tempFileIn.Name(), tempFileOut.Name())
	output, err := cmd.CombinedOutput() // Capture output from the subprocess
	if err != nil {
		return "error", fmt.Errorf("subprocess error: %v, output: %s", err, string(output))
	}
	fmt.Println("Subprocess output:", string(output))

	data_output, err := os.ReadFile(tempFileOut.Name())
	if err != nil {
		return "error", fmt.Errorf("Error reading temp file: %v", err)
	}

	var result PaymentDetails
	if err := json.Unmarshal(data_output, &result); err != nil {
		return "error", fmt.Errorf("unmarshalling JSON: %v \nfrom output %s ", err, output)
	}
	fmt.Printf("Result unmarshalled: %v confirmation %s", result, result.Confirmation)

	return result.Confirmation, err
}

func WithdrawProcess(data_input_file string, data_output_file string) {
	var data PaymentDetails

	fmt.Printf("I'm in the subprocess with PID %d\n", os.Getpid())

	// Read the JSON data from the file
	data_input, err := os.ReadFile(data_input_file)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %v", err))
	}

	if err := json.Unmarshal([]byte(data_input), &data); err != nil {
		panic(fmt.Sprintf("Error unmarshalling JSON: %v", err))
	}

	referenceID := fmt.Sprintf("%s-withdrawal_ProcessPID_%d", data.ReferenceID, os.Getpid())
	data.ReferenceID = referenceID
	bank := BankingService{"bank-api.example.com"}
	confirmation, err := bank.Withdraw(data.SourceAccount, data.Amount, referenceID)
	data.Confirmation = "*" + confirmation + "*"
	if err != nil {
		panic(fmt.Sprintf("Withdrowal error %v", err))
	}

	// Serialize the data to JSON
	data_out, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Serialization failure %v", err)
	}
	fmt.Println("writing output " + string(data_out) + " to " + data_output_file)
	if err := os.WriteFile(data_output_file, data_out, 0644); err != nil {
		fmt.Printf("Error writing temp file: %v", err)
	}
}

// @@@SNIPEND

// @@@SNIPSTART money-transfer-project-template-go-activity-deposit
func Deposit(ctx context.Context, data PaymentDetails) (string, error) {
	log.Printf("Depositing $%d into account %s running on process PID: %dn.\n\n",
		data.Amount,
		data.TargetAccount,
		os.Getpid(),
	)

	referenceID := fmt.Sprintf("%s-deposit", data.ReferenceID)
	bank := BankingService{"bank-api.example.com"}
	// Uncomment the next line and comment the one after that to simulate an unknown failure
	// confirmation, err := bank.DepositThatFails(data.TargetAccount, data.Amount, referenceID)
	confirmation, err := bank.Deposit(data.TargetAccount, data.Amount, referenceID)
	return confirmation, err
}

// @@@SNIPEND

// @@@SNIPSTART money-transfer-project-template-go-activity-refund
func Refund(ctx context.Context, data PaymentDetails) (string, error) {
	log.Printf("Refunding $%v back into account %v.\n\n",
		data.Amount,
		data.SourceAccount,
	)

	referenceID := fmt.Sprintf("%s-refund", data.ReferenceID)
	bank := BankingService{"bank-api.example.com"}
	confirmation, err := bank.Deposit(data.SourceAccount, data.Amount, referenceID)
	return confirmation, err
}

// @@@SNIPEND
