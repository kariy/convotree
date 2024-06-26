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
}

type Branch struct {
	Name string
	HEAD string // Checkpoint ID
}

type ConversationTree struct {
	checkpoints     map[string]Checkpoint
	branches        map[string]Branch
	currentBranch   string
	mu              sync.RWMutex
}

func NewConversationTree() *ConversationTree {
	ct := &ConversationTree{
		checkpoints: make(map[string]Checkpoint),
		branches:    make(map[string]Branch),
	}
	// Create initial checkpoint and "main" branch
	initialID := uuid.New().String()
	ct.checkpoints[initialID] = Checkpoint{ID: initialID, ParentID: ""}
	ct.branches["main"] = Branch{Name: "main", HEAD: initialID}
	ct.currentBranch = "main"
	return ct
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
		ParentID: ct.branches[ct.currentBranch].HEAD,
	}

	ct.checkpoints[id] = checkpoint
	
	// Update the HEAD of the current branch
	branch := ct.branches[ct.currentBranch]
	branch.HEAD = id
	ct.branches[ct.currentBranch] = branch

	return id
}

func (ct *ConversationTree) CreateBranch(branchName, checkpointID string) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if _, exists := ct.branches[branchName]; exists {
		return fmt.Errorf("branch %s already exists", branchName)
	}

	if _, exists := ct.checkpoints[checkpointID]; !exists {
		return fmt.Errorf("checkpoint %s does not exist", checkpointID)
	}

	ct.branches[branchName] = Branch{
		Name: branchName,
		HEAD: checkpointID,
	}

	return nil
}

func (ct *ConversationTree) SwitchBranch(branchName string) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if _, exists := ct.branches[branchName]; !exists {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	ct.currentBranch = branchName
	return nil
}

func (ct *ConversationTree) GetConversationHistory() []Exchange {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	var history []Exchange
	currentID := ct.branches[ct.currentBranch].HEAD

	for currentID != "" {
		checkpoint := ct.checkpoints[currentID]
		history = append([]Exchange{checkpoint.Exchange}, history...)
		currentID = checkpoint.ParentID
	}

	return history
}

func (ct *ConversationTree) GetBranchNames() []string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	names := make([]string, 0, len(ct.branches))
	for name := range ct.branches {
		names = append(names, name)
	}
	return names
}

func (ct *ConversationTree) GetCurrentBranch() string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	return ct.currentBranch
}

func (ct *ConversationTree) GetCheckpoints() []string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	checkpoints := make([]string, 0, len(ct.checkpoints))
	for id := range ct.checkpoints {
		checkpoints = append(checkpoints, id)
	}
	return checkpoints
}

// TODO: This should be replaced with a real AI service
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
		fmt.Printf("Current branch: %s\n", ct.GetCurrentBranch())
		fmt.Print("Enter command (chat/branch/switch/list/quit): ")
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