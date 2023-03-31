package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type authclients struct {
	Fortnite_IOS_Client    string
	Fortnite_PC_Client     string
	Fortnite_SWITCH_Client string
}
type base64AuthClients struct {
	Fortnite_IOS_Client    string
	Fortnite_PC_Client     string
	Fortnite_SWITCH_Client string
}
type OauthToken struct {
	AccessToken      string    `json:"access_token"`
	ExpiresIn        int       `json:"expires_in"`
	ExpiresAt        time.Time `json:"expires_at"`
	TokenType        string    `json:"token_type"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresIn int       `json:"refresh_expires"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	AccountID        string    `json:"account_id"`
	ClientID         string    `json:"client_id"`
	ClientService    string    `json:"client_service"`
	DisplayName      string    `json:"displayName"`
	ProductID        string    `json:"product_id"`
}

var Base64AuthClients = base64AuthClients{
	Fortnite_PC_Client:     "ZWM2ODRiOGM2ODdmNDc5ZmFkZWEzY2IyYWQ4M2Y1YzY6ZTFmMzFjMjExZjI4NDEzMTg2MjYyZDM3YTEzZmM4NGQ=",
	Fortnite_IOS_Client:    "3446cd72694c4a4485d81b77adbb2141:9209d4a5e25a457fb9b07489d313b41a",
	Fortnite_SWITCH_Client: "OThmN2U0MmMyZTNhNGY4NmE3NGViNDNmYmI0MWVkMzk6MGEyNDQ5YTItMDAxYS00NTFlLWFmZWMtM2U4MTI5MDFjNGQ3",
}

var AuthClients = authclients{
	Fortnite_IOS_Client:    "3446cd72694c4a4485d81b77adbb2141",
	Fortnite_PC_Client:     "ec684b8c687f479fadea3cb2ad83f5c6",
	Fortnite_SWITCH_Client: "98f7e42c2e3a4f86a74eb43fbb41ed39",
}

// Gets an Oauth Token for the given parameters and the Authcode, which you can get from this url: https://www.epicgames.com/id/api/redirect?clientId=ec684b8c687f479fadea3cb2ad83f5c6&responseType=code
// Make sure to replace the clientId param with your client.
// Function Parameters:
// code: the auth code you got from the url
// base64Client: must be for the same client that u used for getting the auth code.
// eg1: if true, you will get an eg1 token back. Some endpoints require this type of token. Its smart to use it.
func Get_OauthToken_By_AuthCode(code string, base64Client string, eg1 bool) (*OauthToken, error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	if eg1 {
		data.Set("token_type", "eg1")
	}
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	if base64Client != "" {
		req.Header.Set("Authorization", fmt.Sprintf("basic %s", base64Client))
	} else {
		return nil, fmt.Errorf("please provide an authclient base64 string. You can use the pre defined ones like this: Base64AuthClients.Fortnite_IOS_Client")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	oauthToken := &OauthToken{}

	err = json.Unmarshal(body, &oauthToken)
	if err != nil {
		log.Fatalln(err)
	}

	return oauthToken, nil
}