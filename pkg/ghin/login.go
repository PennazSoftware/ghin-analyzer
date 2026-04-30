package ghin

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// LoginRequest is the request sent to authenticate
type LoginRequest struct {
	User      User   `json:"user"`
	Token     string `json:"token"`
	UserToken string `json:"user_token"`
	Source    string `json:"source"`
}

type User struct {
	Password    string `json:"password"`
	EmailOrGhin string `json:"email_or_ghin"`
	RememberMe  bool   `json:"remember_me"`
}

type LoginResponse struct {
	GolferUser GolferUser `json:"golfer_user"`
}

type Golfers struct {
	GhinNumber                 string `json:"ghin_number"`
	Suffix                     string `json:"suffix"`
	FirstName                  string `json:"first_name"`
	MiddleName                 string `json:"middle_name"`
	LastName                   string `json:"last_name"`
	Prefix                     string `json:"prefix"`
	PlayerName                 string `json:"player_name"`
	Gender                     string `json:"gender"`
	ClubName                   string `json:"club_name"`
	ClubID                     string `json:"club_id"`
	GolfAssociationName        string `json:"golf_association_name"`
	GolfAssociationID          string `json:"golf_association_id"`
	Display                    string `json:"display"`
	DateOfBirth                string `json:"date_of_birth"`
	LowHiDisplay               string `json:"low_hi_display"`
	Email                      string `json:"email"`
	PrimaryClubCountry         string `json:"primary_club_country"`
	PrimaryClubState           string `json:"primary_club_state"`
	PrimaryClubName            string `json:"primary_club_name"`
	PrimaryClubID              int    `json:"primary_club_id"`
	PrimaryGolfAssociationID   int    `json:"primary_golf_association_id"`
	PrimaryGolfAssociationName string `json:"primary_golf_association_name"`
	RevDate                    string `json:"rev_date"`
	Status                     string `json:"status"`
	TechnologyProvider         string `json:"technology_provider"`
	CreatedAt                  string `json:"created_at"`
	SoftCap                    string `json:"soft_cap"`
	HardCap                    string `json:"hard_cap"`
	MessageClubAuthorized      any    `json:"message_club_authorized"`
}
type Subscription struct {
	Active                         bool   `json:"active"`
	SubscriptionAppType            string `json:"subscription_app_type"`
	SubscriptionType               string `json:"subscription_type"`
	InitialSubscriptionDate        string `json:"initial_subscription_date"`
	CurrentSubscriptionStartDate   string `json:"current_subscription_start_date"`
	CurrentSubscriptionEndDate     string `json:"current_subscription_end_date"`
	CurrentSubscriptionRenewalType string `json:"current_subscription_renewal_type"`
}
type FreeTrial struct {
	One8HoleComplete      bool   `json:"18hole_complete"`
	One8HoleDateCompleted string `json:"18hole_date_completed"`
	NineHoleComplete      bool   `json:"9hole_complete"`
	NineHoleDateCompleted string `json:"9hole_date_completed"`
}
type GolferUser struct {
	GolferUserToken         string       `json:"golfer_user_token"`
	GolferID                int          `json:"golfer_id"`
	IsTrial                 bool         `json:"is_trial"`
	GuardianID              any          `json:"guardian_id"`
	GolferUserAcceptedTerms bool         `json:"golfer_user_accepted_terms"`
	GolferCreationDate      time.Time    `json:"golfer_creation_date"`
	Golfers                 []Golfers    `json:"golfers"`
	MinorAccounts           []any        `json:"minor_accounts"`
	Subscription            Subscription `json:"subscription"`
	FreeTrial               FreeTrial    `json:"free_trial"`
}

// Login authenticates against the GHIN service
func (c *Client) Login(ghinNumber string, password string) (LoginResponse, error) {
	var loginResponse LoginResponse

	loginRequest := LoginRequest{
		User: User{
			Password:    password,
			EmailOrGhin: ghinNumber,
			RememberMe:  false,
		},
		Token: "jsGlqh9tkXvkLzDIM7iTnUH3tcB1oG18C8lOSqVumvRaWDCXxWjBFzdKrpsAzyT4o0GJkagmaop5AdoKj2+jP1km5Nd2B712/psGAoOMITJojq1l6aI3J6N18uh1icWj1vFEttute048dW8RveJIkdJJ5Or2dKr7yRmI3LkQ68DZ4gvuLL6s/JqnqAD2jgfQCqWYFJfUVEnrfvZylIfpWs4mSmQ1Tq+3TJ2MWyIU4pOWmeRZ/9pO2KtDTj3IIGov8d+6f9MKVuZELcGUGPBJi70rHr8S8OTIurl+yEqF0Ie7gGbF5FSfPBo70T5A5/L1dvlTzHEAH2DR+CueYhCQ+MmUX1OvHSozKiWqetZ0nr1V53iLC/dOcjveYTktih3PtDcTn6cZtoIdeiNxsC5VW4XaXNyS9YMRgjZ/EIXSku79JUvQL00WL2QEHdt34uiEk/VSr3m0VFwe8LGoWgSiEEQZpK+/ml+67b3JXynoE9uO1mlhOGKpysUd66OKs64rN+HRJyIVgySampR/svUPidm9bQJ1rHKGsoppJY3om93/AxH1mdP8GN+N4krac3PoSOc3npvUuBAleMsiD/FahrFFUrh4fEnCVnsHSc3Gx6xv0Nj0fjzMaC4x5tvWhBz9m1arVGO3Jb8G8qFOKGvUA3KF3dWIl/CoQh2tf0pKtWs=",
	}

	req, err := c.newRequest(http.MethodPost, "https://api2.ghin.com/api/v1/golfer_login.json", loginRequest)
	if err != nil {
		log.Println("ghin client - failed to create new http request for Login()")
		return loginResponse, err
	}

	res, err := c.do(req, &loginResponse)
	if err != nil {
		log.Println("ghin client - failed to perform http request in Login()")
		return loginResponse, err
	}

	if res.StatusCode != http.StatusOK {
		return loginResponse, fmt.Errorf("error: StatusCode=%d: %s", res.StatusCode, res.Status)
	}

	// Save the cookies
	c.credentials.Cookies = res.Cookies()

	return loginResponse, nil
}
