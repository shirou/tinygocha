package game

// Army represents a collection of groups
type Army struct {
	ID     int
	Name   string
	Groups []*Group
	Side   int // 0: A軍, 1: B軍
}

// NewArmy creates a new army
func NewArmy(id int, name string, side int) *Army {
	return &Army{
		ID:     id,
		Name:   name,
		Groups: make([]*Group, 0),
		Side:   side,
	}
}

// AddGroup adds a group to the army
func (a *Army) AddGroup(group *Group) {
	a.Groups = append(a.Groups, group)
}

// Update updates all groups in the army
func (a *Army) Update(deltaTime float64) {
	for _, group := range a.Groups {
		group.Update(deltaTime)
	}
}

// GetAllUnits returns all units in the army
func (a *Army) GetAllUnits() []*Unit {
	var units []*Unit
	for _, group := range a.Groups {
		units = append(units, group.GetAllUnits()...)
	}
	return units
}

// GetAliveUnits returns all alive units in the army
func (a *Army) GetAliveUnits() []*Unit {
	var aliveUnits []*Unit
	for _, unit := range a.GetAllUnits() {
		if unit.IsAlive && !unit.IsRetreating {
			aliveUnits = append(aliveUnits, unit)
		}
	}
	return aliveUnits
}

// GetAliveCount returns the total number of alive units
func (a *Army) GetAliveCount() int {
	return len(a.GetAliveUnits())
}

// GetTotalHealth returns the total health percentage of the army
func (a *Army) GetTotalHealth() float64 {
	units := a.GetAllUnits()
	if len(units) == 0 {
		return 0
	}
	
	totalHealth := 0.0
	for _, unit := range units {
		totalHealth += unit.GetHealthPercentage()
	}
	
	return totalHealth / float64(len(units))
}

// IsDefeated returns true if the army is completely defeated
func (a *Army) IsDefeated() bool {
	for _, group := range a.Groups {
		if !group.IsDefeated() {
			return false
		}
	}
	return true
}

// GetActiveGroups returns groups that are not defeated
func (a *Army) GetActiveGroups() []*Group {
	var activeGroups []*Group
	for _, group := range a.Groups {
		if !group.IsDefeated() {
			activeGroups = append(activeGroups, group)
		}
	}
	return activeGroups
}
