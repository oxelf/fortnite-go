package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
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
		AuthClient: AuthClients.Fortnite_PC_Client,
		Token:      "eg1~eyJraWQiOiJ0RkMyVUloRnBUTV9FYTNxY09kX01xUVQxY0JCbTlrRkxTRGZlSmhzUkc4IiwiYWxnIjoiUFMyNTYifQ.eyJhcHAiOiJwcm9kLWZuIiwic3ViIjoiNDk3NDhjZmU0MDMzNGExMjljYTM2MjAyNTJhYmZiYTIiLCJtdmVyIjpmYWxzZSwiY2xpZCI6ImVjNjg0YjhjNjg3ZjQ3OWZhZGVhM2NiMmFkODNmNWM2IiwiZG4iOiJveGlib3QgMDAxIiwiYW0iOiJhdXRob3JpemF0aW9uX2NvZGUiLCJwZnBpZCI6InByb2QtZm4iLCJwIjoiZU5xMVdOdHUyemdRXC9aXC9BRGVwTHNnMEJQMlRUcEJzZzNRVHJZdEUzWXl5T2JOWVVxU1VwSlwvNzdIVklYUzY0dlVxbysyWmJNdVo0NU04TllHNmVFUXhaSm5YSHJ0SUVsTXJ1MURoUDJnT0F5Z1wvelJTbEI4Y2luQzUzUjRQWkJDclMyTEZadmNcL0RINUZNVTQrVGdlVDJBNHVvbGdmRDM2T0xvYXdTSmV3R2c2R3NTblZQeTdmcEd3XC9ZeGNST0NRejlCczBIeUJCQytGRWw0VFJKSE9sR1BGcHoycmtTWG9nSU1EdGtybm1VVXpGM3llR3MzSkZLc2pBWkl0UVVvMFdcL29uRjhBdXpqc3hIR0Ftd1JtSTFrSXRXWm90cElnWUZ6WWw0MkZwRUJNa0cxc0VvM1NuRUZIK3hEZUhSb0c4emR6S1RvZFhBeElubkF4aTZaUlE5Qm9pSjdTeWJlemQ2WEY2allwcHQwSXpRMnU5Z0FkdGJ2TzNkMUo0K1pcL0t1Q3hBbFlhZFZYRjlKSytwd1g5UUlsaWNoUXdcL0NJbTJqZ0pLUlV6UFdnVDlhaUFoVXhHWnpyaCtWVklESitjZkNrRjNtb0tpM0orWmtMeWhJQUVYclJJSW1iSzV6K2RCSTlTR3p2cFFwd0U3QnBkMDdtakM2RVdVNlwvOUFDTnVJaU5ScW5wRmJZSnlJS1ZYMER3cUZSVVd2em1lTURIV09ETGJUNjMyTnl3d013VlRaNmZpMGp4ZWxFemNEcWlXUW1sNnRnTXFYWFpBeG5JSlhIUWRMK214NEhyY25nR0dOQUE0ZjgrVjIzdGNjZHJQQzQ4Zm4yYVdGalU5M0ErWlViMmhkTzdUSGhpUnkyd3BUXC9WbitmZmI5TDVSQzI2OTZRWWd1dmVoUHdWZUllaEJLOGNIVVEzZERxREFWNWV6QWtBcVBEMkVqdmZHc2VOYkl5QkFzdEdFcG1SQUtJVWt5NWZrNzVLMmtNc2JSRnhIaHlMQm9CVXFoSkticE1mb0xcL2Zhc3FxaW5hS3hXVUxrbnRWNW42VkdBdjBmalU2YWd5a2VCT0dJQmhRVml1XC9UQ2hJb1dscjU4MzRcL1pva2ozbXVlOUlyRlYrNnhPK2tNWUd4MjQ2emZVd09Xc21ZZVMwS3FFdFBIekxVblRHcCsxNk12OWVYS3J1TkdDOTFcL0ZMN09yU21nbXFmZkpNQ3dVOUI1cGp1UkhcL211TjIrSjVuZXdQejAzM1hMaVRxZVwvVWZyS0ZqWXhJY3lBdkRMWGFDS3dqdzd6SFQxVEZPNzQ0Syt6aUxHU2JmV1gwV1wvcktsNGVcLys4K2xwNEFtMjNmZ1Rac1ptZ3lRXC9aZlIzejByNzJhODFsWENhb2ZlTlZibDJUa3lobDUwYXdUNUZyQVNQcUxiczBrXC9DTm42cU45WmN3d2JiZnpFTXh6V0R1K0grWlJSWGFlMGVwbDRyVHZaSUtrZmdDbWE3WWVLXC9zb3BqdnowQzlQMFZ6YWE0XC91RGp3ZVZFd1wvbzJQYktpeVhRWmJ6c09NZ1dXMFZLSVwvRzIzektrMm01V1lUOWJvbldoSmRmQ3JJRVdNalwvU3pNT1hWXC9SNVpBYUVuQmR0YXU0OGJidnBaT2RlK1dWdXRQVE9oaVdMak94T3BkMFd1Vm0raWZoOUx2OG1GSGxFME82MTQ3XC9NSnFkNTlTREhqdzl4XC9NRXRob3g5VDhcL0plVzFcL0tVdHBZOCtuMUY0UVVnOGtoYmJKd0Q5SGxxeUp4YkwrbngxTGxZOStFaGtpa0JxeUhrTTdMNWM0SGNja01yejR4WHhlZE5pWVdqVHpGcGMzKzJ3Vlp2TnFVM2prSFdsWTBtanlySWp6OEYzWFM4MHhxYnpLcW90NjFVYnlodU41YXRZZ2pHNHg2RThHQzZtamRjRHcrU0dXcUVQRnRPdUtqV2ZMVG4zODBMVVVCVElsZG1qQkx1T1NvaDFDWXBrcUd0U3ArcXhkb2htZHBlRlFlU1B6R1RkUFlXbHZ5R2k5NkRRbms1SjVcL2VxNHhHK2VcL2U3Q3FIeWd5STMrUVRzWXkyK0xQb1JPZWQzdmpyUGZaXC9ZQW5RaWx6ZjFiU2xZZ3I3ZUNRNWowNVpjM1UybzNcL1M2b1FtVnZwYW5IMXZKcUtlY1l1c1d3VEtoRk1OR3ExUmJiSGE2TllqcDhlXC9WREM5Vm01ZHZiOUZ2ZCt2d1BrZFwvVEhnPT0iLCJpYWkiOiI0OTc0OGNmZTQwMzM0YTEyOWNhMzYyMDI1MmFiZmJhMiIsInNlYyI6MCwiY2xzdmMiOiJwcm9kLWZuIiwidCI6InMiLCJpYyI6dHJ1ZSwiZXhwIjoxNjgwMjg1MjIwLCJpYXQiOjE2ODAyNzgwMjAsImp0aSI6ImMxMWQ1NmYxMzIxMTQyNGI4YTQ2OTNmZDJkNzA2ZDMwIn0.AI_3VY_CtnX-icOSV5Zi5B7ur-TSJrw2ZPz3ME4chvLs5iCCOVbSS9qEtHdS_qIigM_RE9FrHIO90TTC6g7sgdAq",
		AccountID:  "49748cfe40334a129ca3620252abfba2",
	}
	client, err := NewClient(config)
	if err != nil {
		fmt.Printf("failed to create new xmpp client: %v\n", err)
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

			if msg, err := client.Read(); err != nil {
				fmt.Printf("ERROR %v\n", err)
				break
			} else {

				if strings.HasPrefix(msg, "<presence") {

					presence := &Presence{}
					err := xml.Unmarshal([]byte(msg), presence)
					if err != nil {
						//TODO only log
						continue
					}

					status := &Status{}
					err = json.Unmarshal([]byte(presence.Status), status)
					if err != nil {
						//TODO only log
						continue
					}

					fmt.Printf("<Presence>: type:%s,\n from: %#v,\n status: %#v", presence.Type, presence.From, status)
				}
				if strings.HasPrefix(msg, "<message") {

					message := &Message{}
					err := xml.Unmarshal([]byte(msg), message)
					if err != nil {
						//TODO only log
						continue
					}
					body := &Body{}
					Uerr := json.Unmarshal(message.Body.RawJSON, &body)
					if Uerr != nil {
						fmt.Println("unmarshalling error.")
					}
					switch body.Type {
					case "com.epicgames.social.party.notification.v0.PING":
						fmt.Printf("Ping from: %s\n", body.PingerDN)
						res := client.PartyLookupPing(body.PingerID, client.Config.AccountID)
						fmt.Printf("Party id: %s\n", res.ID)
						fmt.Printf("jid: %s", message.To)
						resp := client.PartySendJoinRequest(message.To, res.ID)
						fmt.Println("resp:" + resp.ID)
					case "com.epicgames.social.party.notification.v0.MEMBER_JOINED":
						fmt.Println("joined: " + body.AccountDN)
					case "com.epicgames.social.party.notification.v0.MEMBER_NEW_CAPTAIN":
						fmt.Printf("new captain: %s, revision: %d\n", body.CaptainId, body.Revision)
						client.Party.Id = body.PartyId
						if body.Revision != 0 {
							client.Party.PartyRevision = body.Revision
						}
						setPlaylist := client.SetPlaylist()
						fmt.Println(setPlaylist)
						setKey := client.SetCustomKey("1234")
						fmt.Println(setKey)
						client.PartyLeave()
					case "com.epicgames.social.party.notification.v0.PARTY_UPDATED":
						fmt.Printf("Party updated from: %s, revision: %d", body.AccountDN, body.Revision)
						if body.Revision != 0 {
							client.Party.PartyRevision = body.Revision
						}
					default:
						continue
					}

				}
			}
		}
	}
}
