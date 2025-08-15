package main

import (
	"Zota/db"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	cfg := db.LoadConfig()

	pool, err := db.ConnectDB(cfg.ConnectionString())
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer pool.Close()

}
