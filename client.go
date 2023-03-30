package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Client struct {
	conn                *websocket.Conn
	id                  string
	Party               Party
	c                   http.Client
	JID                 string
	pingCallback        func(*PartyPing)
	joinCallback        func(*PartyJoin)
	NewCaptainCallback  func(*PartyNewCaptain)
	memberLeftCallback  func(*PartyMemberLeft)
	skinChangedCallback func(string, string)
}
type Party struct {
	Id             string
	ClientRevision int
	PartyRevision  int
	CaptainID      string
	Members        []PartyUser
}

func NewClient(conn *websocket.Conn, auth string) (*Client, error) {
	client := &Client{conn: conn}
	client.c = http.Client{}
	err := client.open(auth)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err = client.sendPing()
				if err != nil {
					_ = client.conn.Close()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return client, nil
}

func (client *Client) open(auth string) error {
	err := client.conn.WriteMessage(websocket.TextMessage, []byte("<open xmlns=\"urn:ietf:params:xml:ns:xmpp-framing\" to=\"prod.ol.epicgames.com\" version=\"1.0\" />"))
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	openResponse, err := client.expectResult("<open")
	if err != nil {
		return err
	}

	open := &Open{}
	err = xml.Unmarshal([]byte(openResponse), open)
	if err != nil {
		return err
	}
	client.id = open.ID

	_, err = client.expectResult("<stream:features")
	if err != nil {
		return err
	}

	err = client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<auth mechanism=\"PLAIN\" xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\">%s</auth>", auth)))
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	_, err = client.expectResult("<success")
	if err != nil {
		return err
	}

	err = client.conn.WriteMessage(websocket.TextMessage, []byte("<open xmlns=\"urn:ietf:params:xml:ns:xmpp-framing\" to=\"prod.ol.epicgames.com\" version=\"1.0\" />"))
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	_, err = client.expectResult("<open")
	if err != nil {
		return err
	}
	_, err = client.expectResult("<stream:features")
	if err != nil {
		return err
	}

	uid := uuid.New()
	id := strings.ReplaceAll(uid.String(), "-", "")
	id = strings.ToUpper(id)
	err = client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<iq id=\"_xmpp_bind1\" type=\"set\"><bind xmlns=\"urn:ietf:params:xml:ns:xmpp-bind\"><resource>V2:Fortnite:WIN::%s</resource></bind></iq>", id)))
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	_, err = client.expectResult("<iq")
	if err != nil {
		return err
	}
	stamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	fmt.Print("Sended Presence with stamp: " + stamp)
	err = client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<presence><status>{\"Status\":\"Moin Leander\",\"bIsPlaying\":true,\"bIsJoinable\":true,\"bHasVoiceSupport\":false,\"ProductName\":\"Fortnite\",\"SessionId\":\"\",\"Properties\":{\"OverrideAppId_s\":\"Fortnite\",\"FortPartySize_i\":1,\"FortSubGame_i\":1,\"InUnjoinableMatch_b\":false}}</status><delay stamp=\"%s\" xmlns=\"urn:xmpp:delay\"/></presence>", stamp)))
	if err != nil {
		return err
	}
	fmt.Println("no error sending presence")
	return nil
}

func (client *Client) Read() (string, error) {
	t, m, err := client.conn.ReadMessage()
	if err != nil {
		return "", err
	}

	message := string(m)
	if t != websocket.TextMessage {
		message, err = client.Read()
		if err != nil {
			return "", err
		}
	}

	return message, nil
}

func (client *Client) expectResult(expect string) (string, error) {
	result, err := client.Read()
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(result, expect) {
		return "", errors.New(fmt.Sprintf("expected %s, but got %s", expect, result))
	}

	return result, nil
}

func (client *Client) sendPing() error {
	//<iq xmlns="jabber:client" id="acbeabf8-b04b-4e94-a044-6d6b8f04514e" type="get"><ping xmlns="urn:xmpp:ping"/></iq>
	ping := Iq{
		Xmlns: "jabber:client",
		ID:    client.id,
		Type:  "get",
		Ping: Ping{
			Xmlns: "urn:xmpp:ping",
		},
	}

	data, err := xml.Marshal(ping)
	if err != nil {
		return err
	}

	err = client.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}
func (client *Client) Disconnect() {
	_ = client.conn.Close()
}
func (c *Client) Listen() {
	for {
		if msg, err := c.Read(); err != nil {
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
				if strings.Contains(presence.From, accountId) {
					c.JID = presence.From
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
					continue
				}
				switch body.Type {
				case "com.epicgames.social.party.notification.v0.PING":
					var party_ping PartyPing
					err := json.Unmarshal(message.Body.RawJSON, &party_ping)
					if err != nil {
						fmt.Println("unmarshalling error.")
						continue
					}
					party_ping.Message = *message
					c.pingCallback(&party_ping)

				case "com.epicgames.social.party.notification.v0.MEMBER_JOINED":

					jsonString := fmt.Sprintf("%s", message.Body.RawJSON)

					var party_join PartyJoin
					err := json.Unmarshal([]byte(jsonString), &party_join)
					if err != nil {
						fmt.Println(message.Body.RawJSON)
						fmt.Println("unmarshalling error.")
						continue
					}
					if party_join.AccountID == accountId {
						if party_join.PartyID != "" {
							c.Party.Id = party_join.PartyID
						}
						c.Party.Members = append(c.Party.Members, PartyUser{
							JoinedAt:  party_join.JoinedAt,
							UpdatedAt: party_join.UpdatedAt,
							AccountId: party_join.AccountID,
							AccountDN: party_join.AccountDN,
							JID:       party_join.Connection.ID,
							Revision:  party_join.Revision,
						})
					} else {
						if party_join.PartyID != "" {
							c.Party.Id = party_join.PartyID
						}
						c.Party.Members = append(c.Party.Members, PartyUser{
							JoinedAt:  party_join.JoinedAt,
							UpdatedAt: party_join.UpdatedAt,
							AccountId: party_join.AccountID,
							AccountDN: party_join.AccountDN,
							JID:       party_join.Connection.ID,
							Revision:  party_join.Revision,
						})
					}
					c.PartyUpdateMemberMeta(metaData)
					c.joinCallback(&party_join)
				case "com.epicgames.social.party.notification.v0.MEMBER_NEW_CAPTAIN":
					var party_new_captain PartyNewCaptain
					err := json.Unmarshal(message.Body.RawJSON, &party_new_captain)
					if err != nil {
						fmt.Println("unmarshalling error.")
						continue
					}
					if party_new_captain.AccountID != "" {
						c.Party.CaptainID = party_new_captain.AccountID
						c.NewCaptainCallback(&party_new_captain)
					} else {
					}
				case "com.epicgames.social.interactions.notification.v2":
					fmt.Println("Social interaction notification.")
					continue
				case "com.epicgames.social.party.notification.v0.PARTY_UPDATED":
					fmt.Printf("Party updated from: %s, revision: %d", body.AccountDN, body.Revision)
					if body.Revision != 0 {
						c.Party.PartyRevision = body.Revision
					}
					continue
				case "com.epicgames.social.party.notification.v0.MEMBER_LEFT":
					var party_member_left PartyMemberLeft
					err := json.Unmarshal(message.Body.RawJSON, &party_member_left)
					if err != nil {
						fmt.Println("unmarshalling error.")
						continue
					}
					c.memberLeftCallback(&party_member_left)
				case "com.epicgames.social.party.notification.v0.MEMBER_STATE_UPDATED":
					var party_member_updated PartyMemberUpdated
					err := json.Unmarshal(message.Body.RawJSON, &party_member_updated)
					if err != nil {
						fmt.Println("unmarshalling error.")
						continue
					}
					if party_member_updated.MemberStateUpdated != nil {
						cosmeticLoadout := UpdatedCosmeticLoadout{}
						fmt.Println(party_member_updated.MemberStateUpdated["Default:AthenaCosmeticLoadout_j"])
						if party_member_updated.MemberStateUpdated["Default:AthenaCosmeticLoadout_j"] == "" {
							fmt.Println("cosmetic loadout empty")
							continue
						} else {

							err := json.Unmarshal([]byte(party_member_updated.MemberStateUpdated["Default:AthenaCosmeticLoadout_j"]), &cosmeticLoadout)
							if err != nil {
								fmt.Println(err)
								continue
							}
							if cosmeticLoadout.AthenaCosmeticLoadout.Character == "None" || cosmeticLoadout.AthenaCosmeticLoadout.Character == "" {

							} else {
								substr := ""
								lastDotIndex := strings.LastIndex(cosmeticLoadout.AthenaCosmeticLoadout.Character, ".")
								if lastDotIndex != -1 && lastDotIndex < len(cosmeticLoadout.AthenaCosmeticLoadout.Character)-1 {
									substr = cosmeticLoadout.AthenaCosmeticLoadout.Character[lastDotIndex+1:]
								}
								c.skinChangedCallback(substr, party_member_updated.AccountID)
							}
						}
					}
				default:
					fmt.Printf("message from type: %s\n", body.Type)
					fmt.Printf("message was: %s", msg)
				}

			}
		}
	}
}
func (c *Client) OnPing(callback func(p *PartyPing)) {
	c.pingCallback = callback
}
func (c *Client) OnJoin(callback func(j *PartyJoin)) {
	c.joinCallback = callback
}
func (c *Client) OnNewCaptain(callback func(c *PartyNewCaptain)) {
	c.NewCaptainCallback = callback
}
func (c *Client) OnMemberLeft(callback func(*PartyMemberLeft)) {
	c.memberLeftCallback = callback
}
func (c *Client) OnSkinChanged(callback func(SkinID string, accountID string)) {
	c.skinChangedCallback = callback
}
