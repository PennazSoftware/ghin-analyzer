package ghin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	scoresByDateTemplate string = "https://api2.ghin.com/api/v1/scores.json?golfer_id={{GHIN_NUMBER}}&offset={{OFFSET}}&limit=365&from_date_played={{DATE_BEGIN}}&to_date_played={{DATE_END}}&statuses=Validated"
	childScoresTemplate  string = "https://api2.ghin.com/api/v1/scores/{{PARENT_SCORE_ID}}/get_child_scores.json?golfer_id={{GHIN_NUMBER}}"
)

// ScoreResponse is the response returned from a call to get scores
type ScoreResponse struct {
	Scores []Score `json:"scores"`
	// TotalCount   int     `json:"total_count"`
	// HighestScore int     `json:"highest_score"`
	// LowestScore  int     `json:"lowest_score"`
}

// Score contains all the information regarding a round of golf
type Score struct {
	ID                     int     `json:"id" dynamodbav:"id"`
	OrderNumber            int     `json:"order_number" dynamodbav:"order_number"`
	ScoreDayOrder          int     `json:"score_day_order" dynamodbav:"score_day_order"`
	Gender                 string  `json:"gender" dynamodbav:"gender"`
	Status                 string  `json:"status" dynamodbav:"status"`
	IsManual               bool    `json:"is_manual" dynamodbav:"is_manual"`
	NumberOfHoles          int     `json:"number_of_holes" dynamodbav:"number_of_holes"`
	NumberOfPlayedHoles    int     `json:"number_of_played_holes" dynamodbav:"number_of_played_holes"`
	GolferID               int     `json:"golfer_id" dynamodbav:"golfer_id"`
	FacilityName           string  `json:"facility_name" dynamodbav:"facility_name"`
	AdjustedGrossScore     int     `json:"adjusted_gross_score" dynamodbav:"adjusted_gross_score"`
	Front9Adjusted         int     `json:"front9_adjusted" dynamodbav:"front9_adjusted"`
	Back9Adjusted          int     `json:"back9_adjusted" dynamodbav:"back9_adjusted"`
	Front9SlopeRating      int     `json:"front9_slope_rating" dynamodbav:"front9_slope_rating"`
	Back9SlopeRating       int     `json:"back9_slope_rating" dynamodbav:"back9_slope_rating"`
	PlayedAt               string  `json:"played_at" dynamodbav:"played_at"`
	CourseID               string  `json:"course_id" dynamodbav:"course_id"`
	CourseName             string  `json:"course_name" dynamodbav:"course_name"`
	Front9CourseName       string  `json:"front9_course_name" dynamodbav:"front9_course_name"`
	Back9CourseName        string  `json:"back9_course_name" dynamodbav:"back9_course_name"`
	Front9CourseRating     float64 `json:"front9_course_rating" dynamodbav:"front9_course_rating"`
	Back9CourseRating      float64 `json:"back9_course_rating" dynamodbav:"back9_course_rating"`
	TeeName                string  `json:"tee_name" dynamodbav:"tee_name"`
	TeeSetID               string  `json:"tee_set_id" dynamodbav:"tee_set_id"`
	TeeSetSide             string  `json:"tee_set_side" dynamodbav:"tee_set_side"`
	Front9TeeName          string  `json:"front9_tee_name" dynamodbav:"front9_tee_name"`
	Back9TeeName           string  `json:"back9_tee_name" dynamodbav:"back9_tee_name"`
	Differential           float64 `json:"differential" dynamodbav:"differential"`
	UnadjustedDifferential float64 `json:"unadjusted_differential" dynamodbav:"unadjusted_differential"`
	ScoreType              string  `json:"score_type" dynamodbav:"score_type"`
	CourseRating           float64 `json:"course_rating" dynamodbav:"course_rating"`
	SlopeRating            int     `json:"slope_rating" dynamodbav:"slope_rating"`
	//Penalty                bool    `json:"penalty"`
	//PenaltyType            `json:"penalty_type"`
	//PenaltyMethod          `json:"penalty_method"`
	//ParentID               `json:"parent_id"`
	ScoreTypeDisplayFull         string       `json:"score_type_display_full" dynamodbav:"score_type_display_full"`
	ScoreTypeDisplayShort        string       `json:"score_type_display_short" dynamodbav:"score_type_display_short"`
	Edited                       bool         `json:"edited" dynamodbav:"edited"`
	PostedAt                     string       `json:"posted_at" dynamodbav:"posted_at"`
	SeasonStartDateAt            string       `json:"season_start_date_at" dynamodbav:"season_start_date_at"`
	SeasonEndDateAt              string       `json:"season_end_date_at" dynamodbav:"season_end_date_at"`
	CourseDisplayValue           string       `json:"course_display_value" dynamodbav:"course_display_value"`
	Used                         bool         `json:"used" dynamodbav:"used"`
	Revision                     bool         `json:"revision" dynamodbav:"revision"`
	Pcc                          float64      `json:"pcc" dynamodbav:"pcc"`
	Adjustments                  []Adjustment `json:"adjustments" dynamodbav:"adjustments"`
	HoleDetails                  []HoleDetail `json:"hole_details" dynamodbav:"hole_details"`
	Statistics                   Statistics   `json:"statistics" dynamodbav:"statistics"`
	Exceptional                  bool         `json:"exceptional" dynamodbav:"exceptional"`
	ScaledUpDifferential         float64      `json:"scaled_up_differential" dynamodbav:"scaled_up_differential"`
	AdjustedScaledUpDifferential float64      `json:"adjusted_scaled_up_differential" dynamodbav:"adjusted_scaled_up_differential"`
	PostedOnHomeCourse           bool         `json:"posted_on_home_course" dynamodbav:"posted_on_home_course"`
	NetScore                     int          `json:"net_score" dynamodbav:"net_score"`
	CourseHandicap               string       `json:"course_handicap" dynamodbav:"course_handicap"`
	GhinCourseDisplayName        string       `json:"ghin_course_display_name" dynamodbav:"ghin_course_display_name"`
}

// Statistics contains overall information for a round of golf
type Statistics struct {
	GirPercent                 string `json:"gir_percent" dynamodbav:"gir_percent"`
	PuttsTotal                 string `json:"putts_total" dynamodbav:"putts_total"`
	ParsPercent                string `json:"pars_percent" dynamodbav:"pars_percent"`
	Par3SAverage               string `json:"par3s_average" dynamodbav:"par3s_average"`
	Par4SAverage               string `json:"par4s_average" dynamodbav:"par4s_average"`
	Par5SAverage               string `json:"par5s_average" dynamodbav:"par5s_average"`
	BogeysPercent              string `json:"bogeys_percent" dynamodbav:"bogeys_percent"`
	MissedLeftPercent          string `json:"missed_left_percent" dynamodbav:"missed_left_percent"`
	MissedLongPercent          string `json:"missed_long_percent" dynamodbav:"missed_long_percent"`
	FairwayHitsPercent         string `json:"fairway_hits_percent" dynamodbav:"fairway_hits_percent"`
	MissedRightPercent         string `json:"missed_right_percent" dynamodbav:"missed_right_percent"`
	MissedShortPercent         string `json:"missed_short_percent" dynamodbav:"missed_short_percent"`
	DoubleBogeysPercent        string `json:"double_bogeys_percent" dynamodbav:"double_bogeys_percent"`
	LastStatsUpdateDate        string `json:"last_stats_update_date" dynamodbav:"last_stats_update_date"`
	LastStatsUpdateType        string `json:"last_stats_update_type" dynamodbav:"last_stats_update_type"`
	BirdiesOrBetterPercent     string `json:"birdies_or_better_percent" dynamodbav:"birdies_or_better_percent"`
	TripleBogeysOrWorsePercent string `json:"triple_bogeys_or_worse_percent" dynamodbav:"triple_bogeys_or_worse_percent"`
}

// HoleDetail contains details about how an individual hole of golf was played
type HoleDetail struct {
	ID                 int  `json:"id" dynamodbav:"id"`
	AdjustedGrossScore int  `json:"adjusted_gross_score" dynamodbav:"adjusted_gross_score"`
	RawScore           int  `json:"raw_score" dynamodbav:"raw_score"`
	HoleNumber         int  `json:"hole_number" dynamodbav:"hole_number"`
	Par                int  `json:"par" dynamodbav:"par"`
	Putts              int  `json:"putts" dynamodbav:"putts"`
	FairwayHit         bool `json:"fairway_hit" dynamodbav:"fairway_hit"`
	GirFlag            bool `json:"gir_flag" dynamodbav:"gir_flag"`
	DriveAccuracy      int  `json:"drive_accuracy" dynamodbav:"drive_accuracy"`
	StrokeAllocation   int  `json:"stroke_allocation" dynamodbav:"stroke_allocation"`
}

// Adjustment relates to score adjustment for a course
type Adjustment struct {
	Type    string  `json:"type" dynamodbav:"type"`
	Value   float64 `json:"value" dynamodbav:"value"`
	Display string  `json:"display" dynamodbav:"display"`
}

// GetScoresByDate retrieves all rounds played by a user within the specified date range
func (c *Client) GetScoresByDate(ghinNumber int, beginDate time.Time, endDate time.Time) ([]Score, error) {
	var scoreResponse ScoreResponse

	url := strings.ReplaceAll(scoresByDateTemplate, "{{GHIN_NUMBER}}", strconv.Itoa(ghinNumber))
	url = strings.ReplaceAll(url, "{{DATE_BEGIN}}", beginDate.Format("2006-01-02"))
	url = strings.ReplaceAll(url, "{{DATE_END}}", endDate.Format("2006-01-02"))

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetScoresByDate()")
		return []Score{}, err
	}

	res, err := c.do(req, &scoreResponse)
	if err != nil {
		// If there are no scores returned then GHIN does this silly thing where instead of returning INTs for
		// fields (like highest_score) they return "-" characters. This causes the marshalling to fail. To work
		// around this we will catch the error here, inspect the message for this special scenario and then return
		// an empty object.  Otherwise we'll return the error.
		if strings.Contains(err.Error(), "cannot unmarshal string into Go struct field ScoreResponse") {
			res.StatusCode = http.StatusOK
			err = nil
		} else {
			log.Printf("ghin client - failed to perform http request GetScoresByDate() -- %+v", err)
			return []Score{}, err
		}
	}

	if res.StatusCode != http.StatusOK {
		return []Score{}, fmt.Errorf("Error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return scoreResponse.Scores, nil
}

// GetScoresByDatePaged retrieves all golf rounds played by a user within the specified date range
func (c *Client) GetScoresByDatePaged(ghinNumber string, beginDate time.Time, endDate time.Time, offset int) (ScoreResponse, error) {
	var scoreResponse ScoreResponse

	url := strings.Replace(scoresByDateTemplate, "{{GHIN_NUMBER}}", ghinNumber, -1)
	url = strings.Replace(url, "{{DATE_BEGIN}}", beginDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{DATE_END}}", endDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{OFFSET}}", strconv.Itoa(offset), -1)

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetScoresByDate()")
		return ScoreResponse{}, err
	}

	res, err := c.do(req, &scoreResponse)
	if err != nil {
		// If there are no scores returned then GHIN does this silly thing where instead of returning INTs for
		// fields (like highest_score) they return "-" characters. This causes the marshalling to fail. To work
		// around this we will catch the error here, inspect the message for this special scenario and then return
		// an empty object.  Otherwise we'll return the error.
		if strings.Contains(err.Error(), "cannot unmarshal string into Go struct field ScoreResponse") {
			res.StatusCode = http.StatusOK
			err = nil
		} else {
			log.Println("ghin client - failed to perform http request GetScoresByDate()")
			return ScoreResponse{}, err
		}
	}

	if res.StatusCode != http.StatusOK {
		return ScoreResponse{}, fmt.Errorf("Error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return scoreResponse, nil
}

// GetChildScores retrieves child scores that make up an 18 hole round
func (c *Client) GetChildScores(ghinNumber string, parentScoreID int) ([]Score, error) {
	var scoreResponse ScoreResponse

	url := strings.Replace(childScoresTemplate, "{{GHIN_NUMBER}}", ghinNumber, -1)
	url = strings.Replace(url, "{{PARENT_SCORE_ID}}", strconv.Itoa(parentScoreID), -1)

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetChildScores()")
		return []Score{}, err
	}

	res, err := c.do(req, &scoreResponse)
	if err != nil {
		log.Println("ghin client - failed to perform http request GetChildScores()")
		return []Score{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []Score{}, fmt.Errorf("Error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return scoreResponse.Scores, nil
}
