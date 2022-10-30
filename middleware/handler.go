package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/chandanacharya1/customer-matching/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres golang driver
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type coordinate struct {
	lat float64
	lng float64
}

// response format
type response struct {
	ID              int64   `json:"id,omitempty"`
	Name            string  `json:"name,omitempty"`
	OperatingRadius int64   `json:"OperatingRadius,omitempty"`
	Rating          float64 `json:"rating,omitempty"`
	Distance        float64 `json:"distance,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// get env variables from .env
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DBNAME")

	//database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, password, port)
	// Open the connection
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}
func CustomerRequest(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	// decode the json request to stock
	err := json.NewDecoder(r.Body).Decode(&customer)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	fmt.Println(customer.Material)
	// call insert stock function and pass the stock
	partners := GetMatchingPartners(customer)
	fmt.Println(len(partners))
	var result []response
	for _, partner := range partners {
		// find distance of each partner from customer's coordinates
		dist := distance(customer.AddressLat, customer.AddressLong, partner.AddressLat, partner.AddressLong, "K")
		// format a response object
		res := response{
			ID:              partner.PartnerID,
			Name:            partner.PartnerName,
			OperatingRadius: partner.OperatingRadius,
			Rating:          partner.Rating,
			Distance:        dist,
		}
		if dist < float64(res.OperatingRadius) {
			result = append(result, res)
		}
	}
	// if partners have same rating, sort by distance
	for i := 0; i < len(result)-1; i++ {
		if result[i].Rating == result[i+1].Rating && result[i].Distance > result[i+1].Distance {
			result[i], result[i+1] = result[i+1], result[i]
		}
	}
	fmt.Println(result)
	// send the response
	json.NewEncoder(w).Encode(result)

}

func GetMatchingPartners(customer models.Customer) []models.Partner {
	db := createConnection()
	defer db.Close()
	sqlStatement := "SELECT p.id, p.name, p.radius, p.rating, a.lattitude, a.longitude " +
		"FROM partner p, materials m, address a, partner_materials pm, partner_address pa " +
		"WHERE p.id = pm.partnerid " +
		"AND m.id = pm.materialid " +
		"AND p.id = pa.partnerid " +
		"AND a.id = pa.addressid " +
		"AND m.name=$1 ORDER BY p.rating DESC"

	rows, err := db.Query(sqlStatement, customer.Material)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	partners := make([]models.Partner, 0)
	for rows.Next() {
		var partner models.Partner
		err := rows.Scan(&partner.PartnerID, &partner.PartnerName, &partner.OperatingRadius, &partner.Rating,
			&partner.AddressLat, &partner.AddressLong)
		if err != nil {
			log.Fatalf("Unable to scan the query result. %v", err)
		}
		partners = append(partners, partner)
	}
	return partners
}

func ListPartners(w http.ResponseWriter, r *http.Request) {
	var partners []models.Partner
	partners = GetAllPartners()
	json.NewEncoder(w).Encode(partners)
}

func GetAllPartners() []models.Partner {
	db := createConnection()
	defer db.Close()
	sqlStatement := "SELECT p.id, p.name, p.radius, p.rating, a.lattitude, a.longitude " +
		"FROM partner p, address a, partner_address pa " +
		"WHERE p.id = pa.partnerid " +
		"AND a.id = pa.addressid "
	/*		"ORDER BY p.rating DESC"*/

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	partners := make([]models.Partner, 0)
	for rows.Next() {
		var partner models.Partner
		err := rows.Scan(&partner.PartnerID, &partner.PartnerName, &partner.OperatingRadius, &partner.Rating,
			&partner.AddressLat, &partner.AddressLong)
		if err != nil {
			log.Fatalf("Unable to scan the query result. %v", err)
		}
		partners = append(partners, partner)
	}
	return partners
}

func GetPartner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	partnerId, err := strconv.Atoi(vars["partnerid"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var partner models.Partner
	partner = GetPartnerFromId(partnerId)
	if len(partner.PartnerName) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		json.NewEncoder(w).Encode(partner)
	}
}

func GetPartnerFromId(id int) models.Partner {
	db := createConnection()
	defer db.Close()
	var partner models.Partner
	sqlStatement := "SELECT p.id, p.name, p.radius, p.rating, a.lattitude, a.longitude " +
		"FROM partner p, address a, partner_address pa " +
		"WHERE p.id = pa.partnerid " +
		"AND a.id = pa.addressid " +
		"AND p.id = $1 " +
		"ORDER BY p.rating DESC"

	err := db.QueryRow(sqlStatement, id).Scan(&partner.PartnerID, &partner.PartnerName, &partner.OperatingRadius,
		&partner.Rating, &partner.AddressLat, &partner.AddressLong)
	if err != nil {
		log.Fatalf("Unable to scan the query result. %v", err)
	}
	return partner
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}
