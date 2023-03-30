package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type BaseRoutes struct {
	FriendsPublicService      string
	PartyPublicService        string
	LightSwitchPublicService  string
	UserSearchService         string
	AccountPublicService      string
	EulaTrackingPublicService string
	EventsPublicService       string
	FortniteContentWebsite    string
	FortnitePublicService     string
	PresencePublicService     string
	AvatarService             string
	StatsProxyPublicService   string
}
type LightSwitchResponse struct {
	Service string `json:"serviceInstanceId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type SendIntentionData struct {
	UserId   string
	ClientId string
}
type IntentionPayload struct {
	Urn string `json:"urn:epic:invite:platformdata_s"`
}
type PartyLookupResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Config    struct {
		Type             string `json:"type"`
		Joinability      string `json:"joinability"`
		Discoverability  string `json:"discoverability"`
		SubType          string `json:"sub_type"`
		MaxSize          int    `json:"max_size"`
		InviteTTL        int    `json:"invite_ttl"`
		JoinConfirmation bool   `json:"join_confirmation"`
		IntentionTTL     int    `json:"intention_ttl"`
	} `json:"config"`
	Members []struct {
		AccountID   string      `json:"account_id"`
		Meta        interface{} `json:"meta"`
		Connections []struct{}  `json:"connections"`
		Revision    int         `json:"revision"`
		UpdatedAt   string      `json:"updated_at"`
		JoinedAt    string      `json:"joined_at"`
		Role        string      `json:"role"`
	} `json:"members"`
	Applicants []interface{} `json:"applicants"`
	Meta       interface{}   `json:"meta"`
}
type ConnectionMeta struct {
	Platform string `json:"urn:epic:conn:platform_s"`
	Type     string `json:"urn:epic:conn:type_s"`
}

type User struct {
	ID             string `json:"id"`
	DisplayName    string `json:"dn"`
	Platform       string `json:"plat"`
	CrossplayPrefs string `json:"CrossplayPreference"`
	SubGame_u      string `json:" SubGame_u"`
}

type Meta struct {
	DisplayName          string           `json:"urn:epic:member:dn_s"`
	JoinRequestUsersJSON JoinRequestUsers `json:"urn:epic:member:joinrequestusers_j"`
}

type Connection struct {
	ID              string         `json:"id"`
	Meta            ConnectionMeta `json:"meta"`
	YieldLeadership bool           `json:"yield_leadership"`
	OfflineTTL      int            `json:"offline_ttl"`
}

type JoinRequestPayload struct {
	Connection Connection `json:"connection"`
	Meta       Meta       `json:"meta"`
}

var BaseRoute = BaseRoutes{
	FriendsPublicService:      "https://friends-public-service-prod.ol.epicgames.com",
	PartyPublicService:        "https://party-service-prod.ol.epicgames.com",
	LightSwitchPublicService:  "https://lightswitch-public-service-prod06.ol.epicgames.com",
	UserSearchService:         "https://user-search-service-prod.ol.epicgames.com",
	AccountPublicService:      "https://account-public-service-prod.ol.epicgames.com",
	EulaTrackingPublicService: "https://eulatracking-public-service-prod-m.ol.epicgames.com",
	EventsPublicService:       "https://events-public-service-live.ol.epicgames.com",
	FortniteContentWebsite:    "https://fortnitecontent-website-prod07.ol.epicgames.com",
	FortnitePublicService:     "https://fortnite-public-service-prod11.ol.epicgames.com",
	PresencePublicService:     "https://presence-public-service-prod.ol.epicgames.com",
	AvatarService:             "https://avatar-service-prod.identity.live.on.epicgames.com",
	StatsProxyPublicService:   "https://statsproxy-public-service-live.ol.epicgames.com",
}

//###################################
//#         Request Handling        #
//###################################

// used for operations that doesnt send a neccessary payload back.
func (c *Client) doNullableRequest(method string, url string, payload []byte, urlEncoded bool) *Error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &Error{ErrorMessage: "Error creating Request."}
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	if urlEncoded {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	res, err := c.c.Do(req)
	if res.StatusCode == 204 || res.StatusCode == 200 {
		c.Party.PartyRevision = c.Party.PartyRevision + 1
		return nil
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return &Error{ErrorMessage: "IO Read Error."}
		}
		EpicError := &Error{}
		err = json.Unmarshal(body, EpicError)
		if err != nil {
			return &Error{
				ErrorMessage: "Error unmarshaling.",
			}
		}
		if res.StatusCode == 409 {
			newRev, convErr := strconv.Atoi(EpicError.EpicMessageVars[1])
			if convErr != nil {
				return &Error{ErrorMessage: "atoi convert error."}
			}
			c.Party.PartyRevision = newRev
			return EpicError
		}
		return EpicError
	}
}

// used for operations that send a neccessary payload back.
func (c *Client) doRequest(method string, url string, payload []byte, urlEncoded bool, out interface{}) (*Error, interface{}) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &Error{ErrorMessage: "Error creating Request."}, nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	if urlEncoded {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	res, err := c.c.Do(req)
	if res.StatusCode == 204 || res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		fmt.Println(body)
		if err != nil {
			return &Error{ErrorMessage: "IO Read Error."}, nil
		}
		if err := json.Unmarshal(body, out); err != nil {
			return nil, Error{ErrorMessage: err.Error()}
		}
		return nil, out
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return &Error{ErrorMessage: "IO Read Error."}, nil
		}
		EpicError := &Error{}
		err = json.Unmarshal(body, EpicError)
		if err != nil {
			return &Error{
				ErrorMessage: "Error unmarshaling.",
			}, nil
		}
		if res.StatusCode == 409 {
			newRev, convErr := strconv.Atoi(EpicError.EpicMessageVars[1])
			if convErr != nil {
				return &Error{ErrorMessage: "atoi convert error."}, nil
			}
			c.Party.PartyRevision = newRev
			return EpicError, nil
		}
		return EpicError, nil
	}
}

//###################################
//#       LightSwitchService        #
//###################################

func (c *Client) LightSwitch_Status_Fortnite() (*LightSwitchResponse, *Error) {
	url := "http://lightswitch-public-service-prod.ol.epicgames.com/lightswitch/api/service/fortnite/status"
	payload := []byte{}
	var lightSwitchResponse LightSwitchResponse
	err, _ := c.doRequest("GET", url, payload, false, &lightSwitchResponse)
	if err != nil {
		return nil, err
	} else {
		return &lightSwitchResponse, nil
	}
}



//###################################
//#            Friends              #
//###################################

func (c *Client) Friends_Add_Or_Accept(friendId string) *Error {
	url := fmt.Sprintf("%s/friends/api/v1/%s/friends/%s", BaseRoute.FriendsPublicService, accountId, friendId)
	payload := []byte{}
	requestError := c.doNullableRequest("POST", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) Friends_Remove_Or_Decline(friendId string) *Error {
	url := fmt.Sprintf("%s/friends/api/v1/%s/friends/%s", BaseRoute.FriendsPublicService, accountId, friendId)
	payload := []byte{}
	requestError := c.doNullableRequest("DELETE", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) Friends_Get_All() ([]Friend, *Error) {
	url := fmt.Sprintf("%s/friends/api/public/friends/%s", BaseRoute.FriendsPublicService, accountId)
	payload := []byte{}
	var friendList []Friend
	err, _ := c.doRequest("GET", url, payload, false, &friendList)
	if err != nil {
		return nil, err
	} else {
		return friendList, nil
	}
}

func (c *Client) Friends_Get_Blocklist() ([]Friend, *Error) {
	url := fmt.Sprintf("%s/friends/api/v1/%s/blocklist", BaseRoute.FriendsPublicService, accountId)
	payload := []byte{}
	var blockList []Friend
	err, _ := c.doRequest("GET", url, payload, false, &blockList)
	if err != nil {
		return nil, err
	} else {
		return *&blockList, nil
	}
}

// Currently not working!!! Problem is the body.
func (c *Client) Friends_Set_Nickname(nickName string, friendID string) *Error {
	uri := fmt.Sprintf("%s/friends/api/v1/%s/friends/%s/alias", BaseRoute.FriendsPublicService, accountId, friendID)
	encodedPayload := url.QueryEscape(nickName)
	bodyBytes := []byte(encodedPayload)
	body := url.Values{}
	body.Set("nick", nickName)
	//bodyString := body.Encode()
	//payload := []byte(bodyString)
	requestError := c.doNullableRequest("PUT", uri, bodyBytes, true)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

//###################################
//#             Party               #
//###################################

func (c *Client) PartySendIntention(userId string) *Error {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/members/%s/intentions/%s", BaseRoute.PartyPublicService, userId, accountId)
	payload, err := json.Marshal(&IntentionPayload{Urn: ""})
	if err != nil {
		return &Error{
			ErrorMessage: "Error marshaling payload.",
		}
	}
	requestError := c.doNullableRequest("POST", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) PartyLookup(partyID string) (*PartyLookupResponse, *Error) {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s", BaseRoute.PartyPublicService, partyID)
	payload, err := json.Marshal(&IntentionPayload{Urn: ""})
	if err != nil {
		return nil, &Error{
			ErrorMessage: "Error marshalling payload.",
		}
	}
	var partyLookup PartyLookupResponse
	requestError, res := c.doRequest("GET", url, payload, false, &partyLookup)
	if requestError != nil {
		if res == nil {
			return nil, requestError
		}
		return nil, requestError
	} else {
		return &partyLookup, nil
	}
}

func PartyLookupPing(c *http.Client, userId string, clientId string) *PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/user/%s/pings/%s/parties", BaseRoute.PartyPublicService, clientId, userId)
	payload, err := json.Marshal(&IntentionPayload{Urn: ""})
	if err != nil {
		fmt.Println("error marschaling")
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.Do(req)
	if requestError != nil {
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	var partyLookupResponses []PartyLookupResponse
	err = json.Unmarshal([]byte(body), &partyLookupResponses)
	if err != nil {
		fmt.Println(err)

	}
	if len(partyLookupResponses) == 0 {
		return nil
	} else {
		return &partyLookupResponses[0]
	}
}

func PartySendInvite(c *http.Client, userId string) PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/user/%s", BaseRoute.PartyPublicService, userId)
	payload := map[string]string{
		"urn:epic:cfg:build-id_s":        "1:3:24395311",
		"urn:epic:conn:platform_s":       "WIN",
		"urn:epic:conn:type_s":           "game",
		"urn:epic:invite:platformdata_s": "",
		"urn:epic:member:dn_s":           "oxelf ay",
	}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payloadbytes))
	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.Do(req)
	if requestError != nil {
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}

func PartySendJoinRequest(c *http.Client, jid string, partyId string) PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/join", BaseRoute.PartyPublicService, partyId, accountId)
	payload := map[string]interface{}{
		"connection": map[string]interface{}{
			"id": jid,
			"meta": map[string]string{
				"urn:epic:conn:platform_s": "WIN",
				"urn:epic:conn:type_s":     "game",
			},
			"yield_leadership": false,
			"offline_ttl":      30,
		},
		"meta": map[string]interface{}{
			"urn:epic:member:dn_s":               "oxibot 001",
			"urn:epic:member:joinrequestusers_j": fmt.Sprintf("{\"users\":[{\"id\":\"%s\", \"dn\":\"oxelf ay\",\"plat\":\"WIN\",\"data\":{\"CrossplayPreference\": \"1\", \"SubGame_u\": \"1\"}}]}", accountId),
		},
	}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadbytes))
	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.Do(req)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}

//	payload = {
//		'delete': deleted_meta,
//		'update': updated_meta,
//		'override': overridden_meta,
//		'revision': revision,
//	}
func (c *Client) PartyUpdateMemberMeta(p *PartyMemberMeta) PartyLookupResponse {
	emoteString := ""
	characterString := ""
	if p.Location == "" {
		p.Location = "PreLobby"
	}
	if p.Emote == "" {
		p.Emote = "None"
		emoteString = ""
	} else {
		emoteString = fmt.Sprintf("/Game/Athena/Items/Cosmetics/Dances/%s.%s", p.Emote, p.Emote)
	}
	if p.CosmeticLoadout.Character == "" {
		p.CosmeticLoadout.Character = "None"
		characterString = "None"
	} else {
		characterString = fmt.Sprintf("/Game/Athena/Items/Cosmetics/Characters/%s.%s", p.CosmeticLoadout.Character, p.CosmeticLoadout.Character)
	}
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, c.Party.Id, accountId)
	payload := map[string]interface{}{
		"delete":   []string{},
		"revision": 0,
		"override": map[string]interface{}{},
		"update": map[string]interface{}{
			"Default:AthenaCosmeticLoadout_j": fmt.Sprintf("{\"AthenaCosmeticLoadout\":{\"characterDef\":\"%s\",\"characterEKey\":\"\",\"backpackDef\":\"None\",\"backpackEKey\":\"\",\"pickaxeDef\":\"/Game/Athena/Items/Cosmetics/Pickaxes/Pickaxe_ID_015_HolidayCandyCane.Pickaxe_ID_015_HolidayCandyCane\",\"pickaxeEKey\":\"\",\"contrailDef\":\"/Game/Athena/Items/Cosmetics/Contrails/Contrail_RedPepper.Contrail_RedPepper\",\"contrailEKey\":\"\",\"scratchpad\":[],\"cosmeticStats\":[{\"statName\":\"TotalVictoryCrowns\",\"statValue\":6},{\"statName\":\"TotalRoyalRoyales\",\"statValue\":0},{\"statName\":\"HasCrown\",\"statValue\":0}]}}", characterString),
			"Default:FrontendEmote_j":         fmt.Sprintf("{\"FrontendEmote\":{\"emoteItemDef\":\"%s\",\"emoteEKey\":\"\",\"emoteSection\":-1}}", emoteString),
		},
	}

	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadbytes))

	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.c.Do(req)
	fmt.Printf("Update meta status code: %d\n", resp.StatusCode)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}

func SetEmote(c *http.Client, partyId string, eID string) PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, partyId, accountId)
	payload := map[string]interface{}{
		"delete":   []string{},
		"revision": 2,
		"override": map[string]interface{}{},
		"update": map[string]interface{}{
			"Default:FrontendEmote_j": fmt.Sprintf("{\"FrontendEmote\":{\"emoteItemDef\":\"/Game/Athena/Items/Cosmetics/Dances/%s.%s\",\"emoteEKey\":\"\",\"emoteSection\":1}}", eID, eID),
		},
	}

	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadbytes))

	if err != nil {
		fmt.Println("error")
	}

	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	request.Header.Add("Content-Type", "application/json")
	resp, requestError := c.Do(request)
	fmt.Printf("Update meta status code: %d\n", resp.StatusCode)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}

func (client *Client) SetCustomKey(newKey string) PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s", BaseRoute.PartyPublicService, accountId)
	payload := map[string]interface{}{
		"revision": client.Party.PartyRevision,
		"meta": map[string]interface{}{
			"delete": []string{},
			"update": map[string]interface{}{
				"Default:CustomMatchKey_s": newKey,
			}},
	}

	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadbytes))

	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := client.c.Do(req)
	fmt.Printf("Update meta status code: %d\n", resp.StatusCode)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}
	var epicError Error
	err = json.Unmarshal([]byte(body), &epicError)
	if err != nil {
		// handle error
	}
	if epicError.EpicErrorCode != "" {
		newRev, convErr := strconv.Atoi(epicError.EpicMessageVars[1])
		if convErr != nil {

		}
		client.Party.PartyRevision = newRev
		client.SetCustomKey(newKey)
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}
func SetReadiness(c *http.Client, partyId string, ready bool) PartyLookupResponse {
	var readyString = ""
	if ready {
		readyString = "Ready"
		fmt.Println(" readiness: ready")
	} else {
		readyString = "NotReady"
		fmt.Println(" readiness: not ready")
	}

	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, partyId, accountId)
	payload := map[string]interface{}{
		"delete":   []string{},
		"revision": 1,
		"override": map[string]interface{}{},
		"update": map[string]interface{}{
			"Default:LobbyState_j": fmt.Sprintf("{\"LobbyState\":{\"inGameReadyCheckStatus\":\"None\",\"gameReadiness\":\"%s\",\"readyInputType\":\"Touch\",\"currentInputType\":\"Touch\",\"hiddenMatchmakingDelayMax\":0,\"hasPreloadedAthena\":true}}", readyString),
		},
	}

	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadbytes))

	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.Do(req)
	fmt.Printf("Update meta status code: %d\n", resp.StatusCode)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}

// "Default:ActivityName_s":"","Default:PlaylistData_j":"{\"PlaylistData\":{\"playlistName\":\"Playlist_ShowdownAlt_Duos\",\"tournamentId\":\"epicgames_Arena_S24_Duos\",\"eventWindowId\":\"Arena_S24_Division6_Duos\",\"regionId\":\"EU\",\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"bGracefullyUpgraded\":false,\"matchmakingRulePreset\":\"RespectParties\"}}","Default:MatchmakingState_j":"{\"MatchmakingState\":{\"currentMatchmakingState\":{\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"requester\":\"INVALID\",\"dataStateId\":12199},\"requestedMatchmakingState\":{\"linkId\":{\"mnemonic\":\"playlist_defaultduo\",\"version\":-1},\"requester\":\"INVALID\",\"dataStateId\":12199},\"coordinatorBroadcast\":\"ReadyForRequests\"}}","Default:ActivityType_s":"BR"
func (c *Client) SetPlaylist() PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s", BaseRoute.PartyPublicService, c.Party.Id)
	payload := map[string]interface{}{
		"revision": c.Party.PartyRevision,
		"meta": map[string]interface{}{
			"delete": []string{},
			"update": map[string]interface{}{
				"Default:PlaylistData_j": "{\"PlaylistData\":{\"playlistName\":\"Playlist_ShowdownAlt_Duos\",\"tournamentId\":\"epicgames_Arena_S24_Duos\",\"eventWindowId\":\"Arena_S24_Division6_Duos\",\"regionId\":\"EU\",\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"bGracefullyUpgraded\":false,\"matchmakingRulePreset\":\"RespectParties\"}}", "Default:MatchmakingState_j": "{\"MatchmakingState\":{\"currentMatchmakingState\":{\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"requester\":\"INVALID\",\"dataStateId\":12199},\"requestedMatchmakingState\":{\"linkId\":{\"mnemonic\":\"playlist_defaultduo\",\"version\":-1},\"requester\":\"INVALID\",\"dataStateId\":12199},\"coordinatorBroadcast\":\"ReadyForRequests\"}}",
			}},
	}

	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadbytes))

	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.c.Do(req)
	fmt.Printf("Update meta status code: %d\n", resp.StatusCode)
	fmt.Printf("update meta reason: %v", resp)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}
	var epicError Error
	err = json.Unmarshal([]byte(body), &epicError)
	if err != nil {
		// handle error
	}
	if epicError.EpicErrorCode != "" {
		newRev, convErr := strconv.Atoi(epicError.EpicMessageVars[1])
		if convErr != nil {

		}
		c.Party.PartyRevision = newRev
		c.SetPlaylist()
	} else {

		response := &PartyLookupResponse{}

		err = json.Unmarshal(body, response)
		if err != nil {
			fmt.Println("error unmarshaling.")
		}

		return *response
	}
	return PartyLookupResponse{}
}
func (c *Client) PartyLeave() PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s", BaseRoute.PartyPublicService, c.Party.Id, accountId)
	payload := map[string]interface{}{
		"connection": map[string]interface{}{
			"id": c.Party.Members[0].JID,
			"meta": map[string]string{
				"urn:epic:conn:platform_s": "WIN",
				"urn:epic:conn:type_s":     "game",
			},
		},
		"meta": map[string]interface{}{
			"urn:epic:member:dn_s":          "oxelf ay",
			"urn:epic:member:type_s":        "game",
			"urn:epic:member:platform_s":    "WIN",
			"urn:epic:member:joinrequest_j": "{\"CrossplayPreference\": \"1\"}",
		},
	}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payloadbytes))
	if err != nil {
		fmt.Println("error")
	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.c.Do(req)
	if requestError != nil {
		fmt.Println("req error: " + requestError.Error())
		fmt.Println("request error.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error.")
	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("error unmarshaling.")
	}

	return *response
}
func (c *Client) Set_Skin(SkinID string) *Error {
	characterString := fmt.Sprintf("/Game/Athena/Items/Cosmetics/Characters/%s.%s", SkinID, SkinID)
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, c.Party.Id, accountId)
	payload := map[string]interface{}{
		"delete":   []string{},
		"revision": c.Party.PartyRevision,
		"override": map[string]interface{}{},
		"update": map[string]interface{}{
			"Default:AthenaCosmeticLoadout_j": fmt.Sprintf("{\"AthenaCosmeticLoadout\":{\"characterDef\":\"%s\",\"characterEKey\":\"\",\"backpackDef\":\"None\",\"backpackEKey\":\"\",\"pickaxeDef\":\"/Game/Athena/Items/Cosmetics/Pickaxes/Pickaxe_ID_015_HolidayCandyCane.Pickaxe_ID_015_HolidayCandyCane\",\"pickaxeEKey\":\"\",\"contrailDef\":\"/Game/Athena/Items/Cosmetics/Contrails/Contrail_RedPepper.Contrail_RedPepper\",\"contrailEKey\":\"\",\"scratchpad\":[],\"cosmeticStats\":[{\"statName\":\"TotalVictoryCrowns\",\"statValue\":6},{\"statName\":\"TotalRoyalRoyales\",\"statValue\":0},{\"statName\":\"HasCrown\",\"statValue\":0}]}}", characterString),
		},
	}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marschaling")
	}
	requestError := c.doNullableRequest("PATCH", url, payloadbytes, false)
	if requestError != nil {
		if requestError.EpicErrorCode == "errors.com.epicgames.social.party.stale_revision" {
			c.Set_Skin(SkinID)
		} else {
			return requestError
		}
	} else {
		return nil
	}
	return &Error{ErrorMessage: "Undefined Error."}
}
