package game

import (
	"fmt"
	"math/rand"

	"github.com/shirou/tinygocha/internal/data"
	gamemath "github.com/shirou/tinygocha/internal/math"
)

// BattleManager manages the battle state and logic
type BattleManager struct {
	ArmyA        *Army
	ArmyB        *Army
	Stage        data.StageConfig
	TerrainData  data.TerrainConfig
	BattleTime   float64
	TimeLimit    float64
	IsActive     bool
	Winner       int // -1: 未決定, 0: A軍勝利, 1: B軍勝利, 2: 引き分け
	
	// Unit ID counter
	nextUnitID int
}

// NewBattleManager creates a new battle manager
func NewBattleManager(stage data.StageConfig, terrainData data.TerrainConfig) *BattleManager {
	return &BattleManager{
		ArmyA:       NewArmy(0, "軍勢A", 0),
		ArmyB:       NewArmy(1, "軍勢B", 1),
		Stage:       stage,
		TerrainData: terrainData,
		BattleTime:  0.0,
		TimeLimit:   stage.TimeLimit,
		IsActive:    false,
		Winner:      -1,
		nextUnitID:  1,
	}
}

// CreatePresetArmy creates a preset army configuration
func (bm *BattleManager) CreatePresetArmy(armyID int, presetType string, dataManager *data.DataManager) error {
	var army *Army
	if armyID == 0 {
		army = bm.ArmyA
	} else {
		army = bm.ArmyB
	}
	
	fmt.Printf("Creating preset army %d (%s)\n", armyID, presetType)
	
	// Get deployment points
	var deploymentPoints []gamemath.Vector2D
	if armyID == 0 {
		deploymentPoints = bm.Stage.GetDeploymentPointsA()
	} else {
		deploymentPoints = bm.Stage.GetDeploymentPointsB()
	}
	
	fmt.Printf("Deployment points for army %d: %v\n", armyID, deploymentPoints)
	
	// Create groups based on preset type
	switch presetType {
	case "バランス型":
		bm.createBalancedArmy(army, deploymentPoints, dataManager)
	case "攻撃重視":
		bm.createOffensiveArmy(army, deploymentPoints, dataManager)
	case "防御重視":
		bm.createDefensiveArmy(army, deploymentPoints, dataManager)
	default:
		bm.createBalancedArmy(army, deploymentPoints, dataManager)
	}
	
	// デバッグ: 作成されたユニット数
	allUnits := army.GetAllUnits()
	fmt.Printf("Army %d created with %d units:\n", armyID, len(allUnits))
	for _, unit := range allUnits {
		fmt.Printf("  Unit ID=%d, Type=%s, Pos=(%.1f,%.1f), AI=%t\n", 
			unit.ID, unit.Type, unit.Position.X, unit.Position.Y, unit.AI != nil)
	}
	
	return nil
}

// createBalancedArmy creates a balanced army composition
func (bm *BattleManager) createBalancedArmy(army *Army, deploymentPoints []gamemath.Vector2D, dataManager *data.DataManager) {
	groupConfigs := []struct {
		leaderType string
		memberType string
		count      int
	}{
		{"infantry", "infantry", 4},
		{"archer", "archer", 3},
		{"mage", "infantry", 2},
	}
	
	for i, config := range groupConfigs {
		if i >= len(deploymentPoints) {
			break
		}
		
		group := bm.createGroup(army.ID, config.leaderType, config.memberType, config.count, deploymentPoints[i], dataManager)
		army.AddGroup(group)
	}
}

// createOffensiveArmy creates an offensive army composition
func (bm *BattleManager) createOffensiveArmy(army *Army, deploymentPoints []gamemath.Vector2D, dataManager *data.DataManager) {
	groupConfigs := []struct {
		leaderType string
		memberType string
		count      int
	}{
		{"cavalry", "cavalry", 2},
		{"archer", "archer", 4},
		{"infantry", "infantry", 3},
	}
	
	for i, config := range groupConfigs {
		if i >= len(deploymentPoints) {
			break
		}
		
		group := bm.createGroup(army.ID, config.leaderType, config.memberType, config.count, deploymentPoints[i], dataManager)
		army.AddGroup(group)
	}
}

// createDefensiveArmy creates a defensive army composition
func (bm *BattleManager) createDefensiveArmy(army *Army, deploymentPoints []gamemath.Vector2D, dataManager *data.DataManager) {
	groupConfigs := []struct {
		leaderType string
		memberType string
		count      int
	}{
		{"heavy_infantry", "heavy_infantry", 3},
		{"infantry", "archer", 4},
		{"mage", "mage", 2},
	}
	
	for i, config := range groupConfigs {
		if i >= len(deploymentPoints) {
			break
		}
		
		group := bm.createGroup(army.ID, config.leaderType, config.memberType, config.count, deploymentPoints[i], dataManager)
		army.AddGroup(group)
	}
}

// createGroup creates a group with specified configuration
func (bm *BattleManager) createGroup(armyID int, leaderType, memberType string, memberCount int, position gamemath.Vector2D, dataManager *data.DataManager) *Group {
	// Get unit configurations
	leaderConfig, err := dataManager.GetUnitConfig(leaderType)
	if err != nil {
		fmt.Printf("Error getting leader config for %s: %v\n", leaderType, err)
		return nil
	}
	
	memberConfig, err := dataManager.GetUnitConfig(memberType)
	if err != nil {
		fmt.Printf("Error getting member config for %s: %v\n", memberType, err)
		return nil
	}
	
	fmt.Printf("Creating group: Leader=%s (HP=%d), Members=%s (HP=%d), Count=%d\n", 
		leaderType, leaderConfig.HP, memberType, memberConfig.HP, memberCount)
	
	// Create leader
	leader := bm.createUnit(UnitType(leaderType), UnitTypeConfig{
		Name:       leaderConfig.Name,
		HP:         leaderConfig.HP,
		Attack:     leaderConfig.Attack,
		Defense:    leaderConfig.Defense,
		Speed:      leaderConfig.Speed,
		Range:      leaderConfig.Range,
		MagicPower: leaderConfig.MagicPower,
		Size:       leaderConfig.Size,  // サイズフィールドを追加
	}, true, armyID)
	leader.Position = position
	leader.Target = position
	
	// Create members
	var members []*Unit
	for i := 0; i < memberCount; i++ {
		member := bm.createUnit(UnitType(memberType), UnitTypeConfig{
			Name:       memberConfig.Name,
			HP:         memberConfig.HP,
			Attack:     memberConfig.Attack,
			Defense:    memberConfig.Defense,
			Speed:      memberConfig.Speed,
			Range:      memberConfig.Range,
			MagicPower: memberConfig.MagicPower,
			Size:       memberConfig.Size,  // サイズフィールドを追加
		}, false, armyID)
		member.Position = position.Add(gamemath.Vector2D{
			X: float64(rand.Intn(40) - 20),
			Y: float64(rand.Intn(40) - 20),
		})
		member.Target = member.Position
		members = append(members, member)
	}
	
	// Create group
	group := NewGroup(len(bm.ArmyA.Groups)+len(bm.ArmyB.Groups), armyID, leader, members)
	
	// Set group IDs for all units
	leader.GroupID = group.ID
	for _, member := range members {
		member.GroupID = group.ID
	}
	
	return group
}

// createUnit creates a new unit with terrain modifiers applied
func (bm *BattleManager) createUnit(unitType UnitType, config UnitTypeConfig, isLeader bool, armyID int) *Unit {
	unit := NewUnit(bm.nextUnitID, unitType, config, isLeader, 0, armyID)
	bm.nextUnitID++
	
	// Apply terrain modifiers
	bm.applyTerrainModifiers(unit)
	
	return unit
}

// applyTerrainModifiers applies terrain effects to a unit
func (bm *BattleManager) applyTerrainModifiers(unit *Unit) {
	// Apply movement modifier
	unit.Speed *= bm.TerrainData.MovementModifier
	
	// Apply defense modifier
	unit.Defense = int(float64(unit.Defense) * bm.TerrainData.DefenseModifier)
	
	// Apply unit type specific bonuses
	switch unit.Type {
	case UnitTypeInfantry:
		unit.AttackPower = int(float64(unit.AttackPower) * bm.TerrainData.InfantryBonus)
	case UnitTypeArcher:
		unit.AttackPower = int(float64(unit.AttackPower) * bm.TerrainData.ArcherBonus)
	case UnitTypeMage:
		unit.AttackPower = int(float64(unit.AttackPower) * bm.TerrainData.MageBonus)
		unit.MagicPower = int(float64(unit.MagicPower) * bm.TerrainData.MageBonus)
	}
}

// StartBattle starts the battle
func (bm *BattleManager) StartBattle() {
	bm.IsActive = true
	bm.BattleTime = 0.0
	bm.Winner = -1
}

// Update updates the battle state
func (bm *BattleManager) Update(deltaTime float64) {
	if !bm.IsActive {
		return
	}
	
	// Update battle time
	bm.BattleTime += deltaTime
	
	// Update armies
	bm.ArmyA.Update(deltaTime)
	bm.ArmyB.Update(deltaTime)
	
	// Update AI behaviors
	bm.updateAI(deltaTime)
	
	// Handle unit collisions
	bm.handleCollisions()
	
	// Process combat
	bm.processCombat()
	
	// Check win conditions
	bm.checkWinConditions()
}

// processCombat handles combat between units
func (bm *BattleManager) processCombat() {
	unitsA := bm.ArmyA.GetAliveUnits()
	unitsB := bm.ArmyB.GetAliveUnits()
	
	// Army A attacks Army B
	for _, unitA := range unitsA {
		if !unitA.CanAttack() {
			continue
		}
		
		// Find closest enemy in range
		var target *Unit
		minDistance := float64(unitA.Range + 1) // Start with out of range
		
		for _, unitB := range unitsB {
			distance := unitA.Position.Distance(unitB.Position)
			if distance <= unitA.Range && distance < minDistance {
				target = unitB
				minDistance = distance
			}
		}
		
		// Attack if target found
		if target != nil {
			unitA.Attack(target)
		}
	}
	
	// Army B attacks Army A
	for _, unitB := range unitsB {
		if !unitB.CanAttack() {
			continue
		}
		
		// Find closest enemy in range
		var target *Unit
		minDistance := float64(unitB.Range + 1)
		
		for _, unitA := range unitsA {
			distance := unitB.Position.Distance(unitA.Position)
			if distance <= unitB.Range && distance < minDistance {
				target = unitA
				minDistance = distance
			}
		}
		
		// Attack if target found
		if target != nil {
			unitB.Attack(target)
		}
	}
}

// checkWinConditions checks if the battle should end
func (bm *BattleManager) checkWinConditions() {
	// Check if time limit reached
	if bm.BattleTime >= bm.TimeLimit {
		bm.IsActive = false
		// Determine winner by remaining health
		healthA := bm.ArmyA.GetTotalHealth()
		healthB := bm.ArmyB.GetTotalHealth()
		
		if healthA > healthB {
			bm.Winner = 0 // Army A wins
		} else if healthB > healthA {
			bm.Winner = 1 // Army B wins
		} else {
			bm.Winner = 2 // Draw
		}
		return
	}
	
	// Check if either army is defeated
	if bm.ArmyA.IsDefeated() && bm.ArmyB.IsDefeated() {
		bm.IsActive = false
		bm.Winner = 2 // Draw
	} else if bm.ArmyA.IsDefeated() {
		bm.IsActive = false
		bm.Winner = 1 // Army B wins
	} else if bm.ArmyB.IsDefeated() {
		bm.IsActive = false
		bm.Winner = 0 // Army A wins
	}
}

// GetWinnerName returns the name of the winner
func (bm *BattleManager) GetWinnerName() string {
	switch bm.Winner {
	case 0:
		return "軍勢A"
	case 1:
		return "軍勢B"
	case 2:
		return "引き分け"
	default:
		return "未決定"
	}
}

// updateAI updates AI behaviors for all units
func (bm *BattleManager) updateAI(deltaTime float64) {
	// Update Army A AI (fight against Army B)
	unitsA := bm.ArmyA.GetAliveUnits()
	unitsB := bm.ArmyB.GetAliveUnits()
	
	// デバッグ: 軍勢の状況
	fmt.Printf("AI Update - Army A: %d units, Army B: %d units\n", len(unitsA), len(unitsB))
	
	for _, unit := range unitsA {
		if unit.AI != nil {
			unit.AI.Update(unit, unitsB, deltaTime)
		}
	}
	
	// Update Army B AI (fight against Army A)
	for _, unit := range unitsB {
		if unit.AI != nil {
			unit.AI.Update(unit, unitsA, deltaTime)
		}
	}
}

// handleCollisions handles collisions between all units
func (bm *BattleManager) handleCollisions() {
	allUnits := append(bm.ArmyA.GetAliveUnits(), bm.ArmyB.GetAliveUnits()...)
	
	// Check collisions between all pairs of units
	for i := 0; i < len(allUnits); i++ {
		for j := i + 1; j < len(allUnits); j++ {
			unit1 := allUnits[i]
			unit2 := allUnits[j]
			
			if unit1.IsCollidingWith(unit2) {
				unit1.ResolveCollision(unit2)
			}
		}
	}
}
