package ghin

// CourseSearchResponse is the response returned from a call to search for courses
type CourseSearchResponse struct {
	Courses []SearchedCourse `json:"courses"`
}

type SearchedCourse struct {
	CourseID             int           `json:"CourseID"`
	CourseStatus         string        `json:"CourseStatus"`
	CourseName           string        `json:"CourseName"`
	GeoLocationLatitude  float64       `json:"GeoLocationLatitude"`
	GeoLocationLongitude float64       `json:"GeoLocationLongitude"`
	FacilityID           int           `json:"FacilityID"`
	FacilityStatus       string        `json:"FacilityStatus"`
	FacilityName         string        `json:"FacilityName"`
	FullName             string        `json:"FullName"`
	Address1             string        `json:"Address1"`
	Address2             string        `json:"Address2"`
	City                 string        `json:"City"`
	State                string        `json:"State"`
	Zip                  string        `json:"Zip"`
	Country              string        `json:"Country"`
	EntCountryCode       int           `json:"EntCountryCode"`
	EntStateCode         int           `json:"EntStateCode"`
	LegacyCRPCourseID    int           `json:"LegacyCRPCourseId"`
	Telephone            string        `json:"Telephone"`
	Email                string        `json:"Email"`
	UpdatedOn            string        `json:"UpdatedOn"`
	Ratings              []interface{} `json:"Ratings"`
}

///////////////////////////////////////////////////////////////////////
