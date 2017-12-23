package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

var clients []Client
var mutex sync.Mutex

type Client struct {
	conn sockjs.Session
	name string
	r    *Role
}

func RemoveClient(toRemove Client) {
	for i, c := range clients {
		if c.conn == toRemove.conn {
			mutex.Lock()
			clients = append(clients[:i], clients[i+1:]...)
			mutex.Unlock()
		}
	}
}

func SendMessageFrom(client Client, message string) {
	for _, c := range clients {
		if c.conn != client.conn {
			c.conn.Send(message)
		}
	}
}

func HandleConnection(conn sockjs.Session) {
	name, _ := conn.Recv()
	name = strings.TrimSpace(name)

	client := Client{
		conn: conn,
		name: name,
	}

	mutex.Lock()
	clients = append(clients, client)
	for _, c := range clients {
		conn.Send("CONNECT:" + c.name + "\n")
	}
	mutex.Unlock()

	connectMsg := "CONNECT:" + name + "\n"
	fmt.Print(connectMsg)
	SendMessageFrom(client, connectMsg)
	for {
		// will listen for message to process ending in newline (\n)
		inMessage, err := conn.Recv()
		if err != nil {
			// User disconnected
			disconnectMsg := "DISCONNECT:" + name + "\n"
			fmt.Print(disconnectMsg)
			SendMessageFrom(client, disconnectMsg)
			RemoveClient(client)
			return
		}

		msg := string(inMessage)
		if strings.HasPrefix(msg, "GO") {
			nplayers := len(clients)
			roles := AssignRoles(nplayers, []int{rtPercival, rtMorgana, rtMerlin})
			mutex.Lock()
			for i, c := range clients {
				c.r = roles[i]
				c.conn.Send("ROLE:" + strconv.Itoa(c.r.roleType) + "\n")
				infoStr := ""
				for i, index := range c.r.info {
					infoStr += clients[index].name
					if i != len(c.r.info)-1 {
						infoStr += ", "
					}
				}
				c.conn.Send("INFO:" + infoStr + "\n")
			}
			mutex.Unlock()
		}
	}
}

func main() {
	fmt.Println("Launching server...")

	opts := sockjs.DefaultOptions
	handler := sockjs.NewHandler("/avalon", opts, HandleConnection)
	http.Handle("/avalon/", handler)
	http.Handle("/", http.FileServer(http.Dir("web/")))
	log.Println("Server started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
