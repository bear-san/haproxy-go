package main

import (
	"errors"
	"fmt"
	"log"
	
	v3 "github.com/bear-san/haproxy-go/dataplane/v3"
)

func main() {
	// Create a client
	client := v3.Client{
		BaseUrl:    "http://localhost:5555",
		Credential: "YWRtaW46YWRtaW4=", // base64 encoded "admin:admin"
	}

	// Example 1: Handling specific error types
	backends, err := client.ListBackends("non-existent-transaction")
	if err != nil {
		switch {
		case v3.IsNotFound(err):
			fmt.Println("Resource not found:", err)
		case v3.IsUnauthorized(err):
			fmt.Println("Authentication failed:", err)
		case v3.IsConflict(err):
			fmt.Println("Conflict occurred:", err)
		case v3.IsUnknownError(err):
			// Get HTTP status code for unknown errors
			statusCode := v3.GetHTTPStatusCode(err)
			fmt.Printf("Unknown error with status %d: %v\n", statusCode, err)
		default:
			fmt.Println("Other error:", err)
		}
		return
	}

	// Example 2: Transaction commit with error handling
	version, err := client.GetVersion()
	if err != nil {
		log.Fatal("Failed to get version:", err)
	}

	transaction, err := client.CreateTransaction(*version)
	if err != nil {
		log.Fatal("Failed to create transaction:", err)
	}

	// ... do some operations ...

	// Commit transaction
	_, err = client.CommitTransaction(*transaction.Id)
	if err != nil {
		if v3.IsCommitFailed(err) {
			fmt.Println("Transaction commit failed:", err)
			// You can extract the transaction ID from CommitFailedError
			var commitErr *v3.CommitFailedError
			if errors.As(err, &commitErr) {
				fmt.Printf("Failed transaction ID: %s\n", commitErr.TransactionID)
			}
		} else {
			fmt.Println("Other error during commit:", err)
		}
		return
	}

	// Example 3: Type assertions for detailed error information
	backend, err := client.GetBackend("non-existent", "some-transaction")
	if err != nil {
		// Type assertion to get detailed error info
		var unknownErr *v3.UnknownError
		if errors.As(err, &unknownErr) {
			fmt.Printf("HTTP Status: %d\n", unknownErr.StatusCode)
			fmt.Printf("Error Message: %s\n", unknownErr.Message)
		}
	}

	_ = backends
	_ = backend
}