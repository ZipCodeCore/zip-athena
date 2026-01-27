package fsrs

import (
	"math"
	"time"
)

// FSRS (Free Spaced Repetition Scheduler) algorithm implementation
// Based on the FSRS-4.5 algorithm for optimizing spaced repetition

const (
	// Default FSRS parameters (optimized values)
	DefaultRequestRetention = 0.9
	MaximumInterval        = 36500 // ~100 years in days
	
	// Rating grades
	RatingAgain = 1
	RatingHard  = 2
	RatingGood  = 3
	RatingEasy  = 4
)

// Parameters represents the FSRS algorithm parameters
type Parameters struct {
	RequestRetention float64
	W                []float64 // 17 weights for the algorithm
}

// Card represents a flashcard with FSRS scheduling information
type Card struct {
	Due           time.Time
	Stability     float64
	Difficulty    float64
	ElapsedDays   int
	ScheduledDays int
	Reps          int
	Lapses        int
	State         int // 0=New, 1=Learning, 2=Review, 3=Relearning
	LastReview    time.Time
}

// ReviewLog represents a single review record
type ReviewLog struct {
	Rating      int
	ScheduledDays int
	ElapsedDays   int
	Review       time.Time
	State        int
}

// SchedulingInfo contains scheduling information for all rating options
type SchedulingInfo struct {
	Card  Card
	ReviewLog ReviewLog
}

// SchedulingCards contains the next scheduling for each rating
type SchedulingCards struct {
	Again SchedulingInfo
	Hard  SchedulingInfo
	Good  SchedulingInfo
	Easy  SchedulingInfo
}

// NewDefaultParameters returns default FSRS parameters
func NewDefaultParameters() Parameters {
	return Parameters{
		RequestRetention: DefaultRequestRetention,
		W: []float64{
			0.4, 0.6, 2.4, 5.8, 4.93, 0.94, 0.86, 0.01, 1.49, 0.14, 0.94, 2.18, 0.05, 0.34, 1.26, 0.29, 2.61,
		},
	}
}

// NewCard creates a new card with default values
func NewCard() Card {
	return Card{
		Due:           time.Now(),
		Stability:     0,
		Difficulty:    0,
		ElapsedDays:   0,
		ScheduledDays: 0,
		Reps:          0,
		Lapses        0,
		State:         0, // New
		LastReview:    time.Time{},
	}
}

// Repeat schedules the next review based on the card's current state and rating
func (p *Parameters) Repeat(card Card, now time.Time) SchedulingCards {
	if card.State == 0 {
		// New card
		card.ElapsedDays = 0
	} else {
		card.ElapsedDays = int(now.Sub(card.LastReview).Hours() / 24)
	}

	card.LastReview = now
	card.Reps++

	s := SchedulingCards{}
	s.Again = p.schedule(card, now, RatingAgain)
	s.Hard = p.schedule(card, now, RatingHard)
	s.Good = p.schedule(card, now, RatingGood)
	s.Easy = p.schedule(card, now, RatingEasy)

	return s
}

func (p *Parameters) schedule(card Card, now time.Time, rating int) SchedulingInfo {
	newCard := card
	
	if card.State == 0 {
		// New card - initialize difficulty
		newCard.Difficulty = p.initDifficulty(rating)
		newCard.Stability = p.initStability(rating)
		newCard.State = 1 // Learning
	} else {
		// Existing card
		if rating == RatingAgain {
			// Failed review
			newCard.Stability = p.nextForgetStability(card.Difficulty, card.Stability, card.ElapsedDays)
			newCard.Difficulty = p.nextDifficulty(card.Difficulty, rating)
			newCard.State = 3 // Relearning
			newCard.Lapses++
		} else {
			// Successful review
			newCard.Stability = p.nextRecallStability(card.Difficulty, card.Stability, card.ElapsedDays, rating)
			newCard.Difficulty = p.nextDifficulty(card.Difficulty, rating)
			newCard.State = 2 // Review
		}
	}

	interval := p.nextInterval(newCard.Stability)
	newCard.ScheduledDays = interval
	newCard.Due = now.AddDate(0, 0, interval)

	return SchedulingInfo{
		Card: newCard,
		ReviewLog: ReviewLog{
			Rating:        rating,
			ScheduledDays: interval,
			ElapsedDays:   card.ElapsedDays,
			Review:        now,
			State:         newCard.State,
		},
	}
}

func (p *Parameters) initStability(rating int) float64 {
	return math.Max(p.W[rating-1], 0.1)
}

func (p *Parameters) initDifficulty(rating int) float64 {
	difficulty := p.W[4] - (float64(rating-3) * p.W[5])
	return p.constrainDifficulty(difficulty)
}

func (p *Parameters) nextInterval(stability float64) int {
	interval := int(math.Round(stability * 9 * (1/p.RequestRetention - 1)))
	if interval < 1 {
		return 1
	}
	if interval > MaximumInterval {
		return MaximumInterval
	}
	return interval
}

func (p *Parameters) nextDifficulty(difficulty float64, rating int) float64 {
	nextD := difficulty - p.W[6]*(float64(rating)-3)
	return p.constrainDifficulty(p.meanReversion(difficulty, nextD))
}

func (p *Parameters) meanReversion(init, current float64) float64 {
	return p.W[7]*init + (1-p.W[7])*current
}

func (p *Parameters) constrainDifficulty(difficulty float64) float64 {
	return math.Min(math.Max(difficulty, 1), 10)
}

func (p *Parameters) nextRecallStability(d, s float64, elapsed int, rating int) float64 {
	hardPenalty := 1.0
	if rating == RatingHard {
		hardPenalty = p.W[15]
	}
	easyBonus := 1.0
	if rating == RatingEasy {
		easyBonus = p.W[16]
	}

	retrievability := math.Exp(math.Log(0.9) * float64(elapsed) / s)
	newS := s * (1 + math.Exp(p.W[8]) *
		(11-d) *
		math.Pow(s, -p.W[9]) *
		(math.Exp((1-retrievability)*p.W[10]) - 1) *
		hardPenalty *
		easyBonus)

	return newS
}

func (p *Parameters) nextForgetStability(d, s float64, elapsed int) float64 {
	retrievability := math.Exp(math.Log(0.9) * float64(elapsed) / s)
	newS := p.W[11] * math.Pow(d, -p.W[12]) *
		(math.Pow(s+1, p.W[13]) - 1) *
		math.Exp((1-retrievability)*p.W[14])
	return math.Max(newS, 0.1)
}

// GetRetrievability calculates the current probability of recall
func (p *Parameters) GetRetrievability(card Card, now time.Time) float64 {
	if card.State == 0 {
		return 0
	}
	elapsedDays := int(now.Sub(card.LastReview).Hours() / 24)
	return math.Exp(math.Log(0.9) * float64(elapsedDays) / card.Stability)
}
