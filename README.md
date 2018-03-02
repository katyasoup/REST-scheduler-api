# REST Scheduler API

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

### Prerequisites

[Go v1.10](https://golang.org/doc/go1.10)  
[Gin](https://github.com/gin-gonic/gin)  
[PostgresSQL v9.6](https://www.postgresql.org/download/)  

#### Recommended
[Postico](https://eggerapps.at/postico/) or other PostgreSQL management tool of your choice  
[Restlet Client](https://chrome.google.com/webstore/detail/restlet-client-rest-api-t/aejoelaoggembcahagimdiliamlcdmfm), [Postman](https://www.getpostman.com/) or other REST client of your choice

### Install
Download and install Go: [Golang Installation Guide](https://golang.org/doc/install)  

Download and install Gin: ```$ go get github.com/gin-gonic/gin```

Download and install PostgresSQL: [PostgresSQL download](https://www.postgresql.org/download/)

Set up your PostgresSQL database:  

- Run the script from ```db_setup.sql``` to create and populate your tables

### Run
- Download and clone the project files: ```$ git clone https://github.com/katyasoup/wiw-challenge.git```
- Adjust database variables in ```database.go``` as needed at lines 9 and 21 
	- If you did not set up a username and password in Postgres, use default user "postgres" and remove password fields from line 9 ```const``` and line 22 ```psqlInfo```
	- ```dbname``` is the name of the database where you created your tables  
- Spin it up! ```$ go run main.go types.go database.go``` The project will be available on port 8080


## Acessing the API endpoints
Access the following information by supplying as needed:  

- an employee ID (```:empID```) or shift ID (```:shiftID```) as an integer;  
- dates (```:startDate``` & ```:endDate```) as yyy-mm-ddThh:mm:ssZ (ex. 2018-02-27T09:00:00Z)  
- for PUT and POST routes, send data as type ```application/x-www-form-url-encoded```

ex. POST:
```
{  
"manager": 7,  
"break": 1,  
"startTime": "2018-04-08T13:00:00Z",  
"endTime": "2018-04-08T17:00:00Z"  
}
```

ex. PUT: (assign employee)  
```
{
"employee": 1,
"id": 20
}
```

ex. PUT: (edit times)  
```
{
"startTime": "2018-03-06T09:00:00Z",
"endTime": "2018-03-06T15:00:00Z",
"id": 16
}
```

### User Stories

- [x] As an employee, I want to know when I am working, by being able to see all of the shifts assigned to me.
	- GET: ``` /myshifts/:empID```
	- ex. [http://localhost:8080/myshifts/2](http://localhost:8080/myshifts/2)
- [x] As an employee, I want to know who I am working with, by being able to see the employees that are working during the same time period as me.
	- GET: ```/roster/:startDate/:endDate```
	- ex. [http://localhost:8080/roster/2018-03-01T09:00:00/2018-03-01T17:00:00Z](http://localhost:8080/roster/2018-03-01T09:00:00/2018-03-01T17:00:00Z)
- [x] As an employee, I want to know how much I worked, by being able to get a summary of hours worked for each week.
	- GET: ```/hours/:empID/:startDate/:endDate```
	- ex. [http://localhost:8080/hours/4/2018-03-01T09:00:00Z/2018-03-07T17:00:00Z](http://localhost:8080/hours/4/2018-03-01T09:00:00Z/2018-03-07T17:00:00Z)
- [x] As an employee, I want to be able to contact my managers, by seeing manager contact information for my shifts.
	- GET: ```/mymanagers/:empID```
	- ex. [http://localhost:8080/mymanagers/1](http://localhost:8080/mymanagers/1)
- [x] As a manager, I want to schedule my employees, by creating shifts for any employee.
	- 	POST: ```/shifts```
- [x] As a manager, I want to see the schedule, by listing shifts within a specific time period.
	- GET: ```/schedule/:start/:end```
	- ex. [http://localhost:8080/schedule/2018-03-01T09:00:00Z/2018-03-04T17:00:00Z](http://localhost:8080/schedule/2018-03-01T09:00:00Z/2018-03-04T17:00:00Z)
- [x] As a manager, I want to be able to change a shift, by updating the time details.
	- 	PUT: ```/shifts```
- [x] As a manager, I want to be able to assign a shift, by changing the employee that will work a shift.
	- 	PUT: ```/shifts/assign```
	-  To see currently unassigned shifts: [http://localhost:8080/shifts/unassigned](http://localhost:8080/shifts/unassigned)
- [x] As a manager, I want to contact an employee, by seeing employee details.
	- GET: ```/employees``` or ```/employees/:id```
	- ex. [http://localhost8080/employees](http://localhost:8080/employees)
	- ex. [http://localhost8080/employees/1](http://localhost:8080/employees/1)
	
## Known Limitations & Next Steps

- Currently, there is no authorization and therefore no way to access data from the current logged in user, which would dictate access rights (employees can read only; managers can read and write). [Pstore](https://github.com/xyproto/pstore) and [AuthBoss](https://github.com/volatiletech/authboss) look like promising auth frameworks to add this critical functionality.
	- The ```basicAuth``` branch of this project contains code for simple username and password validation for manager PUT and POST routes
- Dates are stored in the database as ISO 8601 format; logic is needed to convert these to RFC 2822 format.  
- Ideally, a method updateShift() should be available for both assigning employees and updating time details; currently having issues with allowing null values for shift breaks, employee id, phone, and/or email fields in the PUT request

### Testing

- Investigate in-memory database solution and mocking frameworks
