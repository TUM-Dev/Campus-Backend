package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"

	pb "github.com/TUM-Dev/Campus-Backend/api"

	"github.com/TUM-Dev/Campus-Backend/model/heatmap/dbservice"
)

type JsonEntry struct {
	Intensity float64
	Lat       float64
	Long      float64
	Floor     string
}

const (
	heatmapDB = "./data/sqlite/heatmap.db"
)

// Retrieves all access points from the database
// and stores them in JSON format in `dst` e.g. "data/json/ap.json"
func saveAPsToJSON(dst string, totalLoad int) {
	APs := dbservice.RetrieveAPsOfTUM(true)
	var jsonData []JsonEntry

	for _, ap := range APs {
		currTotLoad := 0
		var intensity float64 = float64(currTotLoad) / float64(totalLoad)
		lat, _ := strconv.ParseFloat(ap.Lat, 64)
		lng, _ := strconv.ParseFloat(ap.Long, 64)
		jsonEntry := JsonEntry{intensity, lat, lng, ap.Floor}
		jsonData = append(jsonData, jsonEntry)
	}

	bytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dst, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	pb.CampusServer
}


func forecast() {
	path := "./forecast.py"
	cmd := exec.Command("python", "-u", path)
	out, err := cmd.Output()

	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(out))
}
