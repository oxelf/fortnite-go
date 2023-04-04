package main

import (
	"fmt"

	fortnitego "github.com/oxelf/fortnite-go"
)

func cmain() {
	//The config for our client.
	config := fortnitego.ClientConfig{
		XMPP:       true,
		Connection: "Fortnite",
		AuthClient: fortnitego.AuthClients.Fortnite_PC_Client}
	oauthResponse, eerr := fortnitego.Get_Token_By_AuthCode("41e5491e8944418b86dfbf39c2fba32f", fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
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
	client.OnPing(func(p *fortnitego.PartyPing) {
		fmt.Println("New Ping from: " + p.PingerDN)
	})
	client.OnJoin(func(j *fortnitego.PartyJoin) {
		fmt.Println("New Join: " + j.AccountDN)
	})
	client.OnMemberLeft(func(m *fortnitego.PartyMemberLeft) {
		fmt.Println("Member left: " + m.AccountID)
	})
	client.OnSkinChanged(func(skin string, userid string) {
		if userid != client.Config.AccountID {
			fmt.Printf("User: %s, new Skin: %s", userid, skin)
		} else {
			fmt.Println("client new skin:" + skin)
		}
	})
	client.OnNewCaptain(func(c *fortnitego.PartyNewCaptain) {
		fmt.Println("new captain: " + c.AccountDN)
	})
	client.OnBlocklistUpdate(func(bu *fortnitego.BlocklistUpdate) {
		fmt.Printf("Blocklist updated. User: %s, New Status: %s\n", bu.AccountID, bu.Status)
	})
	client.OnFriendRequest(func(fr *fortnitego.FriendshipRequest) {
		fmt.Printf("Friendrequest message from: %s, Status: %s\n", fr.From, fr.Status)
	})
	client.OnMemberDisconnected(func(pmd *fortnitego.PartyMemberDisconnected) {
		fmt.Println("member disconnected")
	})
	client.OnMemberExpired(func(pme *fortnitego.PartyMemberExpired) {
		fmt.Println("member expired")
	})
	client.OnMemberKicked(func(pmk *fortnitego.PartyMemberKicked) {
		fmt.Println("Member kicked")
	})
	client.OnMemberRequireConfirmation(func(pmrc *fortnitego.PartyMemberRequireConfirmation) {
		fmt.Println("member needs confirmation")
	})
	client.OnPresence(func(s *fortnitego.Presence) {
		fmt.Println("new presence from: " + s.From)
	})
	//call ```client.Listen()```to start listening to messages and receiving callbacks. This is a blocking operation, call it at the end of your func.
	client.Listen()
}
