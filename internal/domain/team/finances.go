// domain/team/finances.go
package team

import (
	"time"
)

// FinancialManager handles team finances
type FinancialManager struct {
	team *Team
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          string
	Type        TransactionType
	Amount      int64
	Description string
	Date        time.Time
	PlayerID    string // For transfer/wage transactions
}

// TransactionType represents types of financial transactions
type TransactionType string

const (
	TransactionTransferIn  TransactionType = "transfer_in"
	TransactionTransferOut TransactionType = "transfer_out"
	TransactionWages       TransactionType = "wages"
	TransactionTicketSales TransactionType = "ticket_sales"
	TransactionSponsorship TransactionType = "sponsorship"
	TransactionPrizeMoney  TransactionType = "prize_money"
	TransactionOther       TransactionType = "other"
)

// NewFinancialManager creates a financial manager
func NewFinancialManager(team *Team) *FinancialManager {
	return &FinancialManager{team: team}
}

// CanAffordTransfer checks if team can afford a transfer
func (fm *FinancialManager) CanAffordTransfer(fee int64, wages int64) bool {
	if fee > fm.team.Budget {
		return false
	}

	// Check wage budget
	currentWages := fm.GetTotalWages()
	if currentWages+wages > fm.team.WageBudget {
		return false
	}

	return true
}

// GetTotalWages calculates total weekly wages
func (fm *FinancialManager) GetTotalWages() int64 {
	var total int64
	for _, p := range fm.team.Players {
		total += p.Wage
	}
	return total
}

// GetWageBudgetRemaining calculates remaining wage budget
func (fm *FinancialManager) GetWageBudgetRemaining() int64 {
	return fm.team.WageBudget - fm.GetTotalWages()
}

// ProcessMatchRevenue calculates match day income
func (fm *FinancialManager) ProcessMatchRevenue(attendance int, isHome bool) int64 {
	if !isHome {
		return 0 // Away teams typically don't get gate receipts
	}

	// Simple calculation: average ticket price * attendance
	avgTicketPrice := int64(30) // Base price

	// Adjust for stadium utilization
	utilization := float64(attendance) / float64(fm.team.Stadium.Capacity)
	if utilization > 0.9 {
		avgTicketPrice = int64(float64(avgTicketPrice) * 1.2) // Premium pricing
	}

	revenue := avgTicketPrice * int64(attendance)

	// Additional revenue (concessions, parking, etc.)
	revenue = int64(float64(revenue) * 1.3)

	return revenue
}

// CalculateSeasonBudget estimates budget for next season
func (fm *FinancialManager) CalculateSeasonBudget(leaguePosition int, cupProgress string) {
	baseBudget := int64(10000000) // 10M base

	// League position bonus
	if leaguePosition <= 3 {
		baseBudget *= 3
	} else if leaguePosition <= 6 {
		baseBudget *= 2
	} else if leaguePosition <= 10 {
		baseBudget = int64(float64(baseBudget) * 1.5)
	}

	// Cup progress bonus
	switch cupProgress {
	case "winner":
		baseBudget += 5000000
	case "final":
		baseBudget += 3000000
	case "semi":
		baseBudget += 1500000
	case "quarter":
		baseBudget += 500000
	}

	fm.team.Budget = baseBudget
	fm.team.WageBudget = baseBudget / 52 // Weekly wage budget
}
