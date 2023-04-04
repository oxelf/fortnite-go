package main

import (
	"fmt"
	"time"

	fortnitego "github.com/oxelf/fortnite-go"
)

func main() {
	// token, err := fortnitego.Get_Token_By_ClientCredentials(fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("token from client credentials: " + token.AccessToken)
	// return
	// oauthResponse, eerr := fortnitego.Get_Token_By_AuthCode("745ef83005dc42198cb51c6c2648898b", fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	// if eerr != nil {
	// 	fmt.Println(eerr.EpicErrorMessage)
	// 	return
	// }
	// token, err := fortnitego.Get_Token_By_Refresh(oauthResponse.RefreshToken, fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	token, err := fortnitego.Get_Token_By_ClientCredentials(fortnitego.Base64AuthClients.Fortnite_SWITCH_Client, true)
	fmt.Println("token from client credentials: " + token.AccessToken)
	deviceCode, err := fortnitego.Get_DeviceCode(token.AccessToken)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(deviceCode.VerificationUriCompleted)
	time.Sleep(time.Second * 20)
	token, err = fortnitego.Get_Token_By_DeviceCode(deviceCode.DeviceCode, fortnitego.Base64AuthClients.Fortnite_SWITCH_Client, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("token from device code: " + token.AccessToken)
	return
	deviceAuth, err := fortnitego.Create_DeviceAuth(token.AccessToken, token.AccountID)
	if err != nil {
		fmt.Println(err)
		return
	}
	token, err = fortnitego.Get_Token_By_DeviceAuth(*deviceAuth, fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("token by device auth: " + token.AccessToken)
	exchangecode, err := fortnitego.Get_ExchangeToken(token.AccessToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	token, err = fortnitego.Get_Token_By_Exchange(exchangecode.Token, fortnitego.Base64AuthClients.Fortnite_IOS_Client, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("token from exchange: " + token.AccessToken)
}
