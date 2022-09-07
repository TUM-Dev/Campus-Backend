package LRZscraper

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TUM-Dev/Campus-Backend/model/heatmap/DBService"
)

type NetworkLoad struct {
	Datapoints [][]float64 `json:"datapoints"`
	Target     string      `json:"target"`
}

type history [31][24]int

// History data for an Access Point
type History struct {
	name       string
	historyPtr *history
	max        int
	min        int
}

func getNetworkAlias(apName, network string) string {
	return fmt.Sprintf("alias(ap.%s.ssid.%s,%%22%s%%22)", apName, network, network)
}

func BuildURL(apName string, days int, format string) (url string) {
	rendererEndpoint := "http://graphite-kom.srv.lrz.de/render/"
	width := "?width=640"
	height := "height=240"
	title := "title=SSIDs%20(weekly)"
	areaMode := "areaMode=stacked"
	xFormat := "xFormat=%25d.%25m."
	timezone := "tz=CET"
	from := fmt.Sprintf("from=-%ddays", days)
	networks := []string{"eduroam", "lrz", "mwn-events", "@BayernWLAN", "other"}
	aliases := ""
	for i, network := range networks {
		aliases += getNetworkAlias(apName, network)
		if i < len(network)-1 {
			aliases += ","
		}
	}
	target := fmt.Sprintf("target=cactiStyle(group(%s))", aliases)
	fontName := "fontName=Courier"
	format = fmt.Sprintf("format=%s", format)
	url = fmt.Sprintf(`%s%s&%s&%s&%s&%s&%s&%s&%s&%s&%s`,
		rendererEndpoint,
		width,
		height,
		title,
		areaMode,
		xFormat,
		timezone,
		from,
		fontName,
		target,
		format)
	return
}

func GetGraphiteDataForAP(apName string, from int, t *http.Transport) ([]NetworkLoad, error) {
	if !strings.HasPrefix(apName, "apa") {
		log.Fatalf("Name of the Access Point must start with \"apa\"!")
	}

	url := BuildURL(apName, from, "json")

	c := &http.Client{
		Transport: t,
	}

	resp, httpError := c.Get(url)
	if httpError != nil {
		log.Printf("Could not retrieve json data from URL! %q", httpError)
		return nil, httpError
	}
	defer resp.Body.Close()

	var networks []NetworkLoad
	err := json.NewDecoder(resp.Body).Decode(&networks)
	if err != nil {
		// sometimes lrz's server sends a bad JSON response
		log.Printf("Could not decode JSON response: %v", err)
		return nil, err
	}

	return networks, nil
}

func getTotalMaxMin(networks []NetworkLoad) (int, int) {
	totCurr, totMax, totMin := 0.0, 0.0, 0.0
	for _, network := range networks {
		curr, max, min := getCurrMaxMin(network.Target)
		totCurr += curr
		totMax += max
		totMin += min
	}

	mx, mn := int(totMax), int(totMin)
	return mx, mn
}

func getCurrMaxMin(networkLoad string) (curr, max, min float64) {
	fields := strings.Fields(networkLoad) //remove whitespaces

	currStr := strings.Split(fields[1], ":")[1]
	curr, err := strconv.ParseFloat(currStr, 32)
	if err != nil || math.IsNaN(curr) {
		curr = 0.0
	}

	maxStr := strings.Split(fields[2], ":")[1]
	max, err = strconv.ParseFloat(maxStr, 32)
	if err != nil || math.IsNaN(max) {
		max = 0.0
	}

	minStr := strings.Split(fields[3], ":")[1]
	min, err = strconv.ParseFloat(minStr, 32)
	if err != nil || math.IsNaN(min) {
		min = 0.0
	}

	return
}

func GetHistoriesFrom(from int) []History {
	APs := DBService.RetrieveAPsOfTUM(true)
	var wg sync.WaitGroup
	channel := make(chan History)

	t := http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
		MaxConnsPerHost:     50,
		MaxIdleConns:        50,
	}

	start := time.Now()
	log.Println("Start time:", start)
	wg.Add(len(APs))
	for _, ap := range APs {
		go func(apName string) {
			GetLast30DaysForAP(apName, from, channel, &t)
			wg.Done()
		}(ap.Name)
	}

	go func() {
		wg.Wait()
		close(channel)
		elapsed := time.Since(start)
		log.Println("Elapsed time:", elapsed)
	}()

	skippedAPs := make([]string, 0)
	histories := make([]History, 0)

	for history := range channel {
		if history.historyPtr == nil {
			skippedAPs = append(skippedAPs, history.name)
		} else {
			histories = append(histories, history)
		}
	}

	log.Println("Appended total: ", len(histories))
	log.Println("Started storing in DB:", time.Now())
	// storeHistories(histories)
	skipped := len(skippedAPs)
	fmt.Println("Total nr of skipped APs:", skipped, skippedAPs)

	// storeMaxMins(histories)
	return histories
}

func storeMaxMins(histories []History) {
	for _, apHistory := range histories {
		apName := apHistory.name
		max := apHistory.max
		min := apHistory.min
		DBService.UpdateMinMax("Max", apName, max)
		DBService.UpdateMinMax("Min", apName, min)
	}
}

func StoreHistories(histories []History) {
	for _, apHistory := range histories {
		apName := apHistory.name
		history := *apHistory.historyPtr
		storeHistoryOfAP(apName, history)
	}
}

func storeHistoryOfAP(apName string, history history) {
	days := len(history)
	hours := len(history[0])
	if days != 31 || hours != 24 {
		return
	}
	for day := 0; day < days; day++ {
		for hour := 0; hour < hours; hour++ {
			avg := history[day][hour]
			DBService.UpdateHistory(day, hour, avg, apName)
		}
	}
}

func GetTodaysData() map[string][24]int {
	todayAPs := GetHistoriesFrom(1)
	averages := make(map[string][24]int)
	for _, todayAP := range todayAPs {
		averages[todayAP.name] = (*todayAP.historyPtr)[0]
		fmt.Println("Today", todayAP.name, (*todayAP.historyPtr)[0])
	}
	return averages
}

func GetLast30DaysForAP(apName string, from int, c chan History, t *http.Transport) {
	networks, err := GetGraphiteDataForAP(apName, from, t)
	if len(networks) == 0 {
		fmt.Println("len(networks) = 0!")
		c <- History{apName, nil, 0, 0}
		return
	}

	max, min := getTotalMaxMin(networks)
	if err != nil {
		fmt.Printf("Skipping %s. Bad JSON response from server.\n", apName)
		c <- History{apName, nil, 0, 0}
		return
	}

	gesamt := networks[0].Datapoints
	n := len(networks)
	for i := 1; i < n; i++ {
		for j := range gesamt {
			lenOtherNetwork := len(networks[i].Datapoints)
			if j < lenOtherNetwork {
				gesamt[j][0] += networks[i].Datapoints[j][0]
			}
		}
	}

	history := calcHourlyAvgs(gesamt)
	c <- History{apName, &history, max, min}
}

// It calculates hourly averages of network load for the given AP data.
// Returns a 31x24 matrix, where cell (i,j) holds
// the hourly avg. of day i and hour j.
func calcHourlyAvgs(datapoints [][]float64) history {
	var history history //last 30 days + today
	var prevHour, prevDay *int
	day := 31
	cnt := 0
	avg := 0
	n := 0 // network load (Nr of connected devices)

	for _, datapoint := range datapoints {
		n += int(datapoint[0])
		ts := int(datapoint[1])

		t := getTimeFromTimestamp(ts)
		hour := t.Hour()
		currDay := t.Day()

		if prevHour == nil || prevDay == nil {
			prevHour = &hour
			prevDay = &currDay
		}

		if hour != *prevHour {
			*prevHour = hour
			avg = n / cnt
			history[day-1][hour] = avg
			cnt = 0
			n = 0
		} else {
			cnt += 1
		}

		if currDay != *prevDay {
			*prevDay = currDay
			day -= 1
		}
	}

	history[day-1][*prevHour] = avg
	return history
}

// parses a Unix timestamp (i.e. milliseconds from EPOCH) and
// returns it as time.Time
func getTimeFromTimestamp(timestamp int) time.Time {
	ts := fmt.Sprintf("%d", timestamp)
	t, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(t, 0)
	return tm
}
