package ghin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	golfer_template = "https://api2.ghin.com/api/v1/golfers.json?status=Active&from_ghin=true&per_page=25&page=1&golfer_id={{GHIN_NUMBER}}&includeLowHandicapIndex=true&source=GHINcom"
)

type GolfersResponse struct {
	Golfers []Golfer `json:"golfers,omitempty"`
}

type Golfer struct {
	FirstName             string  `json:"first_name,omitempty"`
	LastName              string  `json:"last_name,omitempty"`
	Gender                string  `json:"gender,omitempty"`
	Email                 string  `json:"email,omitempty"`
	PhoneNumber           string  `json:"phone_number,omitempty"`
	Suffix                any     `json:"suffix,omitempty"`
	Prefix                any     `json:"prefix,omitempty"`
	MiddleName            any     `json:"middle_name,omitempty"`
	IsTrial               bool    `json:"is_trial,omitempty"`
	Status                string  `json:"status,omitempty"`
	Ghin                  string  `json:"ghin,omitempty"`
	HandicapIndex         string  `json:"handicap_index,omitempty"`
	AssociationID         int     `json:"association_id,omitempty"`
	AssociationName       string  `json:"association_name,omitempty"`
	ClubName              string  `json:"club_name,omitempty"`
	ClubID                int     `json:"club_id,omitempty"`
	State                 string  `json:"state,omitempty"`
	Country               string  `json:"country,omitempty"`
	LowHi                 string  `json:"low_hi,omitempty"`
	SoftCap               string  `json:"soft_cap,omitempty"`
	HardCap               string  `json:"hard_cap,omitempty"`
	Entitlement           bool    `json:"entitlement,omitempty"`
	ClubAffiliationID     int     `json:"club_affiliation_id,omitempty"`
	IsHomeClub            bool    `json:"is_home_club,omitempty"`
	RevDate               string  `json:"rev_date,omitempty"`
	HiValue               float64 `json:"hi_value,omitempty"`
	HiDisplay             string  `json:"hi_display,omitempty"`
	MessageClubAuthorized any     `json:"message_club_authorized,omitempty"`
	LowHiValue            float64 `json:"low_hi_value,omitempty"`
	LowHiDisplay          string  `json:"low_hi_display,omitempty"`
	LowHiDate             string  `json:"low_hi_date,omitempty"`
	UseScaling            bool    `json:"use_scaling,omitempty"`
	HasDigitalProfile     bool    `json:"has_digital_profile,omitempty"`
}

// GetGolfer retrieves the golfer's information by their GHIN number
func (c *Client) GetGolfer(ghinNumber int) ([]Golfer, error) {
	var golferResponse GolfersResponse

	url := strings.ReplaceAll(golfer_template, "{{GHIN_NUMBER}}", strconv.Itoa(ghinNumber))

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetGolfer()")
		return []Golfer{}, err
	}

	res, err := c.do(req, &golferResponse)
	if err != nil {
		log.Printf("ghin client - failed to perform http request GetGolfer() -- %+v", err)
		return []Golfer{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []Golfer{}, fmt.Errorf("Error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return golferResponse.Golfers, nil
}
