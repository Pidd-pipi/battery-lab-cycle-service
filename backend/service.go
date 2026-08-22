package main

type CycleService struct{ store *CellStore }

func NewCycleService(store *CellStore) *CycleService { return &CycleService{store: store} }
func (s *CycleService) Cells() []Cell                { return s.store.List() }
func (s *CycleService) ChangeStatus(id, status string) (Cell, error) {
	if err := ValidateCellStatus(status); err != nil {
		return Cell{}, err
	}
	return s.store.UpdateStatus(id, status)
}
