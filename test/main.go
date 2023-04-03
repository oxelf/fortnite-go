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
		AuthClient: fortnitego.AuthClients.Fortnite_IOS_Client}
	oauthResponse, eerr := fortnitego.Get_OauthToken_By_AuthCode("3eb5b218305047dab471ec5fa266e516", fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
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
	sessionId := "8c9e7d3608914b0897dd7fa76406eda9"
	fileLocation, _ := client.FileLocationInfo(fmt.Sprint("https://datastorage-public-service-live.ol.epicgames.com/api/v1/access/fnreplaysmetadata/public%2F" + sessionId + ".json"))
	fileUrl := ""
	for _, value := range fileLocation.Files {
		fileUrl = value.ReadLink
		break
	}
	replay, _ := fortnitego.DownloadReplayCDNFile(fileUrl)
	fmt.Sprint(replay)
	//Handling callbacks.
	// client.OnPresence(func(p *fortnitego.Status) {
	// 	if p.SessionID != "" {
	// 		fileLocation, _ := client.FileLocationInfo(fmt.Sprint("https://datastorage-public-service-live.ol.epicgames.com/api/v1/access/fnreplaysmetadata/public%2F" + p.SessionID + ".json"))
	// 		fileUrl := ""
	// 		for _, value := range fileLocation.Files {
	// 			fileUrl = value.ReadLink // return the first key-value pair
	// 			break                    // exit the loop after the first iteration
	// 		}
	// 		fortnitego.DownloadReplayCDNFile(fileUrl)
	// 	}

	// })

	client.Listen()
}
