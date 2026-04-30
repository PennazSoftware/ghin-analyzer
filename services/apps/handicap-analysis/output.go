package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// createCSVFileByHole creates a CSV file for each hole's analysis data.
func createCSVFileByHole(holeAnalysis map[int]HoleAnalysis, description string) error {
	for holeNum, ha := range holeAnalysis {
		filename := fmt.Sprintf("results/holes/hole_%02d_%s.csv", holeNum, description)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		header := []string{"Handicap", "TotalNetScore", "TotalGrossScore", "AvgGrossScore", "AvgNetScore", "Count", "EagleCountGross", "BirdieCountGross", "ParCountGross", "BogeyCountGross", "DoubleBogeyCountGross", "TripleBogeyOrWorseCountGross", "EagleCountNet", "BirdieCountNet", "ParCountNet", "BogeyCountNet", "DoubleBogeyCountNet", "TripleBogeyOrWorseCountNet"}
		if err := writer.Write(header); err != nil {
			return err
		}

		// Write data rows
		for handicap := -10; handicap <= 54; handicap++ {
			if hha, exists := ha.HandicapAnalysis[handicap]; exists {
				row := []string{
					fmt.Sprintf("%d", handicap),
					fmt.Sprintf("%d", hha.TotalNetScore),
					fmt.Sprintf("%d", hha.TotalGrossScore),
					fmt.Sprintf("%.2f", float64(hha.TotalGrossScore)/float64(hha.Count)),
					fmt.Sprintf("%.2f", float64(hha.TotalNetScore)/float64(hha.Count)),
					fmt.Sprintf("%d", hha.Count),
					fmt.Sprintf("%d", hha.EagleCountGross),
					fmt.Sprintf("%d", hha.BirdieCountGross),
					fmt.Sprintf("%d", hha.ParCountGross),
					fmt.Sprintf("%d", hha.BogeyCountGross),
					fmt.Sprintf("%d", hha.DoubleBogeyCountGross),
					fmt.Sprintf("%d", hha.TripleBogeyOrWorseCountGross),
					fmt.Sprintf("%d", hha.EagleCountNet),
					fmt.Sprintf("%d", hha.BirdieCountNet),
					fmt.Sprintf("%d", hha.ParCountNet),
					fmt.Sprintf("%d", hha.BogeyCountNet),
					fmt.Sprintf("%d", hha.DoubleBogeyCountNet),
					fmt.Sprintf("%d", hha.TripleBogeyOrWorseCountNet),
				}
				if err := writer.Write(row); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
