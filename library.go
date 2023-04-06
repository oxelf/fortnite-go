package fortnitego

import (
	"fmt"
	"time"
)

const LibraryServiceUrl = "https://library-service.live.use1a.on.epicgames.com/library"

type PlayTimeQueryResponse struct {
	AccountID          string `json:"accountId"`
	ArtifactID         string `json:"artifactId"`
	TotalTimeInSeconds int    `json:"totalTime"`
}

type LibraryItem struct {
	ResponseMetaData struct {
		StateToken string `json:"stateToken"`
	} `json:"stateToken"`
	Records []LibraryItemRecord `json:"records"`
}

type LibraryItemRecord struct {
	NameSpace       string    `json:"namespace"`
	CatalogItemID   string    `json:"catalogItemId"`
	AppName         string    `json:"appName"`
	ProductID       string    `json:"productId"`
	SandBoxName     string    `json:"sandboxName"`
	SandBoxType     string    `json:"sandboxType"`
	AcquisitionDate time.Time `json:"acquisitionDate"`
}

// LauncherAppClient2 token required.
func (c *Client) Library_Playtime_All() ([]PlayTimeQueryResponse, *Error) {
	url := fmt.Sprintf("%s/api/public/playtime/account/%s/all", LibraryServiceUrl, c.Config.AccountID)
	payload := []byte{}
	var response []PlayTimeQueryResponse
	err, _ := c.doRequest("GET", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

// LauncherAppClient2 token required.
func (c *Client) Library_Playtime_For_Game(artifactId string) (*PlayTimeQueryResponse, *Error) {
	url := fmt.Sprintf("%s/api/public/playtime/account/%s/artifact/%s", LibraryServiceUrl, c.Config.AccountID, artifactId)
	payload := []byte{}
	var response PlayTimeQueryResponse
	err, _ := c.doRequest("GET", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return &response, nil
	}
}

// LauncherAppClient2 token required.
func (c *Client) Library_Items() (*LibraryItem, *Error) {
	url := fmt.Sprintf("%s/api/public/items?includeMetadata=true", LibraryServiceUrl)
	payload := []byte{}
	var response LibraryItem
	err, _ := c.doRequest("GET", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return &response, nil
	}
}

func (c *Client) CheckAccess() (*AccessResponse, *Error) {
	url := fmt.Sprintf("https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/accesscontrol/status")
	payload := []byte{}
	var response AccessResponse
	err, _ := c.doRequest("GET", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return &response, nil
	}
}

func (c *Client) RequestAccess() (*AccessResponse, *Error) {
	url := fmt.Sprintf("https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/storeaccess/v1/request_access/%s", c.Config.AccountID)
	payload := []byte{}
	var response AccessResponse
	err, _ := c.doRequest("POST", url, payload, false, &response)
	if err != nil {
		return nil, err
	} else {
		return &response, nil
	}
}
