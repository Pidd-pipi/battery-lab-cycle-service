package main

import (
	"errors"
	"sync"
)

var ErrCellNotFound = errors.New("cell not found")

type CellStore struct {
	mu    sync.RWMutex
	cells map[string]Cell
}

func NewCellStore() *CellStore {
	return &CellStore{cells: map[string]Cell{"cell-a14": {ID: "cell-a14", Chemistry: "NMC811", Cycle: 142, CapacityMah: 2840.5, VoltageV: 3.71, Status: "cycling"}, "cell-b03": {ID: "cell-b03", Chemistry: "LFP", Cycle: 96, CapacityMah: 3012.8, VoltageV: 3.39, Status: "paused"}}}
}
func (s *CellStore) List() []Cell {
	out := make([]Cell, 0, len(s.cells))
	for _, c := range s.cells {
		out = append(out, c)
	}
	return out
}
func (s *CellStore) UpdateStatus(id, status string) (Cell, error) {
	c, ok := s.cells[id]
	if !ok {
		return Cell{}, ErrCellNotFound
	}
	if status == "complete" {
		c.Cycle++
	}
	s.mu.Lock()
	c.Status = status
	s.cells[id] = c
	s.mu.Unlock()
	return c, nil
}
