package main


import "fmt"

func (c *Chatango) Message(user *User, room *Room, message *Message) {
	fmt.Printf("%s %s %s\n", room.Name, user.Name, message.Body)
	if message.Body == "halo" {
		room.Message("halo juga")
	}
	if message.Body == "out aja" {
		room.Disconnect()
	}

}

func (c *Chatango) PMessage(user *User, private *PrivateMessage, message string) {
	fmt.Printf("%s %s %s\n", private.Name, user.Name, message)
	if message == "halo" {
		private.Message(user.Name, "halo juga")
	}
	if message == "out aja" {
		private.Disconnect()
	}

}

func (c *Chatango) Connect(room *Room) {
	fmt.Printf("connected to %s\n", room.Name)
}

func (c *Chatango) Disconnect(room *Room, err error) {
	fmt.Printf("disconnected to %s %s\n", room.Name, err)
}


func main(){
	ch := NewChatango()

	//	example
	//	ch.EasyStart([]string{"nico-nico", "monosekai", "desertofdead"}, "Name", "")
	ch.EasyStart([]string{"nico-nico"}, "Name", "")

}
