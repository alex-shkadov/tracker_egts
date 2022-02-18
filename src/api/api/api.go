package api

import (
	"database/sql"
	"gorm.io/gorm"
	"tracker/internal/app/models"
	"tracker/internal/app/server"
)

type Api struct {
	OM *server.ObjectManager
	DB *gorm.DB
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

func (api *Api) GetTrackerGPSData(trackerId uint16, dateFrom string, dateTo string) ([]*models.SrPosData, error) {
	tracker := models.Tracker{}

	tx := api.DB.Where("id = ?", trackerId).First(&tracker)

	if tx.Error != nil {

	}

	if tx.RowsAffected == 0 {

	}

	sdrs, err := api.DB.Raw(""+
		"SELECT s.id, ntm, latitude, longitude, mv, bb, spd, alts, dir, odm, satellites "+
		"FROM service_data_records as sdr "+
		"JOIN sr_pos_data as s ON s.service_data_record_id = sdr.id "+
		"WHERE sdr.tracker_id = ? AND s.ntm BETWEEN ? AND ? AND sdr.deleted_at IS NULL AND s.deleted_at IS NULL", trackerId, dateFrom, dateTo).Rows()

	defer sdrs.Close()
	if err != nil {
		panic(err)
	}

	result := []*models.SrPosData{}

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
		var odm sql.NullInt32
		var sat sql.NullByte

		err := sdrs.Scan(&id, &ntm, &lat, &lng, &mv, &bb, &spd, &alts, &dir, &odm, &sat)
		if err != nil {
			return nil, err
		}

		srpd := &models.SrPosData{
			ID:         uint64(id.Int64),
			Ntm:        ntm.Time,
			Latitude:   lat.Float64,
			Longitude:  lng.Float64,
			Mv:         mv.Bool,
			Bb:         bb.Bool,
			Spd:        uint16(spd.Int16),
			Alts:       alts.Int32,
			Dir:        dir.Byte,
			Odm:        uint32(odm.Int32),
			Satellites: uint(sat.Byte),
		}

		result = append(result, srpd)
	}

	return result, nil
}
