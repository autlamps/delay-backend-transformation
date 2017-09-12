package main

import (
	"fmt"

	"encoding/csv"
	"encoding/json"
	"flag"
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
	"regexp"
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
	t := time.Now()
	ID_M = make(map[string]uuid.UUID)
	AGent, err := getAgency()
	emptyErr(err)
	ROent, err := getRoute()
	emptyErr(err)
	CAent, err := getCalender()
	emptyErr(err)
	STent, err := getStops()
	emptyErr(err)
	TRent, err := getTrip()
	emptyErr(err)
	STTent := getStopTrips()

	db := database.CreateCon(DB_URL)
	input.AgIn(AGent, db, ID_M)
	input.RoIn(ROent, db, ID_M)
	input.CaIn(CAent, db, ID_M)
	input.StIn(STent, db, ID_M)
	input.TrIn(TRent, db, ID_M)
	input.SttIn(STTent, db, ID_M)
	fmt.Println(time.Now().Sub(t))
}

func emptyErr(err string) {
	if err != "" {
		fmt.Println(err)
	}
}

// Calls the AT API for Agency List and then returns AGEntities
func getAgency() (update.AGEntities, string) {
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

	return ag.Entities, ag.Error
}

func getRoute() (update.ROEntities, string) {
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

	return ro.Entities, ro.Error
}

func getTrip() (update.TREntities, string) {
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

	return tr.Entities, tr.Error
}

func getCalender() (update.CAEntities, string) {
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

	return ca.Entities, ca.Error
}

func getStops() (update.STEntities, string) {
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

	return st.Entities, st.Error
}

func getStopTrips() update.STTEntities {
	resp, err := http.Get("https://cdn01.at.govt.nz/data/stop_times.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)

	//for {
	//	records, err := r.Read()
	//	if err == io.EOF {
	//		return stt
	//	}
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	for i := 1; i < len(records); i++ {
	//		recsp := strings.Split(records[i], ",")
	//		for k := 0; k < len(recsp); k++ {
	//			trip_id := m[recsp[0]]
	//			arrival_time := recsp[1]
	//			departure_time := recsp[2]
	//			stop_id := m[recsp[3]]
	//			stop_sequence := recsp[4]
	//			stt = update.STTEntity{TripID:trip_id, ArrivalTime:arrival_time, DepatureTime:departure_time, StopID:stop_id, StopSequence:stop_sequence}
	//
	//		}
	//	}
	//}

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

		//at, err := time.Parse("15:04:05", r[1])
		//
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//dt, err := time.Parse("15:04:05", r[2])
		//
		//if err != nil {
		//	log.Fatal(err)
		//}

		seq64, err := strconv.ParseInt(r[4], 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		seq := int(seq64)

		nw := update.STTEntity{
			TripID:       r[0],
			ArrivalTime:  hrCheck(r[1]),
			DepatureTime: hrCheck(r[2]),
			StopID:       r[3],
			StopSequence: seq,
		}

		ste = append(ste, nw)

		i++
	}
	return ste
}

func hrCheck(input string) time.Time {
	matched, err := regexp.MatchString(`[2][4-9]:[\s\S][\s\S]:[\s\S][\s\S]`, input)
	if err != nil {
		log.Fatal(err)
	}
	if matched == true {
		str := strings.Split(input, ":")

		ot, err := strconv.Atoi(str[0])

		if err != nil {
			log.Fatal(err)
		}

		ot = ot - 24
		sot := strconv.Itoa(ot)

		input = fmt.Sprintf("%v:%v:%v", sot, str[1], str[2])

	}

	dt, err := time.Parse("15:04:05", input)

	if err != nil {
		log.Fatal(err)
	}
	return dt
}
