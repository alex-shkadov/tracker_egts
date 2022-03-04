package api

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
	"tracker/internal/app/models"
	"tracker/internal/app/server"
)

type Api struct {
	OM *server.ObjectManager
	DB *gorm.DB
}

type GpsData struct {
	Ntm        time.Time `json:"ntm"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Spd        uint16    `json:"spd"`
	Alts       int32     `json:"alts"`
	Dir        byte      `json:"dir"`
	Satellites uint      `json:"satellites"`
}

func (api *Api) GetTrackers() ([]models.Tracker, error) {
	trackers := []models.Tracker{}

	tx := api.DB.Find(&trackers)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return trackers, nil
}

func (api *Api) GetLastTrackerPosition(trackerId uint16) (*models.SrPosData, error) {
	sdr := api.DB.Raw(""+
		"SELECT s.id, ntm, latitude, longitude, mv, bb, spd, alts, dir, dirh, odm, satellites "+
		"FROM service_data_records as sdr "+
		"JOIN sr_pos_data as s ON s.service_data_record_id = sdr.id "+
		"WHERE sdr.tracker_id = ? AND s.deleted_at IS NULL AND s.ntm < '2100-01-01' AND s.vld = 1 "+
		"ORDER BY s.ntm DESC "+
		"LIMIT 1", trackerId).Row()

	if sdr.Err() != nil && sdr.Err() == sql.ErrNoRows {
		return nil, nil
	}

	var id sql.NullInt64
	var ntm sql.NullTime
	var lat sql.NullFloat64
	var lng sql.NullFloat64
	var mv sql.NullByte
	var bb sql.NullByte
	var spd sql.NullInt16
	var alts sql.NullInt32
	var dir sql.NullByte
	var dirh sql.NullByte
	var odm sql.NullInt32
	var sat sql.NullByte

	err := sdr.Scan(&id, &ntm, &lat, &lng, &mv, &bb, &spd, &alts, &dir, &dirh, &odm, &sat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &models.SrPosData{
		ID:         uint64(id.Int64),
		Ntm:        ntm.Time,
		Latitude:   lat.Float64,
		Longitude:  lng.Float64,
		Mv:         mv.Byte,
		Bb:         bb.Byte,
		Spd:        uint16(spd.Int16),
		Alts:       alts.Int32,
		Dir:        dir.Byte,
		Dirh:       dirh.Byte,
		Odm:        uint32(odm.Int32),
		Satellites: uint(sat.Byte),
	}, nil
}

func (api *Api) GetTrackerGPSData(trackerId uint16, dateFrom string, dateTo string, all bool) ([]*GpsData, error) {
	tracker := models.Tracker{}

	tx := api.DB.Where("id = ?", trackerId).First(&tracker)

	if tx.Error != nil {

	}

	if tx.RowsAffected == 0 {

	}

	sdrs, err := api.DB.Raw(""+
		"SELECT s.id, ntm, latitude, longitude, mv, bb, spd, alts, dir, dirh, odm, satellites "+
		"FROM service_data_records as sdr "+
		"JOIN sr_pos_data as s ON s.service_data_record_id = sdr.id "+
		"WHERE sdr.tracker_id = ? AND s.ntm BETWEEN ? AND ? AND sdr.deleted_at IS NULL AND s.deleted_at IS NULL AND s.ntm < '2100-01-01' AND s.vld = 1 "+
		"ORDER BY s.id", trackerId, dateFrom, dateTo).Rows()

	defer sdrs.Close()
	if err != nil {
		panic(err)
	}

	result := []*GpsData{}

	var prev *GpsData

	var ntmCounter time.Time

	for sdrs.Next() {

		var id sql.NullInt64
		var ntm sql.NullTime
		var lat sql.NullFloat64
		var lng sql.NullFloat64
		var mv sql.NullBool
		var bb sql.NullBool
		var spd sql.NullInt16
		var alts sql.NullInt32
		var dir sql.NullByte
		var dirh sql.NullByte
		var odm sql.NullInt32
		var sat sql.NullByte

		err := sdrs.Scan(&id, &ntm, &lat, &lng, &mv, &bb, &spd, &alts, &dir, &dirh, &odm, &sat)
		if err != nil {
			return nil, err
		}

		if !all {
			if ntmCounter.IsZero() {
				ntmCounter = ntm.Time
			}

			if prev != nil {
				if prev.Spd == 0 && spd.Int16 == 0 {
					continue
				}

				if prev.Latitude == lat.Float64 && prev.Longitude == lng.Float64 && prev.Alts == alts.Int32 {
					continue
				}
			}

			if ntm.Time.Unix() < ntmCounter.Unix()+10 {
				continue
			}

			ntmCounter = ntm.Time
		}

		srpd := &GpsData{
			Ntm:        ntm.Time,
			Latitude:   lat.Float64,
			Longitude:  lng.Float64,
			Spd:        uint16(spd.Int16),
			Alts:       alts.Int32,
			Dir:        dir.Byte,
			Satellites: uint(sat.Byte),
		}

		result = append(result, srpd)

		prev = srpd
	}

	return result, nil
}
