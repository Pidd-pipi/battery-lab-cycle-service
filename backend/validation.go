package main

import "fmt"

var cellStatuses = map[string]bool{"cycling": true, "paused": true, "complete": true}

func ValidateCellStatus(status string) error {
	if !cellStatuses[status] {
		return fmt.Errorf("status must be cycling, paused, or complete")
	}
	return nil
}
