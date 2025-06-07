// domain/team/formation.go
package team

import (
	"github.com/devvspaces/fantasy_league/internal/domain/player"
)

// Formation represents team formation
type Formation string

const (
	Formation442     Formation = "4-4-2"
	Formation433     Formation = "4-3-3"
	Formation451     Formation = "4-5-1"
	Formation352     Formation = "3-5-2"
	Formation532     Formation = "5-3-2"
	Formation4231    Formation = "4-2-3-1"
	Formation4312    Formation = "4-3-1-2"
	FormationDefault           = Formation442
)

// Lineup represents a match lineup
type Lineup struct {
	Formation   Formation
	Starters    []player.PlayerID // 11 players
	Positions   []player.Position // Position for each starter
	Substitutes []player.PlayerID // Bench players
	Captain     player.PlayerID
}

// FormationRequirements defines position requirements
type FormationRequirements map[player.Position]int

// GetPositionRequirements returns required positions for formation
func (f Formation) GetPositionRequirements() FormationRequirements {
	switch f {
	case Formation442:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 4,
			player.PositionMID: 4,
			player.PositionFWD: 2,
		}
	case Formation433:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 4,
			player.PositionMID: 3,
			player.PositionFWD: 3,
		}
	case Formation451:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 4,
			player.PositionMID: 5,
			player.PositionFWD: 1,
		}
	case Formation352:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 3,
			player.PositionMID: 5,
			player.PositionFWD: 2,
		}
	case Formation532:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 5,
			player.PositionMID: 3,
			player.PositionFWD: 2,
		}
	case Formation4231:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 4,
			player.PositionMID: 5, // 2 DM + 3 AM
			player.PositionFWD: 1,
		}
	case Formation4312:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 4,
			player.PositionMID: 4, // 3 CM + 1 AM
			player.PositionFWD: 2,
		}
	default:
		return FormationRequirements{
			player.PositionGK:  1,
			player.PositionDEF: 4,
			player.PositionMID: 4,
			player.PositionFWD: 2,
		}
	}
}

// IsValid checks if formation is valid
func (f Formation) IsValid() bool {
	switch f {
	case Formation442, Formation433, Formation451,
		Formation352, Formation532, Formation4231, Formation4312:
		return true
	default:
		return false
	}
}

// GetFormationStrength calculates formation effectiveness
func (f Formation) GetFormationStrength(matchup Formation) float64 {
	// Rock-paper-scissors style advantages
	advantages := map[Formation]map[Formation]float64{
		Formation442: {
			Formation433: 0.9,
			Formation451: 1.1,
			Formation352: 1.0,
		},
		Formation433: {
			Formation442: 1.1,
			Formation451: 0.9,
			Formation532: 1.1,
		},
		Formation451: {
			Formation433: 1.1,
			Formation442: 0.9,
			Formation352: 1.0,
		},
		Formation352: {
			Formation442: 1.0,
			Formation532: 0.9,
			Formation433: 0.9,
		},
		Formation532: {
			Formation433: 0.9,
			Formation352: 1.1,
			Formation442: 1.0,
		},
	}

	if adv, ok := advantages[f]; ok {
		if mult, ok := adv[matchup]; ok {
			return mult
		}
	}

	return 1.0 // No advantage
}
