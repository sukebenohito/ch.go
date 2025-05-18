package main

import (
	//"fmt"
	"log"
	"net/http"
	"net/url"
	//"regexp"
	"strings"
	"github.com/gorilla/websocket"
)

//Room classs
type PrivateMessage struct {
	Name string
	Server string
	Port string
	Host string
	FirstCommand bool
	Connected bool
	Ws *websocket.Conn
	Mgr *Chatango
}


func NewPrivateMessage(c *Chatango) *PrivateMessage {
	return &PrivateMessage{
		Mgr: c,
		Name: "PrivateMessage",
		FirstCommand: true,
		Server: "c1.chatango.com",
		Port: "8080",
	}
}

func (p *PrivateMessage) SendCommand(args ...string) {
	terminator := ""
	if p.FirstCommand {
		terminator = "\x00"
		p.FirstCommand = false
	} else {
		terminator = "\r\n\x00"
	}

	command := strings.Join(args, ":") + terminator

	err := p.Ws.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("err write:", err)
		return
	}
}

func (p *PrivateMessage) Auth() {

	if p.Mgr.UserName != "" && p.Mgr.Password != "" {
		endpoint := "https://chatango.com/login"

		params := url.Values{}
		params.Add("user_id", p.Mgr.UserName)
		params.Add("password", p.Mgr.Password)
		params.Add("storecookie", "on")
		params.Add("checkerrors", "yes")

		reqURL := endpoint + "?" + params.Encode()
		resp, err := http.Get(reqURL)
		if err != nil {
			log.Println("Request failed:", err)
			p.Disconnect()
			return
		}
		defer resp.Body.Close()

		// Read headers to find auth cookie
		cookies := resp.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "auth.chatango.com" {
				token := cookie.Value
				p.SendCommand("tlogin", token, "2")
				return
			}
		}

		// If cookie not found, maybe try parsing manually from headers
		//rawHeaders := resp.Header["Set-Cookie"]
		//for _, header := range rawHeaders {
		//	re := regexp.MustCompile(`auth\.chatango\.com\s*=\s*([^;]*)`)
		//	match := re.FindStringSubmatch(header)
		//	if len(match) > 1 {
		//		token := match[1]
		//		p.SendCommand("tlogin", token, "2")
		//		return
		//	}
		//}

		// No token found
		p.Disconnect()
	}
}


func(p *PrivateMessage) Connect(){
	p.Host = p.Server + ":" + p.Port
	u := url.URL{Scheme: "ws", Host: p.Host, Path: "/"}
	header := http.Header{}
	header.Add("Origin", "https://st.chatango.com")
	p.Connected = true
	log.Printf("connecting to %s", p.Host)
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	p.Ws = ws
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer p.Ws.Close()

	p.Auth()

	for {
		_, message, err := p.Ws.ReadMessage()
		if err != nil {
			p.Connected  = false
			log.Printf("%s: %s \n", p.Name, err)
			return
		}
		p.Feed(string(message))
	}
}

func (p *PrivateMessage) Ping(){
	p.SendCommand("")
}

func (p *PrivateMessage) Disconnect(){
	p.Ws.Close()
}

func (p *PrivateMessage) Message(username string, message string){
	p.SendCommand("msg", strings.ToLower(username), message)
}

func (p *PrivateMessage) Feed(food string) {
	food = strings.ReplaceAll(food, "\r\n\u0000", "")
	parts := strings.Split(food, ":")
	cmd := "Rcmd_" + parts[0]
	args := parts[1:]
	//log.Printf("%s: %s", cmd, args)
	switch cmd {
		case "Rcmd_msg":
			p.Rcmd_msg(args)
		default :
			//fmt.Printf("%s: %s \n", cmd, args)
	}
}

func (p *PrivateMessage) Rcmd_msg(args []string){
	user := NewUser(args[0])
	//sub := strings.Join(args[5], ":")
	text := _strip_html(args[5])
	p.Mgr.PMessage(user, p, text)
}
