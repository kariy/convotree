package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

type Exchange struct {
	UserInput  string
	AIResponse string
}

type Checkpoint struct {
	ID       string
	Exchange Exchange
	ParentID string
	BranchName string  // New field for branch name
}

type ConversationTree struct {
	checkpoints       map[string]Checkpoint
	currentCheckpoint string
	branchNames       map[string]string  // Map of branch names to checkpoint IDs
	mu                sync.RWMutex
}

func NewConversationTree() *ConversationTree {
	return &ConversationTree{
		checkpoints: make(map[string]Checkpoint),
		branchNames: make(map[string]string),
	}
}

func (ct *ConversationTree) AddExchange(userInput, aiResponse string) string {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	id := uuid.New().String()
	checkpoint := Checkpoint{
		ID: id,
		Exchange: Exchange{
			UserInput:  userInput,
			AIResponse: aiResponse,
		},
		ParentID: ct.currentCheckpoint,
	}

	ct.checkpoints[id] = checkpoint
	ct.currentCheckpoint = id
	return id
}

func (ct *ConversationTree) CreateBranch(fromCheckpointID string, branchName string) (string, error) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if _, exists := ct.checkpoints[fromCheckpointID]; !exists {
		return "", fmt.Errorf("invalid checkpoint ID")
	}

	if _, exists := ct.branchNames[branchName]; exists {
		return "", fmt.Errorf("branch name already exists")
	}

	ct.currentCheckpoint = fromCheckpointID
	ct.branchNames[branchName] = fromCheckpointID

	// Update the checkpoint with the branch name
	checkpoint := ct.checkpoints[fromCheckpointID]
	checkpoint.BranchName = branchName
	ct.checkpoints[fromCheckpointID] = checkpoint

	return fromCheckpointID, nil
}

func (ct *ConversationTree) SwitchBranch(branchName string) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	checkpointID, exists := ct.branchNames[branchName]
	if !exists {
		return fmt.Errorf("branch name not found")
	}

	ct.currentCheckpoint = checkpointID
	return nil
}

func (ct *ConversationTree) GetConversationHistory() []Exchange {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	var history []Exchange
	currentID := ct.currentCheckpoint

	for currentID != "" {
		checkpoint := ct.checkpoints[currentID]
		history = append([]Exchange{checkpoint.Exchange}, history...)
		currentID = checkpoint.ParentID
	}

	return history
}

func (ct *ConversationTree) GetCheckpointIDs() []string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	ids := make([]string, 0, len(ct.checkpoints))
	for id := range ct.checkpoints {
		ids = append(ids, id)
	}
	return ids
}

func (ct *ConversationTree) GetBranchNames() []string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	names := make([]string, 0, len(ct.branchNames))
	for name := range ct.branchNames {
		names = append(names, name)
	}
	return names
}

// Main function and CLI implementation would go here

func getAIResponse(_ []Exchange) (string, error) {
	// Generate a random string of up to 10 words
	words := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon"}
	n := rand.Intn(10) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(words[rand.Intn(len(words))])
		sb.WriteString(" ")
	}
	return strings.TrimSpace(sb.String()), nil
}

func main() {
	ct := NewConversationTree()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter command (chat/branch/switch/list/history/quit): ")
		scanner.Scan()
		command := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch command {
		case "chat":
			fmt.Print("You: ")
			scanner.Scan()
			userInput := scanner.Text()

			aiResponse, err := getAIResponse(ct.GetConversationHistory())
			if err != nil {
				fmt.Printf("Error getting AI response: %v\n", err)
				continue
			}

			fmt.Printf("AI: %s\n", aiResponse)
			newCheckpointID := ct.AddExchange(userInput, aiResponse)
			fmt.Printf("New checkpoint created: %s\n", newCheckpointID)

		case "branch":
			fmt.Print("Enter checkpoint ID to branch from: ")
			scanner.Scan()
			checkpointID := scanner.Text()

			fmt.Print("Enter a name for the new branch: ")
			scanner.Scan()
			branchName := scanner.Text()

			branchedID, err := ct.CreateBranch(checkpointID, branchName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Created branch '%s' from checkpoint: %s\n", branchName, branchedID)
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
			checkpoints := ct.GetCheckpointIDs()
			fmt.Println("Available checkpoints:")
			for _, checkpoint := range checkpoints {
				fmt.Println(checkpoint)
			}

			branches := ct.GetBranchNames()
			fmt.Println("Available branches:")
			for _, branch := range branches {
				fmt.Println(branch)
			}

		case "history":
			history := ct.GetConversationHistory()
			fmt.Println("Conversation history:")
			for _, exchange := range history {
				fmt.Printf("You: %s\nAI: %s\n\n", exchange.UserInput, exchange.AIResponse)
			}

		case "quit":
			return

		default:
			fmt.Println("Invalid command")
		}
	}
}