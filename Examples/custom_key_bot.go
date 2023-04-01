package main

import (
	"fmt"
)

var accountId = "49748cfe40334a129ca3620252abfba2"

var metaData = &PartyMemberMeta{
	Location:         "PreLobby",
	PlayersLeft:      "0",
	MatchStartedTime: "0001-01-01T00:00:00.000Z",
	Emote:            "EID_Worm",
	LobbyState: LobbyState{
		Readiness: true,
		InputType: "MouseAndKeyboard",
	},
	FrontEndMapMarker: FrontEndMapMarker{
		IsSet: false,
		X:     0,
		Y:     0,
	},
	CosmeticLoadout: CosmeticLoadout{
		Character: "CID_029_Athena_Commando_F_Halloween",
		BackPack:  "None",
		Contrails: "None",
		PickAxe:   "DefaultPickaxe",
	},
	Banner: Banner{
		BannerID:      "standardbanner15",
		BannerColorId: "defaultcolor15",
		SeasonLevel:   1,
	},
	BattlePassInfo: BattlePassInfo{
		Level:         12,
		Purchased:     true,
		SelfBoostXP:   0,
		FriendBoostXP: 0,
	},
}

func main() {
	config := ClientConfig{
		XMPP:       true,
		Connection: "Fortnite",
		AuthClient: AuthClients.Fortnite_PC_Client}
	url := GetAuthCodeUrl(AuthClients.Fortnite_IOS_Client)
	fmt.Println(url)
	val, eerr := Get_OauthToken_By_AuthCode("bf0dbfd3231c43809b8c1a02720f2449", Base64AuthClients.Fortnite_IOS_Client, true)
	if eerr != nil {
		fmt.Println(eerr.EpicErrorMessage)
	} else {
		fmt.Println("Acc token: " + val.AccessToken)
		config.Token = val.AccessToken
		config.AccountID = val.AccountID
		Create_DeviceAuth(val.AccessToken, val.AccountID)
	}
	client, cerr := NewClient(config)
	if cerr != nil {
		fmt.Printf("failed to create new xmpp client: %v\n", cerr)
	} else {
		fmt.Println("Program started.")
		client.Friends_Add_Or_Accept("463e92ed9ed047869788c331fac51e9d")
		res, ResErr := client.LightSwitch_Status_Fortnite()
		if ResErr != nil {
			fmt.Println(ResErr)
		} else {
			fmt.Println(res)
		}
		client.OnPing(
			func(p *PartyPing) {
				fmt.Println("New Ping from: " + p.PingerDN)
				res := client.PartyLookupPing(p.PingerID, client.Config.AccountID)
				fmt.Printf("Party id: %s\n", res.ID)
				fmt.Printf("jid: %s", p.Message.To)
				resp := client.PartySendJoinRequest(p.Message.To, res.ID)
				fmt.Println("resp:" + resp.ID)
			},
		)
		client.OnJoin(func(j *PartyJoin) {
			fmt.Println("New Join: " + j.AccountDN)
			partyLookup, err := client.PartyLookup(client.Party.Id)
			if err != nil {
				fmt.Println(partyLookup)
			}
			response := client.PartySendInitialMemberData(metaData)
			fmt.Println(response)
			keyRes := client.SetReadiness(j.PartyID, true)
			fmt.Println(keyRes)
		})
		client.OnMemberLeft(
			func(m *PartyMemberLeft) {
				fmt.Println("Member left: " + m.AccountID)
			},
		)
		client.OnSkinChanged(func(s string, id string) {
			if id != client.Config.AccountID {
				fmt.Println("new Skin: " + s)
				err := client.Set_Skin(s)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("client new skin:" + s)
			}
		})
		client.OnNewCaptain(func(c *PartyNewCaptain) {
			fmt.Println("new captain: " + c.AccountDN)
			setPlaylist := client.SetPlaylist()
			fmt.Println(setPlaylist)
			setKey := client.SetCustomKey("1234")
			fmt.Println(setKey)
			client.PartyLeave()
		})
		client.Listen()
		for {
			client.Listen()
		}
	}
}
