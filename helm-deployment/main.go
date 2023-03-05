package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// Establish DB connection
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)")
	if err != nil {
		fmt.Println("ERROR:", err)
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Connection opened:", db)

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("/ - Aloha!")
		return c.SendString("Aloha, World!")
	})

	log.Fatal(app.Listen(":3000"))

}
