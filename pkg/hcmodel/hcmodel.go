package hcmodel

type Golfer struct {
	GolferID       int    `json:"golferID" dynamodbav:"golferID"`
	FirstName      string `json:"firstName" dynamodbav:"firstName"`
	LastName       string `json:"lastName" dynamodbav:"lastName"`
	CreatedAt      string `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt      string `json:"updatedAt" dynamodbav:"updatedAt"`
	ScoreStartDate string `json:"scoreStartDate,omitempty" dynamodbav:"scoreStartDate,omitempty"`
	ScoreEndDate   string `json:"scoreEndDate,omitempty" dynamodbav:"scoreEndDate,omitempty"`
	LastGhinUpdate string `json:"lastGhinUpdate,omitempty" dynamodbav:"lastGhinUpdate,omitempty"`
	Email          string `json:"email,omitempty" dynamodbav:"email,omitempty"`
}
