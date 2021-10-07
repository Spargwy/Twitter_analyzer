package main

import (
	"dev-team/api"
	"dev-team/storage"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAccountsFromEnv() []string {
	usernames := os.Getenv("USERNAMES")
	usernames = strings.Replace(usernames, " ", "", -1)
	usernamesArray := strings.Split(usernames, ",")
	return usernamesArray
}

func main() {
	connectionString := os.Getenv("CONN")
	var database storage.PostgreSQL
	accounts := GetAccountsFromEnv()
	err := database.Init(connectionString)
	if err != nil {
		log.Fatalf("ERROR IN INIT DB: %v", err)
	}
	log.Print("\nDatabase init\n")

	err = database.CreateTables()
	if err != nil {
		log.Fatal("CAN NOT CREATE TABLES: ", err)
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})
	app.Get("/handler", func(c *fiber.Ctx) error {
		accounts, err := api.GetAllAccountsData(accounts)
		if err != nil {
			log.Printf("ERROR IN GET ACCOUNT DATA: %v\n", err)
			return err
		}
		if err == nil {
			err = database.InsertData(accounts)
			if err != nil {
				log.Println("CAN NOT INSERT DATA: ", err)
				return err
			}
			return c.JSON(accounts)
		}
		return c.SendString("Rate limit exceeded")
	})
	log.Fatal(app.Listen(":8900"))
}
