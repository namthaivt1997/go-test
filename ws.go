package main

import (
	db_ "demo/DB"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lann/ps"
	"io/ioutil"
	"net/http"
	"time"
)
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Room struct {
	Id string `json:"id"`
	Status string `json:"status"`
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type isLogin struct {
	isLogin bool `json:"isLogin"`
}


var (
	upgrader =  websocket.Upgrader{}

)


var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel =
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

func JoinRoom(c echo.Context) error {
	broadcast <- Message{Username:c.QueryParam("username")}

	return c.String(http.StatusOK,"")
}

func login(c echo.Context) error {
	// get data from body
	data, _ := ioutil.ReadAll(c.Request().Body)


	userLogin := db_.User{}

	json.Unmarshal(data,&userLogin)

	fmt.Println(string(data))



	db := db_.ConnectDB()

	sqlStatement := "SELECT id, name, password FROM user WHERE id =" + userLogin.Id

	rows, err := db.Query(sqlStatement)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusCreated, err);
	}
	defer rows.Close()

	user := db_.User{}

	for rows.Next() {

		err2 := rows.Scan(&user.Id, &user.Name, &user.Pw)
		fmt.Println(">>>> user: ",user)
		// Exit if we get an error
		if err2 != nil {
			fmt.Print(err2)
		}
	}

	fmt.Println(userLogin)

	// Throws unauthorized error
	if userLogin.Pw != user.Pw {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]ps.Any{
		"token": t,
		"isLogin": true,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}


func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/app/login", login)

	// Unauthenticated route
	e.GET("/app/accessible", accessible)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)


	e.GET("/user",JoinRoom)
	e.Static("/", "./public")
	e.GET("/ws", hello)
	go handleMessages()

	e.Logger.Fatal(e.Start(":1323"))
}
