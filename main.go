package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
)

var token = "eg1~eyJraWQiOiJ0RkMyVUloRnBUTV9FYTNxY09kX01xUVQxY0JCbTlrRkxTRGZlSmhzUkc4IiwiYWxnIjoiUFMyNTYifQ.eyJhcHAiOiJwcm9kLWZuIiwic3ViIjoiNDk3NDhjZmU0MDMzNGExMjljYTM2MjAyNTJhYmZiYTIiLCJkdmlkIjoiYWIyNDU0YzEwNmU3NDdjMDkxMDA5YTAyODBjNzhlY2MiLCJtdmVyIjpmYWxzZSwiY2xpZCI6ImVjNjg0YjhjNjg3ZjQ3OWZhZGVhM2NiMmFkODNmNWM2IiwiZG4iOiJveGlib3QgMDAxIiwiYW0iOiJhdXRob3JpemF0aW9uX2NvZGUiLCJwZnBpZCI6InByb2QtZm4iLCJwIjoiZU5xMVdOdHUyemdRXC9aXC9BRGVwTHNnMEJQMlRUcEJzZzNRVHJZdEUzWXl5T2JOWVVxU1VwSlwvNzdIVklYUzY0dlVxbysyWmJNdVo0NU04TllHNmVFUXhaSm5YSHJ0SUVsTXJ1MURoUDJnT0F5Z1wvelJTbEI4Y2luQzUzUjRQWkJDclMyTEZadmNcL0RINUZNVTQrVGdlVDJBNHVvbGdmRDM2T0xvYXdTSmV3R2c2R3NTblZQeTdmcEd3XC9ZeGNST0NRejlCczBIeUJCQytGRWw0VFJKSE9sR1BGcHoycmtTWG9nSU1EdGtybm1VVXpGM3llR3MzSkZLc2pBWkl0UVVvMFdcL29uRjhBdXpqc3hIR0Ftd1JtSTFrSXRXWm90cElnWUZ6WWw0MkZwRUJNa0cxc0VvM1NuRUZIK3hEZUhSb0c4emR6S1RvZFhBeElubkF4aTZaUlE5Qm9pSjdTeWJlemQ2WEY2allwcHQwSXpRMnU5Z0FkdGJ2TzNkMUo0K1pcL0t1Q3hBbFlhZFZYRjlKSytwd1g5UUlsaWNoUXdcL0NJbTJqZ0pLUlV6UFdnVDlhaUFoVXhHWnpyaCtWVklESitjZkNrRjNtb0tpM0orWmtMeWhJQUVYclJJSW1iSzV6K2RCSTlTR3p2cFFwd0U3QnBkMDdtakM2RVdVNlwvOUFDTnVJaU5ScW5wRmJZSnlJS1ZYMER3cUZSVVd2em1lTURIV09ETGJUNjMyTnl3d013VlRaNmZpMGp4ZWxFemNEcWlXUW1sNnRnTXFYWFpBeG5JSlhIUWRMK214NEhyY25nR0dOQUE0ZjgrVjIzdGNjZHJQQzQ4Zm4yYVdGalU5M0ErWlViMmhkTzdUSGhpUnkyd3BUXC9WbitmZmI5TDVSQzI2OTZRWWd1dmVoUHdWZUllaEJLOGNIVVEzZERxREFWNWV6QWtBcVBEMkVqdmZHc2VOYkl5QkFzdEdFcG1SQUtJVWt5NWZrNzVLMmtNc2JSRnhIaHlMQm9CVXFoSkticE1mb0xcL2Zhc3FxaW5hS3hXVUxrbnRWNW42VkdBdjBmalU2YWd5a2VCT0dJQmhRVml1XC9UQ2hJb1dscjU4MzRcL1pva2ozbXVlOUlyRlYrNnhPK2tNWUd4MjQ2emZVd09Xc21ZZVMwS3FFdFBIekxVblRHcCsxNk12OWVYS3J1TkdDOTFcL0ZMN09yU21nbXFmZkpNQ3dVOUI1cGp1UkhcL211TjIrSjVuZXdQejAzM1hMaVRxZVwvVWZyS0ZqWXhJY3lBdkRMWGFDS3dqdzd6SFQxVEZPNzQ0Syt6aUxHU2JmV1gwV1wvcktsNGVcLys4K2xwNEFtMjNmZ1Rac1ptZ3lRXC9aZlIzejByNzJhODFsWENhb2ZlTlZibDJUa3lobDUwYXdUNUZyQVNQcUxiczBrXC9DTm42cU45WmN3d2JiZnpFTXh6V0R1K0grWlJSWGFlMGVwbDRyVHZaSUtrZmdDbWE3WWVLXC9zb3BqdnowQzlQMFZ6YWE0XC91RGp3ZVZFd1wvbzJQYktpeVhRWmJ6c09NZ1dXMFZLSVwvRzIzektrMm01V1lUOWJvbldoSmRmQ3JJRVdNalwvU3pNT1hWXC9SNVpBYUVuQmR0YXU0OGJidnBaT2RlK1dWdXRQVE9oaVdMak94T3BkMFd1Vm0raWZoOUx2OG1GSGxFME82MTQ3XC9NSnFkNTlTREhqdzl4XC9NRXRob3g5VDhcL0plVzFcL0tVdHBZOCtuMUY0UVVnOGtoYmJKd0Q5SGxxeUp4YkwrbngxTGxZOStFaGtpa0JxeUhrTTdMNWM0SGNja01yejR4WHhlZE5pWVdqVHpGcGMzKzJ3Vlp2TnFVM2prSFdsWTBtanlySWp6OEYzWFM4MHhxYnpLcW90NjFVYnlodU41YXRZZ2pHNHg2RThHQzZtamRjRHcrU0dXcUVQRnRPdUtqV2ZMVG4zODBMVVVCVElsZG1qQkx1T1NvaDFDWXBrcUd0U3ArcXhkb2htZHBlRlFlU1B6R1RkUFlXbHZ5R2k5NkRRbms1SjVcL2VxNHhHK2VcL2U3Q3FIeWd5STMrUVRzWXkyK0xQb1JPZWQzdmpyUGZaXC9ZQW5RaWx6ZjFiU2xZZ3I3ZUNRNWowNVpjM1UybzNcL1M2b1FtVnZwYW5IMXZKcUtlY1l1c1d3VEtoRk1OR3ExUmJiSGE2TllqcDhlXC9WREM5Vm01ZHZiOUZ2ZCt2d1BrZFwvVEhnPT0iLCJpYWkiOiI0OTc0OGNmZTQwMzM0YTEyOWNhMzYyMDI1MmFiZmJhMiIsInNlYyI6MCwiY2xzdmMiOiJwcm9kLWZuIiwidCI6InMiLCJpYyI6dHJ1ZSwiZXhwIjoxNjgwMDQyMDM2LCJpYXQiOjE2ODAwMzQ4MzYsImp0aSI6Ijg0OWE3OGJiY2EzNzQzMTBhY2RjNDRjMWJjM2Q1NTlhIn0.AC8o3iJUwg1q-4RF3L1qvH2Bp8LuE1XGYq-cLYoOqMzWzora7DCUbRSU7ivCh1ueDbK6mgs-vUTz63xVYzx9pEXk"
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

// Bearer: eg1~eyJraWQiOiJ0RkMyVUloRnBUTV9FYTNxY09kX01xUVQxY0JCbTlrRkxTRGZlSmhzUkc4IiwiYWxnIjoiUFMyNTYifQ.eyJhcHAiOiJwcm9kLWZuIiwic3ViIjoiZWM3ZTc0YTQzNDIyNGIyZTlmMDBkYWQ5OWRmYmVmYTIiLCJtdmVyIjpmYWxzZSwiY2xpZCI6IjM0NDZjZDcyNjk0YzRhNDQ4NWQ4MWI3N2FkYmIyMTQxIiwiZG4iOiJveGVsZiBheSIsImFtIjoiZGV2aWNlX2F1dGgiLCJwZnBpZCI6InByb2QtZm4iLCJwIjoiZU5xMVdHMXoyamdRXC9qK1pKRk1JYlMrYTRVTXVUWHFaU1MrZG8zT1RiOHdpclVGRmxueVNUTUtcL3Y1WDhFcE1DdG9GK0FteTA3XC92c3MwcU05VnA2WkZ5WlhEaHZMTXlSdWJYem1MSjdCSjliRkE5T2dSYWpTeGtcL3g0TlA1OG5XWTdsRHk1QlwveHM4akdGMk5oc1BSYklqWHlZY1BBc1QxdFVobW1NQ1EzU3FKMmtcL1FlNm5uN3Z0a2RPbGdOUjU4M0NHMHRPWGY1WGNGNnk4b0pBZVBZb0oyaGZZcnBIZ3B0UXdtQ2VtNG9XZnJkZ3U0SmNlTVpSbkpIdytQY3VhTXpnUG5KdGVlWmZsTVNjNnFuXC9qcTBXcFFON2xmdU9DZjFQUUF1SmRHTzZiSUMrZlpXYXNLVWtEaGtsNWhTcFwvQlhKM1FTYmtDeFRKckVxbXdpNW1rdnpMTW15VnFadndDN1FTZEMrYmNHM3RUdkMyeU1cLzdqUExQb1VIUHNZQ0p6WlRMSHV5cURoUDJEQ3NIaEpLYnpucXgyUjRiK2ZSMTlyT3BJUWE0NU9jZUVlZEhLZ0NBWDdrdEZ0NGFTb1AyZnVWUWlHQ0F6QVI2WXhUbUZnWDczOGpxZk9XNWxWaVIwWmtrVEIrZkhneDBWUWVKNW9mNkNIRnhKMHBJYWtWUCt3SHFaVUdYOFdrenpIS3lRb04zNDZseEp2WFFzMFYxSzVuU0JmWlE2ZjYxQ1M3MEh5c3laV3dEaFFxeitXbE1Lbmk5U1dOSWhLb2hZVnhSRnFWZjBjbng5UW9zZW5pYVZQZlF2VmdzR1IzWG9ObzFxaDdYQmI0RzFyXC9kXC9cL3dKcktRVUY1aUU4M2ZyMWROWThUNTczZyt3aFFyOEJyNFJXUlZ0K3VuWlJLWHFJbmJmSXBrSDVWSW9wb1prNGFlayt6OHpyazhiYTlmN0FHZUVrbWZmRVFHZTRKQ1V6MEZVYnQrcXBEODFCcVRERVVwcDAwR1U4RE00eHAyRkNneVYyWHFtUXBtRkc4eExtRnVQYzZJVVoyenY1cDVHNml6a1pXbWMwZEMrR0VPVEdpSHNib0cwaGJ4d0tXSTVad05vVlZaV3RaYnlCUXlZalwvcVpwcmdPRmlKaGRqV2dtTUV3QVFnZkwrQUswUmhVbTRLWWpUQm16ekRPU0VtcnhrZjc5SnIyOW9mZTFkSWxMNzhqTm5TYXNxT2xOZlRJY3dzU2FnZ3VjRXJDK21SbjFSTlV0MjZxS2htTkdjZXZRTjFmbmlTWHh3blhGdWx4UlE2cVlcLzFJWk53SkRPY1ZmUzF5WHo0UG9rM2w5bzRVMVVuUmtvSGRDK201SjZvTXZQV2JWc0RHclpzcndaYXpDZGlXalhSeTFkVlp1ZGJwSnUwdlljZ2lXTDdxZ1RCK0dwWWhPUFduaVBZR212eEY4bDF0aVM4ait5K2xYRzNDMUtpa1p5djVBYkx3R1JaVU45clEwSzlkUXR4NEIyVVUxVDJvbVNFRVA3STlHeFk1a0NneWNzbHczQmdlZzFFYVFPOURmSFZrNFNTdzJ3V2lIeDRFUlUyK0lPQWpXRGVKOUhCRnhQdEswa1wvcFVlVkwyUzBhVWZ4MFpiUGVBRlwveFwvSVlNWjY1TVR4TDlRU2VNdUo1c0JOMERsRk9ia05INTV3VkNHeklKVTA3TEJwajZnc3grUERocUxHMjVVWDZiV3FPQlUzRVNiZjZERVVFNlhJSzNwUUN0R1wvWU9id01yWWdBYURRWjkxZVZLc2MzMjM1dWIrR0RTV1c2SFU1Q2ZKb1BwN1RiT3NRY0U2c0s4ZTArUnEyK2F6bGZ0dDJWd3pZZ2J2cms5aWR2WTNEVWxKNUVaanZhV29laFMyOHVaTnlWa3Y3bHFVeUZrY05HUW54dmxjN2FzbVNjaUkrS0o1NXJCcmxQNFkwMWp2andLb3RxRzlaMXp2QXRKSWJXcytcL0NEYWVjR1c2N2hOTWxSZFZ6UkZ2UmlyeEcrQ3JzMkpVZUtzUjBnZDArWGdQRUpyQzgyZzE5VUNaMDJlUlpYVm5kTVhYRDNHMno2eW9Ob3A2dW5lZGRuWkJzcGhWWm5qajRCUHQ1RXNIM1RaRU85MTNuZTNOVCtSVTEzRXE3U0xXTHlmYWtwUFwvOVpZZG1lSHVPMm91VlJxWSs5ZU01S0Vnc0I2YjcyRmhpNFdjWnJOSVJjZHQ0dkN1ZTFMWDczeUZjUnBQUGdmUkhjYnJnPT0iLCJpYWkiOiJlYzdlNzRhNDM0MjI0YjJlOWYwMGRhZDk5ZGZiZWZhMiIsInNlYyI6MSwiY2xzdmMiOiJwcm9kLWZuIiwibHB2IjoxNjc5NzYzMTU1LCJ0IjoicyIsImljIjp0cnVlLCJleHAiOjE2Nzk3NzAzNTUsImlhdCI6MTY3OTc2MzE1NSwianRpIjoiOGY0NGM4NGIyNjIwNDg2Yzk5NDBhMTc5NmVkYmFhOGMifQ.AIUfO2iaZfg39B5jg3Hdy_zGc90vPq2yVY9OcWaIn63UWyzEgl0t6tyNXLEihr58IxJeC2LQR3t6My2ZRcWf0K4G
func main() {
	c := &http.Client{}
	//res := PartyLookupPending(c, "ec7e74a434224b2e9f00dad99dfbefa2")
	//
	//fmt.Println(res)
	// res := PartySendIntention(c, &SendIntentionData{UserId: "e5076458b28a4fa290d9993dc9e35574", ClientId: "ec7e74a434224b2e9f00dad99dfbefa2"})
	// resp := PartyLookupPing(c, "e5076458b28a4fa290d9993dc9e35574", "ec7e74a434224b2e9f00dad99dfbefa2")
	// fmt.Println(resp)
	// fmt.Println(res)
	header := http.Header{}
	header.Add("Sec-WebSocket-Protocol", "xmpp")

	con, _, err := websocket.DefaultDialer.Dial("wss://xmpp-service-prod.ol.epicgames.com", header)
	if err != nil {
		res := fmt.Errorf("failed to connect to xmpp service: %v", err)
		fmt.Println(res)
	}

	auth := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("\u0000%s\u0000%s", accountId,
		token,
	)))
	client, err := NewClient(con, auth)
	if err != nil {
		res := fmt.Errorf("failed to create new xmpp client: %v: %w", err, asynq.SkipRetry)
		fmt.Println(res)
	}
	fmt.Println("Program started.")
	friendsList, friendErr := client.Friends_Get_All()
	if friendErr != nil {
		fmt.Println("err")
	}
	client.Friends_Set_Nickname("Floppa", friendsList[0].AccountId)
	client.OnPing(
		func(p *PartyPing) {
			fmt.Println("New Ping from: " + p.PingerDN)
			res := PartyLookupPing(c, p.PingerID, accountId)
			fmt.Printf("Party id: %s\n", res.ID)
			fmt.Printf("jid: %s", p.Message.To)
			resp := PartySendJoinRequest(c, p.Message.To, res.ID)
			fmt.Println("resp:" + resp.ID)
		},
	)
	client.OnJoin(func(j *PartyJoin) {
		fmt.Println("New Join: " + j.AccountDN)
		partyLookup, err := client.PartyLookup(client.Party.Id)
		if err != nil {
			fmt.Println(partyLookup)
		}
		response := PartyUpdateMemberMeta(c, j.PartyID, metaData)
		fmt.Println(response)
		keyRes := SetReadiness(c, j.PartyID, true)
		fmt.Println(keyRes)
	})
	client.OnMemberLeft(
		func(m *PartyMemberLeft) {
			fmt.Println("Member left: " + m.AccountID)
		},
	)
	client.OnSkinChanged(func(s string, id string) {
		fmt.Println("new Skin: " + s)
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
					res := PartyLookupPing(c, body.PingerID, accountId)
					fmt.Printf("Party id: %s\n", res.ID)
					fmt.Printf("jid: %s", message.To)
					resp := PartySendJoinRequest(c, message.To, res.ID)
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
				case "com.epicgames.social.interactions.notification.v2":
					fmt.Println("Social interaction notification.")
				case "com.epicgames.social.party.notification.v0.PARTY_UPDATED":
					fmt.Printf("Party updated from: %s, revision: %d", body.AccountDN, body.Revision)
					if body.Revision != 0 {
						client.Party.PartyRevision = body.Revision
					}
				default:
					fmt.Printf("message from type: %s\n", body.Type)
					fmt.Printf("message was: %s", msg)
				}

			}
		}
	}

}
