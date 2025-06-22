package graphics

import (
	"math"
)

// AnimationType represents different types of animations
type AnimationType int

const (
	AnimationIdle AnimationType = iota
	AnimationWalk
	AnimationAttack
	AnimationDeath
)

// AnimationState holds the current animation state
type AnimationState struct {
	Type          AnimationType
	Frame         int
	FrameTime     float64
	FrameDuration float64
	TotalFrames   int
	Loop          bool
	Finished      bool
}

// NewAnimationState creates a new animation state
func NewAnimationState(animType AnimationType) *AnimationState {
	state := &AnimationState{
		Type:          animType,
		Frame:         0,
		FrameTime:     0,
		FrameDuration: 0.15, // 150ms per frame
		Loop:          true,
		Finished:      false,
	}
	
	// Set frame count based on animation type
	switch animType {
	case AnimationIdle:
		state.TotalFrames = 4
		state.FrameDuration = 0.5 // Slower for idle
	case AnimationWalk:
		state.TotalFrames = 4
		state.FrameDuration = 0.15
	case AnimationAttack:
		state.TotalFrames = 3
		state.FrameDuration = 0.1
		state.Loop = false
	case AnimationDeath:
		state.TotalFrames = 5
		state.FrameDuration = 0.2
		state.Loop = false
	}
	
	return state
}

// Update updates the animation state
func (as *AnimationState) Update(deltaTime float64) {
	if as.Finished && !as.Loop {
		return
	}
	
	as.FrameTime += deltaTime
	
	if as.FrameTime >= as.FrameDuration {
		as.FrameTime = 0
		as.Frame++
		
		if as.Frame >= as.TotalFrames {
			if as.Loop {
				as.Frame = 0
			} else {
				as.Frame = as.TotalFrames - 1
				as.Finished = true
			}
		}
	}
}

// Reset resets the animation to the beginning
func (as *AnimationState) Reset() {
	as.Frame = 0
	as.FrameTime = 0
	as.Finished = false
}

// SetAnimation changes the current animation type
func (as *AnimationState) SetAnimation(animType AnimationType) {
	if as.Type == animType {
		return
	}
	
	as.Type = animType
	as.Reset()
	
	// Update parameters for new animation type
	switch animType {
	case AnimationIdle:
		as.TotalFrames = 4
		as.FrameDuration = 0.5
		as.Loop = true
	case AnimationWalk:
		as.TotalFrames = 4
		as.FrameDuration = 0.15
		as.Loop = true
	case AnimationAttack:
		as.TotalFrames = 3
		as.FrameDuration = 0.1
		as.Loop = false
	case AnimationDeath:
		as.TotalFrames = 5
		as.FrameDuration = 0.2
		as.Loop = false
	}
}

// GetAnimationOffset returns offset values for animation effects
func (as *AnimationState) GetAnimationOffset() (float64, float64) {
	switch as.Type {
	case AnimationIdle:
		// Gentle bobbing motion
		bob := math.Sin(float64(as.Frame) * math.Pi / 2) * 1.0
		return 0, bob
		
	case AnimationWalk:
		// Walking bounce
		bounce := math.Abs(math.Sin(float64(as.Frame) * math.Pi / 2)) * 2.0
		return 0, -bounce
		
	case AnimationAttack:
		// Forward thrust motion
		thrust := 0.0
		if as.Frame == 1 {
			thrust = 3.0
		}
		return thrust, 0
		
	case AnimationDeath:
		// Falling motion
		fall := float64(as.Frame) * 2.0
		return 0, fall
	}
	
	return 0, 0
}

// GetScaleModifier returns scale modification for animation
func (as *AnimationState) GetScaleModifier() float64 {
	switch as.Type {
	case AnimationAttack:
		if as.Frame == 1 {
			return 1.2 // Slightly larger during attack
		}
	case AnimationDeath:
		// Shrink as dying
		return 1.0 - (float64(as.Frame) / float64(as.TotalFrames) * 0.3)
	}
	
	return 1.0
}

// GetRotationModifier returns rotation modification for animation
func (as *AnimationState) GetRotationModifier() float64 {
	switch as.Type {
	case AnimationDeath:
		// Rotate as falling
		return float64(as.Frame) * math.Pi / 8
	}
	
	return 0.0
}
