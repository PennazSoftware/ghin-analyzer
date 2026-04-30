package main

import "strings"

// calculateNetScoreForHole calculates the net score for a single golf hole.
//
// Parameters:
//
//	grossScore int: The player's actual score for the hole.
//	holeStrokeIndex int: The Stroke Index (difficulty ranking) of the hole (1-18).
//	courseHandicap int: The player's Course Handicap (a whole number).
//
// Returns:
//
//	int: The calculated net score for the hole.
func calculateNetScoreForHole(grossScore, holeStrokeIndex, courseHandicap int) int {
	// Determine the number of strokes the player receives on this specific hole.
	strokesReceivedOnHole := 0

	// Handle "plus" handicaps (where the player adds strokes back)
	if courseHandicap < 0 {
		// A plus handicap means the player gives strokes back, starting from the easiest holes (highest stroke index).
		// We add strokes to their gross score if the hole's Stroke Index is greater than or equal to (18 - abs(courseHandicap) + 1).
		// Example: A +2 handicap gives back 2 strokes on the 17th and 18th hardest holes (Stroke Index 17 and 18).
		// 18 - 2 + 1 = 17. So on holes with SI >= 17, they add a stroke.
		// For a +1 handicap, 18 - 1 + 1 = 18. So on SI 18, they add a stroke.
		if holeStrokeIndex >= (18 - abs(courseHandicap) + 1) {
			strokesReceivedOnHole = 1 // Add 1 stroke for each applicable hole
		}
	} else { // Standard handicaps (where the player subtracts strokes)
		// Players receive one stroke for each hole whose Stroke Index is less than or equal to their Course Handicap.
		if holeStrokeIndex <= courseHandicap {
			strokesReceivedOnHole = 1
		}

		// If the Course Handicap is greater than 18, the player receives more than one stroke on certain holes.
		// They receive one stroke on every hole, plus an additional stroke for each hole where
		// the Stroke Index (resetting from 1) is less than or equal to (Course Handicap - 18).
		if courseHandicap > 18 {
			additionalStrokes := courseHandicap - 18
			if holeStrokeIndex <= additionalStrokes {
				strokesReceivedOnHole += 1
			}
		}
	}

	// Calculate the net score.
	// For "plus" handicaps, strokes are added. For standard handicaps, strokes are subtracted.
	if courseHandicap < 0 {
		return grossScore + strokesReceivedOnHole
	}
	return grossScore - strokesReceivedOnHole
}

// Helper function to get the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// convertCourseHandicap converts a course handicap string (which may include a '+' sign for plus handicaps)
// into an integer. Plus handicaps are represented as negative integers.
func convertCourseHandicap(courseHandicap string) int {
	if strings.Contains(courseHandicap, "+") {
		// Remove the '+' sign and convert to negative integer
		return -1 * atoi(strings.TrimPrefix(courseHandicap, "+"))
	}
	return atoi(courseHandicap)
}

func atoi(s string) int {
	var n int
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			break
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
