package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	scores, err := readAllData()
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		return
	}
	fmt.Printf("Total scores read: %d\n", len(scores))

	// Filter scores by Bear Creek course ID #1526
	bearCreekScores := filterByCourseID("1526", scores)
	fmt.Printf("Total Bear Creek scores: %d\n", len(bearCreekScores))

	// Further filter to only include scores with at least 18 holes played
	bearCreekScores = filterByNumberOfHolesPlayed(18, bearCreekScores)
	fmt.Printf("Bear Creek scores with at least 18 holes played: %d\n", len(bearCreekScores))

	// Build hole analysis
	holeAnalysis := buildHoleAnalysis(bearCreekScores)

	// Output analysis
	for holeNum := 1; holeNum <= 18; holeNum++ {
		if ha, exists := holeAnalysis[holeNum]; exists {
			fmt.Printf("Hole %d (Par %d):\n", ha.HoleNumber, ha.Par)
			for handicap := -10; handicap <= 54; handicap++ {
				if hha, exists := ha.HandicapAnalysis[handicap]; exists {
					avgGross := float64(hha.TotalGrossScore) / float64(hha.Count)
					avgNet := float64(hha.TotalNetScore) / float64(hha.Count)
					fmt.Printf("  Handicap %d: Avg Gross: %.2f, Avg Net: %.2f (from %d scores)\n", handicap, avgGross, avgNet, hha.Count)
				}
			}
		}
	}

	// Create CSV files for each hole's analysis
	err = createCSVFileByHole(holeAnalysis, "first_pass")
	if err != nil {
		fmt.Printf("Error creating CSV files: %v\n", err)
		return
	}
}

// readAllData reads all score data from the JSON files located in the "scores" directory.
func readAllData() ([]Scores, error) {
	var allScores []Scores

	files, err := os.ReadDir("scores")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join("scores", file.Name()))
			if err != nil {
				return nil, err
			}

			var scoreFile ScoreFile
			if err := json.Unmarshal(data, &scoreFile); err != nil {
				return nil, err
			}

			fmt.Printf("Read %d scores from %s\n", len(scoreFile.Scores), file.Name())

			allScores = append(allScores, scoreFile.Scores...)
		}
	}

	return allScores, nil
}
