package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"convotree/internal/ai"
	"convotree/internal/core"

	"github.com/google/uuid"
)

func main() {
	ct := core.NewConversationTree()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Current branch: %s\n", ct.GetCurrentBranch())
		fmt.Print("Enter command (chat/branch/switch/list/quit): ")
		scanner.Scan()
		command := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch command {
		case "chat":
			fmt.Print("You: ")
			scanner.Scan()
			userInput := scanner.Text()

			aiResponse, err := ai.GetAIResponse(ct.GetConversationHistory())
			if err != nil {
				fmt.Printf("Error getting AI response: %v\n", err)
				continue
			}

			fmt.Printf("AI: %s\n", aiResponse)
			newCheckpointID := ct.AddExchange(userInput, aiResponse)
			fmt.Printf("New checkpoint created: %s\n", newCheckpointID)

			fmt.Print("Enter a name for the new branch: ")
			scanner.Scan()
			newBranchName := scanner.Text()

			var id uuid.UUID
			fmt.Print("Enter the checkpoint to branch from (leave blank for current HEAD): ")
			scanner.Scan()
			fromBranchInput := scanner.Text()

			if fromBranchInput == "" {
				branch, err := ct.GetBranch(ct.GetCurrentBranch())
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}
				id = branch.HEAD
			} else {
				var err error
				id, err = uuid.Parse(fromBranchInput)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}
			}
			err1 := ct.CreateBranch(newBranchName, id)

			if err1 != nil {
				fmt.Printf("Error: %v\n", err1)
			} else {
				fmt.Printf("Created branch '%s' from checkpoint '%s'\n", newBranchName, id)
			}

		case "switch":
			fmt.Print("Enter branch name to switch to: ")
			scanner.Scan()
			branchName := scanner.Text()

			err := ct.SwitchBranch(branchName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Switched to branch: %s\n", branchName)
			}

		case "list":
			branches := ct.GetBranchNames()
			fmt.Println("Available branches:")
			for _, branch := range branches {
				fmt.Println(branch)
			}

		case "quit":
			return

		default:
			fmt.Println("Invalid command")
		}
	}
}
