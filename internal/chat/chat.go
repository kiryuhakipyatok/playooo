package chat

import (
	"io"
	"sync"
	"github.com/gofiber/contrib/websocket"
	"log"
)

type Member struct{
	login string
	conn *websocket.Conn
}

type Chat struct {
	mu sync.Mutex
	members map[*Member]bool
}

var (
	chats = make(map[string]*Chat)
	cmu sync.Mutex
)

func NewChat() *Chat {
	return &Chat{
		members: make(map[*Member]bool),
	}
}

func(c *Chat) addMember(m *Member){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.members[m]=true
	log.Printf("new member: %s",m.conn.RemoteAddr())
}

func(c *Chat) removeMember(m *Member){
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.members,m)
}

func HandleWS(ws *websocket.Conn) {
	chatId:=ws.Params("id")
	if chatId == ""{
		ws.Close()
		return
	}
	cmu.Lock()
	chat,ok:=chats[chatId]
	if !ok{
		chat = NewChat()
		chats[chatId]=chat
	}
	cmu.Unlock()
	member:=Member{conn: ws}
	chat.addMember(&member)
	defer func ()  {
		chat.removeMember(&member)
		ws.Close()
		log.Printf("connection closed in chat %s: %s", chatId, ws.RemoteAddr())
		chat.mu.Lock()
		if len(chat.members) == 0{
			delete(chats,chatId)
		}
		chat.mu.Unlock()
	}()
	chat.readLoop(ws)
}

func (c *Chat) broadcast(b []byte, mt int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for member := range c.members {
		go func(ws *websocket.Conn) {
			msg:=append([]byte(member.login + ": "), b...)
			if err := member.conn.WriteMessage(mt, msg); err != nil {
				log.Printf("error writing: %v", err)
				ws.Conn.Close()
				delete(c.members,member)
			}
		}(member.conn)
	}
}

func (c *Chat) readLoop(ws *websocket.Conn) {
	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error reading: %v", err)
			continue
		}
		log.Println(string(msg))
		c.broadcast(msg, mt)
	}
}

