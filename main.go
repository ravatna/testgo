package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

type Data struct {
	ActivePower int `json:"active_power"`
	PowerInput  int `json:"power_input"`
}

func generateMockData(db *sql.DB, numRows int) error {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numRows; i++ {
		activePower := rand.Intn(999) + 1
		powerInput := rand.Intn(999) + 1

		_, err := db.Exec("INSERT INTO power_data (active_power, power_input) VALUES (?, ?)", activePower, powerInput)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	app := fiber.New()

	// Create a MySQL database connection
	dsn := "root:rootpassword@tcp(mysql:3306)/dbname" // "mysql" is the service name in Docker Compose
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS power_data (
		active_power INT,
		power_input INT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// seed data to database (100 rows)
	err = generateMockData(db, 100)
	if err != nil {
		log.Fatal(err)
	}

	// GET /api/query
	app.Get("/api/query", func(c *fiber.Ctx) error {
		start, _ := strconv.Atoi(c.Query("start"))
		end, _ := strconv.Atoi(c.Query("end"))

		rows, err := db.Query("SELECT SUM(active_power), SUM(power_input) FROM power_data WHERE active_power BETWEEN ? AND ?", start, end)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var sumActivePower int
		var sumPowerInput int

		for rows.Next() {
			err := rows.Scan(&sumActivePower, &sumPowerInput)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Create JSON response
		data := Data{
			ActivePower: sumActivePower,
			PowerInput:  sumPowerInput,
		}
		response, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		// response json
		c.Set("Content-Type", "application/json")
		return c.Send(response)
	})

	log.Fatal(app.Listen(":3000"))
}
