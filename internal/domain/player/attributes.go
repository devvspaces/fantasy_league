// domain/player/attributes.go
package player

// Attributes represents player attributes (0-100 scale)
type Attributes struct {
	// Overall
	Quality int // Q: Overall quality

	// Technical
	Keeping  int // Kp: Goalkeeping
	Tackling int // Tk: Defensive ability
	Passing  int // Pa: Passing ability
	Shooting int // Sh: Shooting ability
	Heading  int // He: Heading ability

	// Physical
	Speed   int // Sp: Speed
	Stamina int // St: Stamina/fitness

	// Mental
	Perception  int // Pe: Awareness/vision
	BallControl int // Bc: Technical ability

	// Hidden attributes (affect development and consistency)
	Consistency      int // How consistent performances are
	ImportantMatches int // Performance in big games
	Potential        int // Maximum potential ability
	Ambition         int // Drive to improve
	Professionalism  int // Training attitude
}

// NewDefaultAttributes creates default attributes based on position
func NewDefaultAttributes(position Position) Attributes {
	base := Attributes{
		Quality:          65,
		Consistency:      70,
		ImportantMatches: 70,
		Potential:        75,
		Ambition:         70,
		Professionalism:  70,
	}

	switch position {
	case PositionGK:
		base.Keeping = 70
		base.Tackling = 20
		base.Passing = 50
		base.Shooting = 10
		base.Heading = 30
		base.Speed = 40
		base.Stamina = 70
		base.Perception = 65
		base.BallControl = 30
	case PositionDEF:
		base.Keeping = 20
		base.Tackling = 70
		base.Passing = 55
		base.Shooting = 35
		base.Heading = 65
		base.Speed = 65
		base.Stamina = 75
		base.Perception = 60
		base.BallControl = 50
	case PositionMID:
		base.Keeping = 20
		base.Tackling = 55
		base.Passing = 70
		base.Shooting = 55
		base.Heading = 50
		base.Speed = 70
		base.Stamina = 80
		base.Perception = 70
		base.BallControl = 70
	case PositionFWD:
		base.Keeping = 20
		base.Tackling = 30
		base.Passing = 60
		base.Shooting = 75
		base.Heading = 60
		base.Speed = 75
		base.Stamina = 70
		base.Perception = 65
		base.BallControl = 70
	}

	return base
}

// GetGoalkeeperRating calculates GK overall rating
func (a *Attributes) GetGoalkeeperRating() int {
	return int(
		float64(a.Keeping)*0.5 +
			float64(a.Speed)*0.1 +
			float64(a.Perception)*0.2 +
			float64(a.Stamina)*0.1 +
			float64(a.Passing)*0.1,
	)
}

// GetDefenderRating calculates DEF overall rating
func (a *Attributes) GetDefenderRating() int {
	return int(
		float64(a.Tackling)*0.3 +
			float64(a.Heading)*0.2 +
			float64(a.Speed)*0.15 +
			float64(a.Stamina)*0.15 +
			float64(a.Passing)*0.1 +
			float64(a.Perception)*0.1,
	)
}

// GetMidfielderRating calculates MID overall rating
func (a *Attributes) GetMidfielderRating() int {
	return int(
		float64(a.Passing)*0.25 +
			float64(a.BallControl)*0.2 +
			float64(a.Perception)*0.15 +
			float64(a.Stamina)*0.15 +
			float64(a.Tackling)*0.15 +
			float64(a.Shooting)*0.1,
	)
}

// GetForwardRating calculates FWD overall rating
func (a *Attributes) GetForwardRating() int {
	return int(
		float64(a.Shooting)*0.3 +
			float64(a.BallControl)*0.2 +
			float64(a.Speed)*0.2 +
			float64(a.Heading)*0.15 +
			float64(a.Perception)*0.15,
	)
}

// CanImprove checks if attribute can still improve
func (a *Attributes) CanImprove(attribute string, currentAge int) bool {
	// Physical attributes peak earlier
	physicalPeakAge := 28
	technicalPeakAge := 32

	switch attribute {
	case "Speed", "Stamina":
		return currentAge < physicalPeakAge
	case "Perception", "Passing", "BallControl":
		return currentAge < technicalPeakAge
	default:
		return currentAge < 30
	}
}
