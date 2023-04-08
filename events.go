package fortnitego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const EventsBaseUrl = "https://events-public-service-live.ol.epicgames.com"

type TokenResponse struct {
	Accounts []TokenAccount
}

type TokenAccount struct {
	AccountID string   `json:"accountId"`
	Tokens    []string `json:"tokens"`
}

type EventHistorySession struct {
	ScoreKey struct {
		GameID        string `json:"gameId"`
		EventID       string `json:"eventId"`
		EventWindowId string `json:"eventWindowId"`
	} `json:"scoreKey"`
	TeamID         string       `json:"teamId"`
	TeamAccountIds []string     `json:"teamAccountIds"`
	LiveSessionID  string       `json:"liveSessionId"`
	PointsEarned   int          `json:"pointsEarned"`
	Rank           int          `json:"rank"`
	Percentile     float64      `json:"percentile"`
	SessionHistory []EventMatch `json:"sessionHistory"`
}

type EventMatch struct {
	SessionID    string    `json:"sessionId"`
	EndTime      time.Time `json:"endTime"`
	TrackedStats struct {
		PlacementStatIndex      int `json:"PLACEMENT_STAT_INDEX"`
		TimeAliveStat           int `json:"TIME_ALIVE_STAT"`
		TeamElimsStatIndex      int `json:"TEAM_ELIMS_STAT_INDEX"`
		MatchPlayedStat         int `json:"MATCH_PLAYED_STAT"`
		PlacementTiebreakerStat int `json:"PLACEMENT_TIEBREAKER_STAT"`
		VictoryRoyaleStat       int `json:"VICTORY_ROYALE_STAT"`
	}
}

func EventPlayerTokens(accountID string) (*TokenResponse, *Error) {
	url := fmt.Sprintf("%s/api/v1/players/Fortnite/tokens?teamAccountIds=%s", EventsBaseUrl, accountID)
	c := http.Client{}
	payload := []byte{}
	req, nerr := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if nerr != nil {

	}
	resp, err := c.Do(req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{ErrorMessage: "IO Read Error."}
	}
	if err != nil {
		return nil, &Error{
			ErrorMessage: nerr.Error(),
		}
	} else {
		var response TokenResponse
		rerr := json.Unmarshal(body, &response)
		if rerr != nil {
			return nil, &Error{ErrorMessage: "JSON Read Error"}
		}
		return &response, nil
	}
}

func (c *Client) EventHistory(EventID string) ([]EventHistorySession, *Error) {
	url := fmt.Sprintf("%s/api/v1/events/Fortnite/%s/history/%s", EventsBaseUrl, EventID, c.Config.AccountID)
	payload := []byte{}
	var response []EventHistorySession
	err, _ := c.doRequest("GET", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (c *Client) EventWindowHistory(EventID string, EventWindowID string) ([]EventHistorySession, *Error) {
	url := fmt.Sprintf("%s/api/v1/events/Fortnite/%s/%s/history/%s", EventsBaseUrl, EventID, EventWindowID, c.Config.AccountID)
	payload := []byte{}
	var response []EventHistorySession
	err, _ := c.doRequest("GET", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
