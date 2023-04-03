# fortnite-go
A Package for using the http and xmpp features of fortnite.
Inspired by fortnitepy and fnbr.js

* Installation
  `go get github.com/oxelf/fortnite-go`

* Usage Example
  ```golang 
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
        //On Friendship Request, add the person.
        client.OnFriendRequest(func(fr *fortnitego.FriendshipRequest) {
                client.Friends_Add_Or_Accept(fr.From)
            })
        //call `client.Listen()` to start listening to messages and receiving callbacks. This is a blocking operation, call it at the end of your func.
        client.Listen()
    }


<h2>Links<h2>
  - [Discord] (https://discord.gg/sDt6Zk7bbg)
  - [Docs] (https://pkg.go.dev/github.com/oxelf/fortnite-go)