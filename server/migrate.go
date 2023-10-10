package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	if db := setupDB(); db != nil {
		log.Fatal("db is nil. How did this happen?")
	}
}
