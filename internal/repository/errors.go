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
