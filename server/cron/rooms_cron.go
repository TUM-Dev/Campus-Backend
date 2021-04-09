package cron

import (
	"github.com/TUM-Dev/Campus-Backend/base"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/getsentry/sentry-go"
	"log"
)

const roomFinderAPIURLBase = "https://portal.mytum.de/campus/roomfinder/"
const roomFinderEndpointMaps = "getMaps"
const roomFinderEndpointBuildingMaps = "getBuildingMaps"
const roomFinderEndpointBuildingData = "getBuildingData"
const roomFinderEndpointBuildingDefaultMap = "getBuildingDefaultMap"
const roomFinderEndpointBuildingMap = "getBuildingMap?b_id="

func (c ServiceCron) fetchRooms() {
	log.Printf("fetching rooms from TUMonline.")
	var oldRooms []model.RoomfinderRooms
	c.DB.Find(&oldRooms)
	_, _ = getMaps()
}

func getMaps() (maps []model.RoomfinderMaps, err error) {
	//var res []model.RoomfinderMaps
	resp, err := base.XmlRpcCall(
		roomFinderAPIURLBase,
		roomFinderEndpointMaps,
		struct{ Who string }{}) // todo, what might be a reasonable reply type to unmarshal to?
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	log.Printf("%v", resp)
	return nil, nil
}
