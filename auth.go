package fortnitego

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
type DeviceCode struct {
	UserCode                 string `json:"user_code"`
	DeviceCode               string `json:"device_code"`
	VerificationUri          string `json:"verification_uri"`
	VerificationUriCompleted string `json:"verification_uri_complete"`
	Prompt                   string `json:"prompt"`
	ExpiresInSeconds         int    `json:"expires_in"`
	Interval                 int    `json:"interval"`
	ClientID                 string `json:"client_id"`
}
type DeviceAuth struct {
	DeviceID  string `json:"deviceId"`
	AccountID string `json:"accountId"`
	Secret    string `json:"secret"`
}
type ExchangeToken struct {
	ExpiresInSeconds int    `json:"expiresInSeconds"`
	Token            string `json:"code"`
	AuthClientId     string `json:"creatingClientId"`
}

var Base64AuthClients = base64AuthClients{
	Fortnite_PC_Client:     "ZWM2ODRiOGM2ODdmNDc5ZmFkZWEzY2IyYWQ4M2Y1YzY6ZTFmMzFjMjExZjI4NDEzMTg2MjYyZDM3YTEzZmM4NGQ=",
	Fortnite_IOS_Client:    "MzQ0NmNkNzI2OTRjNGE0NDg1ZDgxYjc3YWRiYjIxNDE6OTIwOWQ0YTVlMjVhNDU3ZmI5YjA3NDg5ZDMxM2I0MWE=",
	Fortnite_SWITCH_Client: "OThmN2U0MmMyZTNhNGY4NmE3NGViNDNmYmI0MWVkMzk6MGEyNDQ5YTItMDAxYS00NTFlLWFmZWMtM2U4MTI5MDFjNGQ3",
}

var AuthClients = authclients{
	Fortnite_IOS_Client:    "3446cd72694c4a4485d81b77adbb2141",
	Fortnite_PC_Client:     "ec684b8c687f479fadea3cb2ad83f5c6",
	Fortnite_SWITCH_Client: "98f7e42c2e3a4f86a74eb43fbb41ed39",
}

func GetAuthCodeUrl(authClient string) string {
	url := fmt.Sprintf("https://www.epicgames.com/id/api/redirect?clientId=%s&responseType=code", authClient)
	return url
}

// Gets an Oauth Token for the given parameters and the Authcode, which you can get from this url: https://www.epicgames.com/id/api/redirect?clientId=ec684b8c687f479fadea3cb2ad83f5c6&responseType=code
// Make sure to replace the clientId param with your client.
// Function Parameters:
// code: the auth code you got from the url
// base64Client: must be for the same client that u used for getting the auth code.
// eg1: if true, you will get an eg1 token back. Some endpoints require this type of token. Its smart to use it.
func Get_Token_By_AuthCode(code string, base64Client string, eg1 bool) (*OauthToken, *Error) {
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
		return nil, &Error{ErrorMessage: "please provide an authclient base64 string. You can use the pre defined ones like this: Base64AuthClients.Fortnite_IOS_Client"}
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
	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		oauthToken := &OauthToken{}
		err = json.Unmarshal(body, &oauthToken)
		if err != nil {
			log.Fatalln(err)
		}

		return oauthToken, nil
	}
}

func Create_DeviceAuth(token string, accountId string) (*DeviceAuth, *Error) {
	fmt.Println("getting device auth")
	uri := fmt.Sprint("https://account-public-service-prod.ol.epicgames.com/account/api/public/account/" + accountId + "/deviceAuth")
	data := url.Values{}

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
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
	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		deviceAuth := &DeviceAuth{}
		err = json.Unmarshal(body, &deviceAuth)
		if err != nil {
			log.Fatalln(err)
		}
		return deviceAuth, nil
	}
}

// base64Client: Needs to be the same client, thats used for the token to create the device auth.
//
// eg1: whetever or not the token should be type eg1.
//
// deviceAuth: A DeviceAuth object is required, as the accountID, the DeviceID and the Secret are required for requesting.
func Get_Token_By_DeviceAuth(deviceAuth DeviceAuth, base64Client string, eg1 bool) (*OauthToken, *Error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	if eg1 {
		data.Set("token_type", "eg1")
	}
	data.Set("grant_type", "device_auth")
	data.Set("device_id", deviceAuth.DeviceID)
	data.Set("secret", deviceAuth.Secret)
	data.Set("account_id", deviceAuth.AccountID)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	if base64Client == "" {
		return nil, &Error{
			ErrorMessage: "Please provide an base64 encoded client.",
		}
	}
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", base64Client))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		token := &OauthToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			log.Fatalln(err)
		}
		return token, nil
	}
}

func Get_ExchangeToken(token string) (*ExchangeToken, *Error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/exchange"
	data := url.Values{}
	req, err := http.NewRequest("GET", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	} else {
		return nil, &Error{ErrorMessage: "please provide an bearer token."}
	}
	req.Header.Add("Content-Type", "application/json")
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
	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		exchangeToken := &ExchangeToken{}
		err = json.Unmarshal(body, &exchangeToken)
		if err != nil {
			log.Fatalln(err)
		}

		return exchangeToken, nil
	}
}

// base64Client: Needs to be the same client, thats used for the token to create the device auth.
//
// eg1: whetever or not the token should be type eg1.
//
// exchangeToken: Your exchange token
func Get_Token_By_Exchange(exchangeToken string, base64Client string, eg1 bool) (*OauthToken, *Error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	if eg1 {
		data.Set("token_type", "eg1")
	}
	if base64Client == "" {
		return nil, &Error{
			ErrorMessage: "Please provide an base64 client string.",
		}
	}
	if exchangeToken == "" {
		return nil, &Error{
			ErrorMessage: "Please provide an exchangetoken",
		}
	}
	data.Set("grant_type", "exchange_code")
	data.Set("exchange_code", exchangeToken)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", base64Client))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		token := &OauthToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			log.Fatalln(err)
		}
		return token, nil
	}
}

// eg1: whetever or not the token should be type eg1.
//
// base64client: The client you want to get an token for.
func Get_Token_By_ClientCredentials(base64Client string, eg1 bool) (*OauthToken, *Error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	if eg1 {
		data.Set("token_type", "eg1")
	}
	if base64Client == "" {
		return nil, &Error{
			ErrorMessage: "Please provide a base64 client.",
		}
	}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", base64Client))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		token := &OauthToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			log.Fatalln(err)
		}
		return token, nil
	}
}

// eg1: whetever or not the token should be type eg1.
//
// refreshtoken: the refresh token you want to exchange for an accesstoken.
//
// base64client: The client you want to get an token for.
func Get_Token_By_Refresh(refreshToken string, base64Client string, eg1 bool) (*OauthToken, *Error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	if eg1 {
		data.Set("token_type", "eg1")
	}
	if base64Client == "" {
		return nil, &Error{
			ErrorMessage: "Please provide a base64 client.",
		}
	}
	if refreshToken == "" {
		return nil, &Error{
			ErrorMessage: "Please provide a refresh token.",
		}
	}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", base64Client))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		token := &OauthToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			log.Fatalln(err)
		}
		return token, nil
	}
}

// eg1: whetever or not the token should be type eg1.
//
// refreshtoken: the deviceCode you want to exchange for an accesstoken.
//
// base64client: The same that was used for creating the devicecode.
func Get_Token_By_DeviceCode(deviceCode string, base64Client string, eg1 bool) (*OauthToken, *Error) {
	uri := "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token"
	data := url.Values{}
	if eg1 {
		data.Set("token_type", "eg1")
	}
	if base64Client == "" {
		return nil, &Error{
			ErrorMessage: "Please provide a base64 client.",
		}
	}
	if deviceCode == "" {
		return nil, &Error{
			ErrorMessage: "Please provide a deviceCode.",
		}
	}
	data.Set("grant_type", "device_code")
	data.Set("device_code", deviceCode)
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", base64Client))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		token := &OauthToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			log.Fatalln(err)
		}
		return token, nil
	}
}
func Get_DeviceCode(token string) (*DeviceCode, *Error) {
	uri := "https://account-public-service-prod03.ol.epicgames.com/account/api/oauth/deviceAuthorization"
	data := url.Values{}

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
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

	if resp.StatusCode >= 400 {
		newerr := &Error{}
		err = json.Unmarshal(body, &newerr)
		if err != nil {
			log.Fatalln(err)
		}
		return nil, newerr
	} else {
		deviceCode := &DeviceCode{}
		err = json.Unmarshal(body, &deviceCode)
		if err != nil {
			log.Fatalln(err)
		}
		return deviceCode, nil
	}
}
