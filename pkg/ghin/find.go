package ghin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	findGolferURLTemplate string = "https://api2.ghin.com/api/v1/golfermethods.asmx/FindGolfer.json?activeOnly=true&username=GHIN2020&password=GHIN2020&club=0&association=0&ghinNumber={{GHIN_NUMBER}}"
)

// FoundGolfer contains information about a golfer found using the FindGolfer methods
type FoundGolfer struct {
	GHINNumber         string `json:"GHINNumber"`
	FirstName          string `json:"FirstName"`
	MiddleName         string `json:"MiddleName"`
	LastName           string `json:"LastName"`
	Prefix             string `json:"Prefix"`
	Suffix             string `json:"Suffix"`
	ClubName           string `json:"ClubName"`
	AssocName          string `json:"AssocName"`
	ServiceName        string `json:"ServiceName"`
	Value              string `json:"Value"`
	TotalDiff          string `json:"TotalDiff"`
	IndexType          string `json:"IndexType"`
	Assoc              string `json:"Assoc"`
	Club               string `json:"Club"`
	ClubID             string `json:"ClubId"`
	Service            string `json:"Service"`
	Address1           string `json:"Address1"`
	Address2           string `json:"Address2"`
	City               string `json:"City"`
	State              string `json:"State"`
	Zip                string `json:"Zip"`
	MemberType         string `json:"MemberType"`
	DateOfBirth        string `json:"DateOfBirth"`
	Email              string `json:"Email"`
	StatusDate         string `json:"StatusDate"`
	PrimaryClubState   string `json:"PrimaryClubState"`
	PrimaryClubCountry string `json:"PrimaryClubCountry"`
	Local              string `json:"Local"`
	Holes              string `json:"Holes"`
	RevDate            string `json:"RevDate"`
	Active             string `json:"Active"`
	TScoreCount        string `json:"TScoreCount"`
	Display            string `json:"Display"`
	Gender             string `json:"Gender"`
	Type3              bool   `json:"Type3"`
	LowHI              string `json:"LowHI"`
	LowHIDisplay       string `json:"LowHIDisplay"`
	ClubType           string `json:"ClubType"`
	PlayerName         string `json:"PlayerName"`
	SearchValue        string `json:"SearchValue"`
	TrendValue         string `json:"TrendValue"`
	TrendTotalDiff     string `json:"TrendTotalDiff"`
	TrendIndexType     string `json:"TrendIndexType"`
	TrendRevDate       string `json:"TrendRevDate"`
	TrendTScoreCount   string `json:"TrendTScoreCount"`
	TrendDisplay       string `json:"TrendDisplay"`
	HiValue            string `json:"HiValue"`
	HiDisplay          string `json:"HiDisplay"`
	LowHiValue         string `json:"LowHiValue"`
	TechnologyProvider string `json:"TechnologyProvider"`
	MembershipPaidTime string `json:"MembershipPaidTime"`
	SoftCap            string `json:"SoftCap"`
	HardCap            string `json:"HardCap"`
}

// FindGolferByGhinNumber finds a golfer given their ghinNumber
func (c *Client) FindGolferByGhinNumber(ghinNumber string) ([]FoundGolfer, error) {
	foundGolfers := []FoundGolfer{}

	url := strings.Replace(findGolferURLTemplate, "{{GHIN_NUMBER}}", ghinNumber, -1)

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request FindGolferByGhinNumber()")
		return foundGolfers, err
	}

	res, err := c.do(req, &foundGolfers)
	if err != nil {
		log.Println("ghin client - failed to perform http request GetCourseHandicaps()")
		return foundGolfers, err
	}

	if res.StatusCode != http.StatusOK {
		return foundGolfers, fmt.Errorf("Error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return foundGolfers, nil
}
