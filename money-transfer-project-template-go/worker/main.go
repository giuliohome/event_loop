package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"money-transfer-project-template-go/app"
)

// Entry point - checks if we are in subprocess mode
func init() {
	// log.Printf("Init command %s", os.Args[0])
	if len(os.Args) == 4 && os.Args[1] == "WithdrawProcess" {
		// log.Printf("Init WithdrawProcess %s", os.Args[1])
		app.WithdrawProcess(os.Args[2], os.Args[3])
		os.Exit(0) // Ensure the subprocess exits after completing its work
	}
}

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, app.MoneyTransferTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(app.MoneyTransfer)
	w.RegisterActivity(app.Withdraw)
	w.RegisterActivity(app.Deposit)
	w.RegisterActivity(app.Refund)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

// @@@SNIPEND
