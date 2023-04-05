package fortnitego

import "time"

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
