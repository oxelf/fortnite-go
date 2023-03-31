package main

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

type Open struct {
	XMLName xml.Name `xml:"open"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	From    string   `xml:"from,attr"`
	ID      string   `xml:"id,attr"`
	Version string   `xml:"version,attr"`
	Lang    string   `xml:"lang,attr"`
}

type Presence struct {
	XMLName xml.Name `xml:"presence"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	From    string   `xml:"from,attr"`
	To      string   `xml:"to,attr"`
	Status  string   `xml:"status"`
	Type    string   `xml:"type,attr"`
	Delay   struct {
		Text  string `xml:",chardata"`
		Xmlns string `xml:"xmlns,attr"`
		Stamp string `xml:"stamp,attr"`
	} `xml:"delay"`
}

type Iq struct {
	XMLName xml.Name `xml:"iq"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	ID      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	Ping    Ping     `xml:"ping"`
}

type Ping struct {
	Text  string `xml:",chardata"`
	Xmlns string `xml:"xmlns,attr"`
}

type Status struct {
	Status           string `json:"Status"`
	BIsPlaying       bool   `json:"bIsPlaying"`
	BIsJoinable      bool   `json:"bIsJoinable"`
	BHasVoiceSupport bool   `json:"bHasVoiceSupport"`
	SessionID        string `json:"SessionId"`
	ProductName      string `json:"ProductName"`
}
type Message struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr"`
	From    string   `xml:"from,attr"`
	CorrID  string   `xml:"corr-id,attr"`
	To      string   `xml:"to,attr"`
	Body    struct {
		RawJSON json.RawMessage `xml:",innerxml"`
		Type    string          `xml:"type"`
	} `xml:"body"`
}

type Body struct {
	Sent       time.Time         `json:"sent"`
	Type       string            `json:"type"`
	NS         string            `json:"ns"`
	Expires    time.Time         `json:"expires"`
	Meta       map[string]string `json:"meta"`
	PingerID   string            `json:"pinger_id"`
	PingerDN   string            `json:"pinger_dn"`
	PingerPL   string            `json:"pinger_pl"`
	PingerPLDN string            `json:"pinger_pl_dn"`
	Revision   int               `json:"revision"`
	PartyId    string            `json:"party_id"`
	CaptainId  string            `json:"captain_id"`
	AccountId  string            `json:"account_id"`
	AccountDN  string            `json:"account_dn"`
}

type FriendshipRequest struct {
	From       string    `json:"from"`
	To         string    `json:"to"`
	Timestamp  time.Time `json:"timestamp"`
	Status     string    `json:"status"`
	RawMessage json.RawMessage
}

type BlocklistUpdate struct {
	OwnerID    string `json:"ownerId"`
	AccountID  string `json:"accountId"`
	Status     string `json:"status"`
	RawMessage json.RawMessage
}
