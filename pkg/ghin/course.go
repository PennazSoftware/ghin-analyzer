package ghin

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	courseHandicapsTemplate string = "https://api2.ghin.com/api/v1/course_handicaps.json?golfer_id={{GHIN_NUMBER}}&course_id={{COURSE_ID}}"
	courseDetailsTemplate   string = "https://api2.ghin.com/api/v1/crsCourseMethods.asmx/GetCourseDetails.json?courseId={{COURSE_ID}}"
)

// CourseHandicapsResponse is the response returned from a call to get course handicaps
type CourseHandicapsResponse struct {
	TeeSets []TeeSet `json:"tee_sets"`
}

// TeeSet contains handicap information for a particular Tee
type TeeSet struct {
	TeeSetID int      `json:"tee_set_id"`
	Name     string   `json:"name"`
	Gender   string   `json:"gender"`
	Ratings  []Rating `json:"ratings"`
}

// Rating contains handicap rating information for a TeeSet
type Rating struct {
	TeeSetSide            string  `json:"tee_set_side"`
	CourseRating          float64 `json:"course_rating"`
	SlopeRating           int     `json:"slope_rating"`
	Par                   int     `json:"par"`
	CourseHandicap        int     `json:"course_handicap"`
	CourseHandicapDisplay string  `json:"course_handicap_display"`
}

// CourseDetailsResponse is the response returned from a call to get course details
type CourseDetailsResponse struct {
	CourseStatus string         `json:"CourseStatus"`
	Facility     Facility       `json:"Facility"`
	Season       Season         `json:"Season"`
	TeeSets      []TeeSetDetail `json:"TeeSets"`
}

// Facility contains information about a golf facility
type Facility struct {
	FacilityID                  int     `json:"FacilityId"`
	FacilityName                string  `json:"FacilityName"`
	FacilityNumber              string  `json:"FacilityNumber"`
	GeoLocationFormattedAddress string  `json:"GeoLocationFormattedAddress"`
	GeoLocationLongitude        float64 `json:"GeoLocationLongitude"`
	GeoLocationLatitude         float64 `json:"GeoLocationLatitude"`
}

// Season contains information about the season associated with a golf course
type Season struct {
	SeasonName      string `json:"SeasonName"`
	SeasonStartDate string `json:"SeasonStartDate"`
	SeasonEndDate   string `json:"SeasonEndDate"`
	IsAllYear       bool   `json:"IsAllYear"`
}

// TeeSetDetail contains details about a tee set for a course
type TeeSetDetail struct {
	Ratings          []DetailsRating `json:"Ratings"`
	Holes            []DetailsHole   `json:"Holes"`
	TeeSetRatingID   int             `json:"TeeSetRatingId"`
	TeeSetRatingName string          `json:"TeeSetRatingName"`
	Gender           string          `json:"Gender"`
	HolesNumber      int             `json:"HolesNumber"`
	TotalYardage     int             `json:"TotalYardage"`
	TotalMeters      int             `json:"TotalMeters"`
}

// DetailsRating contains rating information for a tee set
type DetailsRating struct {
	RatingType   string  `json:"RatingType"`
	CourseRating float64 `json:"CourseRating"`
	SlopeRating  float64 `json:"SlopeRating"`
	BogeyRating  float64 `json:"BogeyRating"`
}

// DetailsHole contains information about a tee set hole
type DetailsHole struct {
	Number     int `json:"Number"`
	HoleID     int `json:"HoleId"`
	Length     int `json:"Length"`
	Par        int `json:"Par"`
	Allocation int `json:"Allocation"`
}

// GetCourseHandicaps retrieves course handicaps for each tee set at a course
func (c *Client) GetCourseHandicaps(ghinNumber string, courseID int) (CourseHandicapsResponse, error) {
	var courseHandicapResponse CourseHandicapsResponse

	url := strings.Replace(courseHandicapsTemplate, "{{GHIN_NUMBER}}", ghinNumber, -1)
	url = strings.Replace(url, "{{COURSE_ID}}", strconv.Itoa(courseID), -1)

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetCourseHandicaps()")
		return courseHandicapResponse, err
	}

	res, err := c.do(req, &courseHandicapResponse)
	if err != nil {
		log.Println("ghin client - failed to perform http request GetCourseHandicaps()")
		return courseHandicapResponse, err
	}

	if res.StatusCode != http.StatusOK {
		return courseHandicapResponse, fmt.Errorf("error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return courseHandicapResponse, nil
}

// GetCourseDetails retrieves course details
func (c *Client) GetCourseDetails(courseID int) (CourseDetailsResponse, error) {
	var courseDetailsResponse CourseDetailsResponse

	url := strings.Replace(courseDetailsTemplate, "{{COURSE_ID}}", strconv.Itoa(courseID), -1)

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request GetCourseDetails()")
		return courseDetailsResponse, err
	}

	res, err := c.do(req, &courseDetailsResponse)
	if err != nil {
		log.Println("ghin client - failed to perform http request GetCourseDetails()")
		return courseDetailsResponse, err
	}

	if res.StatusCode != http.StatusOK {
		return courseDetailsResponse, fmt.Errorf("error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return courseDetailsResponse, nil
}

func (c *Client) SearchCourses(searchName string) (CourseSearchResponse, error) {
	var courseResponse CourseSearchResponse

	url := fmt.Sprintf("https://api2.ghin.com/api/v1/crsCourseMethods.asmx/SearchCourses.json?name=%s", url.QueryEscape(searchName))

	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("ghin client - failed to create new http request SearchCourses()")
		return courseResponse, err
	}

	res, err := c.do(req, &courseResponse)
	if err != nil {
		log.Println("ghin client - failed to perform http request SearchCourses()")
		return courseResponse, err
	}

	if res.StatusCode != http.StatusOK {
		return courseResponse, fmt.Errorf("error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	return courseResponse, nil
}
