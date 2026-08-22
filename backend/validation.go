package main

import (
	"errors"
	"fmt"
)

var cellStatuses = map[string]bool{"cycling": true, "paused": true, "complete": true}

var ErrCellStatusInvalid = errors.New("invalid cell status")

func ValidateCellStatus(status string) error {
	if !cellStatuses[status] {
		return fmt.Errorf("%w: %s", ErrCellStatusInvalid, status)
	}
	return nil
}
