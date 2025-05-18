package main


import "fmt"

func GroupMessage(user *User, room *Room, message *Message) {
	fmt.Printf("%s %s %s\n", room.Name, user.Name, message.Body)
	if message.Body == "halo" {
		room.Message("halo juga")
	}
	if message.Body == "out aja" {
		room.Disconnect()
	}

}

func PMessage(user *User, private *PrivateMessage, message string) {
	fmt.Printf("%s %s %s\n", private.Name, user.Name, message)
	if message == "halo" {
		private.Message(user.Name, "halo juga")
	}
	if message == "out aja" {
		private.Disconnect()
	}

}


func main(){
	ch := NewChatango()
	ch.PMessage,ch.GroupMessage = PMessage, GroupMessage
	//GroupConnect, GroupDisconnect
	//	example
	//	ch.EasyStart([]string{"nico-nico", "monosekai", "desertofdead"}, "Name", "Password")
	ch.EasyStart([]string{"nico-nico"}, "devilsona", "")


}
