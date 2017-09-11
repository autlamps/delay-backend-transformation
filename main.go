package main

import (
	"fmt"

	"encoding/json"
	"github.com/autlamps/delay-backend-transformation/update"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"flag"
	"os"
	"database/sql"
	"github.com/google/uuid"
	"github.com/autlamps/delay-backend-transformation/database"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "gtfs"
)

var API_KEY string

func init() {
	flag.StringVar(&API_KEY, "API_KEY", "", "at api key")
	flag.Parse()

	if API_KEY == "" {
		API_KEY = os.Getenv("API_KEY")
	}

}

func main() {
	ent, err := getAgency()

	if err != "" {
		fmt.Println(err)
	}
	db := database.CreateCon()
	insertDB(ent, db)
}

// Calls the AT API for Agency List and then creates agency classes
func getAgency() (update.AUEntities, string) {
	urlWithKey := fmt.Sprintf("http://api.at.govt.nz/v1/gtfs/agency?api_key=%v", API_KEY)

	resp, err := http.Get(urlWithKey)

	if err != nil {
		log.Fatalf("Failed to call api: %v",err)
	}

	if resp.StatusCode == 403 {
		fmt.Println("API key does not work")
	}

	var ag update.AUAPIResponse

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ag)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(ag.Entities)
	return ag.Entities, ag.Error
}

func insertDB(entities update.AUEntities, db *sql.DB)  {

	for i := 0; i < len(entities); i++ {
		gtfs_agency_id := entities[i].AgencyID
		agency_name := entities[i].AgencyName
		agen_contents_uuid, err := uuid.NewRandom()
		if err != nil {
			fmt.Println("No new UUID")
			log.Fatal(err.Error())
		}

		printls, err := db.Exec("INSERT INTO agency (gtfs_agency_id, agency_name, agency_id) VALUES ($1, $2, $3);", gtfs_agency_id, agency_name, agen_contents_uuid )
		if err != nil {
			log.Fatal(err.Error())
		}
		printls.LastInsertId()
	}
	fmt.Println("Done")
}