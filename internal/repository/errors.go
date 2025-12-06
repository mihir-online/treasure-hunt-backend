package repository

import "fmt"

// DuplicateChestError represents a duplicate chest_id error
type DuplicateChestError struct {
	ChestID string
}

func (e *DuplicateChestError) Error() string {
	return fmt.Sprintf("chest %s already has an owner", e.ChestID)
}

// IsDuplicateChestError checks if an error is a DuplicateChestError
func IsDuplicateChestError(err error) bool {
	_, ok := err.(*DuplicateChestError)
	return ok
}

// DuplicateExplorerError represents a duplicate explorer entry error (player already explored this
// chest)
type DuplicateExplorerError struct {
	ChestID  string
	PlayerID int
}

func (e *DuplicateExplorerError) Error() string {
	return fmt.Sprintf("player %d has already explored chest %s", e.PlayerID, e.ChestID)
}

// IsDuplicateExplorerError checks if an error is a DuplicateExplorerError
func IsDuplicateExplorerError(err error) bool {
	_, ok := err.(*DuplicateExplorerError)
	return ok
}
