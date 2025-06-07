// domain/common/errors.go
package common

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common domain errors
var (
	ErrPlayerNotFound = DomainError{
		Code:    "PLAYER_NOT_FOUND",
		Message: "Player not found",
	}

	ErrTeamNotFound = DomainError{
		Code:    "TEAM_NOT_FOUND",
		Message: "Team not found",
	}

	ErrInvalidFormation = DomainError{
		Code:    "INVALID_FORMATION",
		Message: "Invalid formation",
	}

	ErrInsufficientPlayers = DomainError{
		Code:    "INSUFFICIENT_PLAYERS",
		Message: "Not enough players for lineup",
	}

	ErrPlayerUnavailable = DomainError{
		Code:    "PLAYER_UNAVAILABLE",
		Message: "Player is unavailable",
	}

	ErrMatchAlreadyPlayed = DomainError{
		Code:    "MATCH_ALREADY_PLAYED",
		Message: "Match has already been played",
	}

	ErrInvalidTactics = DomainError{
		Code:    "INVALID_TACTICS",
		Message: "Invalid tactical settings",
	}

	ErrSeasonNotActive = DomainError{
		Code:    "SEASON_NOT_ACTIVE",
		Message: "Season is not active",
	}

	ErrFixtureConflict = DomainError{
		Code:    "FIXTURE_CONFLICT",
		Message: "Fixture scheduling conflict",
	}
)
