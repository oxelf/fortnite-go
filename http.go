package fortnitego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
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

type FileLocationInfo struct {
	Files map[string]struct {
		ReadLink     string `json:"readLink"`
		WriteLink    string `json:"writeLink"`
		Hash         string `json:"hash"`
		LastModified string `json:"lastModified"`
		Size         int    `json:"size"`
		FileLocked   bool   `json:"fileLocked"`
	} `json:"files"`
}

type Meta struct {
	DisplayName          string           `json:"urn:epic:member:dn_s"`
	JoinRequestUsersJSON JoinRequestUsers `json:"urn:epic:member:joinrequestusers_j"`
}
type FileInfo struct {
	ReadLink     string `json:"readLink"`
	WriteLink    string `json:"writeLink"`
	Hash         string `json:"hash"`
	LastModified string `json:"lastModified"`
	Size         int64  `json:"size"`
	FileLocked   bool   `json:"fileLocked"`
}

type FileData struct {
	Files map[string]FileInfo `json:"files"`
}

type AccessResponse struct {
	Play     bool `json:"play"`
	IsBanned bool `json:"isBanned"`
}

type S3Data struct {
	Data FileData `json:"data"`
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

type JoinRequestResponse struct {
	Status  string `json:"status"`
	PartyID string `json:"party_id"`
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
// NOTE: Is only for internal usage.
func (c *Client) doNullableRequest(method string, url string, payload []byte, urlEncoded bool) *Error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &Error{ErrorMessage: "Error creating Request."}
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.Config.Token))
	if urlEncoded {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	res, err := c.c.Do(req)
	if res.StatusCode == 204 || res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return &Error{ErrorMessage: "IO Read Error."}
		}
		fmt.Sprintf(string(body))
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
// NOTE: Is only for internal usage.
func (c *Client) doRequest(method string, url string, payload []byte, urlEncoded bool, out interface{}) (*Error, interface{}) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return &Error{ErrorMessage: "Error creating Request."}, nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.Config.Token))
	if urlEncoded {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	res, err := c.c.Do(req)
	if res.StatusCode == 204 || res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		byteString := string(body)
		fmt.Println(byteString)
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
//#       	   Files     	        #
//###################################

// Get Information on a specific file. Example:
//
//	sessionId := "8c9e7d3608914b0897dd7fa76406eda9"
//	fileLocation, _ := client.FileLocationInfo(fmt.Sprint("https://datastorage-public-service-live.ol.epicgames.com/api/v1/access/fnreplaysmetadata/public%2F" + sessionId + ".json"))
func (c *Client) FileLocationInfo(url string) (*FileLocationInfo, *Error) {
	payload := []byte{}
	var fileLocationInfo FileLocationInfo
	err, _ := c.doRequest("GET", url, payload, false, &fileLocationInfo)
	if err != nil {
		return nil, err
	} else {
		return &fileLocationInfo, nil
	}
}

//###################################
//#       LightSwitchService        #
//###################################

// Get the Status of the Fortnite servers.
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
//#             Party               #
//###################################

// Sends intention to a specific party. I dont really know what its used for.
func (c *Client) PartySendIntention(userId string) *Error {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/members/%s/intentions/%s", BaseRoute.PartyPublicService, userId, c.Config.AccountID)
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

// Get Information on a party using its partyId.
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

// TODO: Struct for this.
func (c *Client) PartyLookupPing(userId string, clientId string) *PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/user/%s/pings/%s/parties", BaseRoute.PartyPublicService, clientId, userId)
	payload, err := json.Marshal(&IntentionPayload{Urn: ""})
	if err != nil {

	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if err != nil {

	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.Config.Token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.c.Do(req)
	if requestError != nil {

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

	}

	var partyLookupResponses []PartyLookupResponse
	err = json.Unmarshal([]byte(body), &partyLookupResponses)
	if err != nil {

	}
	if len(partyLookupResponses) == 0 {
		return nil
	} else {
		return &partyLookupResponses[0]
	}
}

// TODO: Struct or error
func (c *Client) PartySendInvite(userId string) *Error {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/user/%s", BaseRoute.PartyPublicService, userId)
	payload := map[string]string{
		"urn:epic:cfg:build-id_s":        "1:3:",
		"urn:epic:conn:platform_s":       "WIN",
		"urn:epic:conn:type_s":           "game",
		"urn:epic:invite:platformdata_s": "",
		"urn:epic:member:dn_s":           c.Config.DisplayName,
	}
	payloadbytes, _ := json.Marshal(payload)
	ReqError := c.doNullableRequest("GET", url, payloadbytes, false)
	if ReqError != nil {
		return ReqError
	} else {
		return nil
	}
}

func (c *Client) PartySendJoinRequest(jid string, partyId string) (*JoinRequestResponse, *Error) {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/join", BaseRoute.PartyPublicService, partyId, c.Config.AccountID)
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
			"urn:epic:member:dn_s":               c.Config.DisplayName,
			"urn:epic:member:joinrequestusers_j": fmt.Sprintf("{\"users\":[{\"id\":\"%s\", \"dn\":\"%s\",\"plat\":\"WIN\",\"data\":{\"CrossplayPreference\": \"1\", \"SubGame_u\": \"1\"}}]}", c.Config.AccountID, c.Config.DisplayName),
		},
	}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {

	}
	var response JoinRequestResponse
	ReqError, _ := c.doRequest("POST", url, payloadbytes, false, &response)
	if ReqError != nil {
		return nil, ReqError
	} else {
		return &response, nil
	}
}

//	payload = {
//		'delete': deleted_meta,
//		'update': updated_meta,
//		'override': overridden_meta,
//		'revision': revision,
//	}
func (c *Client) PartyUpdateMemberMeta(updatedMeta map[string]interface{}) *Error {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, c.Party.Id, c.Config.AccountID)
	payload := map[string]interface{}{
		"delete":   []string{},
		"revision": c.Party.PartyRevision,
		"override": map[string]interface{}{},
		"update":   updatedMeta,
	}
	payloadbytes, _ := json.Marshal(payload)
	err := c.doNullableRequest("PATCH", url, payloadbytes, false)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Client) PartyUpdateMeta(updatedMeta map[string]interface{}) *Error {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, c.Party.Id, c.Config.AccountID)
	payload := map[string]interface{}{
		"delete":   []string{},
		"revision": c.Party.PartyRevision,
		"override": map[string]interface{}{},
		"update":   updatedMeta,
	}
	payloadbytes, _ := json.Marshal(payload)
	err := c.doNullableRequest("PATCH", url, payloadbytes, false)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Client) SetEmote(eID string) *Error {
	cerr := c.CancelEmote()
	if cerr != nil {
		return cerr
	}
	payload := map[string]interface{}{
		"Default:FrontendEmote_j": fmt.Sprintf("{\"FrontendEmote\":{\"emoteItemDef\":\"/Game/Athena/Items/Cosmetics/Dances/%s.%s\",\"emoteSection\":-1}}", eID, eID),
	}
	err := c.PartyUpdateMemberMeta(payload)
	if err != nil {
		if err.EpicErrorCode == "errors.com.epicgames.social.party.stale_revision" {
			newRevision, _ := strconv.Atoi(err.EpicMessageVars[1])
			c.Party.PartyRevision = newRevision
			c.SetEmote(eID)
		}
		return err
	} else {
		return nil
	}
}

func (c *Client) CancelEmote() *Error {
	payload := map[string]interface{}{
		"Default:FrontendEmote_j": "{\"FrontendEmote\":{\"emoteItemDef\":\"None\",\"emoteEKey\":\"\",\"emoteSection\":-1}}",
	}
	err := c.PartyUpdateMemberMeta(payload)
	if err != nil {
		if err.EpicErrorCode == "errors.com.epicgames.social.party.stale_revision" {
			newRevision, _ := strconv.Atoi(err.EpicMessageVars[1])
			c.Party.PartyRevision = newRevision
			c.CancelEmote()
		}
		return err
	} else {
		return nil
	}
}

func (c *Client) SetCustomKey(newKey string) *Error {
	payload := map[string]interface{}{
		"Default:CustomMatchKey_s": newKey,
	}
	err := c.PartyUpdateMemberMeta(payload)
	if err != nil {
		if err.EpicErrorCode == "errors.com.epicgames.social.party.stale_revision" {
			newRevision, _ := strconv.Atoi(err.EpicMessageVars[1])
			c.Party.PartyRevision = newRevision
			c.SetCustomKey(newKey)
		}
		return err
	} else {
		return nil
	}
}
func (c *Client) SetReadiness(ready bool) *Error {
	var readyString = ""
	if ready {
		readyString = "Ready"
	} else {
		readyString = "NotReady"
	}
	payload := map[string]interface{}{
		"Default:LobbyState_j": fmt.Sprintf("{\"LobbyState\":{\"inGameReadyCheckStatus\":\"None\",\"gameReadiness\":\"SittingOut\",\"readyInputType\":\"Touch\",\"currentInputType\":\"Touch\",\"hiddenMatchmakingDelayMax\":0,\"hasPreloadedAthena\":true}}", readyString),
	}
	err := c.PartyUpdateMemberMeta(payload)
	if err != nil {
		if err.EpicErrorCode == "errors.com.epicgames.social.party.stale_revision" {
			newRevision, _ := strconv.Atoi(err.EpicMessageVars[1])
			c.Party.PartyRevision = newRevision
			c.SetReadiness(ready)
		}
		return err
	} else {
		return nil
	}
}

// "Default:ActivityName_s":"","Default:PlaylistData_j":"{\"PlaylistData\":{\"playlistName\":\"Playlist_ShowdownAlt_Duos\",\"tournamentId\":\"epicgames_Arena_S24_Duos\",\"eventWindowId\":\"Arena_S24_Division6_Duos\",\"regionId\":\"EU\",\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"bGracefullyUpgraded\":false,\"matchmakingRulePreset\":\"RespectParties\"}}","Default:MatchmakingState_j":"{\"MatchmakingState\":{\"currentMatchmakingState\":{\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"requester\":\"INVALID\",\"dataStateId\":12199},\"requestedMatchmakingState\":{\"linkId\":{\"mnemonic\":\"playlist_defaultduo\",\"version\":-1},\"requester\":\"INVALID\",\"dataStateId\":12199},\"coordinatorBroadcast\":\"ReadyForRequests\"}}","Default:ActivityType_s":"BR"
func (c *Client) SetPlaylist() *Error {

	payload := map[string]interface{}{
		"Default:PlaylistData_j": "{\"PlaylistData\":{\"playlistName\":\"Playlist_ShowdownAlt_Duos\",\"tournamentId\":\"epicgames_Arena_S24_Duos\",\"eventWindowId\":\"Arena_S24_Division6_Duos\",\"regionId\":\"EU\",\"linkId\":{\"mnemonic\":\"tournament_epicgames_arena_duos\",\"version\":1},\"bGracefullyUpgraded\":false,\"matchmakingRulePreset\":\"RespectParties\"}}",
	}
	err := c.PartyUpdateMemberMeta(payload)
	if err != nil {
		if err.EpicErrorCode == "errors.com.epicgames.social.party.stale_revision" {
			newRevision, _ := strconv.Atoi(err.EpicMessageVars[1])
			c.Party.PartyRevision = newRevision
			c.SetPlaylist()
		}
		return err
	} else {
		return nil
	}
}

func (c *Client) PartyLeave() PartyLookupResponse {
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s", BaseRoute.PartyPublicService, c.Party.Id, c.Config.AccountID)
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

	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payloadbytes))
	if err != nil {

	}

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.Config.Token))
	req.Header.Add("Content-Type", "application/json")
	resp, requestError := c.c.Do(req)
	if requestError != nil {
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

	}

	response := &PartyLookupResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {

	}

	return *response
}
func (c *Client) Set_Skin(SkinID string) *Error {
	characterString := fmt.Sprintf("/Game/Athena/Items/Cosmetics/Characters/%s.%s", SkinID, SkinID)
	url := fmt.Sprintf("%s/party/api/v1/Fortnite/parties/%s/members/%s/meta", BaseRoute.PartyPublicService, c.Party.Id, c.Config.AccountID)
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
