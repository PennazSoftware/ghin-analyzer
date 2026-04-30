package main

import (
	"PennazSoftware/ghin-analyzer/pkg/awsdynamodb"
	"PennazSoftware/ghin-analyzer/pkg/awssecretsmanager"
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/pkg/hcmodel"
	"PennazSoftware/ghin-analyzer/pkg/pennazemail"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	ghinSecretName     = "Ghin_50155"
	forceGhinRefresh   = true
	defaultSenderEmail = "Dan Penn <headpro@golffuze.com>"

	postingStyleHoleByHole     = "hole_by_hole"
	postingStyleFrontBackTotal = "front_back_total"
	postingStyleTotalOnly      = "total_only"
)

func main() {
	environment := "dev"
	os.Setenv("AWS_PROFILE", "pennaz")
	dbClient := awsdynamodb.New(environment)
	secretsClient := awssecretsmanager.New()

	// Get the GHIN credentials from AWS Secrets Manager
	var ghinSecret awssecretsmanager.GHINSecret
	err := secretsClient.GetSecret(ghinSecretName, &ghinSecret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ghinClient := ghin.NewWithLogin(ghinSecret.Username, ghinSecret.Password)
	golfersMap := make(map[int]hcmodel.Golfer)

	// 2048850, 50155, 8667665
	scores, golfer, err := getGolferScores(2048850, time.Now().AddDate(-1, 0, 0), time.Now(), forceGhinRefresh, dbClient, ghinClient)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	golfersMap[golfer.GolferID] = golfer

	log.Printf("Total Scores found: %d", len(scores))

	results := AnalyzeGHINScores(scores, golfersMap)

	for i := range results {
		if strings.TrimSpace(results[i].FirstName) == "" {
			results[i].FirstName = "Unknown"
		}
		if strings.TrimSpace(results[i].LastName) == "" {
			results[i].LastName = "Unknown"
		}
	}

	for _, r := range results {
		fmt.Printf("Golfer: %s %s (%d)\n", r.FirstName, r.LastName, r.GolferID)
		fmt.Printf("  Handicap Index: %.2f\n", r.HandicapIndex)
		fmt.Printf("  Avg Diff: %.2f\n", r.AvgDifferential)
		fmt.Printf("  Avg Casual Diff: %.2f\n", r.AvgCasualDifferential)
		fmt.Printf("  Avg Tournament Diff: %.2f\n", r.AvgTournamentDifferential)
		fmt.Printf("  Distribution Skew (All/T/C): %.2f / %.2f / %.2f\n", r.SkewnessDifferential, r.SkewnessTournament, r.SkewnessCasual)
		fmt.Printf("  Distribution Skew Gap (Casual - Tournament): %.2f\n", r.SkewnessGap)
		fmt.Printf("  Tournament vs Casual Gap: %.2f\n", r.TournamentVsCasualGap)
		fmt.Printf("  Tournament vs Handicap Gap: %.2f\n", r.TournamentVsHandicapGap)
		fmt.Printf("  Handicap Trajectory Samples/Patterns: %d / %d\n", r.TrajectorySamples, r.TrajectoryPatternCount)
		fmt.Printf("  Handicap Trajectory Rise/Drop Avg: %.2f / %.2f (Pattern Rate: %.1f%%)\n",
			r.AvgPreTournamentRise, r.AvgPostTournamentDrop, r.TrajectoryPatternRate*100)
		fmt.Printf("  Best8 vs Remaining Diff Gap: %.2f\n", r.BestVsRestGap)
		fmt.Printf("  Low Round Frequency (beat HI) Overall/Casual/Tournament: %.1f%% / %.1f%% / %.1f%%\n",
			r.OverallBeatsHandicapRate*100, r.CasualBeatsHandicapRate*100, r.TournamentBeatsHandicapRate*100)
		fmt.Printf("  Blow-Up Round Pattern Rate Overall/Casual/Tournament: %.1f%% / %.1f%% / %.1f%%\n",
			r.BlowUpPatternRate*100, r.BlowUpPatternRateC*100, r.BlowUpPatternRateT*100)
		fmt.Printf("  Avg Posting Delay Days: %.2f\n", r.AvgPostingDelayDays)
		fmt.Printf("  Late Post Count (7+ days): %d\n", r.LatePostCount)
		fmt.Printf("  Late Post Rate: %.1f%%\n", r.LatePostRate*100)
		fmt.Printf("  Late Post Rate Casual: %.1f%%\n", r.LatePostRateCasual*100)
		fmt.Printf("  Late Post Rate Tournament: %.1f%%\n", r.LatePostRateTourney*100)
		fmt.Printf("  Late Post Rate Gap (Casual - Tournament): %.2f\n", r.LatePostRateGap)
		fmt.Printf("  Posting Detail Rate HxH/F&B/Tot(Overall): %.1f%% / %.1f%% / %.1f%%\n",
			r.PostingHoleByHoleRate*100, r.PostingFrontBackTotalRate*100, r.PostingTotalOnlyRate*100)
		fmt.Printf("  Posting Detail Rate HxH/F&B/Tot (Casual): %.1f%% / %.1f%% / %.1f%%\n",
			r.PostingHoleByHoleRateCasual*100, r.PostingFrontBackRateCasual*100, r.PostingTotalOnlyRateCasual*100)
		fmt.Printf("  Posting Detail Rate HxH/F&B/Tot (Tournament): %.1f%% / %.1f%% / %.1f%%\n",
			r.PostingHoleByHoleRateTourney*100, r.PostingFrontBackRateTourney*100, r.PostingTotalOnlyRateTourney*100)
		fmt.Printf("  Total-Only Posting Gap (Casual - Tournament): %.2f\n", r.PostingTotalOnlyRateGap)
		fmt.Printf("  PlayedAt Day-of-Month Distribution (1-31):\n%s\n", r.PlayedAtDistributionGraph)
		fmt.Printf("  Suspicion Score: %.1f\n", r.SuspicionScore)
		fmt.Printf("  Flags: %v\n", r.Flags)
		fmt.Println()
	}

	html, err := writeResultsEmailHTML(results)
	if err != nil {
		fmt.Printf("failed to write results html: %v\n", err)
		return
	}

	// Send email
	emailClient := pennazemail.NewAwsSes()
	err = emailClient.SendEmail([]string{"danpenn@msn.com"},
		fmt.Sprintf("GHIN Suspicion Analysis Results for %s %s (%d)", golfer.FirstName, golfer.LastName, golfer.GolferID),
		html,
		defaultSenderEmail,
	)
	if err != nil {
		fmt.Printf("failed to send email: %v\n", err)
		return
	}

	fmt.Printf("Generated email HTML summary (%d bytes)\n", len(html))
}

func AnalyzeGHINScores(scores []ghin.Score, golfers map[int]hcmodel.Golfer) []GolferAnalysis {
	byGolfer := map[int][]ghin.Score{}
	for _, s := range scores {
		if !isUsableScore(s) {
			continue
		}
		byGolfer[s.GolferID] = append(byGolfer[s.GolferID], s)
	}

	results := make([]GolferAnalysis, 0, len(byGolfer))
	for golferID, golferScores := range byGolfer {
		analysis := analyzeSingleGolfer(golferID, golferScores)
		if golfer, ok := golfers[golferID]; ok {
			analysis.FirstName = golfer.FirstName
			analysis.LastName = golfer.LastName
		}
		results = append(results, analysis)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].SuspicionScore > results[j].SuspicionScore
	})

	return results
}

func analyzeSingleGolfer(golferID int, scores []ghin.Score) GolferAnalysis {
	rounds := make([]scoredRound, 0, len(scores))

	for _, s := range scores {
		playedDate, err := time.Parse("2006-01-02", s.PlayedAt)
		if err != nil {
			// Some scores just have the year and month, so we can try parsing with a fallback format
			playedDate, err = time.Parse("2006-01", s.PlayedAt)
			if err != nil {
				continue
			}
		}

		postedDate, err := parsePostedAtDate(s.PostedAt)
		if err != nil {
			// Some dates just have year and month, so we can try parsing with a fallback format
			postedDate, err = time.Parse("2006-01", s.PostedAt)
			if err != nil {
				continue
			}
		}

		delay := postingDelayDays(playedDate, postedDate)

		rounds = append(rounds, scoredRound{
			Score:         s,
			PlayedDate:    playedDate,
			PostedDate:    postedDate,
			Diff:          roundDifferential(s),
			IsTournament:  isTournamentScore(s),
			BlowupHoles:   countBlowupHoles(s),
			IsBlowUpRound: detectBlowUpRoundPattern(s),
			PostingDelay:  delay,
			IsLatePosted:  delay >= 7,
			PostingStyle:  classifyPostingStyle(s),
		})
	}

	sort.Slice(rounds, func(i, j int) bool {
		return rounds[i].PlayedDate.After(rounds[j].PlayedDate)
	})

	a := GolferAnalysis{
		GolferID:    golferID,
		TotalRounds: len(rounds),
	}

	if len(rounds) == 0 {
		a.Flags = append(a.Flags, "No valid rounds")
		return a
	}
	a.RoundsEndDate = rounds[0].PlayedDate.Format("2006-01-02")
	a.RoundsStartDate = rounds[len(rounds)-1].PlayedDate.Format("2006-01-02")
	a.PlayedAtDistributionGraph = buildPlayedAtDistributionGraph(rounds)

	var allDiffs []float64
	var tournamentDiffs []float64
	var casualDiffs []float64

	var allBlowups []float64
	var tournamentBlowups []float64
	var casualBlowups []float64

	var allDelays []float64
	var tournamentDelays []float64
	var casualDelays []float64

	for _, r := range rounds {
		allDiffs = append(allDiffs, r.Diff)
		allBlowups = append(allBlowups, float64(r.BlowupHoles))
		allDelays = append(allDelays, float64(r.PostingDelay))

		if r.IsLatePosted {
			a.LatePostCount++
		}
		switch r.PostingStyle {
		case postingStyleHoleByHole:
			a.PostingHoleByHoleCount++
		case postingStyleFrontBackTotal:
			a.PostingFrontBackTotalCount++
		default:
			a.PostingTotalOnlyCount++
		}

		if r.IsTournament {
			tournamentDiffs = append(tournamentDiffs, r.Diff)
			tournamentBlowups = append(tournamentBlowups, float64(r.BlowupHoles))
			tournamentDelays = append(tournamentDelays, float64(r.PostingDelay))
			a.TournamentRounds++
			if r.IsBlowUpRound {
				a.BlowUpPatternTourney++
			}
			if r.IsLatePosted {
				a.LatePostCountTourney++
			}
			switch r.PostingStyle {
			case postingStyleHoleByHole:
				a.PostingHoleByHoleCountTourney++
			case postingStyleFrontBackTotal:
				a.PostingFrontBackCountTourney++
			default:
				a.PostingTotalOnlyCountTourney++
			}
		} else {
			casualDiffs = append(casualDiffs, r.Diff)
			casualBlowups = append(casualBlowups, float64(r.BlowupHoles))
			casualDelays = append(casualDelays, float64(r.PostingDelay))
			a.NonTournamentRounds++
			if r.IsBlowUpRound {
				a.BlowUpPatternCasual++
			}
			if r.IsLatePosted {
				a.LatePostCountCasual++
			}
			switch r.PostingStyle {
			case postingStyleHoleByHole:
				a.PostingHoleByHoleCountCasual++
			case postingStyleFrontBackTotal:
				a.PostingFrontBackCountCasual++
			default:
				a.PostingTotalOnlyCountCasual++
			}
		}
		if r.IsBlowUpRound {
			a.BlowUpPatternRounds++
		}
	}

	last20 := make([]float64, 0, min(20, len(rounds)))
	for i := 0; i < len(rounds) && i < 20; i++ {
		last20 = append(last20, rounds[i].Diff)
	}
	a.HandicapIndex = ComputeHandicapIndex(last20)
	traj := analyzeHandicapTrajectory(rounds)
	a.TrajectorySamples = traj.Samples
	a.TrajectoryPatternCount = traj.Patterns
	a.TrajectoryPatternRate = traj.PatternRate
	a.AvgPreTournamentRise = traj.AvgPreTournamentRise
	a.AvgPostTournamentDrop = traj.AvgPostTournamentDrop

	a.AvgDifferential = avg(allDiffs)
	a.AvgTournamentDifferential = avgOrNaN(tournamentDiffs)
	a.AvgCasualDifferential = avgOrNaN(casualDiffs)
	a.StdDevDifferential = stddev(allDiffs)
	a.SkewnessDifferential = skewness(allDiffs)
	a.SkewnessTournament = skewness(tournamentDiffs)
	a.SkewnessCasual = skewness(casualDiffs)
	if !math.IsNaN(a.SkewnessCasual) && !math.IsNaN(a.SkewnessTournament) {
		a.SkewnessGap = a.SkewnessCasual - a.SkewnessTournament
	}

	a.AvgBlowupHolesPerRound = avg(allBlowups)
	a.AvgBlowupHolesTourney = avgOrNaN(tournamentBlowups)
	a.AvgBlowupHolesCasual = avgOrNaN(casualBlowups)

	if !math.IsNaN(a.AvgBlowupHolesCasual) && !math.IsNaN(a.AvgBlowupHolesTourney) {
		a.BlowupGap = a.AvgBlowupHolesCasual - a.AvgBlowupHolesTourney
	}

	a.AvgPostingDelayDays = avg(allDelays)
	a.AvgPostingDelayDaysTourney = avgOrNaN(tournamentDelays)
	a.AvgPostingDelayDaysCasual = avgOrNaN(casualDelays)

	if a.TotalRounds > 0 {
		a.LatePostRate = float64(a.LatePostCount) / float64(a.TotalRounds)
		a.PostingHoleByHoleRate = float64(a.PostingHoleByHoleCount) / float64(a.TotalRounds)
		a.PostingFrontBackTotalRate = float64(a.PostingFrontBackTotalCount) / float64(a.TotalRounds)
		a.PostingTotalOnlyRate = float64(a.PostingTotalOnlyCount) / float64(a.TotalRounds)
	}
	if a.TournamentRounds > 0 {
		a.LatePostRateTourney = float64(a.LatePostCountTourney) / float64(a.TournamentRounds)
		a.PostingHoleByHoleRateTourney = float64(a.PostingHoleByHoleCountTourney) / float64(a.TournamentRounds)
		a.PostingFrontBackRateTourney = float64(a.PostingFrontBackCountTourney) / float64(a.TournamentRounds)
		a.PostingTotalOnlyRateTourney = float64(a.PostingTotalOnlyCountTourney) / float64(a.TournamentRounds)
	}
	if a.NonTournamentRounds > 0 {
		a.LatePostRateCasual = float64(a.LatePostCountCasual) / float64(a.NonTournamentRounds)
		a.BlowUpPatternRateC = float64(a.BlowUpPatternCasual) / float64(a.NonTournamentRounds)
		a.PostingHoleByHoleRateCasual = float64(a.PostingHoleByHoleCountCasual) / float64(a.NonTournamentRounds)
		a.PostingFrontBackRateCasual = float64(a.PostingFrontBackCountCasual) / float64(a.NonTournamentRounds)
		a.PostingTotalOnlyRateCasual = float64(a.PostingTotalOnlyCountCasual) / float64(a.NonTournamentRounds)
	}
	if a.TournamentRounds > 0 {
		a.BlowUpPatternRateT = float64(a.BlowUpPatternTourney) / float64(a.TournamentRounds)
	}
	if a.TotalRounds > 0 {
		a.BlowUpPatternRate = float64(a.BlowUpPatternRounds) / float64(a.TotalRounds)
	}
	a.BlowUpPatternRateGap = a.BlowUpPatternRateC - a.BlowUpPatternRateT

	if !math.IsNaN(a.AvgPostingDelayDaysCasual) && !math.IsNaN(a.AvgPostingDelayDaysTourney) {
		a.LatePostDelayGap = a.AvgPostingDelayDaysCasual - a.AvgPostingDelayDaysTourney
	}
	a.LatePostRateGap = a.LatePostRateCasual - a.LatePostRateTourney
	a.PostingTotalOnlyRateGap = a.PostingTotalOnlyRateCasual - a.PostingTotalOnlyRateTourney

	sorted := append([]float64(nil), last20...)
	sort.Float64s(sorted)
	bestCount := min(8, len(sorted))
	a.Best8Avg = avg(sorted[:bestCount])
	if len(sorted) > bestCount {
		a.RemainingAvg = avg(sorted[bestCount:])
		a.BestVsRestGap = a.RemainingAvg - a.Best8Avg
	} else {
		a.RemainingAvg = math.NaN()
		a.BestVsRestGap = math.NaN()
	}

	if !math.IsNaN(a.HandicapIndex) {
		overallBetter := 0
		for _, d := range allDiffs {
			if d < a.HandicapIndex {
				overallBetter++
			}
		}
		a.OverallBeatsHandicapRate = float64(overallBetter) / float64(len(allDiffs))

		if len(tournamentDiffs) > 0 {
			tBetter := 0
			for _, d := range tournamentDiffs {
				if d < a.HandicapIndex {
					tBetter++
				}
			}
			a.TournamentBeatsHandicapRate = float64(tBetter) / float64(len(tournamentDiffs))
		}

		if len(casualDiffs) > 0 {
			cBetter := 0
			for _, d := range casualDiffs {
				if d < a.HandicapIndex {
					cBetter++
				}
			}
			a.CasualBeatsHandicapRate = float64(cBetter) / float64(len(casualDiffs))
		}
	}

	if len(tournamentDiffs) > 0 && len(casualDiffs) > 0 {
		a.TournamentVsCasualGap = a.AvgCasualDifferential - a.AvgTournamentDifferential
	}
	if len(tournamentDiffs) > 0 && !math.IsNaN(a.HandicapIndex) {
		a.TournamentVsHandicapGap = a.HandicapIndex - a.AvgTournamentDifferential
	}

	a.SuspicionScore, a.Flags = buildSuspicionScore(a)
	return a
}

func isUsableScore(s ghin.Score) bool {
	if s.Status != "Validated" {
		return false
	}
	if s.NumberOfPlayedHoles <= 0 {
		return false
	}
	if roundDifferential(s) == 0 && s.AdjustedGrossScore == 0 {
		return false
	}
	return true
}

func isTournamentScore(s ghin.Score) bool {
	switch s.ScoreType {
	case "T", "C":
		return true
	default:
		return false
	}
}

func roundDifferential(s ghin.Score) float64 {
	// For 9-hole rounds GHIN can populate adjusted_scaled_up_differential.
	// Prefer it when present/non-zero; otherwise fall back to differential.
	if s.AdjustedScaledUpDifferential != 0 {
		return s.AdjustedScaledUpDifferential
	}
	return s.Differential
}

func countBlowupHoles(s ghin.Score) int {
	n := 0
	for _, h := range s.HoleDetails {
		if h.RawScore-h.Par >= 3 {
			n++
		}
	}
	return n
}

func detectBlowUpRoundPattern(s ghin.Score) bool {
	if len(s.HoleDetails) < 9 {
		return false
	}

	severeBlowups := 0
	restHoles := 0
	var restOverPar int

	for _, h := range s.HoleDetails {
		overPar := h.RawScore - h.Par
		if overPar >= 3 {
			severeBlowups++
			continue
		}
		restHoles++
		restOverPar += overPar
	}

	if severeBlowups < 1 || severeBlowups > 2 || restHoles < 6 {
		return false
	}

	restAvgOverPar := float64(restOverPar) / float64(restHoles)
	return restAvgOverPar <= 0.25
}

func classifyPostingStyle(s ghin.Score) string {
	if len(s.HoleDetails) >= min(max(s.NumberOfPlayedHoles, 0), 9) && len(s.HoleDetails) > 0 {
		return postingStyleHoleByHole
	}
	if hasFrontBackSummary(s) {
		return postingStyleFrontBackTotal
	}
	return postingStyleTotalOnly
}

func hasFrontBackSummary(s ghin.Score) bool {
	if s.Front9Adjusted > 0 || s.Back9Adjusted > 0 {
		return true
	}
	if s.Front9SlopeRating > 0 || s.Back9SlopeRating > 0 {
		return true
	}
	if s.Front9CourseRating > 0 || s.Back9CourseRating > 0 {
		return true
	}
	if strings.TrimSpace(s.Front9CourseName) != "" || strings.TrimSpace(s.Back9CourseName) != "" {
		return true
	}
	if strings.TrimSpace(s.Front9TeeName) != "" || strings.TrimSpace(s.Back9TeeName) != "" {
		return true
	}
	return false
}

func parsePostedAtDate(s string) (time.Time, error) {
	// GHIN example: 2025-09-24T01:51:54.938Z
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t.UTC().Truncate(24 * time.Hour), nil
	}
	// Fallback if time portion is ever absent
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unable to parse posted_at: %q", s)
}

func postingDelayDays(playedDate, postedDate time.Time) int {
	delay := int(postedDate.Sub(playedDate).Hours() / 24)
	if delay < 0 {
		return 0
	}
	return delay
}

func analyzeHandicapTrajectory(rounds []scoredRound) handicapTrajectoryMetrics {
	metrics := handicapTrajectoryMetrics{
		AvgPreTournamentRise:  math.NaN(),
		AvgPostTournamentDrop: math.NaN(),
	}

	if len(rounds) < 12 {
		return metrics
	}

	asc := append([]scoredRound(nil), rounds...)
	sort.Slice(asc, func(i, j int) bool {
		return asc[i].PlayedDate.Before(asc[j].PlayedDate)
	})

	rollingHI := make([]float64, len(asc))
	for i := range asc {
		start := i - 19
		if start < 0 {
			start = 0
		}
		diffs := make([]float64, 0, i-start+1)
		for j := start; j <= i; j++ {
			diffs = append(diffs, asc[j].Diff)
		}
		rollingHI[i] = ComputeHandicapIndex(diffs)
	}

	rises := make([]float64, 0)
	drops := make([]float64, 0)

	for i := range asc {
		if !asc[i].IsTournament {
			continue
		}
		if i < 8 || i+3 >= len(asc) {
			continue
		}

		preLong := avg(rollingHI[i-8 : i-3])   // 5 rounds
		preShort := avg(rollingHI[i-3 : i])    // 3 rounds before tournament
		postShort := avg(rollingHI[i+1 : i+4]) // 3 rounds after tournament
		if math.IsNaN(preLong) || math.IsNaN(preShort) || math.IsNaN(postShort) {
			continue
		}

		rise := preShort - preLong
		drop := preShort - postShort
		rises = append(rises, rise)
		drops = append(drops, drop)

		if rise >= 0.70 && drop >= 0.70 {
			metrics.Patterns++
		}
	}

	metrics.Samples = len(rises)
	if metrics.Samples == 0 {
		return metrics
	}

	metrics.AvgPreTournamentRise = avg(rises)
	metrics.AvgPostTournamentDrop = avg(drops)
	metrics.PatternRate = float64(metrics.Patterns) / float64(metrics.Samples)

	return metrics
}

func buildSuspicionScore(a GolferAnalysis) (float64, []string) {
	var score float64
	var flags []string

	if a.TotalRounds < 8 {
		flags = append(flags, "Limited data")
		return 0, flags
	}

	if a.TournamentRounds >= 3 && a.NonTournamentRounds >= 3 {
		switch {
		case a.TournamentVsCasualGap >= 5.0:
			score += 45
			flags = append(flags, "Tournament differentials dramatically better than casual rounds")
		case a.TournamentVsCasualGap >= 3.0:
			score += 28
			flags = append(flags, "Tournament differentials materially better than casual rounds")
		case a.TournamentVsCasualGap >= 2.0:
			score += 15
			flags = append(flags, "Tournament differentials moderately better than casual rounds")
		}
	}

	if a.TournamentRounds >= 3 {
		switch {
		case a.TournamentVsHandicapGap >= 4.0:
			score += 35
			flags = append(flags, "Tournament differentials far better than handicap index")
		case a.TournamentVsHandicapGap >= 2.5:
			score += 20
			flags = append(flags, "Tournament differentials better than handicap index")
		}
	}

	// Low Round Frequency Test:
	// Typical golfers beat their handicap index ~20-25% of the time.
	// Tournament beat-rates in the 50-70%+ range are a strong signal.
	if a.TournamentRounds >= 3 {
		switch {
		case a.TournamentBeatsHandicapRate >= 0.70:
			score += 25
			flags = append(flags, "Low Round Frequency Test: tournament rounds beat handicap at an extremely high rate (70%+)")
		case a.TournamentBeatsHandicapRate >= 0.50:
			score += 15
			flags = append(flags, "Low Round Frequency Test: tournament rounds beat handicap at a suspiciously high rate (50%+)")
		}
	}

	if !math.IsNaN(a.BestVsRestGap) {
		// Variance Between Best 8 and Remaining Rounds:
		// Typical players are often around a 3-4 stroke gap; >6 is a strong signal.
		switch {
		case a.BestVsRestGap > 6.0:
			score += 20
			flags = append(flags, "Best-8 vs remaining-round gap is suspiciously high (>6 strokes)")
		case a.BestVsRestGap >= 4.5:
			score += 12
			flags = append(flags, "Best-8 vs remaining-round gap is above normal range")
		}
	}

	if a.StdDevDifferential >= 6.0 {
		score += 12
		flags = append(flags, "Unusually high scoring variance")
	} else if a.StdDevDifferential >= 4.5 {
		score += 6
		flags = append(flags, "Moderately high scoring variance")
	}

	if a.TotalRounds >= 12 {
		switch {
		case a.OverallBeatsHandicapRate >= 0.45:
			score += 10
			flags = append(flags, "Low Round Frequency Test: overall beat-handicap rate is far above expected (20-25%)")
		case a.OverallBeatsHandicapRate >= 0.35:
			score += 5
			flags = append(flags, "Low Round Frequency Test: overall beat-handicap rate is above expected (20-25%)")
		}
	}

	if a.TournamentRounds >= 3 && a.NonTournamentRounds >= 3 {
		switch {
		case a.BlowupGap >= 1.5:
			score += 18
			flags = append(flags, "Casual rounds contain substantially more blow-up holes than tournament rounds")
		case a.BlowupGap >= 0.8:
			score += 10
			flags = append(flags, "Casual rounds contain more blow-up holes than tournament rounds")
		}
	}

	// Blow-Up Round Pattern:
	// Suspicious when many rounds have just 1-2 severe blow-up holes while the
	// rest of the round is near scratch-level scoring.
	switch {
	case a.BlowUpPatternRate >= 0.30 && a.BlowUpPatternRounds >= 4:
		score += 14
		flags = append(flags, "Frequent blow-up round pattern (1-2 severe holes with otherwise near-scratch play)")
	case a.BlowUpPatternRate >= 0.20 && a.BlowUpPatternRounds >= 3:
		score += 8
		flags = append(flags, "Moderate blow-up round pattern frequency")
	}

	if a.TournamentRounds >= 3 && a.NonTournamentRounds >= 3 {
		switch {
		case a.BlowUpPatternRateGap >= 0.25 && a.BlowUpPatternCasual >= 3:
			score += 10
			flags = append(flags, "Blow-up round pattern appears much more often in casual rounds than tournaments")
		case a.BlowUpPatternRateGap >= 0.15 && a.BlowUpPatternCasual >= 2:
			score += 5
			flags = append(flags, "Blow-up round pattern appears more often in casual rounds than tournaments")
		}
	}

	// Time-Series Handicap Trajectory:
	// Flag repeated sequences where handicap trends up before tournament rounds
	// then drops shortly after.
	if a.TrajectorySamples >= 2 && !math.IsNaN(a.AvgPreTournamentRise) && !math.IsNaN(a.AvgPostTournamentDrop) {
		switch {
		case a.TrajectoryPatternRate >= 0.60 && a.TrajectoryPatternCount >= 2 &&
			a.AvgPreTournamentRise >= 0.80 && a.AvgPostTournamentDrop >= 0.80:
			score += 18
			flags = append(flags, "Handicap repeatedly rises before tournaments and drops soon after")
		case a.TrajectoryPatternRate >= 0.40 && a.TrajectoryPatternCount >= 2 &&
			a.AvgPreTournamentRise >= 0.50 && a.AvgPostTournamentDrop >= 0.50:
			score += 10
			flags = append(flags, "Handicap trajectory shows pre-tournament rises and post-tournament drops")
		}
	}

	// Distribution skew signal:
	// Positive skew means frequent low (better) differentials with a right tail of higher (worse) rounds.
	// A large casual-vs-tournament skew gap can indicate clustering of very low tournament differentials.
	if a.TournamentRounds >= 4 && a.NonTournamentRounds >= 4 &&
		!math.IsNaN(a.SkewnessTournament) && !math.IsNaN(a.SkewnessCasual) {
		switch {
		case a.SkewnessGap >= 1.10:
			score += 10
			flags = append(flags, "Large casual-vs-tournament differential skew gap")
		case a.SkewnessGap >= 0.75:
			score += 6
			flags = append(flags, "Moderate casual-vs-tournament differential skew gap")
		}
	}

	if a.TournamentRounds >= 4 && !math.IsNaN(a.SkewnessTournament) {
		switch {
		case a.SkewnessTournament <= -0.90:
			score += 8
			flags = append(flags, "Tournament differentials are strongly left-skewed")
		case a.SkewnessTournament <= -0.60:
			score += 4
			flags = append(flags, "Tournament differentials are moderately left-skewed")
		}
	}

	// Late-posting signal:
	// Treat this as a secondary indicator, especially when 7+ day late posts
	// happen often, or much more often in casual rounds than tournament rounds.
	switch {
	case a.LatePostRate >= 0.40 && a.LatePostCount >= 5:
		score += 15
		flags = append(flags, "High rate of rounds posted 7+ days after play")
	case a.LatePostRate >= 0.25 && a.LatePostCount >= 4:
		score += 8
		flags = append(flags, "Moderate rate of rounds posted 7+ days after play")
	}

	switch {
	case a.AvgPostingDelayDays >= 10 && a.LatePostCount >= 4:
		score += 8
		flags = append(flags, "Average posting delay is unusually high")
	case a.AvgPostingDelayDays >= 7 && a.LatePostCount >= 3:
		score += 4
		flags = append(flags, "Average posting delay is elevated")
	}

	if a.TournamentRounds >= 3 && a.NonTournamentRounds >= 3 {
		switch {
		case a.LatePostRateGap >= 0.35 && a.LatePostCountCasual >= 4:
			score += 10
			flags = append(flags, "Casual rounds are posted late much more often than tournament rounds")
		case a.LatePostRateGap >= 0.20 && a.LatePostCountCasual >= 3:
			score += 5
			flags = append(flags, "Casual rounds are posted late more often than tournament rounds")
		}
	}

	// Posting-detail signal:
	// Treat repeated total-only posting as a secondary indicator, especially when
	// detailed hole-by-hole entries are rare or casual rounds are much more likely
	// to omit scoring detail.
	switch {
	case a.PostingTotalOnlyRate >= 0.60 && a.PostingTotalOnlyCount >= 6:
		score += 10
		flags = append(flags, "High share of scores are posted as total-only with no hole detail")
	case a.PostingTotalOnlyRate >= 0.40 && a.PostingTotalOnlyCount >= 4:
		score += 5
		flags = append(flags, "Moderate share of scores are posted as total-only")
	}

	if a.TournamentRounds >= 3 && a.NonTournamentRounds >= 3 {
		switch {
		case a.PostingTotalOnlyRateGap >= 0.35 && a.PostingTotalOnlyCountCasual >= 4:
			score += 8
			flags = append(flags, "Casual rounds are posted as total-only much more often than tournament rounds")
		case a.PostingTotalOnlyRateGap >= 0.20 && a.PostingTotalOnlyCountCasual >= 3:
			score += 4
			flags = append(flags, "Casual rounds are posted as total-only more often than tournament rounds")
		}
	}

	if a.TournamentRounds < 3 {
		score *= 0.65
		flags = append(flags, "Low tournament sample size")
	}

	return score, flags
}

func ComputeHandicapIndex(differentials []float64) float64 {
	if len(differentials) == 0 {
		return math.NaN()
	}

	cp := append([]float64(nil), differentials...)
	sort.Float64s(cp)

	n := len(cp)
	switch {
	case n >= 20:
		return avg(cp[:8])
	case n == 19:
		return avg(cp[:7])
	case n == 18:
		return avg(cp[:6])
	case n >= 15:
		return avg(cp[:5])
	case n >= 12:
		return avg(cp[:4])
	case n >= 9:
		return avg(cp[:3])
	case n >= 7:
		return avg(cp[:2])
	default:
		return cp[0]
	}
}

func avg(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	var sum float64
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func avgOrNaN(xs []float64) float64 {
	return avg(xs)
}

func stddev(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	m := avg(xs)
	var ss float64
	for _, x := range xs {
		d := x - m
		ss += d * d
	}
	return math.Sqrt(ss / float64(len(xs)))
}

func skewness(xs []float64) float64 {
	if len(xs) < 3 {
		return math.NaN()
	}

	m := avg(xs)
	s := stddev(xs)
	if s == 0 || math.IsNaN(s) {
		return 0
	}

	var m3 float64
	for _, x := range xs {
		d := x - m
		m3 += d * d * d
	}
	m3 /= float64(len(xs))

	return m3 / (s * s * s)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildPlayedAtDistributionGraph(rounds []scoredRound) string {
	if len(rounds) == 0 {
		return "  (no rounds)"
	}

	start := rounds[len(rounds)-1].PlayedDate
	end := rounds[0].PlayedDate
	startMonth := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)

	monthCounts := make(map[string][31]int)
	for _, r := range rounds {
		key := r.PlayedDate.Format("2006-01")
		current := monthCounts[key]
		day := r.PlayedDate.Day() // 1..31
		current[day-1]++
		monthCounts[key] = current
	}

	var b strings.Builder
	for m := startMonth; !m.After(endMonth); m = m.AddDate(0, 1, 0) {
		label := m.Format("2006-01")
		dayCounts := monthCounts[label]
		monthTotal := 0

		b.WriteString("  ")
		b.WriteString(label)
		b.WriteString(" | ")
		for i := 0; i < 31; i++ {
			n := dayCounts[i]
			monthTotal += n
			if n == 0 {
				b.WriteByte(' ')
				continue
			}
			if n > 9 {
				b.WriteByte('+')
				continue
			}
			b.WriteByte(byte('0' + n))
		}
		b.WriteString(fmt.Sprintf(" (%d)", monthTotal))
		b.WriteByte('\n')
	}
	b.WriteString("            1234567890123456789012345678901\n")
	b.WriteString("                     Day of Month")
	return strings.TrimSuffix(b.String(), "\n")
}

// func assignNamesFromInputPath(results []GolferAnalysis, inputPath string) {
// 	first, last, id, ok := parseNameAndIDFromFilename(filepath.Base(inputPath))
// 	if !ok {
// 		return
// 	}
// 	for i := range results {
// 		if results[i].GolferID == id {
// 			results[i].FirstName = first
// 			results[i].LastName = last
// 		}
// 	}
// }

// func parseNameAndIDFromFilename(fileName string) (firstName, lastName string, golferID int64, ok bool) {
// 	base := strings.TrimSuffix(fileName, filepath.Ext(fileName))
// 	parts := strings.Split(base, "_")
// 	if len(parts) < 3 {
// 		return "", "", 0, false
// 	}
// 	idPart := parts[len(parts)-1]
// 	id, err := strconv.ParseInt(idPart, 10, 64)
// 	if err != nil {
// 		return "", "", 0, false
// 	}
// 	last := parts[0]
// 	first := parts[1]
// 	return sanitizeNameToken(first), sanitizeNameToken(last), id, true
// }

// func sanitizeNameToken(s string) string {
// 	s = strings.TrimSpace(s)
// 	s = strings.ReplaceAll(s, "-", " ")
// 	s = strings.ReplaceAll(s, ".", " ")
// 	s = strings.Join(strings.Fields(s), " ")
// 	if s == "" {
// 		return "Unknown"
// 	}
// 	return s
// }

func slugToken(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return "unknown"
	}
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		if r == ' ' || r == '_' || r == '-' {
			b.WriteRune('_')
		}
	}
	out := strings.Trim(b.String(), "_")
	if out == "" {
		return "unknown"
	}
	return out
}
