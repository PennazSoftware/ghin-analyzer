package ghin

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// API contains the externally accessible methods
type API interface {
	Login(ghinNumber string, password string) (LoginResponse, error)
	GetGolfer(ghinNumber int) ([]Golfer, error)
	GetRevisions(ghinNumber string, beginDate time.Time, endDate time.Time) (RevisionResponse, error)
	GetCourseHandicaps(ghinNumber string, courseID int) (CourseHandicapsResponse, error)
	GetCourseDetails(courseID int) (CourseDetailsResponse, error)
	GetScoresByDate(ghinNumber int, beginDate time.Time, endDate time.Time) ([]Score, error)
	GetScoresByDatePaged(ghinNumber string, beginDate time.Time, endDate time.Time, offset int) (ScoreResponse, error)
	GetChildScores(ghinNumber string, parentScoreID int) ([]Score, error)
	FindGolferByGhinNumber(ghinNumber string) ([]FoundGolfer, error)
	SearchCourses(searchName string) (CourseSearchResponse, error)
}

const (
	httpTimeout = 20 * time.Second
)

// Client is the struct that embodies the AWS DynamoDB client
type Client struct {
	ghinClient  *http.Client
	credentials Credentials
}

// verifying if the Client struct is indeed implementing the AwsDynamoDb interface
var _ API = (*Client)(nil)

// Credentials contains properties for authentication
type Credentials struct {
	Username string
	Password string
	Cookies  []*http.Cookie
}

// New creates a new GHIN client
func New() *Client {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	return &Client{
		ghinClient: client,
	}
}

// New creates a new GHIN client
func NewWithLogin(username string, password string) *Client {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	ghinClient := &Client{
		ghinClient: client,
	}

	loginResponse, err := ghinClient.Login(username, password)
	if err != nil {
		log.Printf("failed to login. Err: %+v", err)
		return nil
	}

	ghinClient.credentials.Username = username
	ghinClient.credentials.Password = password

	if len(loginResponse.GolferUser.Golfers) > 0 {
		log.Printf("Successfully logged-in to GHIN as %s %s", loginResponse.GolferUser.Golfers[0].FirstName, loginResponse.GolferUser.Golfers[0].LastName)
	}

	return ghinClient
}

func (c *Client) newRequest(method string, url string, body interface{}) (*http.Request, error) {
	log.Printf("GHIN Client newRequest(): Method: %s, URL: %s", method, url)

	var buf io.ReadWriter
	var req *http.Request
	var err error

	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			log.Printf("failed to json encode body. Err: %+v", err)
			return nil, err
		}
		req, err = http.NewRequest(method, url, buf)
		if err != nil {
			log.Printf("failed on http.NewRequest() with body. Method: %s, URL: %s, Err: %+v", method, url, err)
			return nil, err
		}

	} else {
		req, err = http.NewRequest(method, url, buf)
		if err != nil {
			log.Printf("failed on http.NewRequest(). Method: %s, URL: %s, Err: %+v", method, url, err)
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")

	// Set Cookies
	for _, cookie := range c.credentials.Cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {

	// fmt.Println("================== REQUEST HEADERS ==================")
	// for k, v := range req.Header {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }

	resp, err := c.ghinClient.Do(req)
	if err != nil {
		log.Println("failed to call httpClient.Do")
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		// Make a copy of the response body so we can read it for debugging/logging purposes
		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)

		// Log the response body for debugging purposes
		//log.Printf("Response Status: %s", resp.Status)
		//log.Printf("Response Body: %s", buf.String())

		err = json.NewDecoder(tee).Decode(v)
		if err != nil {
			log.Printf("failed to decode response body. Err: %+v", err)
			log.Printf("response body: %s", buf.String())
			return resp, err
		}
	}

	// fmt.Println("================== RESPONSE HEADERS ==================")
	// for k, v := range resp.Header {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }

	return resp, err
}
