// Application constants

// FSRS states
export const FSRS_STATES = {
    NEW: 0,
    LEARNING: 1,
    REVIEW: 2,
    RELEARNING: 3
};

// Difficulty levels
export const DIFFICULTY_LEVELS = {
    TRIVIAL: 1,
    VERY_EASY: 2,
    EASY: 3,
    MEDIUM: 4,
    HARD: 5,
    VERY_HARD: 6,
    EXPERT: 7
};

// Default category
export const DEFAULT_CATEGORY = 'java';

// Helper function to normalize category to lowercase
export function normalizeCategory(category) {
    return (category || DEFAULT_CATEGORY).toLowerCase();
}
