package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	invalidName := true
	var name string
	for invalidName {
		invalidName = false
		var err error
		name, err = conn.Recv()
		if err != nil {
			// Client disconnected
			return
		}
		name = strings.TrimSpace(name)

		mutex.Lock()
		for _, c := range clients {
			if name == c.name {
				conn.Send("INVALID\n")
				invalidName = true
				break
			}
		}
		mutex.Unlock()
	}

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
	log.Println(name + " connected")
	fmt.Print(connectMsg)
	SendMessageFrom(client, connectMsg)
	for {
		// will listen for message to process ending in newline (\n)
		inMessage, err := conn.Recv()
		if err != nil {
			// User disconnected
			disconnectMsg := "DISCONNECT:" + name + "\n"
			log.Println(name + " disconnected")
			fmt.Print(disconnectMsg)
			SendMessageFrom(client, disconnectMsg)
			RemoveClient(client)
			return
		}

		msg := string(inMessage)
		if strings.HasPrefix(msg, "GO:") {
			log.Println("New game")
			SendMessageFrom(client, msg)
			nplayers := len(clients)
			special := strings.Split(msg, ":")[1]
			specials := make([]int, 0, rtMorgana)
			if strings.Contains(special, ",") {
				for _, val := range strings.Split(special, ",") {
					r, err := strconv.Atoi(val)
					if err == nil && r != -1 {
						specials = append(specials, r)
					}
				}
			}
			roles := AssignRoles(nplayers, specials)
			mutex.Lock()
			for i, c := range clients {
				c.r = roles[i]
				c.conn.Send("ROLE:" + strconv.Itoa(c.r.roleType) + "\n")
				log.Println(c.name + ": " + roleNames[c.r.roleType])
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
	f, err := os.OpenFile("server_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Server started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
