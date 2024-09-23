package main

import (
	"log"

	"github.com/aishwaryaBuildd/go_boiler_plate/server"
	"github.com/aishwaryaBuildd/go_boiler_plate/store/mysql"
)

func main() {
	db, err := mysql.NewDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	store := mysql.NewStore(db)
	r := server.NewServer(store)
	r.Run(":8080")
}
