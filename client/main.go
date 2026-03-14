package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AntonioMartinezLopez/enginsight/jrpc"
)

func main() {
	opts := jrpc.ClientOptions{}

	client, err := jrpc.NewClient("http://localhost:8080/rpc", opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter a line: ")
		if !scanner.Scan() {
			break
		}
		message := strings.TrimSpace(scanner.Text())

		result, err := client.Count(ctx, message)

		if err != nil {
			log.Fatalf("Count request failed: %v", err)
		}

		fmt.Printf("Message: %q, Character count: %d\n", message, result)
	}
}
