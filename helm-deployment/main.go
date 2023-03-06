package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

// https://github.com/gofiber/recipes/blob/master/mysql/main.go
// Database settings
const (
	host     = "localhost"
	port     = 3306 // Default port
	user     = "root"
	password = "root"
	dbname   = "classicmodels"
)

// Connect function
func Connect() error {
	var err error
	// Use DSN string to open
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(mysql:%v)/%s", user, password, port, dbname))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

type Employee struct {
	EmployeeNumber sql.NullInt64  `json:"employeeNumber"`
	LastName       sql.NullString `json:"lastName"`
	FirstName      sql.NullString `json:"firstName"`
	Extension      sql.NullString `json:"extension"`
	Email          sql.NullString `json:"email"`
	OfficeCode     sql.NullString `json:"officeCode"`
	ReportsTo      sql.NullInt64  `json:"reportsTo"`
	JobTitle       sql.NullString `json:"jobTitle"`
}

type Employees struct {
	Employees []Employee `json:"employees"`
}

func main() {

	// Establish DB connection
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("/ - Aloha!")

		rows, err := db.Query("SELECT * FROM employees LIMIT 10")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Employees{}

		for rows.Next() {
			employee := Employee{}
			if err := rows.Scan(&employee.EmployeeNumber, &employee.LastName, &employee.FirstName, &employee.Extension, &employee.Email, &employee.OfficeCode, &employee.ReportsTo, &employee.JobTitle); err != nil {
				return err // Exit if we get an error
			}

			// Append Employee to Employees
			result.Employees = append(result.Employees, employee)
		}
		return c.JSON(result)
	})

	// Get all records from MySQL
	app.Get("/employee", func(c *fiber.Ctx) error {
		// Get Employee list from database
		rows, err := db.Query("SELECT * FROM employees order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		// Return Employees in JSON format
		return c.JSON(rows)
	})

	log.Fatal(app.Listen(":3000"))

}
