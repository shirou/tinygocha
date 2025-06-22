package game

import (
	"fmt"
	
	"github.com/shirou/tinygocha/internal/graphics"
	"github.com/shirou/tinygocha/internal/math"
)

// UnitType represents different types of units
type UnitType string

const (
	UnitTypeInfantry UnitType = "infantry"
	UnitTypeArcher   UnitType = "archer"
	UnitTypeMage     UnitType = "mage"
)

// Unit represents an individual unit in the game
type Unit struct {
	ID           int
	Type         UnitType
	Name         string
	HP           int
	MaxHP        int
	AttackPower  int
	Defense      int
	Speed        float64
	Range        float64
	MagicPower   int
	Size         float64  // ユニットの大きさ（衝突判定用）
	Position     math.Vector2D
	Target       math.Vector2D
	IsLeader     bool
	IsAlive      bool
	IsRetreating bool
	GroupID      int
	ArmyID       int
	
	// Combat state
	LastAttackTime float64
	AttackCooldown float64
	
	// Animation state
	Animation *graphics.AnimationState
	
	// AI behavior
	AI *AIBehavior
}

// NewUnit creates a new unit with the given configuration
func NewUnit(id int, unitType UnitType, config UnitTypeConfig, isLeader bool, groupID, armyID int) *Unit {
	unit := &Unit{
		ID:             id,
		Type:           unitType,
		Name:           config.Name,
		HP:             config.HP,
		MaxHP:          config.HP,
		AttackPower:    config.Attack,
		Defense:        config.Defense,
		Speed:          config.Speed,
		Range:          config.Range,
		MagicPower:     config.MagicPower,
		Size:           config.Size,  // サイズを設定
		Position:       math.Vector2D{},
		Target:         math.Vector2D{},
		IsLeader:       isLeader,
		IsAlive:        true,
		IsRetreating:   false,
		GroupID:        groupID,
		ArmyID:         armyID,
		LastAttackTime: 0,
		AttackCooldown: 1.0, // 1 second cooldown
		Animation:      graphics.NewAnimationState(graphics.AnimationIdle),
		AI:             NewAIBehavior(unitType),
	}
	
	// デバッグ: ユニット作成確認
	fmt.Printf("Created Unit ID=%d, Type=%s, HP=%d/%d, Alive=%t, Army=%d, Size=%.1f\n", 
		unit.ID, unit.Type, unit.HP, unit.MaxHP, unit.IsAlive, unit.ArmyID, unit.Size)
	
	return unit
}

// Update updates the unit's state
func (u *Unit) Update(deltaTime float64) {
	if !u.IsAlive {
		// Set death animation if not already set
		if u.Animation.Type != graphics.AnimationDeath {
			u.Animation.SetAnimation(graphics.AnimationDeath)
		}
		u.Animation.Update(deltaTime)
		return
	}
	
	// Update attack cooldown
	if u.LastAttackTime > 0 {
		u.LastAttackTime -= deltaTime
		if u.LastAttackTime < 0 {
			u.LastAttackTime = 0
		}
	}
	
	// Determine animation based on state
	isMoving := u.Position.Distance(u.Target) > u.GetCollisionRadius()  // 衝突半径を考慮した移動判定
	
	if u.LastAttackTime > u.AttackCooldown * 0.7 { // Recently attacked
		if u.Animation.Type != graphics.AnimationAttack {
			u.Animation.SetAnimation(graphics.AnimationAttack)
		}
	} else if isMoving {
		if u.Animation.Type != graphics.AnimationWalk {
			u.Animation.SetAnimation(graphics.AnimationWalk)
		}
	} else {
		if u.Animation.Type != graphics.AnimationIdle {
			u.Animation.SetAnimation(graphics.AnimationIdle)
		}
	}
	
	// Update animation
	u.Animation.Update(deltaTime)
	
	// Move towards target if not at target
	if isMoving {
		direction := u.Target.Sub(u.Position).Normalize()
		movement := direction.Mul(u.Speed * deltaTime)
		u.Position = u.Position.Add(movement)
	}
}

// MoveTo sets the unit's target position
func (u *Unit) MoveTo(target math.Vector2D) {
	u.Target = target
}

// CanAttack checks if the unit can attack
func (u *Unit) CanAttack() bool {
	return u.IsAlive && u.LastAttackTime <= 0
}

// Attack performs an attack on the target unit
func (u *Unit) Attack(target *Unit) int {
	if !u.CanAttack() || !target.IsAlive {
		return 0
	}
	
	// Check range (攻撃範囲 + 両方の衝突半径を考慮)
	distance := u.Position.Distance(target.Position)
	effectiveRange := u.Range + u.GetCollisionRadius() + target.GetCollisionRadius()
	if distance > effectiveRange {
		return 0
	}
	
	// Trigger attack animation
	u.Animation.SetAnimation(graphics.AnimationAttack)
	
	// Calculate damage
	baseDamage := u.AttackPower
	if u.Type == UnitTypeMage {
		baseDamage += u.MagicPower
	}
	
	// Apply defense
	damage := baseDamage - target.Defense
	if damage < 1 {
		damage = 1 // Minimum damage
	}
	
	// Apply damage
	target.TakeDamage(damage)
	
	// Set cooldown
	u.LastAttackTime = u.AttackCooldown
	
	return damage
}

// TakeDamage applies damage to the unit
func (u *Unit) TakeDamage(damage int) {
	if !u.IsAlive {
		return
	}
	
	u.HP -= damage
	if u.HP <= 0 {
		u.HP = 0
		u.IsAlive = false
	}
}

// StartRetreating makes the unit start retreating
func (u *Unit) StartRetreating(exitPoint math.Vector2D) {
	u.IsRetreating = true
	u.Target = exitPoint
}

// GetHealthPercentage returns the unit's health as a percentage
func (u *Unit) GetHealthPercentage() float64 {
	if u.MaxHP == 0 {
		return 0
	}
	return float64(u.HP) / float64(u.MaxHP)
}

// GetCollisionRadius returns the collision radius for this unit
func (u *Unit) GetCollisionRadius() float64 {
	// サイズに基づいて衝突半径を計算（基本半径 * サイズ倍率）
	baseRadius := 3.0  // 基本半径を10.0から3.0に縮小
	return baseRadius * u.Size
}

// GetSightRange returns the sight range for this unit
func (u *Unit) GetSightRange() float64 {
	// デフォルトで500m（5000px）の知覚範囲
	// 実際の実装では、ユニット設定から取得する
	return 5000.0
}

// IsCollidingWith checks if this unit is colliding with another unit
func (u *Unit) IsCollidingWith(other *Unit) bool {
	if !u.IsAlive || !other.IsAlive {
		return false
	}
	
	distance := u.Position.Distance(other.Position)
	combinedRadius := u.GetCollisionRadius() + other.GetCollisionRadius()
	
	return distance < combinedRadius
}

// ResolveCollision resolves collision with another unit by pushing them apart
func (u *Unit) ResolveCollision(other *Unit) {
	if !u.IsAlive || !other.IsAlive {
		return
	}
	
	distance := u.Position.Distance(other.Position)
	combinedRadius := u.GetCollisionRadius() + other.GetCollisionRadius()
	
	if distance < combinedRadius && distance > 0 {
		// 重なりを解消するために押し出す
		overlap := combinedRadius - distance
		direction := other.Position.Sub(u.Position).Normalize()
		
		// 両方のユニットを半分ずつ押し出す
		pushDistance := overlap * 0.5
		u.Position = u.Position.Sub(direction.Mul(pushDistance))
		other.Position = other.Position.Add(direction.Mul(pushDistance))
	}
}
