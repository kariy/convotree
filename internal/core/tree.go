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
	ID       uuid.UUID
	Exchange Exchange
	ParentID uuid.UUID
}

type Branch struct {
	Name string
	HEAD uuid.UUID // Checkpoint ID
}

type ConversationTree struct {
	checkpoints   map[uuid.UUID]Checkpoint
	branches      map[string]Branch
	currentBranch string
	mu            sync.RWMutex
}

func NewConversationTree() *ConversationTree {
	ct := &ConversationTree{
		checkpoints: make(map[uuid.UUID]Checkpoint),
		branches:    make(map[string]Branch),
	}

	ct.branches["main"] = Branch{Name: "main"}
	ct.currentBranch = "main"

	return ct
}

func (ct *ConversationTree) AddExchange(userInput, aiResponse string) uuid.UUID {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	id := uuid.New()
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

func (ct *ConversationTree) CreateBranch(branchName string, checkpointID uuid.UUID) error {
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

	for {
		if currentID == uuid.Nil {
			break
		}

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

func (ct *ConversationTree) GetBranch(branchName string) (*Branch, error) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	branch, exists := ct.branches[branchName]
	if !exists {
		return nil, fmt.Errorf("branch %s does not exist", branchName)
	}

	return &branch, nil
}
func (ct *ConversationTree) GetCurrentBranch() string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	return ct.currentBranch
}


func (ct *ConversationTree) GetCheckpoints() []uuid.UUID {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	checkpoints := make([]uuid.UUID, 0, len(ct.checkpoints))
	for id := range ct.checkpoints {
		checkpoints = append(checkpoints, id)
	}
	return checkpoints
}

