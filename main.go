package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "Katie"
	dbname = "wiw-challenge"
)

// Global DB handles
var db *sql.DB
var err error

func stringToInt64(str string) int64 {
	id, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		panic(err)
	}
	return id
}

func main() {
	// db setup
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database!!!")

	routes := gin.Default()

	// SHIFT routes:

	// get all shifts
	routes.GET("/shifts", func(c *gin.Context) {
		results := getAllShifts()
		c.JSON(200, results)
	})

	// get single shift by id
	routes.GET("/shifts/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getShiftByID(id)
		c.JSON(200, result)
	})

	// get all shifts for single employee
	routes.GET("/myshifts/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		results := getMyShifts(id)
		c.JSON(200, results)
	})

	// get all shifts for date range
	routes.GET("/schedule/:start/:end", func(c *gin.Context) {
		start := c.Params.ByName("start")
		end := c.Params.ByName("end")
		results := getSchedule(start, end)
		c.JSON(200, results)
	})

	// add new shift
	routes.POST("/shifts", func(c *gin.Context) {
		var shift Shift
		c.BindJSON(&shift)
		result := createShift(shift)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			c.JSON(201, gin.H{"success": result})
		}
	})

	// change shift time
	routes.PUT("/shifts", func(c *gin.Context) {
		var shift Shift
		c.BindJSON(&shift)
		result := editShiftTime(shift)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			c.JSON(201, gin.H{"success": result})
		}
	})

	// add employee to shift
	routes.PUT("/shifts/assign", func(c *gin.Context) {
		var shift Shift
		c.BindJSON(&shift)
		result := scheduleEmployee(shift)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			c.JSON(201, gin.H{"success": result})
		}
	})

	// USER routes:

	// get all users
	routes.GET("/users", func(c *gin.Context) {
		results := getAllUsers()
		c.JSON(200, results)
	})

	// get single user by id
	routes.GET("/users/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getUserByID(id)
		c.JSON(200, result)
	})

	// get all users with role of employee
	routes.GET("/employees", func(c *gin.Context) {
		results := getAllEmployees()
		c.JSON(200, results)
	})

	// get single employee by id
	routes.GET("/employees/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getEmployeeByID(id)
		c.JSON(200, result)
	})

	// get all users with role of manager
	routes.GET("/managers", func(c *gin.Context) {
		results := getAllManagers()
		c.JSON(200, results)
	})

	// get single manager by id
	routes.GET("/managers/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getManagerByID(id)
		c.JSON(200, result)
	})

	// get coworkers for date range
	routes.GET("/roster/:start/:end", func(c *gin.Context) {
		start := c.Params.ByName("start")
		end := c.Params.ByName("end")
		results := getCoworkers(start, end)
		c.JSON(200, results)
	})

	routes.Run()
}
