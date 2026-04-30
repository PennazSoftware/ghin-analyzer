package main

import (
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"time"
)

type PlayerScores struct {
	Scores []GHINScore `json:"scores"`
}

type GHINScore struct {
	ID                           int64          `json:"id"`
	OrderNumber                  int            `json:"order_number"`
	ScoreDayOrder                int            `json:"score_day_order"`
	Gender                       string         `json:"gender"`
	Status                       string         `json:"status"`
	IsManual                     bool           `json:"is_manual"`
	NumberOfHoles                int            `json:"number_of_holes"`
	NumberOfPlayedHoles          int            `json:"number_of_played_holes"`
	GolferID                     int64          `json:"golfer_id"`
	FacilityName                 string         `json:"facility_name"`
	AdjustedGrossScore           int            `json:"adjusted_gross_score"`
	Front9Adjusted               *int           `json:"front9_adjusted"`
	Back9Adjusted                *int           `json:"back9_adjusted"`
	PostedOnHomeCourse           bool           `json:"posted_on_home_course"`
	PlayedAt                     string         `json:"played_at"`
	Front9SlopeRating            *float64       `json:"front9_slope_rating"`
	Back9SlopeRating             *float64       `json:"back9_slope_rating"`
	CourseID                     string         `json:"course_id"`
	CourseName                   string         `json:"course_name"`
	Front9CourseName             *string        `json:"front9_course_name"`
	Back9CourseName              *string        `json:"back9_course_name"`
	Front9CourseRating           *float64       `json:"front9_course_rating"`
	Back9CourseRating            *float64       `json:"back9_course_rating"`
	TeeName                      string         `json:"tee_name"`
	TeeSetID                     string         `json:"tee_set_id"`
	TeeSetSide                   string         `json:"tee_set_side"`
	Front9TeeName                *string        `json:"front9_tee_name"`
	Back9TeeName                 *string        `json:"back9_tee_name"`
	Differential                 float64        `json:"differential"`
	UnadjustedDifferential       float64        `json:"unadjusted_differential"`
	ScaledUpDifferential         *float64       `json:"scaled_up_differential"`
	AdjustedScaledUpDifferential *float64       `json:"adjusted_scaled_up_differential"`
	ScoreType                    string         `json:"score_type"`
	Penalty                      *float64       `json:"penalty"`
	PenaltyType                  *string        `json:"penalty_type"`
	PenaltyMethod                *string        `json:"penalty_method"`
	ParentID                     *int64         `json:"parent_id"`
	CourseRating                 float64        `json:"course_rating"`
	SlopeRating                  int            `json:"slope_rating"`
	ScoreTypeDisplayFull         string         `json:"score_type_display_full"`
	ScoreTypeDisplayShort        string         `json:"score_type_display_short"`
	Edited                       bool           `json:"edited"`
	PostedAt                     string         `json:"posted_at"`
	SeasonStartDateAt            string         `json:"season_start_date_at"`
	SeasonEndDateAt              string         `json:"season_end_date_at"`
	ChallengeAvailable           any            `json:"challenge_available"`
	NetScore                     int            `json:"net_score"`
	CourseHandicap               string         `json:"course_handicap"`
	CourseDisplayValue           string         `json:"course_display_value"`
	GHINCourseNameDisplay        string         `json:"ghin_course_name_display"`
	Used                         bool           `json:"used"`
	Revision                     bool           `json:"revision"`
	PCC                          float64        `json:"pcc"`
	Adjustments                  []any          `json:"adjustments"`
	HoleDetails                  []GHINHole     `json:"hole_details"`
	Statistics                   GHINStatistics `json:"statistics"`
	Exceptional                  bool           `json:"exceptional"`
	MessageClubAuthorized        any            `json:"message_club_authorized"`
	IsRecent                     bool           `json:"is_recent"`
	NetScoreDifferential         float64        `json:"net_score_differential"`
	ShortCourse                  any            `json:"short_course"`
}

type GHINHole struct {
	ID                   int64   `json:"id"`
	AdjustedGrossScore   int     `json:"adjusted_gross_score"`
	RawScore             int     `json:"raw_score"`
	HoleNumber           int     `json:"hole_number"`
	Par                  int     `json:"par"`
	Putts                *int    `json:"putts"`
	FairwayHit           *bool   `json:"fairway_hit"`
	GIRFlag              *bool   `json:"gir_flag"`
	DriveAccuracy        *string `json:"drive_accuracy"`
	StrokeAllocation     int     `json:"stroke_allocation"`
	ApproachShotAccuracy *string `json:"approach_shot_accuracy"`
	XHole                bool    `json:"x_hole"`
	MostLikelyScore      *int    `json:"most_likely_score"`
}

type GHINStatistics struct {
	PuttsTotal                               string `json:"putts_total"`
	OnePuttOrBetterPercent                   string `json:"one_putt_or_better_percent"`
	TwoPuttPercent                           string `json:"two_putt_percent"`
	ThreePuttOrWorsePercent                  string `json:"three_putt_or_worse_percent"`
	TwoPuttOrBetterPercent                   string `json:"two_putt_or_better_percent"`
	UpAndDownsTotal                          string `json:"up_and_downs_total"`
	Par3sAverage                             string `json:"par3s_average"`
	Par4sAverage                             string `json:"par4s_average"`
	Par5sAverage                             string `json:"par5s_average"`
	ParsPercent                              string `json:"pars_percent"`
	BirdiesOrBetterPercent                   string `json:"birdies_or_better_percent"`
	BogeysPercent                            string `json:"bogeys_percent"`
	DoubleBogeysPercent                      string `json:"double_bogeys_percent"`
	TripleBogeysOrWorsePercent               string `json:"triple_bogeys_or_worse_percent"`
	FairwayHitsPercent                       string `json:"fairway_hits_percent"`
	MissedLeftPercent                        string `json:"missed_left_percent"`
	MissedRightPercent                       string `json:"missed_right_percent"`
	MissedLongPercent                        string `json:"missed_long_percent"`
	MissedShortPercent                       string `json:"missed_short_percent"`
	GIRPercent                               string `json:"gir_percent"`
	MissedLeftApproachShotAccuracyPercent    string `json:"missed_left_approach_shot_accuracy_percent"`
	MissedRightApproachShotAccuracyPercent   string `json:"missed_right_approach_shot_accuracy_percent"`
	MissedLongApproachShotAccuracyPercent    string `json:"missed_long_approach_shot_accuracy_percent"`
	MissedShortApproachShotAccuracyPercent   string `json:"missed_short_approach_shot_accuracy_percent"`
	MissedGeneralApproachShotAccuracyPercent string `json:"missed_general_approach_shot_accuracy_percent"`
	LastStatsUpdateDate                      string `json:"last_stats_update_date"`
	LastStatsUpdateType                      string `json:"last_stats_update_type"`
}

type GolferAnalysis struct {
	GolferID  int
	FirstName string
	LastName  string

	HandicapIndex float64

	TotalRounds               int
	TournamentRounds          int
	NonTournamentRounds       int
	RoundsStartDate           string
	RoundsEndDate             string
	PlayedAtDistributionGraph string

	AvgDifferential           float64
	AvgTournamentDifferential float64
	AvgCasualDifferential     float64
	StdDevDifferential        float64
	SkewnessDifferential      float64
	SkewnessTournament        float64
	SkewnessCasual            float64
	SkewnessGap               float64

	Best8Avg      float64
	RemainingAvg  float64
	BestVsRestGap float64

	TournamentBeatsHandicapRate float64
	OverallBeatsHandicapRate    float64
	CasualBeatsHandicapRate     float64

	TournamentVsCasualGap   float64
	TournamentVsHandicapGap float64

	TrajectorySamples      int
	TrajectoryPatternCount int
	TrajectoryPatternRate  float64
	AvgPreTournamentRise   float64
	AvgPostTournamentDrop  float64

	AvgBlowupHolesPerRound float64
	AvgBlowupHolesTourney  float64
	AvgBlowupHolesCasual   float64
	BlowupGap              float64
	BlowUpPatternRounds    int
	BlowUpPatternRate      float64
	BlowUpPatternTourney   int
	BlowUpPatternRateT     float64
	BlowUpPatternCasual    int
	BlowUpPatternRateC     float64
	BlowUpPatternRateGap   float64

	AvgPostingDelayDays        float64
	AvgPostingDelayDaysTourney float64
	AvgPostingDelayDaysCasual  float64
	LatePostCount              int
	LatePostRate               float64
	LatePostCountTourney       int
	LatePostRateTourney        float64
	LatePostCountCasual        int
	LatePostRateCasual         float64
	LatePostDelayGap           float64
	LatePostRateGap            float64

	PostingHoleByHoleCount        int
	PostingFrontBackTotalCount    int
	PostingTotalOnlyCount         int
	PostingHoleByHoleRate         float64
	PostingFrontBackTotalRate     float64
	PostingTotalOnlyRate          float64
	PostingHoleByHoleCountTourney int
	PostingFrontBackCountTourney  int
	PostingTotalOnlyCountTourney  int
	PostingHoleByHoleRateTourney  float64
	PostingFrontBackRateTourney   float64
	PostingTotalOnlyRateTourney   float64
	PostingHoleByHoleCountCasual  int
	PostingFrontBackCountCasual   int
	PostingTotalOnlyCountCasual   int
	PostingHoleByHoleRateCasual   float64
	PostingFrontBackRateCasual    float64
	PostingTotalOnlyRateCasual    float64
	PostingTotalOnlyRateGap       float64

	SuspicionScore float64
	Flags          []string
}

type scoredRound struct {
	Score         ghin.Score
	PlayedDate    time.Time
	PostedDate    time.Time
	Diff          float64
	IsTournament  bool
	BlowupHoles   int
	IsBlowUpRound bool
	PostingDelay  int
	IsLatePosted  bool
	PostingStyle  string
}

type handicapTrajectoryMetrics struct {
	Samples               int
	Patterns              int
	PatternRate           float64
	AvgPreTournamentRise  float64
	AvgPostTournamentDrop float64
}
