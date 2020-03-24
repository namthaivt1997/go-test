package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	//"fmt"
	//"github.com/casbin/casbin"
	//"github.com/labstack/echo/v4"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//casbin_mw "github.com/labstack/echo-contrib/casbin"
)


// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Student struct {
	Id string `json:"studen_id"`
	Name string `json:"student_name"`
}

type Students struct {
	Students []Student `json:"Student"`
}

var (
	upgrader1 = websocket.Upgrader{}
)

var con *sql.DB
// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func hello1(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}


func ConnectDB() *sql.DB  {
	db, err := sql.Open("mysql", "devteam:*devteam@123@tcp(192.168.11.49:3306)/testing")
	if err != nil {
		panic(err.Error())
	}
	//defer db.Close()

	return db

}

func getData(c echo.Context) error {


	//w := echo.NewResponse()
	con := ConnectDB()

	sqlStatement := "SELECT studen_id,student_name FROM student order by studen_id"

	rows, err := con.Query(sqlStatement)
	fmt.Println(rows)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		//return c.JSON(http.StatusCreated, u);
	}

	defer rows.Close()
	result := Students{}
	fmt.Println(">>>",rows)

	for rows.Next() {
		st := Student{}

		err2 := rows.Scan(&st.Id, &st.Name)
		// Exit if we get an error
		if err2 != nil {
			fmt.Print(err2)
		}
		result.Students = append(result.Students, st)
	}




	fmt.Println(">>>>",result.Students)


	//return c.JSON(http.StatusOK, response)
	return c.JSON(http.StatusOK, result.Students)
}



func main() {



	e := echo.New()
	g := e.Group("/admin")


	cookieGroup := e.Group("/cookie")

	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	g.GET("/db",getData)
	g.Use(middleware.Logger())
	g.Use(middleware.Recover())
	g.Static("/", "../public")
	g.GET("/ws", hello)

	cookieGroup.GET("/login", Login)
	g.GET("/main", func(c echo.Context) error {
		return c.String(http.StatusOK,"OK")
	})
	e.POST("/cats", addCat)

	e.Start(":8000")
}

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "thai" && password == "123" {
		cookie :=&http.Cookie{}

		//this is the same
		// cookie := new(http.Cookie)

		cookie.Name = "SessionID"
		cookie.Value = "1111"

		c.SetCookie(cookie)

	}
	return c.String(http.StatusOK, "cookie set")
}

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("read false",err)
	}

	err = json.Unmarshal(b, &cat);
	if err != nil {
		log.Println("json false",err)
	}

	log.Println("cat",cat)
	return c.String(http.StatusOK,"we got your cat")
}
