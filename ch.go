//author agunq.e@gmail.com

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"math"
	"math/rand"
	"time"
	"strings"
	"strconv"
	"regexp"
	"github.com/gorilla/websocket"
)

type ServerWeight struct {
	ID     string
	Weight int
}

var tsweights = []ServerWeight{
	{"5", 75}, {"6", 75}, {"7", 75}, {"8", 75}, {"16", 75}, {"17", 75}, {"18", 75},
	{"9", 95}, {"11", 95}, {"12", 95}, {"13", 95}, {"14", 95}, {"15", 95}, {"19", 110},
	{"23", 110}, {"24", 110}, {"25", 110}, {"26", 110}, {"28", 104}, {"29", 104},
	{"30", 104}, {"31", 104}, {"32", 104}, {"33", 104}, {"35", 101}, {"36", 101},
	{"37", 101}, {"38", 101}, {"39", 101}, {"40", 101}, {"41", 101}, {"42", 101},
	{"43", 101}, {"44", 101}, {"45", 101}, {"46", 101}, {"47", 101}, {"48", 101},
	{"49", 101}, {"50", 101}, {"52", 110}, {"53", 110}, {"55", 110}, {"57", 110},
	{"58", 110}, {"59", 110}, {"60", 110}, {"61", 110}, {"62", 110}, {"63", 110},
	{"64", 110}, {"65", 110}, {"66", 110}, {"68", 95}, {"71", 116}, {"72", 116},
	{"73", 116}, {"74", 116}, {"75", 116}, {"76", 116}, {"77", 116}, {"78", 116},
	{"79", 116}, {"80", 116}, {"81", 116}, {"82", 116}, {"83", 116}, {"84", 116},
}

var totalWeight = 7034

func _getServer(group string) string {
	group = strings.ReplaceAll(group, "_", "q")
	group = strings.ReplaceAll(group, "-", "q")
	end := len(group)
	if end > 5 {
		end = 5
	}
	fnv, _ := strconv.ParseInt(group[:end], 36, 64)

	if len(group) <= 6 {
		return ""
	}
	start := 6
	end = start + 3
	if end > len(group) {
		end = len(group)
	}

	lnvStr := group[start:end]

	var lnv int64
	if lnvStr != "" {
		lnv, _ = strconv.ParseInt(lnvStr, 36, 64)
		if lnv < 1000 {
			lnv = 1000
		}
	} else {
		lnv = 1000
	}

	num := float64(fnv%lnv) / float64(lnv)

	//totalWeight := 0
	//for _, sw := range tsweights {
	//	totalWeight += sw.Weight
	//}

	cumfreq := 0.0
	for _, sw := range tsweights {
		cumfreq += float64(sw.Weight) / float64(totalWeight)
		if num <= cumfreq {
			return fmt.Sprintf("s%s.chatango.com", sw.ID)
		}
	}

	panic(fmt.Sprintf("Couldn't find host server for room %s", group))
}

func _getAnonID(n string, ssid string) string {
	if n == "" || ssid == "" {
		return ""
	}
	id := ""
	for i := 0; i < 4; i++ {
		a, _ := strconv.Atoi(n[i:i+1])
		b, _ := strconv.Atoi(ssid[i+4 : i+5])
		sum := a + b
		id += strconv.Itoa(sum % 10)
	}
	return id
}

func _strip_html(msg string) string {
	htmlRegex := regexp.MustCompile(`<\/?[^>]*>`)
	strippedMsg := htmlRegex.ReplaceAllString(msg, "")
	return strippedMsg
}

func _clean_message(msg string) (string, string, string) {
	nRegex := regexp.MustCompile(`<n(.*?)\/>`)
	nMatch := nRegex.FindStringSubmatch(msg)
	var n string
	if len(nMatch) > 1 {
		n = nMatch[1]
	}

	fRegex := regexp.MustCompile(`<f(.*?)>`)
	fMatch := fRegex.FindStringSubmatch(msg)
	var f string
	if len(fMatch) > 1 {
		f = fMatch[1]
	}

	msg = nRegex.ReplaceAllString(msg, "")
	msg = fRegex.ReplaceAllString(msg, "")
	msg = _strip_html(msg)
	msg = strings.ReplaceAll(msg, "&lt;", "<")
	msg = strings.ReplaceAll(msg, "&gt;", ">")
	msg = strings.ReplaceAll(msg, "&quot;", "\"")
	msg = strings.ReplaceAll(msg, "&apos;", "'")
	msg = strings.ReplaceAll(msg, "&amp;", "&")

	return msg, n, f
}

func _parseFont(f string) (string, string, int) {
	if f != "" {
		sizeColorFontFace := strings.SplitN(f, "=", 2)
		sizeColor := strings.TrimSpace(sizeColorFontFace[0])
		fontFace := sizeColorFontFace[1]

		sizeRegex := regexp.MustCompile(`x(\d\d|\d)`)
		sizeMatch := sizeRegex.FindStringSubmatch(sizeColor)
		var size int
		if len(sizeMatch) > 1 {
			size, _ = strconv.Atoi(sizeMatch[1])
		} else {
			size = 0
		}

		col := sizeRegex.ReplaceAllString(sizeColor, "")
		if col == "" {
			col = "000"
		}

		face := fontFace[1 : len(fontFace)-1]
		if face == "" {
			face = "0"
		}

		return col, face, size
	} else {
		return "000", "0", 10
	}
}

func _genUid() string {
	min := math.Pow10(15)
	max := math.Pow10(16)
	rand.Seed(time.Now().UnixNano())
	num := int(min + rand.Float64()*(max-min+1))
	return fmt.Sprintf("%d", num)
}


//User Class
type User struct{
	Name string
	NameColor string;
	FontFace string;
	FontSize int;
	FontColor string;
}

var users = make(map[string]*User)

func NewUser(name string) *User {
	lowerName := strings.ToLower(name)

	if user, exists := users[lowerName]; exists {
		return user
	}

	user := &User{
		Name: name,
	}
	users[lowerName] = user

	return user
}


//Message classs
type Message struct{
	User *User
	Time string
	Puid string
	Body string
	Ip string
	Channel string
}

func NewMessage(user *User, body string, time string, puid string, ip string, channel string) *Message {
	return &Message{
		User: user,
		Body: body,
		Time: time,
		Puid: puid,
		Ip: ip,
		Channel: channel,
	}
}

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
	r.Host = r.Server + ":" + r.Port
	u := url.URL{Scheme: "wss", Host: r.Host, Path: "/"}
	//log.Printf("connecting to %s", r.Host)
	header := http.Header{}
	header.Add("Origin", "https://st.chatango.com")
	r.Mgr.Connect(r)
	r.Connected = true
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	r.Ws = ws
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer r.Ws.Close()

	r.Auth()

	for {
		_, message, err := r.Ws.ReadMessage()
		if err != nil {
			r.Connected  = false
			r.Mgr.Disconnect(r, err)
			return
		}
		r.Feed(string(message))
	}
}

func (r *Room) Message(msg string){
	r.SendCommand("bm", "t12r", r.Channel, msg )
}

func (r *Room) Ping(){
	r.SendCommand("")
}

func (r *Room) Disconnect(){
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

	r.Mgr.Message(user, r, msg)
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


type Chatango struct{
	UserName string
	Password string
	PrivateMessage *PrivateMessage
	User *User
	RoomList map[string]*Room
	Running bool
}


func NewChatango() *Chatango{
	return &Chatango{
		Running: true,
	}
}

func (c *Chatango) EasyStart(rooms []string, username string, password string) {
	c.UserName = username
	c.Password = password
	c.RoomList = make(map[string]*Room)
	c.User = NewUser(username)
	for _, roomName := range rooms {
		c.RoomList[roomName] = NewRoom(roomName, c)
		go c.RoomList[roomName].Connect()
	}

	if c.UserName != "" && c.Password != ""{
		c.PrivateMessage = NewPrivateMessage(c)
		go c.PrivateMessage.Connect()
	}

	interval := time.Second * 10
	ticker := time.NewTicker(interval)
	tickerChannel := ticker.C
	for c.Running {
		<-tickerChannel

		var activeConnections []string

		// Perform the ping task
		for _, room := range c.RoomList {
			if room.Connected == true{
				room.Ping()
				activeConnections = append(activeConnections, room.Name)
			}
		}
		if c.PrivateMessage.Connected == true{
			c.PrivateMessage.Ping()
			activeConnections = append(activeConnections, c.PrivateMessage.Name)
		}

		if len(activeConnections) == 0 {
			log.Println("No active connections left, stopping the ping loop.")
			c.Running = false
		}
	}
}
