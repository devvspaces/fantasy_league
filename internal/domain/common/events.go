// domain/common/events.go
package common

import (
	"time"
)

// EventType represents the type of domain event
type EventType string

const (
	// Match events
	EventMatchScheduled EventType = "match.scheduled"
	EventMatchStarted   EventType = "match.started"
	EventMatchCompleted EventType = "match.completed"
	EventGoalScored     EventType = "match.goal_scored"
	EventCardIssued     EventType = "match.card_issued"

	// Player events
	EventPlayerInjured    EventType = "player.injured"
	EventPlayerRecovered  EventType = "player.recovered"
	EventPlayerSuspended  EventType = "player.suspended"
	EventPlayerTrained    EventType = "player.trained"
	EventPlayerProgressed EventType = "player.progressed"

	// Team events
	EventLineupSet        EventType = "team.lineup_set"
	EventTacticsChanged   EventType = "team.tactics_changed"
	EventFormationChanged EventType = "team.formation_changed"

	// Season events
	EventSeasonStarted     EventType = "season.started"
	EventSeasonCompleted   EventType = "season.completed"
	EventFixturesGenerated EventType = "season.fixtures_generated"
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	GetID() string
	GetType() EventType
	GetOccurredAt() time.Time
	GetAggregateID() string
}

// BaseEvent provides common event fields
type BaseEvent struct {
	ID          string
	Type        EventType
	OccurredAt  time.Time
	AggregateID string
}

func (e BaseEvent) GetID() string            { return e.ID }
func (e BaseEvent) GetType() EventType       { return e.Type }
func (e BaseEvent) GetOccurredAt() time.Time { return e.OccurredAt }
func (e BaseEvent) GetAggregateID() string   { return e.AggregateID }

// Match Events
type MatchScheduledEvent struct {
	BaseEvent
	HomeTeamID  string
	AwayTeamID  string
	ScheduledAt time.Time
}

type MatchCompletedEvent struct {
	BaseEvent
	HomeScore int
	AwayScore int
	Stats     map[string]interface{}
}

type GoalScoredEvent struct {
	BaseEvent
	MatchID  string
	PlayerID string
	TeamID   string
	Minute   int
	AssistBy string
}

// Player Events
type PlayerInjuredEvent struct {
	BaseEvent
	PlayerID     string
	InjuryType   string
	ExpectedDays int
}

type PlayerTrainedEvent struct {
	BaseEvent
	PlayerID       string
	TrainingType   string
	AttributeGains map[string]int
}

// Team Events
type LineupSetEvent struct {
	BaseEvent
	TeamID    string
	MatchID   string
	PlayerIDs []string
	Formation string
}

// Season Events
type SeasonStartedEvent struct {
	BaseEvent
	SeasonID  string
	LeagueID  string
	StartDate time.Time
	Teams     []string
}
