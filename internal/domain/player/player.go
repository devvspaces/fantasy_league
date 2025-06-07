// domain/player/player.go
package player

import (
	"time"
)

// PlayerID represents a unique player identifier
type PlayerID string

// Position represents player positions
type Position string

const (
	PositionGK  Position = "GK"
	PositionDEF Position = "DEF"
	PositionMID Position = "MID"
	PositionFWD Position = "FWD"
)

// Status represents player availability
type Status string

const (
	StatusAvailable Status = "available"
	StatusInjured   Status = "injured"
	StatusSuspended Status = "suspended"
	StatusOnLoan    Status = "on_loan"
	StatusRetired   Status = "retired"
)

// Player represents a football player
type Player struct {
	ID          PlayerID
	FirstName   string
	LastName    string
	Nickname    string
	DateOfBirth time.Time
	Nationality string

	// Physical characteristics
	Height int // in cm
	Weight int // in kg

	// Playing information
	Position      Position
	PreferredFoot string // "left", "right", "both"
	ShirtNumber   int
	ContractUntil time.Time
	MarketValue   int64 // in currency units
	Wage          int64 // weekly wage

	// Current state
	Status  Status
	Fitness float64 // 0-100
	Morale  float64 // 0-100
	Form    float64 // 0-100

	// Attributes
	Attributes Attributes

	// Career stats
	CareerStats CareerStats

	// Team affiliation
	CurrentTeamID string

	// Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CareerStats tracks player's career statistics
type CareerStats struct {
	TotalMatches     int
	TotalGoals       int
	TotalAssists     int
	TotalYellowCards int
	TotalRedCards    int
	TotalCleanSheets int // for goalkeepers
	SeasonStats      []SeasonStats
}

// SeasonStats tracks stats for a specific season
type SeasonStats struct {
	SeasonID      string
	TeamID        string
	Matches       int
	Goals         int
	Assists       int
	YellowCards   int
	RedCards      int
	CleanSheets   int
	AverageRating float64
}

// NewPlayer creates a new player
func NewPlayer(id PlayerID, firstName, lastName string, position Position, dateOfBirth time.Time) *Player {
	return &Player{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: dateOfBirth,
		Position:    position,
		Status:      StatusAvailable,
		Fitness:     100,
		Morale:      75,
		Form:        70,
		Attributes:  NewDefaultAttributes(position),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Age calculates the player's current age
func (p *Player) Age() int {
	now := time.Now()
	years := now.Year() - p.DateOfBirth.Year()
	if now.YearDay() < p.DateOfBirth.YearDay() {
		years--
	}
	return years
}

// FullName returns the player's full name
func (p *Player) FullName() string {
	if p.Nickname != "" {
		return p.Nickname
	}
	return p.FirstName + " " + p.LastName
}

// IsAvailable checks if player can play
func (p *Player) IsAvailable() bool {
	return p.Status == StatusAvailable && p.Fitness >= 70
}

// CanPlayPosition checks if player can play in a given position
func (p *Player) CanPlayPosition(pos Position) bool {
	// Exact match
	if p.Position == pos {
		return true
	}

	// Some flexibility for similar positions
	switch p.Position {
	case PositionDEF:
		return pos == PositionDEF || (pos == PositionMID && p.Attributes.Passing > 60)
	case PositionMID:
		return pos == PositionMID || pos == PositionDEF || pos == PositionFWD
	case PositionFWD:
		return pos == PositionFWD || (pos == PositionMID && p.Attributes.Passing > 60)
	default:
		return false
	}
}

// GetOverallRating calculates overall rating based on position
func (p *Player) GetOverallRating() int {
	switch p.Position {
	case PositionGK:
		return p.Attributes.GetGoalkeeperRating()
	case PositionDEF:
		return p.Attributes.GetDefenderRating()
	case PositionMID:
		return p.Attributes.GetMidfielderRating()
	case PositionFWD:
		return p.Attributes.GetForwardRating()
	default:
		return p.Attributes.Quality
	}
}

// UpdateMatchStats updates player statistics after a match
func (p *Player) UpdateMatchStats(goals, assists, yellowCards, redCards int, rating float64) {
	p.CareerStats.TotalMatches++
	p.CareerStats.TotalGoals += goals
	p.CareerStats.TotalAssists += assists
	p.CareerStats.TotalYellowCards += yellowCards
	p.CareerStats.TotalRedCards += redCards

	// Update form based on performance
	p.updateForm(rating)
}

// updateForm adjusts player form based on recent performance
func (p *Player) updateForm(matchRating float64) {
	// Form is weighted average of recent performances
	weight := 0.3 // Weight of new performance
	p.Form = p.Form*(1-weight) + (matchRating*10)*weight

	// Ensure form stays within bounds
	if p.Form > 100 {
		p.Form = 100
	} else if p.Form < 0 {
		p.Form = 0
	}
}
