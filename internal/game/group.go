package game

import (
	"math"

	gamemath "github.com/shirou/tinygocha/internal/math"
)

// FormationType represents different formation types
type FormationType int

const (
	CircleFormation FormationType = iota
	// Future: LineFormation, WedgeFormation, etc.
)

// Formation defines the formation parameters
type Formation struct {
	Type    FormationType
	Radius  float64
	Spacing float64
}

// Group represents a group of units with a leader
type Group struct {
	ID        int
	Leader    *Unit
	Members   []*Unit
	Formation Formation
	ArmyID    int
	
	// Formation state
	targetPosition gamemath.Vector2D
}

// NewGroup creates a new group
func NewGroup(id, armyID int, leader *Unit, members []*Unit) *Group {
	return &Group{
		ID:      id,
		Leader:  leader,
		Members: members,
		Formation: Formation{
			Type:    CircleFormation,
			Radius:  50.0,
			Spacing: 20.0,
		},
		ArmyID:         armyID,
		targetPosition: leader.Position,
	}
}

// Update updates the group and maintains formation
func (g *Group) Update(deltaTime float64) {
	if g.Leader == nil || !g.Leader.IsAlive {
		g.handleLeaderDeath()
		return
	}
	
	// Update leader first
	g.Leader.Update(deltaTime)
	
	// Update formation target based on leader position
	// リーダーが移動中の場合は目標位置、そうでなければ現在位置を使用
	if g.Leader.Position.Distance(g.Leader.Target) > 5.0 {
		g.targetPosition = g.Leader.Target
	} else {
		g.targetPosition = g.Leader.Position
	}
	
	// Update members and maintain formation
	g.updateFormation()
	
	// Update all members
	for _, member := range g.Members {
		if member.IsAlive {
			member.Update(deltaTime)
		}
	}
}

// updateFormation maintains the group's formation
func (g *Group) updateFormation() {
	if g.Leader == nil || !g.Leader.IsAlive {
		return
	}
	
	switch g.Formation.Type {
	case CircleFormation:
		g.updateCircleFormation()
	}
}

// updateCircleFormation arranges members in a circle around the leader
func (g *Group) updateCircleFormation() {
	aliveMembers := g.getAliveMembers()
	if len(aliveMembers) == 0 {
		return
	}
	
	angleStep := 2 * math.Pi / float64(len(aliveMembers))
	
	for i, member := range aliveMembers {
		if member.IsRetreating {
			continue
		}
		
		angle := float64(i) * angleStep
		offsetX := math.Cos(angle) * g.Formation.Radius
		offsetY := math.Sin(angle) * g.Formation.Radius
		
		formationPos := g.targetPosition.Add(gamemath.Vector2D{
			X: offsetX,
			Y: offsetY,
		})
		
		member.MoveTo(formationPos)
	}
}

// getAliveMembers returns all alive members
func (g *Group) getAliveMembers() []*Unit {
	var alive []*Unit
	for _, member := range g.Members {
		if member.IsAlive && !member.IsRetreating {
			alive = append(alive, member)
		}
	}
	return alive
}

// handleLeaderDeath handles the case when the leader dies
func (g *Group) handleLeaderDeath() {
	// Make all members retreat
	for _, member := range g.Members {
		if member.IsAlive && !member.IsRetreating {
			// Set retreat target to screen edge (simplified)
			exitPoint := gamemath.Vector2D{X: -100, Y: member.Position.Y}
			if member.ArmyID == 1 { // Army B retreats to right
				exitPoint.X = 1124 // Screen width + 100
			}
			member.StartRetreating(exitPoint)
		}
	}
}

// MoveGroup moves the entire group to a new position
func (g *Group) MoveGroup(target gamemath.Vector2D) {
	if g.Leader != nil && g.Leader.IsAlive {
		g.Leader.MoveTo(target)
	}
}

// GetAllUnits returns all units in the group (leader + members)
func (g *Group) GetAllUnits() []*Unit {
	units := []*Unit{}
	if g.Leader != nil {
		units = append(units, g.Leader)
	}
	units = append(units, g.Members...)
	return units
}

// GetAliveCount returns the number of alive units in the group
func (g *Group) GetAliveCount() int {
	count := 0
	if g.Leader != nil && g.Leader.IsAlive {
		count++
	}
	for _, member := range g.Members {
		if member.IsAlive {
			count++
		}
	}
	return count
}

// IsDefeated returns true if the group is completely defeated
func (g *Group) IsDefeated() bool {
	return g.GetAliveCount() == 0
}
