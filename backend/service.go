package main

type CycleService struct{ store *CellStore }

func NewCycleService(store *CellStore) *CycleService { return &CycleService{store: store} }
func (s *CycleService) Cells() []Cell                { return s.store.List() }
func (s *CycleService) ChangeStatus(id, status string) (Cell, error) {
	if err := ValidateCellStatus(status); err != nil {
		return Cell{}, err
	}
	cell, err := s.store.UpdateStatus(id, status)
	if err != nil {
		return Cell{}, err
	}
	if status == "complete" {
		cell.Cycle++
	}
	return cell, nil
}
