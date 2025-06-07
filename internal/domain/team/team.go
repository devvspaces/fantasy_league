// domain/team/team.go
package team

import (
	"fmt"
	"time"

	"github.com/devvspaces/fantasy_league/internal/domain/common"
	"github.com/devvspaces/fantasy_league/internal/domain/player"
)

// TeamID represents a unique team identifier
type TeamID string

// Team represents a football team
type Team struct {
	ID        TeamID
	Name      string
	ShortName string
	Founded   int
	Stadium   Stadium

	// Squad
	Players     []player.Player
	Captain     *player.PlayerID
	ViceCaptain *player.PlayerID

	// Tactical setup
	Formation Formation
	Tactics   TeamTactics

	// Staff
	ManagerName string

	// Financials
	Budget     int64
	WageBudget int64

	// Performance
	CurrentForm []MatchResult // Last 5 matches
	SeasonStats TeamSeasonStats

	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Stadium represents team's home ground
type Stadium struct {
	Name      string
	Capacity  int
	City      string
	Country   string
	PitchType string // "grass", "artificial"
}

// MatchResult represents a recent match outcome
type MatchResult struct {
	MatchID      string
	Opponent     string
	IsHome       bool
	GoalsFor     int
	GoalsAgainst int
	Result       string // "W", "D", "L"
}

// TeamSeasonStats tracks seasonal performance
type TeamSeasonStats struct {
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	Points         int
	LeaguePosition int
}

// NewTeam creates a new team
func NewTeam(id TeamID, name string, stadium Stadium) *Team {
	return &Team{
		ID:        id,
		Name:      name,
		ShortName: name[:3], // Simple default
		Stadium:   stadium,
		Formation: FormationDefault,
		Tactics:   DefaultTactics(),
		Players:   []player.Player{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddPlayer adds a player to the squad
func (t *Team) AddPlayer(p player.Player) error {
	// Check squad size limit
	if len(t.Players) >= 30 {
		return fmt.Errorf("squad size limit reached")
	}

	// Check if player already exists
	for _, existing := range t.Players {
		if existing.ID == p.ID {
			return fmt.Errorf("player already in squad")
		}
	}

	t.Players = append(t.Players, p)
	t.UpdatedAt = time.Now()
	return nil
}

// RemovePlayer removes a player from the squad
func (t *Team) RemovePlayer(playerID player.PlayerID) error {
	for i, p := range t.Players {
		if p.ID == playerID {
			// Remove player
			t.Players = append(t.Players[:i], t.Players[i+1:]...)

			// Clear captain if needed
			if t.Captain != nil && *t.Captain == playerID {
				t.Captain = nil
			}
			if t.ViceCaptain != nil && *t.ViceCaptain == playerID {
				t.ViceCaptain = nil
			}

			t.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("player not found in squad")
}

// GetPlayer retrieves a player by ID
func (t *Team) GetPlayer(playerID player.PlayerID) (*player.Player, error) {
	for _, p := range t.Players {
		if p.ID == playerID {
			return &p, nil
		}
	}
	return nil, common.ErrPlayerNotFound
}

// GetAvailablePlayers returns players available for selection
func (t *Team) GetAvailablePlayers() []player.Player {
	available := []player.Player{}
	for _, p := range t.Players {
		if p.IsAvailable() {
			available = append(available, p)
		}
	}
	return available
}

// GetPlayersByPosition returns players who can play in a position
func (t *Team) GetPlayersByPosition(pos player.Position) []player.Player {
	players := []player.Player{}
	for _, p := range t.Players {
		if p.CanPlayPosition(pos) {
			players = append(players, p)
		}
	}
	return players
}

// ValidateLineup checks if a lineup is valid
func (t *Team) ValidateLineup(lineup Lineup) error {
	// Check if we have 11 players
	if len(lineup.Starters) != 11 {
		return common.ErrInsufficientPlayers
	}

	// Check if all players are available
	for _, playerID := range lineup.Starters {
		player, err := t.GetPlayer(playerID)
		if err != nil {
			return err
		}
		if !player.IsAvailable() {
			return common.ErrPlayerUnavailable
		}
	}

	// Check formation requirements
	if !lineup.Formation.IsValid() {
		return common.ErrInvalidFormation
	}

	// Validate positions match formation
	return t.validateFormationPositions(lineup)
}

// validateFormationPositions ensures players are in correct positions
func (t *Team) validateFormationPositions(lineup Lineup) error {
	requiredPositions := lineup.Formation.GetPositionRequirements()

	// Count positions in lineup
	positionCount := make(map[player.Position]int)
	for i, playerID := range lineup.Starters {
		p, _ := t.GetPlayer(playerID)
		assignedPos := lineup.Positions[i]

		if !p.CanPlayPosition(assignedPos) {
			return fmt.Errorf("player %s cannot play %s", p.FullName(), assignedPos)
		}

		positionCount[assignedPos]++
	}

	// Check requirements met
	for pos, required := range requiredPositions {
		if positionCount[pos] != required {
			return fmt.Errorf("formation requires %d %s, got %d", required, pos, positionCount[pos])
		}
	}

	return nil
}

// GetTeamStrength calculates overall team strength
func (t *Team) GetTeamStrength() float64 {
	if len(t.Players) == 0 {
		return 0
	}

	totalStrength := 0.0
	count := 0

	// Get best 11 players
	for _, p := range t.GetBestEleven() {
		totalStrength += float64(p.GetOverallRating())
		count++
	}

	if count == 0 {
		return 0
	}

	return totalStrength / float64(count)
}

// GetBestEleven returns the strongest possible lineup
func (t *Team) GetBestEleven() []player.Player {
	available := t.GetAvailablePlayers()
	if len(available) < 11 {
		return available
	}

	// Simple selection: highest rated players per position
	bestEleven := []player.Player{}

	// Get required positions for current formation
	requirements := t.Formation.GetPositionRequirements()

	for pos, count := range requirements {
		candidates := []player.Player{}
		for _, p := range available {
			if p.CanPlayPosition(pos) {
				candidates = append(candidates, p)
			}
		}

		// Sort by rating and take required count
		// (simplified - in real implementation would use proper sorting)
		for i := 0; i < count && i < len(candidates); i++ {
			bestEleven = append(bestEleven, candidates[i])
		}
	}

	return bestEleven
}

// UpdateForm adds a match result to recent form
func (t *Team) UpdateForm(result MatchResult) {
	t.CurrentForm = append([]MatchResult{result}, t.CurrentForm...)
	if len(t.CurrentForm) > 5 {
		t.CurrentForm = t.CurrentForm[:5]
	}
}

// GetFormString returns form as string (e.g., "WWLDW")
func (t *Team) GetFormString() string {
	form := ""
	for _, result := range t.CurrentForm {
		form += result.Result
	}
	return form
}
