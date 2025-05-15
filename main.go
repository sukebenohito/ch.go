package main


import "fmt"


func (c *_Chatango) Message(user *_User, room *_Room, message *_Message) {
	fmt.Printf("%s %s %s\n", room.Name, user.Name, message.Body)
	if message.Body == "halo" {
		room.Message("halo juga")
	}
	if message.Body == "out aja" {
		room.Disconnect()
	}
}

func (c *_Chatango) Connect(room *_Room) {
	fmt.Printf("connected to %s\n", room.Name)
}

func (c *_Chatango) Disconnect(room *_Room, err error) {
	fmt.Printf("disconnected to %s %s\n", room.Name, err)
}


func main(){
	ch := Chatango()

//	example
//	ch.EasyStart([]string{"nico-nico", "monosekai", "desertofdead"}, "BotName", "")
	ch.EasyStart([]string{"nico-nico"}, "BotName", "")

}


