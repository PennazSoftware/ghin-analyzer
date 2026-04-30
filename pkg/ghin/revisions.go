package ghin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	revisionURLTemplate = "https://api2.ghin.com/api/v1/golfers/{{GHIN_NUMBER}}/handicap_history.json?revCount=0&date_begin={{DATE_BEGIN}}&date_end={{DATE_END}}&source=GHINcom"
)

// RevisionResponse is the response returned from a call to get handicap revision scores
type RevisionResponse struct {
	Revisions []Revision `json:"handicap_revisions"`
}

// Revision contains information about a handicap revision update
type Revision struct {
	GHINNumber             string `json:"GHINNumber"`
	Assoc                  string `json:"Assoc"`
	Club                   string `json:"Club"`
	Service                string `json:"Service"`
	RevDate                string `json:"RevDate"`
	Display                string `json:"Display"`
	Value                  string `json:"Value"`
	LowHIDisplay           string `json:"LowHIDisplay"`
	LowHI                  string `json:"LowHI"`
	HIBeforeSoftCapDisplay string `json:"HIBeforeSoftCapDisplay"`
	HIBeforeSoftCap        string `json:"HIBeforeSoftCap"`
}

// GetRevisions retrieves the handicap revisions for the specified golfer
func (c *Client) GetRevisions(ghinNumber string, beginDate time.Time, endDate time.Time) (RevisionResponse, error) {
	var revisionResponse RevisionResponse

	url := strings.Replace(revisionURLTemplate, "{{GHIN_NUMBER}}", ghinNumber, -1)
	url = strings.Replace(url, "{{DATE_BEGIN}}", beginDate.Format("2006-01-02"), -1)
	url = strings.Replace(url, "{{DATE_END}}", endDate.Format("2006-01-02"), -1)

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetRevisions()")
		return revisionResponse, err
	}

	res, err := c.do(req, &revisionResponse)
	if err != nil {
		log.Println("ghin client - failed to perform http request GetRevisions()")
		return revisionResponse, err
	}

	if res.StatusCode != http.StatusOK {
		return revisionResponse, fmt.Errorf("Error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return revisionResponse, nil
}
