package main

import (
	"flag"
	"log"
)


func main() {
	deleteFlag := flag.Bool("delete", false, "Set to true to delete the 'usf-courses' collection")
	flag.Parse()
	
	db, err := Start(*deleteFlag)
	if err != nil{
		log.Fatalf("Error starting program: %v\n", err)
	}

	StartUserInterface(db)
}