package main

import (
	"fmt"

	fortnitego "github.com/oxelf/fortnite-go"
)

func main() {
	//The config for our client.
	config := fortnitego.ClientConfig{
		XMPP:       true,
		Connection: "Fortnite",
		AuthClient: fortnitego.AuthClients.Fortnite_PC_Client}
	oauthResponse, eerr := fortnitego.Get_OauthToken_By_AuthCode("41e5491e8944418b86dfbf39c2fba32f", fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	if eerr != nil {
		fmt.Println(eerr.EpicErrorMessage)
	}
	//Setting the Token and AccountID for the client config.
	config.Token = oauthResponse.AccessToken
	config.AccountID = oauthResponse.AccountID

	//creating the client.
	client, cerr := fortnitego.NewClient(config)
	if cerr != nil {
		fmt.Printf("failed to create new xmpp client: %v\n", cerr)
	}
	fmt.Println("Program started.")
	//Handling callbacks.
	client.OnPing(
		func(p *fortnitego.PartyPing) {
			fmt.Println("New Ping from: " + p.PingerDN)
			//Looking up the ping for the partyid and our JID.
			res := client.PartyLookupPing(p.PingerID, client.Config.AccountID)
			//Sending a Join Request. If we received an ping before, we need to send a join request to join.
			client.PartySendJoinRequest(p.Message.To, res.ID)
		},
	)
	client.OnJoin(func(j *fortnitego.PartyJoin) {
		fmt.Println("New Join: " + j.AccountDN)
		if j.AccountID == client.Config.AccountID {
			//We set the skin, because if we dont set it, our client wont be shown as player in the party.
			client.Set_Skin("CID_029_Athena_Commando_F_Halloween")
			//Just to do something, we make the bot ready up.
			client.SetReadiness(j.PartyID, true)
		}

	})
	client.OnMemberLeft(
		func(m *fortnitego.PartyMemberLeft) {
			fmt.Println("Member left: " + m.AccountID)
		},
	)
	client.OnSkinChanged(func(skin string, id string) {
		//Just to demonstrate, if someone changes skin, and the accountid is not the id of our bot, We change our skin to the same skin as the other person.
		if id != client.Config.AccountID {
			fmt.Println("new Skin: " + skin)
			//Set the skin that the other user has set.
			client.Set_Skin(skin)
		} else {
			fmt.Println("client new skin:" + skin)
		}
	})
	client.OnNewCaptain(func(c *fortnitego.PartyNewCaptain) {
		fmt.Println("new captain: " + c.AccountDN)
		//check if our bot is the new captain
		if c.AccountID == client.Config.AccountID {
			setPlaylist := client.SetPlaylist()
			fmt.Println(setPlaylist)
			setKey := client.SetCustomKey("1234")
			fmt.Println(setKey)
			client.PartyLeave()
		}
	})
	//call ```client.Listen()```to start listening to messages and receiving callbacks. This is a blocking operation, call it at the end of your func.
	client.Listen()
}
