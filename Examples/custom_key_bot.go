package main

import (
	"fmt"

	fortnitego "github.com/oxelf/fortnite-go"
)

var accountId = "49748cfe40334a129ca3620252abfba2"

var metaData = &fortnitego.PartyMemberMeta{
	Location:         "PreLobby",
	PlayersLeft:      "0",
	MatchStartedTime: "0001-01-01T00:00:00.000Z",
	Emote:            "EID_Worm",
	LobbyState: fortnitego.LobbyState{
		Readiness: true,
		InputType: "MouseAndKeyboard",
	},
	FrontEndMapMarker: fortnitego.FrontEndMapMarker{
		IsSet: false,
		X:     0,
		Y:     0,
	},
	CosmeticLoadout: fortnitego.CosmeticLoadout{
		Character: "CID_029_Athena_Commando_F_Halloween",
		BackPack:  "None",
		Contrails: "None",
		PickAxe:   "DefaultPickaxe",
	},
	Banner: fortnitego.Banner{
		BannerID:      "standardbanner15",
		BannerColorId: "defaultcolor15",
		SeasonLevel:   1,
	},
	BattlePassInfo: fortnitego.BattlePassInfo{
		Level:         12,
		Purchased:     true,
		SelfBoostXP:   0,
		FriendBoostXP: 0,
	},
}

func main() {
	config := fortnitego.ClientConfig{
		XMPP:       true,
		Connection: "Fortnite",
		AuthClient: fortnitego.AuthClients.Fortnite_PC_Client}
	url := fortnitego.GetAuthCodeUrl(fortnitego.AuthClients.Fortnite_IOS_Client)
	fmt.Println(url)
	val, eerr := fortnitego.Get_OauthToken_By_AuthCode("3cb997fb0a5844bb90ba49492abbc298", fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	if eerr != nil {
		fmt.Println(eerr.EpicErrorMessage)
	} else {
		fmt.Println("Acc token: " + val.AccessToken)
		config.Token = val.AccessToken
		config.AccountID = val.AccountID
		fortnitego.Create_DeviceAuth(val.AccessToken, val.AccountID)
	}
	client, cerr := fortnitego.NewClient(config)
	if cerr != nil {
		fmt.Printf("failed to create new xmpp client: %v\n", cerr)
	} else {
		fmt.Println("Program started.")
		client.OnPresence(func(p *fortnitego.Status) {
			if p.SessionID != "" {
				fileLocation, _ := client.FileLocationInfo(fmt.Sprint("https://datastorage-public-service-live.ol.epicgames.com/api/v1/access/fnreplaysmetadata/public%2F" + p.SessionID + ".json"))
				fileUrl := ""
				for _, value := range fileLocation.Files {
					fileUrl = value.ReadLink // return the first key-value pair
					break                    // exit the loop after the first iteration
				}
				fortnitego.DownloadReplayCDNFile(fileUrl)
			}

		})

		res, ResErr := client.LightSwitch_Status_Fortnite()
		if ResErr != nil {
			fmt.Println(ResErr)
		} else {
			fmt.Println(res)
		}
		client.OnFriendRequest(func(fr *fortnitego.FriendshipRequest) {
			client.Friends_Add_Or_Accept(fr.From)
		})
		client.OnPing(
			func(p *fortnitego.PartyPing) {
				fmt.Println("New Ping from: " + p.PingerDN)
				res := client.PartyLookupPing(p.PingerID, client.Config.AccountID)
				fmt.Printf("Party id: %s\n", res.ID)
				fmt.Printf("jid: %s", p.Message.To)
				resp := client.PartySendJoinRequest(p.Message.To, res.ID)
				fmt.Println("resp:" + resp.ID)
			},
		)
		client.OnJoin(func(j *fortnitego.PartyJoin) {
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
			func(m *fortnitego.PartyMemberLeft) {
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
		client.OnNewCaptain(func(c *fortnitego.PartyNewCaptain) {
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
