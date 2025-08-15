package main

import (
	"Zota/db"
	"Zota/initializer"
	"Zota/model"
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

	tableError := initializer.InitSchema(pool)
	if tableError != nil {
		return
	}

	store := model.NewStore(pool)
	if storeError := store.Set("Key", "Apple"); storeError != nil {
		log.Fatal(storeError)
	}

	var getValue, _ = store.Get("Key")
	log.Default().Println("Key:", getValue)

	err = store.Delete("Key")
	if err != nil {
		return
	}

	defer pool.Close()

}
