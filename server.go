package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_"go/doc"
	"net/http"
	//_"encoding/json"
)
type User struct {
	id     string `json:"id"`
	name   string `json:"name"`
	password string `json: "password"`
}
type Users struct {
	Users []User `json:"employee"`
}

var con *sql.DB

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:9128617542@tcp(127.0.0.1:3306)/dev")
	if err != nil {
		fmt.Println(">>>",err.Error())
	} else {
		fmt.Println("db is connected")


	}
	//defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Println("hhhh",err.Error())
	}
	return db
}



func ShowUsers(c echo.Context) error {
	con := ConnectDB()
	sqlStatement := "SELECT id, name, password FROM user order by id"

	rows, err := con.Query(sqlStatement)
	fmt.Println(rows)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusCreated, err);
	}
	defer rows.Close()
	result := Users{}

	for rows.Next() {
		user := User{}

		err2 := rows.Scan(&user.id, &user.name, &user.password)
		fmt.Println(">>>> user: ",user)
		// Exit if we get an error
		if err2 != nil {
			fmt.Print(err2)
		}
		result.Users = append(result.Users, user)
	}

	fmt.Println(result)
	return c.JSON(http.StatusOK, result)
}

func main() {


	e := echo.New()

	e.GET("/app", ShowUsers)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}
