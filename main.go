package main

import (
	"flag"
	"log"
	"os"

	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/autlamps/delay-backend-transformation/database"
	"github.com/autlamps/delay-backend-transformation/input"
	"github.com/autlamps/delay-backend-transformation/update"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
)

var API_KEY string
var DB_URL string
var HTTPForbiden = errors.New("HTTP 403 - Failed")

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
	version, err := getVersion()
	fmt.Println("Current GTFS verison: ", version)

	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	ag := make(chan update.AGReturn)
	ro := make(chan update.ROReturn)
	ca := make(chan update.CAReturn)
	st := make(chan update.STReturn)
	tr := make(chan update.TRReturn)
	stt := make(chan update.STTReturn)

	go getAgency(ag, version)
	go getRoute(ro, version)
	go getCalender(ca, version)
	go getStops(st, version)
	go getTrip(tr, version)
	go getStopTrips(stt, version)

	AGret := <-ag
	if AGret.Error != nil {
		log.Fatal(AGret.Error)
	}
	AGent := AGret.Entities

	ROret := <-ro
	if ROret.Error != nil {
		log.Fatal(ROret.Error)
	}
	ROent := ROret.Entities

	CAret := <-ca
	if CAret.Error != nil {
		log.Fatal(CAret.Error)
	}
	CAent := CAret.Entities

	STret := <-st
	if STret.Error != nil {
		log.Fatal(STret.Error)
	}
	STent := STret.Entities

	TRret := <-tr
	if TRret.Error != nil {
		log.Fatal(TRret.Error)
	}
	TRent := TRret.Entities

	STTret := <-stt
	if STTret.Error != nil {
		log.Fatal(STTret.Error)
	}
	STTent := STTret.Entities

	db := database.CreateCon(DB_URL)

	is := input.InService{
		Db:         db,
		AgencyMap:  nil,
		ServiceMap: nil,
		RouteMap:   nil,
		TripMap:    nil,
		StopMap:    nil,
	}

	is.Init()

	is.AgIn(AGent)
	fmt.Println("Done Agency: ", time.Now().Sub(t))
	is.RoIn(ROent)
	fmt.Println("Done Routes: ", time.Now().Sub(t))
	is.CaIn(CAent)
	fmt.Println("Done Calender: ", time.Now().Sub(t))
	is.StIn(STent)
	fmt.Println("Done Stops: ", time.Now().Sub(t))
	is.TrIn(TRent)
	fmt.Println("Done Trips: ", time.Now().Sub(t))
	is.SttIn(STTent)
	fmt.Println("Done StopTime: ", time.Now().Sub(t))

}

func getVersion() (ver string, err error) {
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/versions?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		return "", err
	}

	if resp.StatusCode == 403 {
		return "", HTTPForbiden
	}

	var vr update.VRAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&vr)

	if err != nil {
		return "", err
	}

	return vr.ATVerDet[0].ATVersion, nil
}

// Calls the AT API for Agency List and then returns AGEntities
func getAgency(c chan<- update.AGReturn, ver string) {
	var toChan update.AGReturn
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/agency?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	if resp.StatusCode == 403 {
		toChan.Error = HTTPForbiden
		c <- toChan
	}

	var ag update.AGAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ag)

	if err != nil {
		toChan.Error = err
		c <- toChan
	}
	fmt.Println("Get Agencies done: ", time.Now().Sub(t))

	toChan.Entities = ag.Entities

	c <- toChan
}

// Calls the AT API for the Route list and then returns
// ROEntities
func getRoute(c chan update.ROReturn, ver string) {
	var toChan update.ROReturn
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/routes?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	if resp.StatusCode == 403 {
		toChan.Error = HTTPForbiden
		c <- toChan
	}

	var ro update.ROAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ro)

	final := update.ROEntities{}

	for _, a := range ro.Entities {
		if strings.Contains(a.RouteID, ver) {
			final = append(final, a)
		}

	}

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	fmt.Println("Get Routes done: ", time.Now().Sub(t))

	toChan.Entities = final
	c <- toChan
}

func getTrip(c chan update.TRReturn, ver string) {
	var toChan update.TRReturn
	t := time.Now()
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/trips?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	if resp.StatusCode == 403 {
		toChan.Error = HTTPForbiden
		c <- toChan
	}

	var tr update.TRAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tr)

	final := update.TREntities{}

	for _, a := range tr.Entities {
		if strings.Contains(a.TripID, ver) {
			final = append(final, a)
		}

	}

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	toChan.Entities = final

	fmt.Println("Get Trips done: ", time.Now().Sub(t))

	c <- toChan
}

func getCalender(c chan update.CAReturn, ver string) {
	t := time.Now()
	var toChan update.CAReturn
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/calendar?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	if resp.StatusCode == 403 {
		fmt.Println("API key does not work (Calender)")
	}

	var ca update.CAAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ca)

	final := update.CAEntities{}

	for _, a := range ca.Entities {
		if strings.Contains(a.ServiceID, ver) {
			final = append(final, a)
		}

	}

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	toChan.Entities = final

	fmt.Println("Get Calendars done: ", time.Now().Sub(t))
	c <- toChan
}

func getStops(c chan update.STReturn, ver string) {
	t := time.Now()
	var toChan update.STReturn
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/stops?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	if resp.StatusCode == 403 {
		toChan.Error = HTTPForbiden
		c <- toChan
	}

	var st update.STAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&st)

	final := update.STEntities{}

	for _, a := range st.Entities {
		if strings.Contains(a.StopID, ver) {
			final = append(final, a)
		}

	}
	if err != nil {
		toChan.Error = err
		c <- toChan
	}

	toChan.Entities = final

	fmt.Println("Get Stops done: ", time.Now().Sub(t))
	c <- toChan
}

func getStopTrips(c chan update.STTReturn, ver string) {
	t := time.Now()
	var toChan update.STTReturn
	resp, err := http.Get("https://cdn01.at.govt.nz/data/stop_times.txt")
	fmt.Println("Got StopTime file: ", time.Now().Sub(t))

	if err != nil {
		toChan.Error = err
		c <- toChan
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
			toChan.Error = err
			c <- toChan
		}

		if i == 0 {
			i++
			continue
		}

		if !strings.Contains(r[0], ver) {
			i++
			continue
		}

		seq64, err := strconv.ParseInt(r[4], 10, 64)

		if err != nil {
			toChan.Error = err
			c <- toChan
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

	toChan.Entities = ste
	fmt.Println("Get StopTimes done: ", time.Now().Sub(t))
	c <- toChan
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
