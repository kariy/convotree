package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"convotree/internal/ai"
	"convotree/internal/core"
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

		case "branch":
			fmt.Print("Enter a name for the new branch: ")
			scanner.Scan()
			newBranchName := scanner.Text()

			fmt.Print("Enter the name of the branch to branch from (default: current branch): ")
			scanner.Scan()
			fromBranch := scanner.Text()
			if fromBranch == "" {
				fromBranch = ct.GetCurrentBranch()
			}

			err := ct.CreateBranch(newBranchName, fromBranch)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Created branch '%s' from '%s'\n", newBranchName, fromBranch)
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
