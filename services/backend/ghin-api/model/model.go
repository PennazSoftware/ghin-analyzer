package model

import "PennazSoftware/ghin-analyzer/pkg/ghin"

// ResponseBase is the base response struct that all other responses should embed
type ResponseBase struct {
	// in: body
	// Status specifies whether the request was successful or not. If the value is 'ERROR' then consult the 'errorMessage' field for additional information about the error. Values: SUCCESS | ERROR
	Status string `json:"status,omitempty"`
	// ErrorMessage contains a descriptive explanation of the error in the situation where an error occurs
	// in: body
	ErrorMessage string `json:"errorMessage,omitempty"`
	// RequestID is a unique identifier for the request and is helpful for debugging
	// in: body
	RequestID string `json:"requestID,omitempty"`
	// Timestamp is the time of the request
	// in: body
	Timestamp string `json:"timestamp,omitempty"`
	// Page specifies the key necessary for results that are paged. If a value is provided, submitting subsequent requests with this value along the 'page' parameter will return the next set of results (i.e., &page=pagevalue)
	// in: body
	Page string `json:"page,omitempty"`
	// Cached indicates whether the response was returned from cache or if it was a live response from GHIN. Values: true | false
	Cached bool `json:"cached,omitempty"`
}

// ErrorResponse is a basic response containing only the base and can be used for all error response conditions
type ErrorResponse struct {
	ResponseBase
}

// CourseSearchResponse is the response returned from a call to search for courses
type CourseSearchResponse struct {
	ResponseBase
	Courses []ghin.SearchedCourse `json:"courses,omitempty"`
}

// CourseDetailsResponse is the response returned from a call to retrieve course details
type CourseDetailsResponse struct {
	ResponseBase
	Course ghin.CourseDetailsResponse `json:"course"`
}

// TeesResponse is the response returned from a call to retrieve tees for a course
type TeesResponse struct {
	ResponseBase
	Tees []ghin.TeeSet `json:"tees,omitempty"`
}

// CourseHandicapsResponse is the response returned from a call to retrieve a golfer's course handicaps
type CourseHandicapsResponse struct {
	ResponseBase
	TeeSets []ghin.TeeSet `json:"teeSets,omitempty"`
}

// GolferResponse is the response returned from a call to retrieve a golfer's information
type GolferResponse struct {
	ResponseBase
	Golfer ghin.Golfer `json:"golfer"`
}

// RevisionsResponse is the response returned from a call to retrieve a golfer's handicap revision history
type RevisionsResponse struct {
	ResponseBase
	Revisions []ghin.Revision `json:"revisions,omitempty"`
}
