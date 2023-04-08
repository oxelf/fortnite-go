package fortnitego

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Replay struct {
	ReplayId              string       `json:"ReplayName"`
	LengthInMS            int          `json:"LengthInMS"`
	NetworkVersion        int          `json:"NetworkVersion"`
	ChangeList            int          `json:"ChangeList"`
	FriendlyName          string       `json:"FriendlyName"`
	TimeStamp             time.Time    `json:"Timestamp"`
	IsLive                bool         `json:"bIsLive"`
	Compressed            bool         `json:"bCompressed"`
	DesiredDelayInSeconds int          `json:"DesiredDelayInSeconds"`
	CheckPoints           []CheckPoint `json:"Checkpoints"`
	DataChunks            []DataChunk  `json:"DataChunks"`
	Events                []CheckPoint `json:"Events"`
}
type CheckPoint struct {
	ID       string `json:"Id"`
	Group    string `json:"Group"`
	MetaData string `json:"Metadata"`
	Time1    int    `json:"Time1"`
	Time2    int    `json:"Time2"`
}

type DataChunk struct {
	ID          string `json:"Id"`
	Time1       int    `json:"Time1"`
	Time2       int    `json:"Time2"`
	SizeInBytes int    `json:"SizeInBytes"`
}

const ReplayBaseDataUrl = "https://datastorage-public-service-live.ol.epicgames.com/api/v1/access/fnreplays/public"

// ###################################
// #       	   Replays     	        #
// ###################################
// https://datastorage-public-service-live.ol.epicgames.com/api/v1/access/fnreplaysmetadata/public
func GetReplayMetadata(url string) (*Replay, *Error) {
	c := http.Client{}
	payload := []byte{}
	req, nerr := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if nerr != nil {

	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, &Error{
			ErrorMessage: nerr.Error(),
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &Error{ErrorMessage: "IO Read Error."}
		}
		var replay Replay
		rerr := json.Unmarshal(body, &replay)
		if rerr != nil {

		}
		return &replay, nil
	}
}

func GetReplayEventOrChunk(url string) (*[]byte, *Error) {
	c := http.Client{}
	payload := []byte{}

	req, nerr := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if nerr != nil {

	}
	req.Header.Add("User-Agent", "Tournament replay downloader")
	resp, err := c.Do(req)
	if err != nil {
		return nil, &Error{
			ErrorMessage: nerr.Error(),
		}
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, &Error{ErrorMessage: "IO Read Error."}
		}
		return &body, nil
	}
}
