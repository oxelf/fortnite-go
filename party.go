package main

import (
	"time"
)

type PartyMemberMeta struct {
	Location string
	//Set 0 if you dont want to assign a specific number.
	PlayersLeft string
	Emote       string
	//Use 0001-01-01T00:00:00.000Z as default value
	MatchStartedTime  string
	LobbyState        LobbyState
	FrontEndMapMarker FrontEndMapMarker
	CosmeticLoadout   CosmeticLoadout
	Banner            Banner
	BattlePassInfo    BattlePassInfo
}
type LobbyState struct {
	Readiness bool
	InputType string
}
type FrontEndMapMarker struct {
	X     int
	Y     int
	IsSet bool
}
type CosmeticLoadout struct {
	Character string `json:"characterDef"`
	BackPack  string `json:"backpackDef"`
	PickAxe   string `json:"pickaxeDef"`
	Contrails string `json:"contrailDef"`
}
type UpdatedCosmeticLoadout struct {
	AthenaCosmeticLoadout struct {
		Character     string        `json:"characterDef"`
		CharacterEKey string        `json:"characterEKey"`
		BackPack      string        `json:"backpackDef"`
		BackPackEKey  string        `json:"backpackEKey"`
		PickAxe       string        `json:"pickaxeDef"`
		PickAxeEKey   string        `json:"pickaxeEKey"`
		Contrails     string        `json:"contrailDef"`
		ContrailEKey  string        `json:"contrailEKey"`
		ScratchPad    []interface{} `json:"scratchpad"`
		CosmeticStats []struct {
			StatName  string `json:"statName"`
			StatValue int    `json:"statValue"`
		} `json:"cosmeticStats"`
	} `json:"AthenaCosmeticLoadout"`
}
type Banner struct {
	BannerID      string
	BannerColorId string
	SeasonLevel   int
}
type BattlePassInfo struct {
	Purchased     bool
	Level         int
	SelfBoostXP   int
	FriendBoostXP int
}
type PartyPing struct {
	Sent       time.Time         `json:"sent"`
	Type       string            `json:"type"`
	Namespace  string            `json:"ns"`
	PingerID   string            `json:"pinger_id"`
	PingerDN   string            `json:"pinger_dn"`
	PingerPL   string            `json:"pinger_pl"`
	PingerPLDN string            `json:"pinger_pl_dn"`
	Expires    time.Time         `json:"expires"`
	Meta       map[string]string `json:"meta"`
	Message    Message
}
type PartyUser struct {
	AccountId string
	AccountDN string
	JoinedAt  time.Time
	UpdatedAt time.Time
	Skin      string
	Pickaxe   string
	Revision  int
	JID       string
}
type PartyJoin struct {
	Sent       time.Time `json:"sent"`
	Type       string    `json:"type"`
	Connection struct {
		ID              string            `json:"id"`
		Meta            map[string]string `json:"meta"`
		ConnectedAt     time.Time         `json:"connected_at"`
		UpdatedAt       time.Time         `json:"updated_at"`
		YieldLeadership bool              `json:"yield_leadership"`
	} `json:"connection"`
	Revision           int    `json:"revision"`
	Namespace          string `json:"ns"`
	PartyID            string `json:"party_id"`
	AccountID          string `json:"account_id"`
	AccountDN          string `json:"account_dn"`
	MemberStateUpdated struct {
		JoinRequestUsersJ string `json:"urn:epic:member:joinrequestusers_j"`
		DNS               string `json:"urn:epic:member:dn_s"`
	} `json:"member_state_updated"`
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type JoinRequestUsers struct {
	Users []struct {
		ID          string            `json:"id"`
		DisplayName string            `json:"dn"`
		Platform    string            `json:"plat"`
		Data        map[string]string `json:"data"`
	} `json:"users"`
}
type PartyNewCaptain struct {
	Sent      string `json:"sent"`
	Type      string `json:"type"`
	Revision  int    `json:"revision"`
	Ns        string `json:"ns"`
	PartyID   string `json:"party_id"`
	AccountID string `json:"account_id"`
	AccountDN string `json:"account_dn"`
	JoinedAt  string `json:"joined_at"`
	UpdatedAt string `json:"updated_at"`
}
type PartyMemberUpdated struct {
	Sent                  time.Time              `json:"sent"`
	Type                  string                 `json:"type"`
	Revision              int                    `json:"revision"`
	Ns                    string                 `json:"ns"`
	PartyID               string                 `json:"party_id"`
	AccountID             string                 `json:"account_id"`
	AccountDN             string                 `json:"account_dn"`
	MemberStateRemoved    []string               `json:"member_state_removed"`
	MemberStateUpdated    map[string]string      `json:"member_state_updated"`
	MemberStateOverridden map[string]interface{} `json:"member_state_overridden"`
	JoinedAt              time.Time              `json:"joined_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
}
type PartyMemberLeft struct {
	PartyID   string `json:"party_id"`
	AccountID string `json:"account_id"`
}
