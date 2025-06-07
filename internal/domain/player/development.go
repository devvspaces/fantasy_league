// domain/player/development.go
package player

import (
	"math"
	"math/rand"
)

// DevelopmentManager handles player growth and decline
type DevelopmentManager struct {
	rand *rand.Rand
}

// NewDevelopmentManager creates a development manager
func NewDevelopmentManager() *DevelopmentManager {
	return &DevelopmentManager{
		rand: rand.New(rand.NewSource(42)), // Use seeded random for consistency
	}
}

// TrainingType represents different training focuses
type TrainingType string

const (
	TrainingGeneral   TrainingType = "general"
	TrainingTechnical TrainingType = "technical"
	TrainingPhysical  TrainingType = "physical"
	TrainingTactical  TrainingType = "tactical"
	TrainingSetPieces TrainingType = "set_pieces"
)

// TrainingResult contains the outcome of training
type TrainingResult struct {
	AttributeChanges map[string]int
	FitnessChange    float64
	MoraleChange     float64
}

// ProcessTraining applies training effects to a player
func (dm *DevelopmentManager) ProcessTraining(player *Player, trainingType TrainingType, intensity float64) TrainingResult {
	result := TrainingResult{
		AttributeChanges: make(map[string]int),
	}

	// Base improvement chance
	improvementChance := dm.calculateImprovementChance(player)

	// Apply training based on type
	switch trainingType {
	case TrainingTechnical:
		dm.trainTechnical(player, improvementChance, &result)
	case TrainingPhysical:
		dm.trainPhysical(player, improvementChance, &result)
	case TrainingTactical:
		dm.trainTactical(player, improvementChance, &result)
	case TrainingSetPieces:
		dm.trainSetPieces(player, improvementChance, &result)
	default:
		dm.trainGeneral(player, improvementChance, &result)
	}

	// Fitness impact
	result.FitnessChange = -5 * intensity

	// Morale impact (players like moderate training)
	if intensity < 0.7 {
		result.MoraleChange = 2
	} else if intensity > 0.9 {
		result.MoraleChange = -3
	}

	return result
}

// ProcessNaturalDevelopment handles age-based attribute changes
func (dm *DevelopmentManager) ProcessNaturalDevelopment(player *Player) {
	age := player.Age()

	// Young players improve naturally
	if age < 23 {
		dm.youngPlayerDevelopment(player)
	} else if age > 30 {
		dm.veteranDecline(player)
	}

	// Update overall quality based on attributes
	player.Attributes.Quality = player.GetOverallRating()
}

// calculateImprovementChance determines likelihood of improvement
func (dm *DevelopmentManager) calculateImprovementChance(player *Player) float64 {
	age := player.Age()

	// Base chance by age
	var baseChance float64
	if age < 21 {
		baseChance = 0.8
	} else if age < 25 {
		baseChance = 0.6
	} else if age < 28 {
		baseChance = 0.4
	} else if age < 32 {
		baseChance = 0.2
	} else {
		baseChance = 0.05
	}

	// Modify by potential
	potentialMod := float64(player.Attributes.Potential) / 100
	baseChance *= potentialMod

	// Modify by professionalism
	profMod := float64(player.Attributes.Professionalism) / 100
	baseChance *= (0.5 + 0.5*profMod)

	// Modify by morale
	moraleMod := player.Morale / 100
	baseChance *= (0.8 + 0.2*moraleMod)

	return baseChance
}

// trainTechnical focuses on technical skills
func (dm *DevelopmentManager) trainTechnical(player *Player, chance float64, result *TrainingResult) {
	attrs := []string{"Passing", "BallControl", "Shooting"}

	for _, attr := range attrs {
		if dm.rand.Float64() < chance {
			improvement := dm.calculateImprovement(player, attr)
			if improvement > 0 {
				result.AttributeChanges[attr] = improvement
				dm.applyAttributeChange(player, attr, improvement)
			}
		}
	}
}

// trainPhysical focuses on physical attributes
func (dm *DevelopmentManager) trainPhysical(player *Player, chance float64, result *TrainingResult) {
	attrs := []string{"Speed", "Stamina", "Heading"}

	for _, attr := range attrs {
		if dm.rand.Float64() < chance*0.8 { // Harder to improve physical
			improvement := dm.calculateImprovement(player, attr)
			if improvement > 0 {
				result.AttributeChanges[attr] = improvement
				dm.applyAttributeChange(player, attr, improvement)
			}
		}
	}
}

// trainTactical focuses on mental attributes
func (dm *DevelopmentManager) trainTactical(player *Player, chance float64, result *TrainingResult) {
	attrs := []string{"Perception", "Tackling"}

	for _, attr := range attrs {
		if dm.rand.Float64() < chance {
			improvement := dm.calculateImprovement(player, attr)
			if improvement > 0 {
				result.AttributeChanges[attr] = improvement
				dm.applyAttributeChange(player, attr, improvement)
			}
		}
	}
}

// trainSetPieces focuses on set piece attributes
func (dm *DevelopmentManager) trainSetPieces(player *Player, chance float64, result *TrainingResult) {
	if player.Position == PositionGK {
		// Goalkeepers improve keeping
		if dm.rand.Float64() < chance {
			improvement := dm.calculateImprovement(player, "Keeping")
			if improvement > 0 {
				result.AttributeChanges["Keeping"] = improvement
				player.Attributes.Keeping += improvement
			}
		}
	} else {
		// Others improve heading and shooting
		attrs := []string{"Heading", "Shooting"}
		for _, attr := range attrs {
			if dm.rand.Float64() < chance*0.7 {
				improvement := dm.calculateImprovement(player, attr)
				if improvement > 0 {
					result.AttributeChanges[attr] = improvement
					dm.applyAttributeChange(player, attr, improvement)
				}
			}
		}
	}
}

// trainGeneral provides balanced training
func (dm *DevelopmentManager) trainGeneral(player *Player, chance float64, result *TrainingResult) {
	// Small chance to improve any attribute
	allAttrs := []string{"Keeping", "Tackling", "Passing", "Shooting", "Heading",
		"Speed", "Stamina", "Perception", "BallControl"}

	// Pick 2-3 random attributes
	numAttrs := 2 + dm.rand.Intn(2)
	for i := 0; i < numAttrs; i++ {
		attr := allAttrs[dm.rand.Intn(len(allAttrs))]
		if dm.rand.Float64() < chance*0.5 {
			improvement := dm.calculateImprovement(player, attr)
			if improvement > 0 {
				result.AttributeChanges[attr] = improvement
				dm.applyAttributeChange(player, attr, improvement)
			}
		}
	}
}

// calculateImprovement determines attribute improvement amount
func (dm *DevelopmentManager) calculateImprovement(player *Player, attribute string) int {
	current := dm.getAttributeValue(player, attribute)

	// Harder to improve higher attributes
	if current >= 95 {
		return 0
	} else if current >= 90 {
		if dm.rand.Float64() < 0.1 {
			return 1
		}
	} else if current >= 80 {
		if dm.rand.Float64() < 0.3 {
			return 1
		}
	} else {
		// Normal improvement
		return 1 + dm.rand.Intn(2)
	}

	return 0
}

// youngPlayerDevelopment handles natural growth for young players
func (dm *DevelopmentManager) youngPlayerDevelopment(player *Player) {
	// Physical growth
	if player.Age() < 21 {
		if dm.rand.Float64() < 0.3 {
			player.Attributes.Speed = int(math.Min(float64(player.Attributes.Speed+1), 100))
			player.Attributes.Stamina = int(math.Min(float64(player.Attributes.Stamina+1), 100))
		}
	}

	// Natural improvement based on potential
	if dm.rand.Float64() < float64(player.Attributes.Potential)/200 {
		// Random attribute improvement
		attrs := []string{"Passing", "BallControl", "Perception", "Tackling"}
		attr := attrs[dm.rand.Intn(len(attrs))]
		dm.applyAttributeChange(player, attr, 1)
	}
}

// veteranDecline handles age-related decline
func (dm *DevelopmentManager) veteranDecline(player *Player) {
	age := player.Age()
	declineRate := float64(age-30) * 0.05

	// Physical decline
	if dm.rand.Float64() < declineRate {
		player.Attributes.Speed = int(math.Max(float64(player.Attributes.Speed-1), 30))
	}
	if dm.rand.Float64() < declineRate*0.8 {
		player.Attributes.Stamina = int(math.Max(float64(player.Attributes.Stamina-1), 30))
	}

	// But experience can still improve
	if dm.rand.Float64() < 0.2 {
		player.Attributes.Perception = int(math.Min(float64(player.Attributes.Perception+1), 100))
	}
}

// Helper methods

func (dm *DevelopmentManager) getAttributeValue(player *Player, attribute string) int {
	switch attribute {
	case "Keeping":
		return player.Attributes.Keeping
	case "Tackling":
		return player.Attributes.Tackling
	case "Passing":
		return player.Attributes.Passing
	case "Shooting":
		return player.Attributes.Shooting
	case "Heading":
		return player.Attributes.Heading
	case "Speed":
		return player.Attributes.Speed
	case "Stamina":
		return player.Attributes.Stamina
	case "Perception":
		return player.Attributes.Perception
	case "BallControl":
		return player.Attributes.BallControl
	default:
		return 50
	}
}

func (dm *DevelopmentManager) applyAttributeChange(player *Player, attribute string, change int) {
	switch attribute {
	case "Keeping":
		player.Attributes.Keeping = int(math.Min(math.Max(float64(player.Attributes.Keeping+change), 0), 100))
	case "Tackling":
		player.Attributes.Tackling = int(math.Min(math.Max(float64(player.Attributes.Tackling+change), 0), 100))
	case "Passing":
		player.Attributes.Passing = int(math.Min(math.Max(float64(player.Attributes.Passing+change), 0), 100))
	case "Shooting":
		player.Attributes.Shooting = int(math.Min(math.Max(float64(player.Attributes.Shooting+change), 0), 100))
	case "Heading":
		player.Attributes.Heading = int(math.Min(math.Max(float64(player.Attributes.Heading+change), 0), 100))
	case "Speed":
		player.Attributes.Speed = int(math.Min(math.Max(float64(player.Attributes.Speed+change), 0), 100))
	case "Stamina":
		player.Attributes.Stamina = int(math.Min(math.Max(float64(player.Attributes.Stamina+change), 0), 100))
	case "Perception":
		player.Attributes.Perception = int(math.Min(math.Max(float64(player.Attributes.Perception+change), 0), 100))
	case "BallControl":
		player.Attributes.BallControl = int(math.Min(math.Max(float64(player.Attributes.BallControl+change), 0), 100))
	}
}
