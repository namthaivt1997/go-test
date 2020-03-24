package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}


var (
	upgrader =  websocket.Upgrader{}

)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel


var user = make(chan string,2)

func hello(c echo.Context) error {

	ws, err := upgrader.Upgrade( c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	clients[ws] = true

	for {


		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, ws)
			break
		}
		fmt.Printf("%s\n", msg)
		str := string(msg)

		broadcast <- Message{Message:str}
		// Write
		//err = ws.WriteMessage(websocket.TextMessage, []byte(<-user))
		//if err != nil {
		//	fmt.Println(err)
		//}
		//err = ws.WriteMessage(websocket.TextMessage, []byte(<-user))
		//if err != nil {
		//	fmt.Println(err)
		//}

	}
	return c.String(http.StatusOK,"")
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Username + msg.Message))
			if err != nil {
				fmt.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func CreatUser(c echo.Context) error {
	broadcast <- Message{Username:c.QueryParam("username")}

	return c.String(http.StatusOK,"")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/user",CreatUser)
	e.Static("/", "./public")
	e.GET("/ws", hello)
	go handleMessages()

	e.Logger.Fatal(e.Start(":1323"))
}
