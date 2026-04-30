package main

type HoleAnalysis struct {
	HoleNumber       int                          `json:"hole_number"`
	Par              int                          `json:"par"`
	HandicapAnalysis map[int]HoleHandicapAnalysis `json:"handicap_analysis,omitempty"`
}

type HoleHandicapAnalysis struct {
	Handicap                     int `json:"handicap"`
	TotalNetScore                int `json:"total_net_score"`
	TotalGrossScore              int `json:"total_gross_score"`
	Count                        int `json:"count"`
	EagleCountGross              int `json:"eagle_count_gross"`
	BirdieCountGross             int `json:"birdie_count_gross"`
	ParCountGross                int `json:"par_count_gross"`
	BogeyCountGross              int `json:"bogey_count_gross"`
	DoubleBogeyCountGross        int `json:"double_bogey_count_gross"`
	TripleBogeyOrWorseCountGross int `json:"triple_bogey_or_worse_count_gross"`
	EagleCountNet                int `json:"eagle_count_net"`
	BirdieCountNet               int `json:"birdie_count_net"`
	ParCountNet                  int `json:"par_count_net"`
	BogeyCountNet                int `json:"bogey_count_net"`
	DoubleBogeyCountNet          int `json:"double_bogey_count_net"`
	TripleBogeyOrWorseCountNet   int `json:"triple_bogey_or_worse_count_net"`
}

type ScoreFile struct {
	Scores       []Scores `json:"scores,omitempty"`
	TotalCount   int      `json:"total_count,omitempty"`
	HighestScore int      `json:"highest_score,omitempty"`
	LowestScore  int      `json:"lowest_score,omitempty"`
	Average      float64  `json:"average,omitempty"`
}
type HoleDetails struct {
	ID                   int64 `json:"id,omitempty"`
	AdjustedGrossScore   int   `json:"adjusted_gross_score,omitempty"`
	RawScore             int   `json:"raw_score,omitempty"`
	HoleNumber           int   `json:"hole_number,omitempty"`
	Par                  int   `json:"par,omitempty"`
	Putts                any   `json:"putts,omitempty"`
	FairwayHit           any   `json:"fairway_hit,omitempty"`
	GirFlag              any   `json:"gir_flag,omitempty"`
	DriveAccuracy        any   `json:"drive_accuracy,omitempty"`
	StrokeAllocation     int   `json:"stroke_allocation,omitempty"`
	ApproachShotAccuracy any   `json:"approach_shot_accuracy,omitempty"`
	XHole                bool  `json:"x_hole,omitempty"`
	MostLikelyScore      any   `json:"most_likely_score,omitempty"`
}
type Statistics struct {
	PuttsTotal                               string `json:"putts_total,omitempty"`
	OnePuttOrBetterPercent                   string `json:"one_putt_or_better_percent,omitempty"`
	TwoPuttPercent                           string `json:"two_putt_percent,omitempty"`
	ThreePuttOrWorsePercent                  string `json:"three_putt_or_worse_percent,omitempty"`
	TwoPuttOrBetterPercent                   string `json:"two_putt_or_better_percent,omitempty"`
	UpAndDownsTotal                          string `json:"up_and_downs_total,omitempty"`
	Par3SAverage                             string `json:"par3s_average,omitempty"`
	Par4SAverage                             string `json:"par4s_average,omitempty"`
	Par5SAverage                             string `json:"par5s_average,omitempty"`
	ParsPercent                              string `json:"pars_percent,omitempty"`
	BirdiesOrBetterPercent                   string `json:"birdies_or_better_percent,omitempty"`
	BogeysPercent                            string `json:"bogeys_percent,omitempty"`
	DoubleBogeysPercent                      string `json:"double_bogeys_percent,omitempty"`
	TripleBogeysOrWorsePercent               string `json:"triple_bogeys_or_worse_percent,omitempty"`
	FairwayHitsPercent                       string `json:"fairway_hits_percent,omitempty"`
	MissedLeftPercent                        string `json:"missed_left_percent,omitempty"`
	MissedRightPercent                       string `json:"missed_right_percent,omitempty"`
	MissedLongPercent                        string `json:"missed_long_percent,omitempty"`
	MissedShortPercent                       string `json:"missed_short_percent,omitempty"`
	GirPercent                               string `json:"gir_percent,omitempty"`
	MissedLeftApproachShotAccuracyPercent    string `json:"missed_left_approach_shot_accuracy_percent,omitempty"`
	MissedRightApproachShotAccuracyPercent   string `json:"missed_right_approach_shot_accuracy_percent,omitempty"`
	MissedLongApproachShotAccuracyPercent    string `json:"missed_long_approach_shot_accuracy_percent,omitempty"`
	MissedShortApproachShotAccuracyPercent   string `json:"missed_short_approach_shot_accuracy_percent,omitempty"`
	MissedGeneralApproachShotAccuracyPercent string `json:"missed_general_approach_shot_accuracy_percent,omitempty"`
	LastStatsUpdateDate                      string `json:"last_stats_update_date,omitempty"`
	LastStatsUpdateType                      string `json:"last_stats_update_type,omitempty"`
}
type Scores struct {
	ID                           int           `json:"id,omitempty"`
	OrderNumber                  int           `json:"order_number,omitempty"`
	ScoreDayOrder                int           `json:"score_day_order,omitempty"`
	Gender                       string        `json:"gender,omitempty"`
	Status                       string        `json:"status,omitempty"`
	IsManual                     bool          `json:"is_manual,omitempty"`
	NumberOfHoles                int           `json:"number_of_holes,omitempty"`
	NumberOfPlayedHoles          int           `json:"number_of_played_holes,omitempty"`
	GolferID                     int           `json:"golfer_id,omitempty"`
	FacilityName                 string        `json:"facility_name,omitempty"`
	AdjustedGrossScore           int           `json:"adjusted_gross_score,omitempty"`
	Front9Adjusted               any           `json:"front9_adjusted,omitempty"`
	Back9Adjusted                any           `json:"back9_adjusted,omitempty"`
	PostedOnHomeCourse           bool          `json:"posted_on_home_course,omitempty"`
	PlayedAt                     string        `json:"played_at,omitempty"`
	Front9SlopeRating            any           `json:"front9_slope_rating,omitempty"`
	Back9SlopeRating             any           `json:"back9_slope_rating,omitempty"`
	CourseID                     string        `json:"course_id,omitempty"`
	CourseName                   string        `json:"course_name,omitempty"`
	Front9CourseName             any           `json:"front9_course_name,omitempty"`
	Back9CourseName              any           `json:"back9_course_name,omitempty"`
	Front9CourseRating           any           `json:"front9_course_rating,omitempty"`
	Back9CourseRating            any           `json:"back9_course_rating,omitempty"`
	TeeName                      string        `json:"tee_name,omitempty"`
	TeeSetID                     string        `json:"tee_set_id,omitempty"`
	TeeSetSide                   string        `json:"tee_set_side,omitempty"`
	Front9TeeName                any           `json:"front9_tee_name,omitempty"`
	Back9TeeName                 any           `json:"back9_tee_name,omitempty"`
	Differential                 float64       `json:"differential,omitempty"`
	UnadjustedDifferential       float64       `json:"unadjusted_differential,omitempty"`
	ScaledUpDifferential         any           `json:"scaled_up_differential,omitempty"`
	AdjustedScaledUpDifferential any           `json:"adjusted_scaled_up_differential,omitempty"`
	ScoreType                    string        `json:"score_type,omitempty"`
	Penalty                      any           `json:"penalty,omitempty"`
	PenaltyType                  any           `json:"penalty_type,omitempty"`
	PenaltyMethod                any           `json:"penalty_method,omitempty"`
	ParentID                     any           `json:"parent_id,omitempty"`
	CourseRating                 float64       `json:"course_rating,omitempty"`
	SlopeRating                  int           `json:"slope_rating,omitempty"`
	ScoreTypeDisplayFull         string        `json:"score_type_display_full,omitempty"`
	ScoreTypeDisplayShort        string        `json:"score_type_display_short,omitempty"`
	Edited                       bool          `json:"edited,omitempty"`
	PostedAt                     string        `json:"posted_at,omitempty"`
	SeasonStartDateAt            string        `json:"season_start_date_at,omitempty"`
	SeasonEndDateAt              string        `json:"season_end_date_at,omitempty"`
	ChallengeAvailable           any           `json:"challenge_available,omitempty"`
	NetScore                     int           `json:"net_score,omitempty"`
	CourseHandicap               string        `json:"course_handicap,omitempty"`
	CourseDisplayValue           string        `json:"course_display_value,omitempty"`
	GhinCourseNameDisplay        string        `json:"ghin_course_name_display,omitempty"`
	Used                         bool          `json:"used,omitempty"`
	Revision                     bool          `json:"revision,omitempty"`
	Pcc                          float64       `json:"pcc,omitempty"`
	Adjustments                  []any         `json:"adjustments,omitempty"`
	HoleDetails                  []HoleDetails `json:"hole_details,omitempty"`
	Statistics                   Statistics    `json:"statistics,omitempty"`
	Exceptional                  bool          `json:"exceptional,omitempty"`
	MessageClubAuthorized        any           `json:"message_club_authorized,omitempty"`
	IsRecent                     bool          `json:"is_recent,omitempty"`
	NetScoreDifferential         float64       `json:"net_score_differential,omitempty"`
	ShortCourse                  any           `json:"short_course,omitempty"`
}
