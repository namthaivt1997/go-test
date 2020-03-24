package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net"
	"net/url"
)

var (
	upgrader1    = websocket.Upgrader{}
)

func CreatClinet(c echo.Context) error {
	websocket.
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "./")
	e.GET("/client", CreatClinet)
	e.Start(":8000")

}


