package main

import (
	json2 "encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	api2 "tracker/api/api"
	db2 "tracker/internal/app/db"
	"tracker/internal/app/models"
)

func json(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func main() {

	port := os.Getenv("port")

	db := db2.DBConnect()
	api := &api2.Api{DB: db}

	http.HandleFunc("/trackers_locations", func(w http.ResponseWriter, r *http.Request) {
		json(w)

		positions := map[uint16]*models.SrPosData{}
		trackerIds := r.URL.Query()["id[]"]
		for _, _id := range trackerIds {
			id, _ := strconv.Atoi(_id)

			posData, err := api.GetLastTrackerPosition(uint16(id))
			if err != nil {
				panic(err)
			}

			positions[uint16(id)] = posData

		}

		js, err2 := json2.Marshal(positions)
		if err2 != nil {
			panic(err2)
		}

		w.Write([]byte(js))

	})

	http.HandleFunc("/tracker_data", func(w http.ResponseWriter, r *http.Request) {
		json(w)

		trackerId := r.URL.Query().Get("id")
		dateFrom := r.URL.Query().Get("dateFrom")
		dateTo := r.URL.Query().Get("dateTo")

		tId, _ := strconv.Atoi(trackerId)

		if tId == 0 || dateFrom == "" || dateTo == "" {
			w.WriteHeader(400)
			w.Write([]byte("{}"))
			return
		}

		posData, err := api.GetTrackerGPSData(uint16(tId), dateFrom, dateTo)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("{}"))
			return
		}

		js, err2 := json2.Marshal(posData)
		if err2 != nil {
			panic(err2)
		}

		w.Write([]byte(js))
	})

	http.HandleFunc("/trackers", func(w http.ResponseWriter, r *http.Request) {
		json(w)

		trackers, err := api.GetTrackers()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("{}"))
			return
		}

		response := map[uint16]models.Tracker{}
		for _, tracker := range trackers {
			response[tracker.ID] = tracker
		}

		js, err2 := json2.Marshal(response)
		if err2 != nil {
			panic(err2)
		}

		w.Write([]byte(js))
	})

	log.Println("Start server on port", port)
	log.Fatalln(http.ListenAndServe("0.0.0.0:"+port, nil))

}
