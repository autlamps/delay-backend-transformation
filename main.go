package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/autlamps/delay-backend-transformation/database"
	"github.com/autlamps/delay-backend-transformation/input"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var API_KEY string
var DB_URL string
var ID_M map[string]uuid.UUID

func init() {
	flag.StringVar(&API_KEY, "API_KEY", "", "at api key")
	flag.StringVar(&DB_URL, "DB_URL", "", "DB URL")
	flag.Parse()

	if API_KEY == "" {
		API_KEY = os.Getenv("API_KEY")
	}
	if DB_URL == "" {
		DB_URL = os.Getenv("DB_URL")
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	t := time.Now()
	ID_M = make(map[string]uuid.UUID)

	ag := make(chan update.AGEntities)
	ro := make(chan update.ROEntities)
	ca := make(chan update.CAEntities)
	st := make(chan update.STEntities)
	tr := make(chan update.TREntities)
	stt := make(chan update.STTEntities)

	go getAgency(ag)
	go getRoute(ro)
	go getCalender(ca)
	go getStops(st)
	go getTrip(tr)
	go getStopTrips(stt)

	AGent := <-ag
	ROent := <-ro
	CAent := <-ca
	STent := <-st
	TRent := <-tr
	STTent := <-stt

	//STTent := getStopTrips()
	//fmt.Println(time.Now().Sub(t))

	db := database.CreateCon(DB_URL)
	input.AgIn(AGent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
	input.RoIn(ROent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
	input.CaIn(CAent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
	input.StIn(STent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
	input.TrIn(TRent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
	input.SttIn(STTent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Calls the AT API for Agency List and then returns AGEntities
func getAgency(c chan update.AGEntities) {
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/agency?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		log.Fatalf("Failed to call api: %v", err)
	}

	if resp.StatusCode == 403 {
		log.Fatal("API key does not work (Agency)")
	}

	var ag update.AGAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ag)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get Agencies done: ", time.Now().Sub(t))
	c <- ag.Entities
}

// Calls the AT API for the Route list and then returns ROEntities
func getRoute(c chan update.ROEntities) {
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/routes?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		log.Fatalf("Failed to call api: %v", err)
	}

	if resp.StatusCode == 403 {
		log.Fatal("API key does not work (Route)")
	}

	var ro update.ROAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ro)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Get Routes done: ", time.Now().Sub(t))
	c <- ro.Entities
}

func getTrip(c chan update.TREntities) {
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/trips?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		log.Fatalf("Failed to call api: %v", err)
	}

	if resp.StatusCode == 403 {
		log.Fatal("API key does not work (Trip)")
	}

	var tr update.TRAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Get Trips done: ", time.Now().Sub(t))
	c <- tr.Entities
}

func getCalender(c chan update.CAEntities) {
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/calendar?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 403 {
		fmt.Println("API key does not work (Calender)")
	}

	var ca update.CAAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ca)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get Calendars done: ", time.Now().Sub(t))
	c <- ca.Entities
}

func getStops(c chan update.STEntities) {
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/stops?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 403 {
		fmt.Println("API key does not work (stops)")
	}

	var st update.STAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&st)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Get Stops done: ", time.Now().Sub(t))
	c <- st.Entities
}

func getStopTrips(c chan update.STTEntities) {
	t := time.Now()
	resp, err := http.Get("https://cdn01.at.govt.nz/data/stop_times.txt")
	fmt.Println("Got StopTime file: ", time.Now().Sub(t))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)

	var ste update.STTEntities

	i := 0

	for {
		r, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		if i == 0 {
			i++
			continue
		}

		seq64, err := strconv.ParseInt(r[4], 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		seq := int(seq64)

		nw := update.STTEntity{
			TripID:       r[0],
			ArrivalTime:  ParseTime(r[1]),
			DepatureTime: ParseTime(r[2]),
			StopID:       r[3],
			StopSequence: seq,
		}

		ste = append(ste, nw)

		i++
	}
	fmt.Println("Get StopTimes done: ", time.Now().Sub(t))
	c <- ste
}

// TODO: move this into its a package that makes more sense
// ParseTime parses GTFS time which can be greater than 24 hours to signify trips occurring over midnight and therefore
// multiple days. However golang and psql dont recognize this so we must convert it to usable time
func ParseTime(st string) time.Time {
	ib := []byte(st)

	// Check to see if the first char of our string is greater than the char 1 then see if
	// the second char is greater than the char 3. If so our time is greater than 24 hours
	matched := ib[0] > 49 && ib[1] > 51

	if matched {
		str := strings.Split(st, ":")

		ot, err := strconv.Atoi(str[0])

		if err != nil {
			log.Fatal(err)
		}

		ot = ot - 24
		sot := strconv.Itoa(ot)

		st = fmt.Sprintf("%v:%v:%v", sot, str[1], str[2])
	}

	dt, err := time.Parse("15:04:05", st)

	if err != nil {
		log.Fatal(err)
	}
	return dt
}
