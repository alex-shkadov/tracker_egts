package server

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"tracker/internal/app/models"
)

type ObjectManager struct {
	db *gorm.DB
}

func NewObjectManager(db *gorm.DB) *ObjectManager {
	return &ObjectManager{db: db}
}

func (om *ObjectManager) GetTracker(imei string) (*models.Tracker, error) {
	var track models.Tracker

	tx := om.db.Where("imei like ?", imei).First(&track)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &track, nil
}

func (om *ObjectManager) SaveTracker(imei string) (*models.Tracker, error) {

	track := &models.Tracker{
		Imei:            imei,
		Title:           imei,
		TransportNumber: imei,
		Description:     imei,
		IsActive:        true,
	}

	result := om.db.Create(track)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected > 0 {
		return track, nil
	}

	return nil, nil
}

func (om *ObjectManager) SavePacket() *models.ServiceDataRecord {
	return nil
}

func (om *ObjectManager) SaveSDR(pId uint16, rNum uint16, oid uint32, tracker *models.Tracker) (*models.ServiceDataRecord, error) {
	sdr := &models.ServiceDataRecord{
		PacketId:         pId,
		RecordNumber:     rNum,
		ObjectIdentifier: oid,
		Tracker:          *tracker,
	}

	result := om.db.Create(sdr)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected > 0 {
		return sdr, nil
	}

	return nil, nil
}

func (om *ObjectManager) SaveSrPosData(sdr *models.ServiceDataRecord, data map[string]interface{}) (*models.SrPosData, error) {

	//speed := data["SPD"].(uint16)
	var altitude int32
	longitude := data["Longitude"].(float64)
	latitude := data["Latitude"].(float64)
	if data["LAHS"] == "1" {
		latitude = -latitude
	}
	if data["LOHS"] == "1" {
		longitude = -longitude
	}

	if data["ALTS"] == "1" {
		altitude = altitude - int32(data["ALT"].(uint32))
	} else {
		altitude = altitude + int32(data["ALT"].(uint32))
	}

	mv := 0
	bb := 0
	vld := 0
	if data["MV"].(string) == "1" {
		mv = 1
	}
	if data["BB"].(string) == "1" {
		bb = 1
	}
	if data["VLD"].(string) == "1" {
		vld = 1
	}

	srd := &models.SrPosData{
		Longitude: longitude,
		Latitude:  latitude,
		Ntm:       data["NTM"].(time.Time),
		Mv:        byte(mv),
		Bb:        byte(bb),
		Spd:       data["SPD"].(uint16),
		Alts:      altitude,
		Dir:       data["DIR"].(byte),
		Vld:       byte(vld),
		Dirh:      data["DIRH"].(byte),
		Odm:       data["ODM"].(uint32),
		SDR:       *sdr,
	}

	sqlText := `
SELECT s.id
FROM service_data_records as sdr
JOIN sr_pos_data as s ON s.service_data_record_id = sdr.id
WHERE sdr.tracker_id = :trackerId
	AND s.ntm = ':ntm'
	AND s.spd = :speed
	AND s.alts = :alts
	AND s.latitude = :latitude
	AND s.longitude = :longitude
	AND s.dir = :dir
  AND sdr.deleted_at IS NULL
LIMIT 1
`
	sqlText = strings.Replace(sqlText, ":trackerId", fmt.Sprint(sdr.TrackerId), 1)
	sqlText = strings.Replace(sqlText, ":ntm", srd.Ntm.Format("2006-01-02 15:04:05"), 1)
	sqlText = strings.Replace(sqlText, ":speed", fmt.Sprint(srd.Spd), 1)
	sqlText = strings.Replace(sqlText, ":alts", fmt.Sprint(srd.Alts), 1)
	sqlText = strings.Replace(sqlText, ":latitude", fmt.Sprint(srd.Latitude), 1)
	sqlText = strings.Replace(sqlText, ":longitude", fmt.Sprint(srd.Longitude), 1)
	sqlText = strings.Replace(sqlText, ":dir", fmt.Sprint(srd.Dir), 1)

	var exists int
	om.db.Raw(sqlText).Scan(&exists)

	if exists > 0 {
		return nil, errors.New("duplicate")
	}

	result := om.db.Create(srd)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected > 0 {
		return srd, nil
	}

	return nil, nil
}
