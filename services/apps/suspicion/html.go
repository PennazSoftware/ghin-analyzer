package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"
)

func writeResultsEmailHTML(results []GolferAnalysis) (string, error) {
	if len(results) == 0 {
		return "", fmt.Errorf("no results to render")
	}

	funcMap := template.FuncMap{
		"mul100": func(v float64) float64 { return v * 100 },
		"join":   func(xs []string, sep string) string { return strings.Join(xs, sep) },
	}
	tpl := template.Must(template.New("email").Funcs(funcMap).Parse(resultsEmailTemplate))
	var buf bytes.Buffer

	data := struct {
		GeneratedAt string
		Results     []GolferAnalysis
	}{
		GeneratedAt: time.Now().Format("Jan 2, 2006 3:04 PM MST"),
		Results:     results,
	}
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const resultsEmailTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body style="margin:0;padding:0;background:#f2f4f7;font-family:Arial,'Helvetica Neue',Helvetica,sans-serif;color:#1f2933;">
  <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="background:#f2f4f7;padding:20px 10px;">
    <tr>
      <td align="center">
        <table role="presentation" width="760" cellspacing="0" cellpadding="0" style="width:100%;max-width:760px;background:#ffffff;border-radius:10px;overflow:hidden;border:1px solid #d9e2ec;">
          <tr>
            <td style="background:#0f5132;color:#ffffff;padding:18px 20px;">
              <div style="font-size:22px;font-weight:700;line-height:1.2;">GHIN Handicap Suspicion Analysis Summary</div>
              <div style="font-size:13px;opacity:0.9;margin-top:4px;">Generated {{.GeneratedAt}}</div>
			  <div style="font-size:12px;opacity:0.9;margin-top:4px;">The handicap suspicion analysis evaluates golfers' scores for potential anomalies that may indicate unusual patterns or behaviors. In other words, it helps identify potential sandbaggers.</div>
              <div style="font-size:12px;opacity:0.95;margin-top:8px;">Full descriptions for all fields are available at <a href="https://www.pennaz.com/suspicion_scoring.html" style="color:#d9f99d;text-decoration:underline;">https://www.pennaz.com/suspicion_scoring.html</a>.</div>
            </td>
          </tr>
          <tr>
            <td style="padding:18px 20px 10px 20px;">
              <div style="font-size:16px;font-weight:700;color:#102a43;margin-bottom:10px;">Summary</div>
              <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;border:1px solid #e4e7eb;">
                <tr style="background:#f7fafc;">
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">Golfer</th>
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">ID</th>
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">Suspicion Score</th>
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">HI</th>
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">T/C Gap</th>
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">Rounds (T/C/Total)</th>
                  <th align="left" style="padding:8px 10px;font-size:12px;color:#486581;border-bottom:1px solid #e4e7eb;">Date Range</th>
                </tr>
                {{range .Results}}
                <tr>
                  <td style="padding:8px 10px;font-size:13px;border-bottom:1px solid #e4e7eb;">{{.FirstName}} {{.LastName}}</td>
                  <td style="padding:8px 10px;font-size:13px;border-bottom:1px solid #e4e7eb;">{{.GolferID}}</td>
                  <td style="padding:8px 10px;font-size:13px;font-weight:700;border-bottom:1px solid #e4e7eb;">{{printf "%.1f" .SuspicionScore}}</td>
                  <td style="padding:8px 10px;font-size:13px;border-bottom:1px solid #e4e7eb;">{{printf "%.2f" .HandicapIndex}}</td>
                  <td style="padding:8px 10px;font-size:13px;border-bottom:1px solid #e4e7eb;">{{printf "%.2f" .TournamentVsCasualGap}}</td>
                  <td style="padding:8px 10px;font-size:13px;border-bottom:1px solid #e4e7eb;">{{.TournamentRounds}} / {{.NonTournamentRounds}} / {{.TotalRounds}}</td>
                  <td style="padding:8px 10px;font-size:13px;border-bottom:1px solid #e4e7eb;">{{.RoundsStartDate}} to {{.RoundsEndDate}}</td>
                </tr>
                {{end}}
              </table>
            </td>
          </tr>
          {{range .Results}}
          <tr>
            <td style="padding:14px 20px 16px 20px;border-top:1px solid #e4e7eb;">
              <div style="font-size:15px;font-weight:700;color:#102a43;">{{.FirstName}} {{.LastName}} ({{.GolferID}})</div>
              <table role="presentation" width="100%" cellspacing="0" cellpadding="0" style="border-collapse:collapse;margin-top:8px;border:1px solid #e4e7eb;">
                <colgroup>
                  <col style="width:280px;min-width:280px;">
                  <col>
                </colgroup>
                <tr><td style="padding:6px 8px;font-size:12px;color:#486581;background:#f7fafc;">Metric</td><td style="padding:6px 8px;font-size:12px;color:#486581;background:#f7fafc;">Value</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Suspicion Score</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f" .SuspicionScore}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Rounds Analyzed (Total / Casual / Tournament)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{.TotalRounds}} / {{.NonTournamentRounds}} / {{.TournamentRounds}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Date Range Analyzed</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{.RoundsStartDate}} to {{.RoundsEndDate}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Handicap Index</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .HandicapIndex}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Avg Differential / Casual / Tournament</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .AvgDifferential}} / {{printf "%.2f" .AvgCasualDifferential}} / {{printf "%.2f" .AvgTournamentDifferential}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Tournament vs Casual Gap</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .TournamentVsCasualGap}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Tournament vs Handicap Gap</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .TournamentVsHandicapGap}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Best 8 vs Rest Gap (recent 20)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .BestVsRestGap}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Beat HI Rate (Overall / Casual / Tournament)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f%%" (mul100 .OverallBeatsHandicapRate)}} / {{printf "%.1f%%" (mul100 .CasualBeatsHandicapRate)}} / {{printf "%.1f%%" (mul100 .TournamentBeatsHandicapRate)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Skew (All / T / C)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .SkewnessDifferential}} / {{printf "%.2f" .SkewnessTournament}} / {{printf "%.2f" .SkewnessCasual}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Skew Gap (Casual - Tournament)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .SkewnessGap}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Handicap Trajectory Samples / Patterns / Rate</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{.TrajectorySamples}} / {{.TrajectoryPatternCount}} / {{printf "%.1f%%" (mul100 .TrajectoryPatternRate)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Handicap Trajectory Avg Rise / Drop</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .AvgPreTournamentRise}} / {{printf "%.2f" .AvgPostTournamentDrop}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Blow-Up Holes Avg (Overall / T / C)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .AvgBlowupHolesPerRound}} / {{printf "%.2f" .AvgBlowupHolesTourney}} / {{printf "%.2f" .AvgBlowupHolesCasual}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Blow-Up Pattern Rate (Overall / C / T)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f%%" (mul100 .BlowUpPatternRate)}} / {{printf "%.1f%%" (mul100 .BlowUpPatternRateC)}} / {{printf "%.1f%%" (mul100 .BlowUpPatternRateT)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Late Posting (Avg Days / Rate / C Rate / T Rate)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.2f" .AvgPostingDelayDays}} / {{printf "%.1f%%" (mul100 .LatePostRate)}} / {{printf "%.1f%%" (mul100 .LatePostRateCasual)}} / {{printf "%.1f%%" (mul100 .LatePostRateTourney)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Posting Detail Counts HxH/F&B/Tot (Overall)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{.PostingHoleByHoleCount}} / {{.PostingFrontBackTotalCount}} / {{.PostingTotalOnlyCount}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Posting Detail Rates HxH/F&B/Tot (Overall)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f%%" (mul100 .PostingHoleByHoleRate)}} / {{printf "%.1f%%" (mul100 .PostingFrontBackTotalRate)}} / {{printf "%.1f%%" (mul100 .PostingTotalOnlyRate)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Posting Detail Rates HxH/F&B/Tot (Casual)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f%%" (mul100 .PostingHoleByHoleRateCasual)}} / {{printf "%.1f%%" (mul100 .PostingFrontBackRateCasual)}} / {{printf "%.1f%%" (mul100 .PostingTotalOnlyRateCasual)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Posting Detail Rates HxH/F&B/Tot (Tournament)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f%%" (mul100 .PostingHoleByHoleRateTourney)}} / {{printf "%.1f%%" (mul100 .PostingFrontBackRateTourney)}} / {{printf "%.1f%%" (mul100 .PostingTotalOnlyRateTourney)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Total-Only Posting Gap (Casual - Tournament)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{printf "%.1f%%" (mul100 .PostingTotalOnlyRateGap)}}</td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">PlayedAt Day-of-Month Distribution (1-31)</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;"><pre style="margin:0;font-family:Menlo,Consolas,'Courier New',monospace;font-size:12px;line-height:1.35;white-space:pre;">{{.PlayedAtDistributionGraph}}</pre></td></tr>
                <tr><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">Flags</td><td style="padding:6px 8px;font-size:13px;border-top:1px solid #e4e7eb;">{{if .Flags}}<ul style="margin:0;padding-left:18px;">{{range .Flags}}<li style="margin:0 0 4px 0;">{{.}}</li>{{end}}</ul>{{else}}None{{end}}</td></tr>
              </table>
            </td>
          </tr>
          {{end}}
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`
