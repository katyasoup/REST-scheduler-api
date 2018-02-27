package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func stringToInt64(str string) int64 {
	id, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		panic(err)
	}
	return id
}

func main() {

	OpenDatabase()
	defer CloseDatabase()

	routes := gin.Default()

	routes.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hey there, thanks for checking out my project! -Katie")
	})
	// // EMPLOYEE user stories:

	// As an employee, I want to know when I am working, by being able to see all of the shifts assigned to me:
	routes.GET("/myshifts/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		results := getShiftsByEmployee(id)
		c.JSON(200, results)
	})

	// As an employee, I want to know who I am working with, by being able to see the
	// employees that are working during the same time period as me:
	routes.GET("/roster/:start/:end", func(c *gin.Context) {
		start := c.Params.ByName("start")
		end := c.Params.ByName("end")
		results := getEmployeeRostersByDateRange(start, end)
		c.JSON(200, results)
	})

	// As an employee, I want to know how much I worked, by being able to get a summary of hours worked for each week:
	// // TODO: add math for subtracting break time from total hours
	// // TODO: only total hours for dates in past
	routes.GET("/hours/:id/:start/:end", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		start := c.Params.ByName("start")
		end := c.Params.ByName("end")
		results := getShiftsByEmployeeInDateRange(id, start, end)
		var totalHours int

		for _, shift := range results {
			start, err := time.Parse(time.RFC3339, shift.Start)
			end, err := time.Parse(time.RFC3339, shift.End)
			if err != nil {
				panic(err)
			}
			shiftHours := end.Hour() - start.Hour()
			fmt.Printf("Shift length: %d ", shiftHours)

			// add hours for each shift in date range to total hours
			totalHours += shiftHours
			fmt.Printf("Total hours: %d ", totalHours)
		}
		summary := Hours{results, totalHours}
		c.JSON(200, summary)
	})

	// As an employee, I want to be able to contact my managers, by seeing manager contact information for my shifts:
	routes.GET("/mymanagers/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		results := getManagerRostersByDateRange(id)
		c.JSON(200, results)
	})

	// // MANAGER user stories:

	// As a manager, I want to see the schedule, by listing shifts within a specific time period:
	// // NOTE route listed as "/schedule/:start/:end" because "/shifts/:start/:end" conflicts with "/shifts/:id"
	routes.GET("/schedule/:start/:end", func(c *gin.Context) {
		start := c.Params.ByName("start")
		end := c.Params.ByName("end")
		results := getShiftsByDateRange(start, end)
		c.JSON(200, results)
	})

	// As a manager, I want to schedule my employees, by creating shifts for any employee:
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

	// As a manager, I want to be able to assign a shift, by changing the employee that will work a shift:
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

	// As a manager, I want to be able to change a shift, by updating the time details:
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

	// As a manager, I want to contact an employee, by seeing employee details:
	// // get all employees
	routes.GET("/employees", func(c *gin.Context) {
		results := getAllEmployees()
		c.JSON(200, results)
	})

	// // get single employee by id
	routes.GET("/employees/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getEmployeeByID(id)
		c.JSON(200, result)
	})

	routes.Run()
}
