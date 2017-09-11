package main

import (
	"fmt"

	"encoding/json"
	"flag"
	"github.com/autlamps/delay-backend-transformation/database"
	"github.com/autlamps/delay-backend-transformation/input"
	"github.com/autlamps/delay-backend-transformation/update"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
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

	db := database.CreateCon(DB_URL)
	input.AgIn(AGent, db, ID_M)
	input.RoIn(ROent, db, ID_M)
	input.CaIn(CAent, db, ID_M)
	input.StIn(STent, db, ID_M)
	input.TrIn(TRent, db, ID_M)
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

	fmt.Println(ag.Entities)
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
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/calendar?api_key=%v", API_KEY)

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
