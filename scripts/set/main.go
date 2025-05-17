package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// nolint
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <key> <value>")
		os.Exit(1)
	}

	key := os.Args[1]
	value := os.Args[2]

	payload := map[string]string{
		"key":   key,
		"value": value,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error creating JSON payload: %v\n", err)
		os.Exit(1)
	}

	url := "http://localhost:8080/v1/write"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Error: %s (Status code: %d)\n", string(body), resp.StatusCode)
		os.Exit(1)
	}

	fmt.Printf("Successfully set %s=%s\n", key, value)
}
