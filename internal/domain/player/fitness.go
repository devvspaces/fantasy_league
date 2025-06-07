// domain/player/fitness.go
package player

import (
	"math"
)

// FitnessManager handles player fitness calculations
type FitnessManager struct {
	fatigueRate     float64
	recoveryRate    float64
	injuryThreshold float64
}

// NewFitnessManager creates a fitness manager
func NewFitnessManager() *FitnessManager {
	return &FitnessManager{
		fatigueRate:     0.15, // Base fatigue per minute played
		recoveryRate:    10.0, // Base recovery per day
		injuryThreshold: 40.0, // Below this fitness, injury risk increases
	}
}

// CalculateMatchFatigue calculates fitness loss from a match
func (fm *FitnessManager) CalculateMatchFatigue(player *Player, minutesPlayed int, matchIntensity float64) float64 {
	if minutesPlayed == 0 {
		return 0
	}

	// Base fatigue
	fatigue := float64(minutesPlayed) * fm.fatigueRate

	// Intensity modifier
	fatigue *= matchIntensity

	// Stamina reduces fatigue
	staminaMod := 1.5 - (float64(player.Attributes.Stamina) / 100)
	fatigue *= staminaMod

	// Age factor
	age := player.Age()
	if age > 30 {
		fatigue *= 1.0 + (float64(age-30) * 0.05)
	}

	// Position factor
	switch player.Position {
	case PositionGK:
		fatigue *= 0.6
	case PositionDEF:
		fatigue *= 0.85
	case PositionMID:
		fatigue *= 1.15
	case PositionFWD:
		fatigue *= 1.0
	}

	return math.Min(fatigue, 60) // Cap maximum fatigue
}

// CalculateDailyRecovery calculates fitness recovery per day
func (fm *FitnessManager) CalculateDailyRecovery(player *Player, trainingIntensity float64) float64 {
	// Base recovery
	recovery := fm.recoveryRate

	// Stamina bonus
	recovery += float64(player.Attributes.Stamina) / 20

	// Age factor
	age := player.Age()
	if age < 23 {
		recovery *= 1.2
	} else if age > 30 {
		recovery *= 0.9 - (float64(age-30) * 0.02)
	}

	// Training intensity reduces recovery
	recovery *= (2 - trainingIntensity)

	// Professionalism helps recovery
	recovery *= 1 + (float64(player.Attributes.Professionalism) / 200)

	return recovery
}

// CalculateInjuryRisk calculates injury probability
func (fm *FitnessManager) CalculateInjuryRisk(player *Player) float64 {
	risk := 0.0

	// Low fitness increases risk
	if player.Fitness < fm.injuryThreshold {
		risk += (fm.injuryThreshold - player.Fitness) / 100
	}

	// Age factor
	age := player.Age()
	if age > 30 {
		risk += float64(age-30) * 0.01
	}

	// Recent injury history would increase risk
	// (would need injury history tracking)

	return math.Min(risk, 0.5) // Cap at 50% risk
}

// ApplyMatchFitness updates player fitness after a match
func (fm *FitnessManager) ApplyMatchFitness(player *Player, minutesPlayed int, intensity float64) {
	fatigue := fm.CalculateMatchFatigue(player, minutesPlayed, intensity)
	player.Fitness = math.Max(0, player.Fitness-fatigue)
}

// ApplyDailyRecovery updates player fitness with daily recovery
func (fm *FitnessManager) ApplyDailyRecovery(player *Player, trainingIntensity float64) {
	recovery := fm.CalculateDailyRecovery(player, trainingIntensity)
	player.Fitness = math.Min(100, player.Fitness+recovery)
}
