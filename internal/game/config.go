package game

// UnitTypeConfig represents unit configuration (re-exported from data package)
type UnitTypeConfig struct {
	Name       string
	HP         int
	Attack     int
	Defense    int
	Speed      float64
	Range      float64
	MagicPower int
	Size       float64  // ユニットの大きさ（衝突判定用）
}
