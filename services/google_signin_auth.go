package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rudyjcruz831/mathSheets/model"
	"github.com/rudyjcruz831/mathSheets/util/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func auth(code string) (*model.Response, *model.Users, *errors.MathSheetsError) {
	tok, tradeCVDErr := getTokenFromAuthCode(code)
	if tradeCVDErr != nil {
		return nil, nil, tradeCVDErr
	}

	tokenID, tradeCVDErr := getUserInfoFromToken(tok)
	if tradeCVDErr != nil {
		return nil, nil, tradeCVDErr
	}

	var tokID model.TokenID
	if err := json.Unmarshal(tokenID, &tokID); err != nil {
		tradeCVDErr := errors.NewInternalServerError("not able to parse token")
		return nil, nil, tradeCVDErr
	}

	u := model.Users{
		// ID:        tokID.ID,
		Email:     tokID.Email,
		FirstName: tokID.FirstName,
		LastName:  tokID.LastName,
	}

	response := model.Response{
		AccessToken: tok.AccessToken,
		ID:          tokID.ID,
		FirstName:   tokID.FirstName,
		LastName:    tokID.LastName,
		Email:       tokID.Email,
		TokenType:   tok.TokenType,
	}

	// fmt.Println(u)
	// fmt.Println(response)
	fmt.Println(response.TokenType)

	return &response, &u, nil
}

func getConfig() *oauth2.Config {
	// Get the Client Id and Client secret stored in enviroment variables
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	// fmt.Println(clientID, clientSecret)

	// Build auth configuration instance
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"profile", "email", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
	// Don't know why but we need this to make it work
	conf.RedirectURL = "postmessage"

	return conf
}

func getTokenFromAuthCode(authCode string) (*oauth2.Token, *errors.MathSheetsError) {
	conf := getConfig()

	// Exchange consumable authorization code for refresh token
	tok, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		tradeCVDErr := errors.UnauthorizedError("getting token form authCode ---" + err.Error())
		return nil, tradeCVDErr
	}

	return tok, nil
}

func getUserInfoFromToken(token *oauth2.Token) ([]byte, *errors.MathSheetsError) {
	conf := getConfig()

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		tradeCVDErr := errors.UnauthorizedError("getting user from token ---" + err.Error())
		return nil, tradeCVDErr
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tradeCVDErr := errors.NewInternalServerError("Something went wrong on our end")
		return nil, tradeCVDErr
	}

	return data, nil
}
