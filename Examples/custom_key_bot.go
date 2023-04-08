package main

import (
	"fmt"
	"time"

	fortnitego "github.com/oxelf/fortnite-go"
)

var accountId = "49748cfe40334a129ca3620252abfba2"

func main() {
	//The config for our client.
	config := fortnitego.ClientConfig{
		XMPP:       true,
		Connection: "Fortnite",
		AuthClient: fortnitego.AuthClients.Fortnite_PC_Client}
	url := fortnitego.GetAuthCodeUrl(fortnitego.AuthClients.Fortnite_PC_Client)
	fmt.Println(url)
	val, eerr := fortnitego.Get_Token_By_AuthCode("250128179f084dcf95de74552f54f5e4", fortnitego.Base64AuthClients.Fortnite_PC_Client, true)
	if eerr != nil {
		fmt.Println(eerr.EpicErrorMessage)
	} else {
		fmt.Println("Acc token: " + val.AccessToken)
		config.Token = val.AccessToken
		config.AccountID = val.AccountID
		config.DisplayName = val.DisplayName
		fortnitego.Create_DeviceAuth(val.AccessToken, val.AccountID)
	}
	client, cerr := fortnitego.NewClient(config)
	if cerr != nil {
		fmt.Printf("failed to create new xmpp client: %v\n", cerr)
	} else {
		fmt.Println("Program started.")
		res, ResErr := client.LightSwitch_Status_Fortnite()
		if ResErr != nil {
			fmt.Println(ResErr)
		} else {
			fmt.Println(res)
		}
		client.OnFriendRequest(func(fr *fortnitego.FriendshipRequest) {
			client.Friend_Add_Or_Accept(fr.From)
		})
		client.OnPing(
			func(p *fortnitego.PartyPing) {
				fmt.Println("New Ping from: " + p.PingerDN)
				res := client.PartyLookupPing(p.PingerID, client.Config.AccountID)
				fmt.Printf("Party id: %s\n", res.ID)
				fmt.Printf("jid: %s", p.Message.To)
				resp, err := client.PartySendJoinRequest(p.Message.To, res.ID)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(resp.Status + "\n")
				}
			},
		)
		client.OnJoin(func(j *fortnitego.PartyJoin) {
			fmt.Println("New Join: " + j.AccountDN)
			response := client.Set_Skin("CID_029_Athena_Commando_F_Halloween")
			fmt.Println(response)
			client.CancelEmote()
			keyRes := client.SetReadiness(false)
			fmt.Println(keyRes)
			time.Sleep(time.Millisecond * 1000)
			err := client.SetEmote("EID_Worm")
			if err != nil {
				fmt.Println(err)
			}
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
			err := client.SetEmote("EID_Worm")
			if err != nil {
				fmt.Println(err)
			}
			setPlaylist := client.SetPlaylist()
			fmt.Println(setPlaylist)
			setKey := client.SetCustomKey("1234")
			fmt.Println(setKey)
			client.CancelEmote()
			client.SetEmote("EID_Worm")
		})
	}
	//call ```client.Listen()```to start listening to messages and receiving callbacks. This is a blocking operation, call it at the end of your func.
	client.Listen()
}
