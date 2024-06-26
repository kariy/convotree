package core

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
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