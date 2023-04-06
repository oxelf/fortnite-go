package fortnitego

import (
	"fmt"
	"time"
)

type Friend struct {
	AccountId string    `json:"accountId"`
	Status    string    `json:"status"`
	Direction string    `json:"direction"`
	Created   time.Time `json:"created"`
	Favorite  bool      `json:"favorite"`
}

type OutgoingFriend struct {
	AccountID string    `json:"accountId"`
	Groups    []string  `json:"groups"`
	Mutuals   int       `json:"mutual"`
	Alias     string    `json:"alias"`
	Note      string    `json:"note"`
	Favorite  string    `json:"favorite"`
	Created   time.Time `json:"created"`
}

type FriendsSummary struct {
	Friends       []OutgoingFriend  `json:"friends"`
	Suggested     []SuggestedFriend `json:"suggested"`
	BlockList     []AccountID       `json:"blocklist"`
	Settings      FriendsSettings   `json:"settings"`
	LimitsReached FriendsLimits     `json:"limitsReached"`
}

type SuggestedFriend struct {
	AccountID   string `json:"accountId"`
	Mutuals     int    `json:"mutuals"`
	Connections struct {
		Epic struct {
			ID          string `json:"id"`
			SortFactors struct {
				X int       `json:"x"`
				Y int       `json:"y"`
				K time.Time `json:"k"`
				L time.Time `json:"l"`
			} `json:"sortFactors"`
		} `json:"epic"`
	} `json:"connections"`
}

type AccountID struct {
	AccountID string `json:"accountId"`
}

type FriendsSettings struct {
	AcceptInvites string `json:"acceptInvites"`
	MutualPrivacy string `json:"mutualPrivacy"`
}

type FriendsLimits struct {
	Incoming bool `json:"incoming"`
	Outgoing bool `json:"outgoing"`
	Accepted bool `json:"accepted"`
}

//###################################
//#            Friends              #
//###################################

func (c *Client) Friend_Add_Or_Accept(friendId string) *Error {
	url := fmt.Sprintf("%s/friends/api/v1/%s/friends/%s", BaseRoute.FriendsPublicService, c.Config.AccountID, friendId)
	payload := []byte{}
	requestError := c.doNullableRequest("POST", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) Friend_Remove_Or_Decline(friendId string) *Error {
	url := fmt.Sprintf("%s/friends/api/v1/%s/friends/%s", BaseRoute.FriendsPublicService, c.Config.AccountID, friendId)
	payload := []byte{}
	requestError := c.doNullableRequest("DELETE", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) Friend_Block(friendId string) *Error {
	url := fmt.Sprintf("%s/friends/api/v1/%s/blocklist/%s", BaseRoute.FriendsPublicService, c.Config.AccountID, friendId)
	payload := []byte{}
	requestError := c.doNullableRequest("POST", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) Friend_Unblock(friendId string) *Error {
	url := fmt.Sprintf("%s/friends/api/v1/%s/blocklist/%s", BaseRoute.FriendsPublicService, c.Config.AccountID, friendId)
	payload := []byte{}
	requestError := c.doNullableRequest("DELETE", url, payload, false)
	if requestError != nil {
		return requestError
	} else {
		return nil
	}
}

func (c *Client) Friends_Get_All() ([]Friend, *Error) {
	url := fmt.Sprintf("%s/friends/api/public/friends/%s", BaseRoute.FriendsPublicService, c.Config.AccountID)
	payload := []byte{}
	var friendList []Friend
	err, _ := c.doRequest("GET", url, payload, false, &friendList)
	if err != nil {
		return nil, err
	} else {
		return friendList, nil
	}
}

func (c *Client) Friends_Get_Incoming() ([]OutgoingFriend, *Error) {
	url := fmt.Sprintf("%s/friends/api/v1/%s/incoming", BaseRoute.FriendsPublicService, c.Config.AccountID)
	payload := []byte{}
	var friendList []OutgoingFriend
	err, _ := c.doRequest("GET", url, payload, false, &friendList)
	if err != nil {
		return nil, err
	} else {
		return friendList, nil
	}
}

func (c *Client) Friends_Get_Outgoing() ([]OutgoingFriend, *Error) {
	url := fmt.Sprintf("%s/friends/api/v1/%s/outgoing", BaseRoute.FriendsPublicService, c.Config.AccountID)
	payload := []byte{}
	var friendList []OutgoingFriend
	err, _ := c.doRequest("GET", url, payload, false, &friendList)
	if err != nil {
		return nil, err
	} else {
		return friendList, nil
	}
}

func (c *Client) Friends_Get_Blocklist() ([]Friend, *Error) {
	url := fmt.Sprintf("%s/friends/api/v1/%s/blocklist", BaseRoute.FriendsPublicService, c.Config.AccountID)
	payload := []byte{}
	var blockList []Friend
	err, _ := c.doRequest("GET", url, payload, false, &blockList)
	if err != nil {
		return nil, err
	} else {
		return *&blockList, nil
	}
}

func (c *Client) Friends_Get_Summary() (*FriendsSummary, *Error) {
	url := fmt.Sprintf("%s/friends/api/v1/%s/summary", BaseRoute.FriendsPublicService, c.Config.AccountID)
	payload := []byte{}
	var summary FriendsSummary
	err, _ := c.doRequest("GET", url, payload, false, &summary)
	if err != nil {
		return nil, err
	} else {
		return &summary, nil
	}
}

// Currently not working!!! Problem is the body.
// func (c *Client) Friends_Set_Nickname(nickName string, friendID string) *Error {
// 	uri := fmt.Sprintf("%s/friends/api/v1/%s/friends/%s/alias", BaseRoute.FriendsPublicService, c.Config.AccountID, friendID)
// 	body := url.Values{}
// 	body.Set("nick", nickName)
// 	stringBody := url.QueryEscape(nickName)
// 	requestError := c.doNullableRequestUrlEncoded("PUT", uri, stringBody)
// 	if requestError != nil {
// 		return requestError
// 	} else {
// 		return nil
// 	}
// }
