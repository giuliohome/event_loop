See Temporal Go SDK [issue](https://github.com/temporalio/sdk-go/issues/1651).

Go sample adapted from [here](https://github.com/temporalio/money-transfer-project-template-go).

### Activity
```go
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
		return "error", fmt.Errorf("Serialization error: %v", err)
	}

	ex, err := os.Executable()
	cmd := exec.Command(ex, string(data))
	output, err := cmd.CombinedOutput() // Capture output from the subprocess
	if err != nil {
		return "error", fmt.Errorf("subprocess error: %v, output: %s", err, string(output))
	}
	fmt.Println("Subprocess output:", string(output))
	var result PaymentDetails
	if err := json.Unmarshal(output, &result); err != nil {
		return "error", fmt.Errorf("Error unmarshalling JSON: %v \nfrom output %s ", err, output)
	}

	return result.confirmation, err
}
```

### Separate SubProcess
```go
func WithdrawProcess(data_input string) () {
	var data PaymentDetails

	if err := json.Unmarshal([]byte(data_input), &data); err != nil {
		panic(fmt.Sprintf("Error unmarshalling JSON: %v", err))
	}

	referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID) + "_ProcessPID_" + strconv.Itoa(os.Getpid())
	bank := BankingService{"bank-api.example.com"}
	confirmation, err := bank.Withdraw(data.SourceAccount, data.Amount, referenceID)
	data.confirmation  = confirmation
	data.ReferenceID = referenceID
	if err != nil {
		panic(fmt.Sprintf("Withdrowal error %v", err))
	}

	// Serialize the data to JSON
	data_out, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Serialization failure %v", err)
	}
	fmt.Printf( string(data_out) )
}
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
