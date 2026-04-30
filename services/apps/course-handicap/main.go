package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Score struct {
	FacilityName                 string   `json:"facility_name"`
	AdjustedGrossScore           int      `json:"adjusted_gross_score"`
	Differential                 float64  `json:"differential"`
	AdjustedScaledUpDifferential *float64 `json:"adjusted_scaled_up_differential"`
	CourseRating                 float64  `json:"course_rating"`
	SlopeRating                  int      `json:"slope_rating"`
	Used                         bool     `json:"-"`
}

func getEffectiveDifferential(score Score) float64 {
	if score.AdjustedScaledUpDifferential != nil && *score.AdjustedScaledUpDifferential != 0 {
		return *score.AdjustedScaledUpDifferential
	}
	return score.Differential
}

type Data struct {
	RevisionScores struct {
		Scores []Score `json:"scores"`
	} `json:"revision_scores"`
}

func main() {
	file, err := os.Open("./data/canessa.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var data Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Println("JSON data successfully read and parsed.")

	facilities := getUniqueFacilityNames(data.RevisionScores.Scores)
	if len(facilities) == 0 {
		log.Fatalf("No facilities found in the data.")
	}

	for _, facilityName := range facilities {
		filteredScores := filterScoresByFacility(data.RevisionScores.Scores, facilityName)
		sort.Slice(filteredScores, func(i, j int) bool {
			return filteredScores[i].Differential < filteredScores[j].Differential
		})

		usedScores := filteredScores
		if len(filteredScores) > 8 {
			usedScores = filteredScores[:8]
		}

		handicap := calculateHandicap(usedScores)
		markedScores := markUsedScores(filteredScores, usedScores)

		fmt.Printf("Facility: %s\n", facilityName)
		fmt.Printf("Calculated Handicap: %.2f\n", handicap)
		fmt.Println("Scores:")
		for _, score := range markedScores {
			usedIndicator := "No"
			if score.Used {
				usedIndicator = "Yes"
			}
			fmt.Printf("- Adjusted Gross Score: %d, Differential: %.2f, Used in Index: %s\n", score.AdjustedGrossScore, score.Differential, usedIndicator)
		}
	}
}

func calculateHandicap(scores []Score) float64 {
	if len(scores) == 0 {
		return 0
	}

	sort.Slice(scores, func(i, j int) bool {
		return getEffectiveDifferential(scores[i]) < getEffectiveDifferential(scores[j])
	})

	n := len(scores)
	if n > 8 {
		scores = scores[:8]
	}

	totalDifferential := 0.0
	for _, score := range scores {
		totalDifferential += getEffectiveDifferential(score)
	}

	averageDifferential := totalDifferential / float64(len(scores))
	return averageDifferential * 0.96 // USGA Handicap formula
}

func filterScoresByFacility(scores []Score, facilityName string) []Score {
	var filteredScores []Score
	for _, score := range scores {
		if score.FacilityName == facilityName {
			filteredScores = append(filteredScores, score)
		}
	}
	return filteredScores
}

func markUsedScores(scores []Score, usedScores []Score) []Score {
	usedMap := make(map[int]bool)
	for _, score := range usedScores {
		usedMap[score.AdjustedGrossScore] = true
	}

	for i := range scores {
		if usedMap[scores[i].AdjustedGrossScore] {
			scores[i].Used = true
		} else {
			scores[i].Used = false
		}
	}
	return scores
}

func getUniqueFacilityNames(scores []Score) []string {
	facilitySet := make(map[string]struct{})
	for _, score := range scores {
		facilitySet[score.FacilityName] = struct{}{}
	}

	var facilities []string
	for facility := range facilitySet {
		facilities = append(facilities, facility)
	}
	return facilities
}
