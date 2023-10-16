package main

import (
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	if db := utils.SetupDB(); db != nil {
		log.Fatal("db is nil. How did this happen?")
	}
}
