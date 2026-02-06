package fsrs

import (
	"math"
	"testing"
	"time"
)

// Test 1: NewDefaultParameters returns valid parameters
func TestNewDefaultParameters(t *testing.T) {
	params := NewDefaultParameters()

	if params.RequestRetention != DefaultRequestRetention {
		t.Errorf("Expected RequestRetention %v, got %v", DefaultRequestRetention, params.RequestRetention)
	}

	if len(params.W) != 17 {
		t.Errorf("Expected 17 weights, got %d", len(params.W))
	}

	// Check that all weights are reasonable values
	for i, w := range params.W {
		if math.IsNaN(w) || math.IsInf(w, 0) {
			t.Errorf("Weight %d is invalid: %v", i, w)
		}
	}
}

// Test 2: NewCard creates card with default values
func TestNewCard(t *testing.T) {
	card := NewCard()

	if card.State != 0 {
		t.Errorf("Expected new card state 0, got %d", card.State)
	}

	if card.Reps != 0 {
		t.Errorf("Expected 0 reps, got %d", card.Reps)
	}

	if card.Lapses != 0 {
		t.Errorf("Expected 0 lapses, got %d", card.Lapses)
	}

	if card.Stability != 0 {
		t.Errorf("Expected 0 stability, got %v", card.Stability)
	}
}

// Test 3: First review of new card initializes difficulty
func TestFirstReviewInitializesDifficulty(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	scheduleCards := params.Repeat(card, now)

	// All ratings should initialize difficulty
	if scheduleCards.Again.Card.Difficulty == 0 {
		t.Error("Again rating should initialize difficulty")
	}
	if scheduleCards.Good.Card.Difficulty == 0 {
		t.Error("Good rating should initialize difficulty")
	}
	if scheduleCards.Easy.Card.Difficulty == 0 {
		t.Error("Easy rating should initialize difficulty")
	}
}

// Test 4: First review sets stability
func TestFirstReviewSetsStability(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	scheduleCards := params.Repeat(card, now)

	if scheduleCards.Good.Card.Stability <= 0 {
		t.Error("Good rating should set positive stability")
	}
}

// Test 5: Rating affects difficulty progression
func TestRatingAffectsDifficulty(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	scheduleCards := params.Repeat(card, now)

	// Easy should have lower difficulty than Again
	easyDiff := scheduleCards.Easy.Card.Difficulty
	againDiff := scheduleCards.Again.Card.Difficulty

	if easyDiff >= againDiff {
		t.Errorf("Easy difficulty (%v) should be less than Again difficulty (%v)", easyDiff, againDiff)
	}
}

// Test 6: Successful review increases stability
func TestSuccessfulReviewIncreasesStability(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	// First review
	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card

	initialStability := card.Stability

	// Second review after interval
	now = now.AddDate(0, 0, card.ScheduledDays)
	scheduleCards = params.Repeat(card, now)
	card = scheduleCards.Good.Card

	if card.Stability <= initialStability {
		t.Errorf("Stability should increase after successful review. Initial: %v, After: %v", initialStability, card.Stability)
	}
}

// Test 7: Failed review (Again) increases lapses
func TestFailedReviewIncreasesLapses(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	// First review - good
	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card

	initialLapses := card.Lapses

	// Second review - fail
	now = now.AddDate(0, 0, 1)
	scheduleCards = params.Repeat(card, now)
	card = scheduleCards.Again.Card

	if card.Lapses != initialLapses+1 {
		t.Errorf("Expected lapses to increase by 1. Initial: %d, After: %d", initialLapses, card.Lapses)
	}
}

// Test 8: Interval calculation returns positive values
func TestIntervalCalculationPositive(t *testing.T) {
	params := NewDefaultParameters()

	stabilities := []float64{0.1, 1.0, 5.0, 10.0, 50.0}

	for _, s := range stabilities {
		interval := params.nextInterval(s)
		if interval < 1 {
			t.Errorf("Interval for stability %v should be at least 1, got %d", s, interval)
		}
		if interval > MaximumInterval {
			t.Errorf("Interval for stability %v exceeded maximum: got %d", s, interval)
		}
	}
}

// Test 9: Difficulty is constrained between 1 and 10
func TestDifficultyConstrained(t *testing.T) {
	params := NewDefaultParameters()

	testValues := []float64{-5.0, 0.0, 0.5, 1.0, 5.0, 10.0, 15.0, 100.0}

	for _, d := range testValues {
		constrained := params.constrainDifficulty(d)
		if constrained < 1.0 || constrained > 10.0 {
			t.Errorf("Difficulty %v should be constrained to [1, 10], got %v", d, constrained)
		}
	}
}

// Test 10: Reps increment with each review
func TestRepsIncrement(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	expectedReps := 0
	for i := 0; i < 5; i++ {
		expectedReps++
		scheduleCards := params.Repeat(card, now)
		card = scheduleCards.Good.Card

		if card.Reps != expectedReps {
			t.Errorf("After %d reviews, expected %d reps, got %d", i+1, expectedReps, card.Reps)
		}

		now = now.AddDate(0, 0, card.ScheduledDays)
	}
}

// Test 11: State transitions correctly
func TestStateTransitions(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	// New card (state 0)
	if card.State != 0 {
		t.Errorf("New card should have state 0, got %d", card.State)
	}

	// First review - moves to Learning (state 1)
	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card
	if card.State != 1 {
		t.Errorf("After first review, card should be in Learning state (1), got %d", card.State)
	}

	// Successful review - moves to Review (state 2)
	now = now.AddDate(0, 0, card.ScheduledDays)
	scheduleCards = params.Repeat(card, now)
	card = scheduleCards.Good.Card
	if card.State != 2 {
		t.Errorf("After successful review, card should be in Review state (2), got %d", card.State)
	}

	// Failed review - moves to Relearning (state 3)
	now = now.AddDate(0, 0, card.ScheduledDays)
	scheduleCards = params.Repeat(card, now)
	card = scheduleCards.Again.Card
	if card.State != 3 {
		t.Errorf("After failed review, card should be in Relearning state (3), got %d", card.State)
	}
}

// Test 12: Due date is set correctly
func TestDueDateSet(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card

	expectedDue := now.AddDate(0, 0, card.ScheduledDays)

	// Allow for small time differences due to calculation
	diff := card.Due.Sub(expectedDue).Abs()
	if diff > time.Second {
		t.Errorf("Due date not set correctly. Expected around %v, got %v", expectedDue, card.Due)
	}
}

// Test 13: Elapsed days calculated correctly
func TestElapsedDaysCalculation(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	// First review
	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card
	card.LastReview = now

	// Second review 5 days later
	now = now.AddDate(0, 0, 5)
	scheduleCards = params.Repeat(card, now)

	if scheduleCards.Good.Card.ElapsedDays != 5 {
		t.Errorf("Expected 5 elapsed days, got %d", scheduleCards.Good.Card.ElapsedDays)
	}
}

// Test 14: Hard rating provides intermediate interval
func TestHardRatingIntermediateInterval(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	scheduleCards := params.Repeat(card, now)

	againInterval := scheduleCards.Again.Card.ScheduledDays
	hardInterval := scheduleCards.Hard.Card.ScheduledDays
	goodInterval := scheduleCards.Good.Card.ScheduledDays
	easyInterval := scheduleCards.Easy.Card.ScheduledDays

	// Hard should be at least as long as Again and no longer than Good
	if hardInterval < againInterval || hardInterval > goodInterval {
		t.Errorf("Hard interval (%d) should be between or equal to Again (%d) and Good (%d)", hardInterval, againInterval, goodInterval)
	}

	// Good should be less than or equal to Easy
	if goodInterval > easyInterval {
		t.Errorf("Good interval (%d) should be less than or equal to Easy (%d)", goodInterval, easyInterval)
	}
}

// Test 15: GetRetrievability returns value between 0 and 1
func TestGetRetrievability(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	// First review to set up card
	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card

	// Check retrievability at different times
	testDays := []int{0, 1, 5, 10, 30}
	for _, days := range testDays {
		testTime := card.LastReview.AddDate(0, 0, days)
		retrievability := params.GetRetrievability(card, testTime)

		if retrievability < 0 || retrievability > 1 {
			t.Errorf("Retrievability should be between 0 and 1, got %v for %d days", retrievability, days)
		}
	}
}

// Test 16: Retrievability decreases over time
func TestRetrievabilityDecreasesOverTime(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	// First review
	scheduleCards := params.Repeat(card, now)
	card = scheduleCards.Good.Card

	r1 := params.GetRetrievability(card, card.LastReview.AddDate(0, 0, 1))
	r2 := params.GetRetrievability(card, card.LastReview.AddDate(0, 0, 5))
	r3 := params.GetRetrievability(card, card.LastReview.AddDate(0, 0, 10))

	if r1 <= r2 || r2 <= r3 {
		t.Errorf("Retrievability should decrease over time: r1=%v, r2=%v, r3=%v", r1, r2, r3)
	}
}

// Test 17: Easy bonus increases stability more
func TestEasyBonusIncreasesStability(t *testing.T) {
	params := NewDefaultParameters()
	card1 := NewCard()
	card2 := NewCard()
	now := time.Now()

	// First review - both good
	scheduleCards1 := params.Repeat(card1, now)
	card1 = scheduleCards1.Good.Card

	scheduleCards2 := params.Repeat(card2, now)
	card2 = scheduleCards2.Good.Card

	// Second review - one good, one easy
	now = now.AddDate(0, 0, card1.ScheduledDays)

	scheduleCards1 = params.Repeat(card1, now)
	goodCard := scheduleCards1.Good.Card

	scheduleCards2 = params.Repeat(card2, now)
	easyCard := scheduleCards2.Easy.Card

	if easyCard.Stability <= goodCard.Stability {
		t.Errorf("Easy card stability (%v) should be greater than Good card stability (%v)", easyCard.Stability, goodCard.Stability)
	}
}

// Test 18: Hard penalty decreases stability gain
func TestHardPenaltyDecreasesStabilityGain(t *testing.T) {
	params := NewDefaultParameters()
	card1 := NewCard()
	card2 := NewCard()
	now := time.Now()

	// First review - both good
	scheduleCards1 := params.Repeat(card1, now)
	card1 = scheduleCards1.Good.Card

	scheduleCards2 := params.Repeat(card2, now)
	card2 = scheduleCards2.Good.Card

	// Second review - one good, one hard
	now = now.AddDate(0, 0, card1.ScheduledDays)

	scheduleCards1 = params.Repeat(card1, now)
	goodCard := scheduleCards1.Good.Card

	scheduleCards2 = params.Repeat(card2, now)
	hardCard := scheduleCards2.Hard.Card

	if hardCard.Stability >= goodCard.Stability {
		t.Errorf("Hard card stability (%v) should be less than Good card stability (%v)", hardCard.Stability, goodCard.Stability)
	}
}

// Test 19: ReviewLog captures correct information
func TestReviewLogCaptures(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	scheduleCards := params.Repeat(card, now)

	log := scheduleCards.Good.ReviewLog

	if log.Rating != RatingGood {
		t.Errorf("Expected rating %d, got %d", RatingGood, log.Rating)
	}

	if log.ScheduledDays != scheduleCards.Good.Card.ScheduledDays {
		t.Errorf("Log scheduled days (%d) doesn't match card (%d)", log.ScheduledDays, scheduleCards.Good.Card.ScheduledDays)
	}

	if log.Review != now {
		t.Errorf("Review time not captured correctly")
	}
}

// Test 20: Long-term spaced repetition increases intervals
func TestLongTermSpacedRepetition(t *testing.T) {
	params := NewDefaultParameters()
	card := NewCard()
	now := time.Now()

	intervals := make([]int, 0, 10)

	// Simulate 10 successful reviews
	for i := 0; i < 10; i++ {
		scheduleCards := params.Repeat(card, now)
		card = scheduleCards.Good.Card
		intervals = append(intervals, card.ScheduledDays)
		now = now.AddDate(0, 0, card.ScheduledDays)
	}

	// Check that intervals generally increase
	increasingCount := 0
	for i := 1; i < len(intervals); i++ {
		if intervals[i] > intervals[i-1] {
			increasingCount++
		}
	}

	// Most intervals should be increasing (allow some variation)
	if increasingCount < 6 {
		t.Errorf("Expected most intervals to increase, only %d out of %d increased. Intervals: %v", increasingCount, len(intervals)-1, intervals)
	}
}
