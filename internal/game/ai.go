package game

import (
	"fmt"
	stdmath "math"
)

// AIBehavior represents AI behavior state for a unit
type AIBehavior struct {
	TargetEnemy      *Unit
	PreferredRange   float64 // 理想的な戦闘距離
	AggressionLevel  float64 // 攻撃性 (0.0-1.0)
	LastDecisionTime float64
	DecisionCooldown float64 // 判断間隔（秒）
	
	// 行動状態
	CurrentAction    AIAction
	ActionStartTime  float64
	ActionDuration   float64
}

// AIAction represents different AI actions
type AIAction int

const (
	AIActionIdle     AIAction = iota // 待機
	AIActionApproach                 // 接近
	AIActionRetreat                  // 後退
	AIActionAttack                   // 攻撃
	AIActionHold                     // 位置保持
)

// NewAIBehavior creates a new AI behavior based on unit type
func NewAIBehavior(unitType UnitType) *AIBehavior {
	ai := &AIBehavior{
		DecisionCooldown: 0.1, // 0.1秒間隔で判断（高速化）
		LastDecisionTime: 0,
		CurrentAction:    AIActionIdle,
	}
	
	// ユニット種別に応じた設定（新スケール対応）
	switch unitType {
	case UnitTypeInfantry:
		ai.PreferredRange = 15.0  // 1.5m = 15px
		ai.AggressionLevel = 0.7
	case UnitTypeArcher:
		ai.PreferredRange = 600.0 // 60m = 600px（射程80mの75%）
		ai.AggressionLevel = 0.5
	case UnitTypeMage:
		ai.PreferredRange = 480.0 // 48m = 480px（射程60mの80%）
		ai.AggressionLevel = 0.4
	case "heavy_infantry":
		ai.PreferredRange = 20.0  // 2m = 20px
		ai.AggressionLevel = 0.8
	case "cavalry":
		ai.PreferredRange = 25.0  // 2.5m = 25px
		ai.AggressionLevel = 0.9
	default:
		ai.PreferredRange = 15.0  // デフォルト
		ai.AggressionLevel = 0.6
	}
	
	return ai
}

// Update updates the AI behavior
func (ai *AIBehavior) Update(unit *Unit, enemies []*Unit, deltaTime float64) {
	if !unit.IsAlive || unit.IsRetreating {
		return
	}
	
	// 判断クールダウンチェック
	ai.LastDecisionTime += deltaTime
	if ai.LastDecisionTime < ai.DecisionCooldown {
		return
	}
	
	ai.LastDecisionTime = 0
	
	// デバッグ: リーダーのみログ出力
	if unit.IsLeader {
		fmt.Printf("AI Update: Unit %d, Enemies: %d\n", unit.ID, len(enemies))
	}
	
	// 敵の探索・選択
	ai.selectTarget(unit, enemies)
	
	if ai.TargetEnemy == nil || !ai.TargetEnemy.IsAlive {
		ai.CurrentAction = AIActionIdle
		if unit.IsLeader {
			fmt.Printf("Unit %d: No target\n", unit.ID)
		}
		return
	}
	
	// 距離ベースの行動決定
	distance := unit.Position.Distance(ai.TargetEnemy.Position)
	ai.decideAction(unit, distance)
	
	// デバッグ: 行動決定の確認
	if unit.IsLeader {
		fmt.Printf("Unit %d: Target=%d, Distance=%.2f, Action=%s\n", 
			unit.ID, ai.TargetEnemy.ID, distance, ai.GetActionName())
	}
	
	// 行動実行
	ai.executeAction(unit, distance)
}

// selectTarget selects the best target enemy
func (ai *AIBehavior) selectTarget(unit *Unit, enemies []*Unit) {
	var bestTarget *Unit
	bestScore := -1.0
	
	// デバッグ: 敵軍の詳細情報
	if unit.IsLeader {
		fmt.Printf("Unit %d (Army %d) selecting target from %d enemies:\n", unit.ID, unit.ArmyID, len(enemies))
		validEnemies := 0
		for i, enemy := range enemies {
			isValid := enemy.IsAlive && !enemy.IsRetreating
			if isValid {
				validEnemies++
			}
			fmt.Printf("  Enemy[%d]: ID=%d, Army=%d, Alive=%t, Retreating=%t, Pos=(%.1f,%.1f), Valid=%t\n", 
				i, enemy.ID, enemy.ArmyID, enemy.IsAlive, enemy.IsRetreating, enemy.Position.X, enemy.Position.Y, isValid)
		}
		fmt.Printf("  Valid enemies: %d/%d\n", validEnemies, len(enemies))
	}
	
	for _, enemy := range enemies {
		if !enemy.IsAlive || enemy.IsRetreating {
			continue
		}
		
		distance := unit.Position.Distance(enemy.Position)
		
		// 知覚範囲チェック - 範囲外の敵は無視
		sightRange := unit.GetSightRange()
		if distance > sightRange {
			continue
		}
		
		// スコア計算（距離、敵の体力、優先度を考慮）
		score := ai.calculateTargetScore(unit, enemy, distance)
		
		// デバッグ: スコア詳細（リーダーのみ）
		if unit.IsLeader {
			fmt.Printf("    Enemy ID=%d: Distance=%.1f, SightRange=%.1f, Score=%.2f\n", enemy.ID, distance, sightRange, score)
		}
		
		if score > bestScore {
			bestScore = score
			bestTarget = enemy
		}
	}
	
	ai.TargetEnemy = bestTarget
	
	if unit.IsLeader {
		if bestTarget != nil {
			fmt.Printf("Unit %d selected target: ID=%d (score: %.2f)\n", unit.ID, bestTarget.ID, bestScore)
		} else {
			fmt.Printf("Unit %d: No valid target found!\n", unit.ID)
		}
	}
}

// calculateTargetScore calculates target priority score
func (ai *AIBehavior) calculateTargetScore(unit *Unit, enemy *Unit, distance float64) float64 {
	// 基本スコア
	score := 1000.0  // 基本スコアを大幅に増加
	
	// 距離による減点（近い敵を優先、ただし極端に遠い敵も除外しない）
	score -= distance * 0.05  // 距離による減点をさらに緩和
	
	// 敵の体力による加点（体力が少ない敵を優先）
	healthPercent := enemy.GetHealthPercentage()
	score += (1.0 - healthPercent) * 30.0
	
	// リーダーボーナス
	if enemy.IsLeader {
		score += 50.0
	}
	
	// 射程内の敵にボーナス
	if distance <= unit.Range {
		score += 100.0
	}
	
	// ユニット種別による優先度
	switch enemy.Type {
	case UnitTypeMage:
		score += 20.0 // 魔術師を優先
	case UnitTypeArcher:
		score += 15.0 // 弓兵を優先
	case UnitTypeInfantry:
		score += 10.0
	}
	
	return score
}

// decideAction decides what action to take based on distance
func (ai *AIBehavior) decideAction(unit *Unit, distance float64) {
	// 衝突半径を考慮した実効距離を計算
	effectiveDistance := distance - unit.GetCollisionRadius() - ai.TargetEnemy.GetCollisionRadius()
	
	// 攻撃可能距離内かチェック（実効距離で判定）
	if effectiveDistance <= unit.Range && unit.CanAttack() {
		ai.CurrentAction = AIActionAttack
		return
	}
	
	// 理想的な距離と比較（実効距離で判定）
	if effectiveDistance > ai.PreferredRange * 1.2 {
		// 遠すぎる場合は接近
		ai.CurrentAction = AIActionApproach
	} else if effectiveDistance < ai.PreferredRange * 0.8 && ai.isRangedUnit(unit) {
		// 近すぎる場合は後退（遠距離ユニットのみ）
		ai.CurrentAction = AIActionRetreat
	} else if effectiveDistance <= unit.Range {
		// 射程内だが攻撃できない場合は位置保持
		ai.CurrentAction = AIActionHold
	} else {
		// その他の場合は接近
		ai.CurrentAction = AIActionApproach
	}
}

// executeAction executes the decided action
func (ai *AIBehavior) executeAction(unit *Unit, distance float64) {
	switch ai.CurrentAction {
	case AIActionApproach:
		ai.moveTowardsTarget(unit, 1.0) // 敵に向かって移動
		
	case AIActionRetreat:
		ai.moveAwayFromTarget(unit, 1.0) // 敵から離れる
		
	case AIActionAttack:
		// 攻撃は Unit.Attack で自動実行される
		
	case AIActionHold:
		// 現在位置を保持（移動しない）
		unit.Target = unit.Position
		
	case AIActionIdle:
		// 何もしない
		unit.Target = unit.Position
	}
}

// moveTowardsTarget moves unit towards the target enemy
func (ai *AIBehavior) moveTowardsTarget(unit *Unit, intensity float64) {
	if ai.TargetEnemy == nil {
		return
	}
	
	direction := ai.TargetEnemy.Position.Sub(unit.Position).Normalize()
	
	// 敵に向かって移動（衝突半径を考慮した理想距離まで）
	currentDistance := unit.Position.Distance(ai.TargetEnemy.Position)
	collisionBuffer := unit.GetCollisionRadius() + ai.TargetEnemy.GetCollisionRadius()
	targetDistance := ai.PreferredRange * 0.9 + collisionBuffer // 理想距離 + 衝突バッファ
	
	if currentDistance > targetDistance {
		// 理想距離まで接近（より大きな移動距離）
		moveDistance := stdmath.Min(currentDistance - targetDistance, 50.0) // 最大50ピクセル移動
		targetPos := unit.Position.Add(direction.Mul(moveDistance * intensity))
		unit.MoveTo(targetPos)
	} else {
		// 既に理想距離内にいる場合は、敵に向かって少し移動
		targetPos := unit.Position.Add(direction.Mul(20.0 * intensity)) // 移動距離を増加
		unit.MoveTo(targetPos)
	}
}

// moveAwayFromTarget moves unit away from the target enemy
func (ai *AIBehavior) moveAwayFromTarget(unit *Unit, intensity float64) {
	if ai.TargetEnemy == nil {
		return
	}
	
	direction := unit.Position.Sub(ai.TargetEnemy.Position).Normalize()
	
	// 理想的な距離まで後退（衝突半径を考慮）
	currentDistance := unit.Position.Distance(ai.TargetEnemy.Position)
	collisionBuffer := unit.GetCollisionRadius() + ai.TargetEnemy.GetCollisionRadius()
	targetDistance := ai.PreferredRange * 1.1 + collisionBuffer // 理想距離 + 衝突バッファ
	moveDistance := targetDistance - currentDistance
	
	if moveDistance > 0 {
		targetPos := unit.Position.Add(direction.Mul(moveDistance * intensity))
		
		// 画面外に出ないようにクランプ
		targetPos.X = stdmath.Max(50, stdmath.Min(974, targetPos.X))
		targetPos.Y = stdmath.Max(100, stdmath.Min(700, targetPos.Y))
		
		unit.MoveTo(targetPos)
	}
}

// isRangedUnit checks if the unit is a ranged unit
func (ai *AIBehavior) isRangedUnit(unit *Unit) bool {
	return unit.Type == UnitTypeArcher || unit.Type == UnitTypeMage
}

// GetActionName returns human-readable action name for debugging
func (ai *AIBehavior) GetActionName() string {
	switch ai.CurrentAction {
	case AIActionIdle:
		return "待機"
	case AIActionApproach:
		return "接近"
	case AIActionRetreat:
		return "後退"
	case AIActionAttack:
		return "攻撃"
	case AIActionHold:
		return "保持"
	default:
		return "不明"
	}
}
