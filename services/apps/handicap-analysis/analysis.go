package main

func buildHoleAnalysis(scores []Scores) map[int]HoleAnalysis {
	holeAnalysis := make(map[int]HoleAnalysis)

	for _, score := range scores {
		courseHandicap := convertCourseHandicap(score.CourseHandicap)

		for _, hole := range score.HoleDetails {
			ha, exists := holeAnalysis[hole.HoleNumber]
			if !exists {
				ha = HoleAnalysis{
					HoleNumber:       hole.HoleNumber,
					Par:              hole.Par,
					HandicapAnalysis: make(map[int]HoleHandicapAnalysis),
				}
			}

			hha, exists := ha.HandicapAnalysis[courseHandicap]
			if !exists {
				hha = HoleHandicapAnalysis{
					Handicap: courseHandicap,
				}
			}

			hha.TotalGrossScore += hole.RawScore
			hha.TotalNetScore += calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap)
			hha.Count++
			hha.EagleCountGross += boolToInt(hole.RawScore == ha.Par-2)
			hha.BirdieCountGross += boolToInt(hole.RawScore == ha.Par-1)
			hha.ParCountGross += boolToInt(hole.RawScore == ha.Par)
			hha.BogeyCountGross += boolToInt(hole.RawScore == ha.Par+1)
			hha.DoubleBogeyCountGross += boolToInt(hole.RawScore == ha.Par+2)
			hha.TripleBogeyOrWorseCountGross += boolToInt(hole.RawScore >= ha.Par+3)
			hha.EagleCountNet += boolToInt(calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap) == ha.Par-2)
			hha.BirdieCountNet += boolToInt(calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap) == ha.Par-1)
			hha.ParCountNet += boolToInt(calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap) == ha.Par)
			hha.BogeyCountNet += boolToInt(calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap) == ha.Par+1)
			hha.DoubleBogeyCountNet += boolToInt(calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap) == ha.Par+2)
			hha.TripleBogeyOrWorseCountNet += boolToInt(calculateNetScoreForHole(hole.RawScore, hole.StrokeAllocation, courseHandicap) >= ha.Par+3)

			ha.HandicapAnalysis[courseHandicap] = hha
			holeAnalysis[hole.HoleNumber] = ha
		}
	}

	return holeAnalysis
}
