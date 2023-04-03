package fortnitego

import (
	b64 "encoding/base64"
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
	conn                      *websocket.Conn
	id                        string
	Party                     Party
	c                         http.Client
	JID                       string
	Config                    ClientConfig
	pingCallback              func(*PartyPing)
	joinCallback              func(*PartyJoin)
	NewCaptainCallback        func(*PartyNewCaptain)
	memberLeftCallback        func(*PartyMemberLeft)
	friendshipRequestCallback func(*FriendshipRequest)
	blocklistUpdateCallback   func(*BlocklistUpdate)
	kickedCallback            func(*PartyMemberKicked)
	expiredCallback           func(*PartyMemberExpired)
	disconnectedCallback      func(*PartyMemberDisconnected)
	needConfirmationCallback  func(*PartyMemberRequireConfirmation)
	skinChangedCallback       func(string, string)
	presenceCallback          func(*Presence)
}
type Party struct {
	Id             string
	ClientRevision int
	PartyRevision  int
	CaptainID      string
	Members        []PartyUser
}

type ClientConfig struct {
	//whetever the client should connect via xmpp or not. Cant be changed after initializing the client.
	XMPP bool
	//Possible values: "Fortnite", "Launcher". The Type of XMPP Connection you want to use. Launcher is just listening to presences.
	Connection string
	//Choose one of the pre defined types from fortnitego.Tokentypes . Take a look at that struct for more information on which types are possible. I recommend Fortnite-Pc for Fortnite XMPP Connection and Fortnite-IOS for Launcher Connection.
	AuthClient string
	//Your access token. We prefer eg1 token, because some servives might not accept other tokens than eg1
	Token string
	//We need your account Id to open a connection with XMPP.
	AccountID string
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.AccountID == "" {
		return nil, fmt.Errorf("please provide an AccountId")
	}
	if config.Token == "" {
		return nil, fmt.Errorf("please provide a authentication token, best case an eg1")
	}
	if config.XMPP {
		if config.Connection == "" {
			return nil, fmt.Errorf("please provide the type of xmpp connection you want to use. Either Fortnite or Launcher")
		}
	}
	if config.AuthClient == "" {
		config.AuthClient = AuthClients.Fortnite_IOS_Client
	}
	client := &Client{Config: config}
	client.c = http.Client{}
	if !client.Config.XMPP {
		return client, nil
	}
	header := http.Header{}
	header.Add("Sec-WebSocket-Protocol", "xmpp")

	con, _, err := websocket.DefaultDialer.Dial("wss://xmpp-service-prod.ol.epicgames.com", header)
	if err != nil {
		res := fmt.Errorf("failed to connect to xmpp service: %v", err)
		fmt.Println(res)
	}
	client.conn = con
	auth := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("\u0000%s\u0000%s", config.AccountID,
		config.Token,
	)))

	err = client.open(auth)
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
	err = client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<iq id=\"_xmpp_bind1\" type=\"set\"><bind xmlns=\"urn:ietf:params:xml:ns:xmpp-bind\"><resource>V2:%s:WIN::%s</resource></bind></iq>", client.Config.Connection, id)))
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	_, err = client.expectResult("<iq")
	if err != nil {
		return err
	}
	PresenceError := client.SendPresence("Battle Royale Lobby - 1/16")
	if PresenceError != nil {
		return err
	}

	return nil
}

func (client *Client) SendPresence(Status string) error {
	stamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	err := client.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<presence><status>{\"Status\":\"%s\",\"bIsPlaying\":true,\"bIsJoinable\":true,\"bHasVoiceSupport\":false,\"ProductName\":\"Fortnite\",\"SessionId\":\"\",\"Properties\":{\"OverrideAppId_s\":\"Fortnite\",\"FortPartySize_i\":1,\"FortSubGame_i\":1,\"InUnjoinableMatch_b\":false}}</status><delay stamp=\"%s\" xmlns=\"urn:xmpp:delay\"/></presence>", Status, stamp)))
	if err != nil {
		return err
	} else {
		return nil
	}
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
		if result == "<failure xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\"><not-authorized/><text xml:lang='en'>Password not verified</text></failure>" {
			return "", fmt.Errorf("an error with the provided token or account id occured. check if they are correct, and if they are, try using an eg1 token")
		}
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
				if strings.Contains(presence.From, c.Config.AccountID) {
					c.JID = presence.From
				}
				fmt.Printf("<Presence>: type:%s,\n from: %#v,\n status: %#v", presence.Type, presence.From, status)
				if c != nil {
					presence.MStatus = *status
					c.presenceCallback(presence)
					continue
				}
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
					continue
				}
				switch body.Type {
				case "USER_BLOCKLIST_UPDATE":
					var blocklistUpdate BlocklistUpdate
					err := json.Unmarshal(message.Body.RawJSON, &blocklistUpdate)
					if err != nil {
						continue
					}
					blocklistUpdate.RawMessage = message.Body.RawJSON
					c.blocklistUpdateCallback(&blocklistUpdate)
				case "com.epicgames.social.party.notification.v0.PING":
					var party_ping PartyPing
					err := json.Unmarshal(message.Body.RawJSON, &party_ping)
					if err != nil {
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
					if party_join.AccountID == c.Config.AccountID {
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
				case "com.epicgames.social.party.notification.v0.MEMBER_KICKED":
					var partyMemberKicked PartyMemberKicked
					partyMemberKicked.RawMessage = message.Body.RawJSON
					c.kickedCallback(&partyMemberKicked)
				case "com.epicgames.social.party.notification.v0.MEMBER_EXPIRED":
					var partyMemberExpired PartyMemberExpired
					partyMemberExpired.RawMessage = message.Body.RawJSON
					c.expiredCallback(&partyMemberExpired)
				case "com.epicgames.social.party.notification.v0.MEMBER_DISCONNECTED":
					var partyMemberDisconnected PartyMemberDisconnected
					partyMemberDisconnected.RawMessage = message.Body.RawJSON
					c.disconnectedCallback(&partyMemberDisconnected)
				case "com.epicgames.social.party.notification.v0.MEMBER_REQUIRE_CONFIRMATION":
					var partyMemberRequireConfirmation PartyMemberRequireConfirmation
					c.needConfirmationCallback(&partyMemberRequireConfirmation)
				case "com.epicgames.social.party.notification.v0.MEMBER_LEFT":
					var party_member_left PartyMemberLeft
					err := json.Unmarshal(message.Body.RawJSON, &party_member_left)
					if err != nil {
						fmt.Println("unmarshalling error.")
						continue
					}
					party_member_left.RawMessage = message.Body.RawJSON
					c.memberLeftCallback(&party_member_left)
				case "FRIENDSHIP_REQUEST":
					var friendshipRequest FriendshipRequest
					err := json.Unmarshal(message.Body.RawJSON, &friendshipRequest)
					if err != nil {
						fmt.Println("unmarshalling error.")
						continue
					}
					friendshipRequest.RawMessage = message.Body.RawJSON
					c.friendshipRequestCallback(&friendshipRequest)
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

func (c *Client) OnFriendRequest(callback func(*FriendshipRequest)) {
	c.friendshipRequestCallback = callback
}
func (c *Client) OnBlocklistUpdate(callback func(*BlocklistUpdate)) {
	c.blocklistUpdateCallback = callback
}
func (c *Client) OnPresence(callback func(*Presence)) {
	c.presenceCallback = callback
}
func (c *Client) OnMemberKicked(callback func(*PartyMemberKicked)) {
	c.kickedCallback = callback
}
func (c *Client) OnMemberExpired(callback func(*PartyMemberExpired)) {
	c.expiredCallback = callback
}
func (c *Client) OnMemberDisconnected(callback func(*PartyMemberDisconnected)) {
	c.disconnectedCallback = callback
}
func (c *Client) OnMemberRequireConfirmation(callback func(*PartyMemberRequireConfirmation)) {
	c.needConfirmationCallback = callback
}
func (c *Client) OnSkinChanged(callback func(SkinID string, accountID string)) {
	c.skinChangedCallback = callback
}
