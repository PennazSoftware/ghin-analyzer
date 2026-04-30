package main

// filterByCourseID filters the scores by the given course ID.
func filterByCourseID(courseID string, scores []Scores) []Scores {
	var filtered []Scores
	for _, score := range scores {
		if score.CourseID == courseID {
			filtered = append(filtered, score)
		}
	}
	return filtered
}

func filterByNumberOfHolesPlayed(minHoles int, scores []Scores) []Scores {
	var filtered []Scores
	for _, score := range scores {
		if score.NumberOfPlayedHoles >= minHoles {
			filtered = append(filtered, score)
		}
	}
	return filtered
}
