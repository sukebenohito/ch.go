//author agunq.e@gmail.com

package main

import (

	"log"
	"net/http"
	"net/url"
	"time"
	"strings"
	"net"
	"github.com/gorilla/websocket"
)

//Room classs
type Room struct {
	Name string
	Uid string
	Server string
	Port string
	Host string
	Channel string
	FirstCommand bool
	Connected bool
	Ws *websocket.Conn
	Mgr *Chatango
}

func NewRoom(name string, c *Chatango) *Room {
	return &Room{
		Mgr: c,
		Name: name,
		Uid: _genUid(),
		FirstCommand: true,
		Server: _getServer(name),
		Port: "8081",
	}
}

func (r *Room) SendCommand(args ...string) {
	terminator := ""
	if r.FirstCommand {
		terminator = "\x00"
		r.FirstCommand = false
	} else {
		terminator = "\r\n\x00"
	}

	command := strings.Join(args, ":") + terminator

	err := r.Ws.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Println("err write:", err)
		return
	}
}

func (r *Room) Auth(){
	if r.Mgr.UserName != "" && r.Mgr.Password != "" {
		r.SendCommand("bauth", r.Name, r.Uid, r.Mgr.UserName, r.Mgr.Password)
	} else if r.Mgr.UserName != "" {
		r.SendCommand("bauth", r.Name, r.Uid)
		r.SendCommand("blogin", r.Mgr.UserName)
	} else {
		r.SendCommand("bauth", r.Name, r.Uid)
	}
}


func(r *Room) Connect(){
	for {
		r.Host = r.Server + ":" + r.Port
		u := url.URL{Scheme: "wss", Host: r.Host, Path: "/"}
		log.Println("Connecting to", r.Host)
		header := http.Header{}
		header.Add("Origin", "https://st.chatango.com")
		ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			log.Println("Dial error:", err)
			time.Sleep(5 * time.Second)
			continue // <-- Retry connection on dial error
		}

		r.Mgr.GroupConnect(r)
		r.Ws = ws
		r.Connected = true
		r.Auth()


		for {
			_, message, err := r.Ws.ReadMessage()
			if err != nil {
				log.Println("Error:", err)

				if closeErr, ok := err.(*websocket.CloseError); ok {
					log.Println("Code Error:", closeErr.Code)
					if closeErr.Code == websocket.CloseAbnormalClosure { // code 1006
						log.Println("Code 1006: Reconnecting...")
						break
					}
				}

				if opErr, ok := err.(*net.OpError); ok {
					if r.Connected {
						log.Println("Network-level error (OpError):", opErr)
						break
					}
				}

				r.Connected = false
				r.Ws.Close()
				return // <-- Exit if not code 1006
			}

			r.Feed(string(message))
		}

		log.Println("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func (r *Room) Message(msg string){
	r.SendCommand("bm", "t12r", r.Channel, msg )
}

func (r *Room) Ping(){
	if r.Connected == true {
		r.SendCommand("")
	}
}

func (r *Room) Disconnect(){
	r.Connected = false
	r.Ws.Close()
}

func (r *Room) Rcmd_b(args []string){
	time := args[0];
	name := args[1];
	puid := args[3];
	_msg, n, f := _clean_message(strings.Join(args[9:], ":"))
	color, face, size := _parseFont(f)

	if name == "" {
		name = "#" + args[2]
		if name == "#" {
			name = "!anon" + _getAnonID(n, args[3])
		}
	}

	user := NewUser(name)
	ip := args[6]
	channel := args[7]
	r.Channel = channel

	user.NameColor = n
	user.FontFace = face
	user.FontSize = size
	user.FontColor = color

	msg := NewMessage(user, _msg, time, puid, ip, channel)

	r.Mgr.GroupMessage(user, r, msg)
}

func (r Room) Feed(food string) {
	parts := strings.Split(food, ":")
	cmd := "Rcmd_" + parts[0]
	args := parts[1:]

	switch cmd {
		case "Rcmd_b":
			r.Rcmd_b(args)
		default :
			//fmt.Printf("%s: %s \n", cmd, args)
	}
}
