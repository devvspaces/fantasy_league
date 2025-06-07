// domain/team/squad.go
package team

import (
	"sort"

	"github.com/devvspaces/fantasy_league/internal/domain/player"

	"github.com/devvspaces/fantasy_league/internal/domain/common"
)

// SquadManager handles squad-related operations
type SquadManager struct {
	team *Team
}

// NewSquadManager creates a squad manager
func NewSquadManager(team *Team) *SquadManager {
	return &SquadManager{team: team}
}

// GetSquadDepth analyzes squad depth by position
func (sm *SquadManager) GetSquadDepth() map[player.Position][]player.Player {
	depth := make(map[player.Position][]player.Player)

	for _, p := range sm.team.Players {
		depth[p.Position] = append(depth[p.Position], p)
	}

	// Sort by rating
	for pos := range depth {
		sort.Slice(depth[pos], func(i, j int) bool {
			return depth[pos][i].GetOverallRating() > depth[pos][j].GetOverallRating()
		})
	}

	return depth
}

// GetSquadAge calculates average squad age
func (sm *SquadManager) GetSquadAge() float64 {
	if len(sm.team.Players) == 0 {
		return 0
	}

	totalAge := 0
	for _, p := range sm.team.Players {
		totalAge += p.Age()
	}

	return float64(totalAge) / float64(len(sm.team.Players))
}

// GetSquadValue calculates total squad value
func (sm *SquadManager) GetSquadValue() int64 {
	var total int64
	for _, p := range sm.team.Players {
		total += p.MarketValue
	}
	return total
}

// GetWageBill calculates total weekly wages
func (sm *SquadManager) GetWageBill() int64 {
	var total int64
	for _, p := range sm.team.Players {
		total += p.Wage
	}
	return total
}

// GetYouthProspects returns players under 21
func (sm *SquadManager) GetYouthProspects() []player.Player {
	prospects := []player.Player{}
	for _, p := range sm.team.Players {
		if p.Age() < 21 {
			prospects = append(prospects, p)
		}
	}
	return prospects
}

// GetVeterans returns players over 30
func (sm *SquadManager) GetVeterans() []player.Player {
	veterans := []player.Player{}
	for _, p := range sm.team.Players {
		if p.Age() > 30 {
			veterans = append(veterans, p)
		}
	}
	return veterans
}

// GetInjuredPlayers returns all injured players
func (sm *SquadManager) GetInjuredPlayers() []player.Player {
	injured := []player.Player{}
	for _, p := range sm.team.Players {
		if p.Status == player.StatusInjured {
			injured = append(injured, p)
		}
	}
	return injured
}

// GetSuspendedPlayers returns all suspended players
func (sm *SquadManager) GetSuspendedPlayers() []player.Player {
	suspended := []player.Player{}
	for _, p := range sm.team.Players {
		if p.Status == player.StatusSuspended {
			suspended = append(suspended, p)
		}
	}
	return suspended
}

// RecommendLineup suggests best lineup for formation
func (sm *SquadManager) RecommendLineup(formation Formation) (*Lineup, error) {
	available := sm.team.GetAvailablePlayers()
	requirements := formation.GetPositionRequirements()

	lineup := &Lineup{
		Formation:   formation,
		Starters:    []player.PlayerID{},
		Positions:   []player.Position{},
		Substitutes: []player.PlayerID{},
	}

	// Track used players
	used := make(map[player.PlayerID]bool)

	// Fill positions by priority: GK, DEF, MID, FWD
	for _, pos := range []player.Position{
		player.PositionGK,
		player.PositionDEF,
		player.PositionMID,
		player.PositionFWD,
	} {
		count := requirements[pos]
		candidates := sm.getBestCandidates(available, pos, used)

		for i := 0; i < count && i < len(candidates); i++ {
			lineup.Starters = append(lineup.Starters, candidates[i].ID)
			lineup.Positions = append(lineup.Positions, pos)
			used[candidates[i].ID] = true
		}
	}

	// Check if we have enough players
	if len(lineup.Starters) < 11 {
		return nil, common.ErrInsufficientPlayers
	}

	// Fill substitutes
	for _, p := range available {
		if !used[p.ID] && len(lineup.Substitutes) < 7 {
			lineup.Substitutes = append(lineup.Substitutes, p.ID)
		}
	}

	// Set captain (highest leadership/experience)
	if captain := sm.selectCaptain(lineup.Starters); captain != nil {
		lineup.Captain = *captain
	}

	return lineup, nil
}

// getBestCandidates returns best players for position
func (sm *SquadManager) getBestCandidates(available []player.Player, pos player.Position, used map[player.PlayerID]bool) []player.Player {
	candidates := []player.Player{}

	for _, p := range available {
		if !used[p.ID] && p.CanPlayPosition(pos) {
			candidates = append(candidates, p)
		}
	}

	// Sort by position-specific rating
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].GetOverallRating() > candidates[j].GetOverallRating()
	})

	return candidates
}

// selectCaptain chooses captain from starters
func (sm *SquadManager) selectCaptain(starters []player.PlayerID) *player.PlayerID {
	if sm.team.Captain != nil {
		// Check if current captain is starting
		for _, id := range starters {
			if id == *sm.team.Captain {
				return sm.team.Captain
			}
		}
	}

	// Otherwise pick most experienced starter
	var bestPlayer *player.Player
	var bestScore float64

	for _, id := range starters {
		if p, err := sm.team.GetPlayer(id); err == nil {
			// Score based on age, experience, and leadership
			score := float64(p.Age()) + float64(p.CareerStats.TotalMatches)/10
			if score > bestScore {
				bestScore = score
				bestPlayer = p
			}
		}
	}

	if bestPlayer != nil {
		return &bestPlayer.ID
	}

	return nil
}
