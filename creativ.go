package fortnitego

import "fmt"

type Creativ_Favorites struct {
	Results []string `json:"results"`
	HasMore bool     `json:"hasMore"`
}

type Creativ_History struct {
	Results []string `json:"results"`
	HasMore bool     `json:"hasMore"`
}

type CreativFavoriteResponse struct {
	SortDate string `json:"sortDate"`
	LinkData struct {
		MapCode          string   `json:"mnemonic"`
		LinkType         string   `json:"linkType"`
		IsActive         bool     `json:"active"`
		Version          int      `json:"version"`
		ModerationStatus string   `json:"moderationStatus"`
		AccountID        string   `json:"accountId"`
		CreatorName      string   `json:"creatorName"`
		DescriptionTags  []string `json:"descriptionTags"`
		MetaData         struct {
			QuickSilverID string `json:"quicksilver_id"`
			ImageUrl      string `json:"image_url"`
			ImageUrls     struct {
				Url_S string `json:"url_s"`
				Url_M string `json:"url_m"`
				Url   string `json:"url"`
			} `json:"image_urls"`
			Title       string `json:"title"`
			Locale      string `json:"locale"`
			MatchMaking struct {
				SelectedJoinInProgressType int    `json:"selectedJoinInProgressType"`
				PlayersPerTeam             int    `json:"playersPerTeam"`
				MaximumNumberOfPlayers     int    `json:"maximumNumberOfPlayers"`
				OverridePlaylist           string `json:"override_Playlist"`
				PlayerCount                int    `json:"playerCount"`
				MMSType                    string `json:"mmsType"`
				MMSPrivacy                 string `json:"mmsPrivacy"`
				NumberOfTeams              int    `json:"numberOfTeams"`
				AllowJoinInProgress        bool   `json:"bAllowJoinInProgress"`
				MinimumNumberOfPlayers     int    `json:"minimumNumberOfPlayers"`
				JoinInProgessTeam          int    `json:"joinInProgressTeam"`
			} `json:"matchmaking"`
			VideoVUID          string `json:"video_vuid"`
			GeneratedImageUrls struct {
				Url_S      string `json:"url_s"`
				Url_M      string `json:"url_m"`
				Url        string `json:"url"`
				Compressed struct {
					Url_S string `json:"url_s"`
					Url_M string `json:"url_m"`
					Url   string `json:"url"`
				} `json:"compressed"`
			} `json:"generated_image_urls"`
			Mode       string `json:"mode"`
			VideoUrl   string `json:"video_url"`
			TagLine    string `json:"tagline"`
			IslandType string `json:"islandType"`
			DynamicXP  struct {
				UniqueGameVersion int    `json:"uniqueGameVersion"`
				CalibrationPhase  string `json:"calibrationPhase"`
			} `json:"dynamicXp"`
			SupportCode  string `json:"supportCode"`
			Introduction string `json:"introduction"`
		} `json:"metadata"`
	} `json:"linkData"`
	IsFavorite bool `json:"isFavorite"`
}

const FortnitePublicService = "https://fngw-mcp-gc-livefn.ol.epicgames.com"

// Same as Creativ Favorites, does not return the history, just empty results field.
// func (c *Client) Creativ_History() (*Creativ_History, *Error) {
// 	url := fmt.Sprintf("https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/game/v2/creative/history/%s", c.Config.AccountID)
// 	payload := []byte{}
// 	var creativHistory Creativ_History
// 	err, _ := c.doRequest("GET", url, payload, false, &creativHistory)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		return &creativHistory, nil
// 	}
// }

// This should work, but returns nothing.
// func (c *Client) Creativ_Favorites() (*Creativ_Favorites, *Error) {
// 	url := fmt.Sprintf("https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/game/v2/creative/favorites/%s?limit=30", c.Config.AccountID)
// 	payload := []byte{}
// 	var creativFavorites Creativ_Favorites
// 	err, _ := c.doRequest("GET", url, payload, false, &creativFavorites)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		return &creativFavorites, nil
// 	}
// }

func (c *Client) Creativ_Add_Favorite(mapCode string) (*CreativFavoriteResponse, *Error) {
	url := fmt.Sprintf("%s/fortnite/api/game/v2/creative/favorites/%s/%s", FortnitePublicService, c.Config.AccountID, mapCode)
	payload := []byte{}
	var creativFavorite CreativFavoriteResponse
	err, _ := c.doRequest("PUT", url, payload, false, &creativFavorite)
	if err != nil {
		return nil, err
	} else {
		return &creativFavorite, nil
	}
}

func (c *Client) Creativ_Remove_Favorite(mapCode string) (*CreativFavoriteResponse, *Error) {
	url := fmt.Sprintf("%s/fortnite/api/game/v2/creative/favorites/%s/%s", FortnitePublicService, c.Config.AccountID, mapCode)
	payload := []byte{}
	var creativFavorite CreativFavoriteResponse
	err := c.doNullableRequest("DELETE", url, payload, false)
	if err != nil {
		return nil, err
	} else {
		return &creativFavorite, nil
	}
}
