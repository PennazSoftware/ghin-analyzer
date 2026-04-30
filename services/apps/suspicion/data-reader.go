package main

import (
	"PennazSoftware/ghin-analyzer/pkg/awsdynamodb"
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/pkg/hcmodel"
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// getGolferScores retrieves a golfer's scores from DynamoDB if they exist, otherwise it fetches them from the GHIN API
func getGolferScores(golferID int, startDate time.Time, endDate time.Time, forceGhinRefresh bool, dbClient awsdynamodb.AwsDynamoDb, ghinClient ghin.API) ([]ghin.Score, hcmodel.Golfer, error) {
	// Get the golfer from the database
	var scores []ghin.Score

	golfer, err := dbClient.GetGolfer(golferID)
	if err == nil {
		if !forceGhinRefresh {
			// Get all the scores we have from the database for the date range requested
			var pageKey map[string]types.AttributeValue = nil
			for {
				tempScores, newPageKey, err := dbClient.GetScoresByGolferPlayedAtRange(golferID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), pageKey)
				if err != nil {
					log.Printf("failed to get scores from database - %+v", err)
					return nil, golfer, err
				}

				scores = append(scores, tempScores...)

				if newPageKey == nil {
					break
				}
				pageKey = newPageKey
			}

			// Get any new scores from GHIN that are not in the database for the date range requested
			if golfer.LastGhinUpdate == endDate.Format("2006-01-02") {
				log.Printf("golfer found in database with score date range that covers requested date range...returning scores from database...")
				return scores, golfer, nil
			}

			// If the golfer's last GHIN update is before the end date of the requested date range, we need to fetch the new scores from GHIN and save them to the database
			log.Printf("golfer found in database with score date range that does not cover requested date range...fetching scores from GHIN...")
			startDate, err = time.Parse("2006-01-02", golfer.LastGhinUpdate)
			if err != nil {
				log.Printf("failed to parse golfer's last GHIN update date - %+v", err)
				return scores, golfer, err
			}
		}

		ghinScores, err := ghinClient.GetScoresByDate(golferID, startDate, endDate)
		if err != nil {
			log.Printf("failed to get scores from ghin - %+v", err)
			return nil, golfer, err
		}

		// Save the scores to the database
		for _, score := range ghinScores {
			_, err = dbClient.CreateScore(score)
			if err != nil {
				log.Printf("failed to create score (id=%d) in database - %+v", score.ID, err)
			}
		}

		scores = append(scores, ghinScores...)

		oldestScoreDate, latestScoreDate := findOldestAndLatestScoreDate(scores)
		golfer.ScoreEndDate = latestScoreDate

		if forceGhinRefresh {
			golfer.ScoreStartDate = oldestScoreDate
		}

		golfer.LastGhinUpdate = time.Now().Format("2006-01-02")
		_, err = dbClient.UpdateGolfer(golfer)
		if err != nil {
			log.Printf("failed to update golfer in database - %+v", err)
			return scores, golfer, err
		}

		return scores, golfer, nil
	}

	log.Printf("no golfer found in database...retrieving golfer from GHIN...")

	golfers, err := ghinClient.GetGolfer(golferID)
	if err != nil {
		log.Printf("failed to get golfer from ghin - %+v", err)
		return nil, hcmodel.Golfer{}, err
	}

	if len(golfers) == 0 {
		log.Printf("no golfer found with id %d", golferID)
		return nil, hcmodel.Golfer{}, nil
	}

	newGolfer := hcmodel.Golfer{
		GolferID:  hcutil.StringToInt(golfers[0].Ghin),
		FirstName: golfers[0].FirstName,
		LastName:  golfers[0].LastName,
	}

	scores, err = ghinClient.GetScoresByDate(golferID, startDate, endDate)
	if err != nil {
		log.Printf("failed to get scores from ghin - %+v", err)
		return scores, newGolfer, err
	}

	// Save the golfer to the database
	oldestScoreDate, latestScoreDate := findOldestAndLatestScoreDate(scores)
	newGolfer.ScoreStartDate = oldestScoreDate
	newGolfer.ScoreEndDate = latestScoreDate
	newGolfer.LastGhinUpdate = time.Now().Format("2006-01-02")

	// Save the scores to the database
	for _, score := range scores {
		_, err = dbClient.CreateScore(score)
		if err != nil {
			log.Printf("failed to create score (id=%d) in database - %+v", score.ID, err)
		}
	}

	_, err = dbClient.CreateGolfer(newGolfer)
	if err != nil {
		log.Printf("failed to create golfer in database - %+v", err)
		return scores, newGolfer, err
	}

	return scores, newGolfer, err
}

func findOldestAndLatestScoreDate(scores []ghin.Score) (string, string) {
	if len(scores) == 0 {
		return "", ""
	}

	oldestDate := scores[0].PlayedAt
	latestDate := scores[0].PlayedAt

	for _, score := range scores {
		if score.PlayedAt < oldestDate {
			oldestDate = score.PlayedAt
		}

		if score.PlayedAt > latestDate {
			latestDate = score.PlayedAt
		}
	}

	return oldestDate, latestDate
}
