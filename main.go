package main

import (
	"fmt"
	"net/http"
	"math"
	"strconv"
	"time"
	// framework echo
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// usedd to connect to database
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// var users []User
var seq   = 1
type User struct {
		Id  		bson.ObjectId       `json:"_id"   bson:"_id,omitempty"`
	    Name        string         		`json:"name"`
	    Status      string         		`json:"status"`
	    Lat         string         		`json:"lat"`
	    Long        string         		`json:"long"`
	    Created_at  time.Time      		`json:"created_at"`
	    Updated_at  time.Time      		`json:"updated_at"` 
}

type Message struct {
	Message  string  `json:"message"`
}

type MessageNearest struct {
	Message  string  `json:"message"`
	Nearest string   `json:"nearest"`
}

func createUser(c echo.Context) error {
	u := &User{
		Name: c.FormValue("name"),
		Status: c.FormValue("status"),
		Lat: c.FormValue("long"),
		Long: c.FormValue("lat"),
	}
	session, err := mgo.Dial("localhost")
    if err != nil {
            return err
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    connection := session.DB("helloWorld").C("users")
    err = connection.Insert(&User{Name: u.Name, Status: u.Status, Lat: u.Lat, Long: u.Long, Created_at: time.Now(), Updated_at: time.Now()})
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func getUsers(c echo.Context) error {
	session, err := mgo.Dial("localhost")
    if err != nil {
            return err
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    connection := session.DB("helloWorld").C("users")
    var results []User
    err = connection.Find(nil).All(&results)
    if err != nil {
        return err
    } 
    return c.JSON(http.StatusCreated, results)
}

func updateUser(c echo.Context) error {
	idUser := bson.ObjectIdHex(c.Param("id"))
	u := &User{
		Name: c.FormValue("name"),
		Status: c.FormValue("status"),
		Lat: c.FormValue("long"),
		Long: c.FormValue("lat"),
	}
	session, err := mgo.Dial("localhost")
    if err != nil {
            return err
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    connection := session.DB("helloWorld").C("users")
    colQuerier := bson.M{"_id": idUser}
    fmt.Println(colQuerier)
    change := bson.M{"$set": bson.M{"name": u.Name,"status": u.Status,"lat": u.Lat,"long": u.Long,"updated_at": time.Now()}}
    err = connection.Update(colQuerier, change)
    if err != nil {
    	return c.JSON(http.StatusCreated, err)
    }
    msg := &Message{
		Message: "success",
	}
    return c.JSON(http.StatusCreated, msg)
}

func deleteUser(c echo.Context) error {
	idUser := bson.ObjectIdHex(c.Param("id"))
	session, err := mgo.Dial("localhost")
    if err != nil {
            return err
    }
    defer session.Close()
    // session.SetMode(mgo.Monotonic, true)
    connection := session.DB("helloWorld").C("users")
    err = connection.Remove(bson.M{"_id": idUser})
	if err != nil {
        return c.JSON(http.StatusCreated, err)
    }
	msg := &Message{
		Message: "success",
	}
    return c.JSON(http.StatusCreated, msg)
}

func nearest(c echo.Context) error {
	idUser := bson.ObjectIdHex(c.Param("id"))
	session, err := mgo.Dial("localhost")
    if err != nil {
            return err
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)
    connection := session.DB("helloWorld").C("users")
    var results []User
    err = connection.Find(nil).All(&results)
    if err != nil {
        return err
    } else {
    	var indexUser int
    	var latUser string
    	var longUser string
    	var nameUser string
    	var nearUser string
    	dist := 99999999.99999999
    	for  index, element := range results {
    		if element.Id == idUser {
    			indexUser = index
    			latUser = element.Lat
    			longUser = element.Long
    			nameUser = element.Name
    		}
		}
		for  index, element := range results {
    		if index != indexUser {
    			x1, err1 := strconv.ParseFloat(longUser, 64)
    			x2, err2 := strconv.ParseFloat(element.Long, 64)
    			y1, err3 := strconv.ParseFloat(latUser, 64)
    			y2, err4 := strconv.ParseFloat(element.Lat, 64)
    			if( err1 == nil && err2 == nil && err3 ==nil && err4 == nil) {
    				diff := math.Sqrt(math.Pow((x1-x2),2) + math.Pow((y1-y2),2))
    				if diff < dist {
    					dist = diff
    					nearUser = element.Name
    				}
    			}
    		}
		}
		fmt.Println(nameUser)
		fmt.Println(nearUser)
		msg := &MessageNearest{
			Message: "success",
			Nearest: nearUser,
		}
    	return c.JSON(http.StatusCreated, msg)
    }
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/add-user", createUser)
	e.GET("/users", getUsers)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	e.GET("/near/:id", nearest)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}