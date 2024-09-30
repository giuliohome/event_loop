See Temporal Go SDK [issue](https://github.com/temporalio/sdk-go/issues/1651).

Go sample adapted from [here](https://github.com/temporalio/money-transfer-project-template-go).

![immagine](https://github.com/user-attachments/assets/7a1741dd-948e-45bb-9975-f7a7df0b497f)

![immagine](https://github.com/user-attachments/assets/1d9c8f83-0305-45c4-8c1b-3cef04fddec2)



### Activity
```go
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
	cmd := exec.Command(ex, string(data))
	output, err := cmd.CombinedOutput() // Capture output from the subprocess
	if err != nil {
		return "error", fmt.Errorf("subprocess error: %v, output: %s", err, string(output))
	}
	fmt.Println("Subprocess output:", string(output))
	var result PaymentDetails
	if err := json.Unmarshal(output, &result); err != nil {
		return "error", fmt.Errorf("unmarshalling JSON: %v \nfrom output %s ", err, output)
	}
	fmt.Printf("Result unmarshalled: %v confirmation %s", result, result.Confirmation)

	return result.Confirmation, err
}
```

### Separate SubProcess
```go
func WithdrawProcess(data_input string) {
	var data PaymentDetails

	if err := json.Unmarshal([]byte(data_input), &data); err != nil {
		panic(fmt.Sprintf("Error unmarshalling JSON: %v", err))
	}
	/* do not write to output except for the final result
	)*/

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
	fmt.Println(string(data_out))
}
// @@@SNIPEND
```

### SubProcess Executor
```go
// Entry point - checks if we are in subprocess mode
func init() {
	if len(os.Args) > 1  {
		app.WithdrawProcess(os.Args[1])
		os.Exit(0) // Ensure the subprocess exits after completing its work
	}
}

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {
```

![immagine](https://github.com/user-attachments/assets/97c4c2e3-b6fa-4224-90a0-505a9dcc2875)

