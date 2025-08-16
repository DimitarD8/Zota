package main

import (
	"Zota/db"
	"Zota/initializer"
	"Zota/model"
	"fmt"
	"log"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	cfg := db.LoadConfig()

	pool, err := db.ConnectDB(cfg.ConnectionString())
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	tableError := initializer.InitSchema(pool)
	if tableError != nil {
		return
	}

	store := model.NewStore(pool)
	if storeError := store.Put("Key", "Apple"); storeError != nil {
		log.Fatal(storeError)
	}

	if storeError := store.Put("1", "Cat"); storeError != nil {
		log.Fatal(storeError)
	}

	if storeError := store.Put("2", "Dog"); storeError != nil {
		log.Fatal(storeError)
	}

	if storeError := store.Put("3", "3"); storeError != nil {
		log.Fatal(storeError)
	}

	var getValue, _ = store.Get("2")
	log.Default().Println("Key:2 Value:", getValue)

	err = store.Delete("2")
	if err != nil {
		return
	}

	var currentState, _ = store.Dumb()
	for k, v := range currentState {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			data, err := store.Dumb()
			if err != nil {
			}
			for k, v := range data {
				fmt.Printf("key[%s] value[%s]\n", k, v)
			}
		}
	}
	defer pool.Close()

}
