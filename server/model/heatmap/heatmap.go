package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"time"

	
	"google.golang.org/protobuf/types/known/emptypb"

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

func (s *server) GetAccessPoint(ctx context.Context, in *pb.APRequest) (*pb.AccessPoint, error) {
	name := in.Name
	ts := in.Timestamp
	log.Printf("Received request for AP with name: %s and timestamp: %s", name, in.Timestamp)

	db := dbservice.InitDB(heatmapDB)

	day, hr := getDayAndHourFromTimestamp(ts)
	ap := dbservice.GetHistoryForSingleAP(name, day, hr)
	db.Close()

	load, _ := strconv.Atoi(ap.Load)

	return &pb.AccessPoint{
		Name:      ap.Name,
		Lat:       ap.Lat,
		Long:      ap.Long,
		Intensity: int64(load),
		Max:       int64(ap.Max),
		Min:       int64(ap.Min),
	}, nil
}

func getDayAndHourFromTimestamp(timestamp string) (int, int) {
	ts := strings.Split(timestamp, " ")
	date := ts[0]
	yearMonthDay := strings.Split(date, "-")
	day, err := strconv.Atoi(yearMonthDay[2])
	if err != nil {
		day = 0
	}
	hr, err := strconv.Atoi(ts[1])
	if err != nil {
		hr = 0
	}
	today := time.Now().Day()
	day = day - today
	return day, hr
}

type location struct {
	lat  string
	long string
}

var locations map[string]location

func (s *server) ListAccessPoints(in *pb.APRequest, stream pb.Campus_ListAccessPointsServer) error {
	ts := in.Timestamp
	day, hr := getDayAndHourFromTimestamp(ts)

	apList := dbservice.GetHistoryForAllAPs(day, hr)

	log.Printf("Sending %d APs ...", len(apList))

	locations = make(map[string]location)
	accessPoints := dbservice.RetrieveAPsOfTUM(true)
	for _, ap := range accessPoints {
		locations[ap.Name] = location{ap.Lat, ap.Long}
	}

	for _, ap := range apList {
		location := locations[ap.Name]

		load, _ := strconv.Atoi(ap.Load)

		if err := stream.Send(
			&pb.APResponse{
				Accesspoint: &pb.AccessPoint{
					Name:      ap.Name,
					Lat:       location.lat,
					Long:      location.long,
					Intensity: int64(load),
					Max:       int64(ap.Max),
					Min:       int64(ap.Min),
				},
			}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ListAllAPNames(in *emptypb.Empty, stream pb.Campus_ListAllAPNamesServer) error {
	names := dbservice.GetAllNames()
	for _, name := range names {
		if err := stream.Send(
			&pb.APName{
				Name: name,
			}); err != nil {
			return err
		}
	}
	return nil
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
