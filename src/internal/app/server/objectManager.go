package server

import (
	"gorm.io/gorm"
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

func (om *ObjectManager) SaveSDR(rNum uint16, oid uint32, tracker *models.Tracker) (*models.ServiceDataRecord, error) {
	sdr := &models.ServiceDataRecord{
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

	srd := &models.SrPosData{
		Longitude: longitude,
		Latitude:  latitude,
		Ntm:       data["NTM"].(time.Time),
		Mv:        data["MV"].(bool),
		Bb:        data["BB"].(bool),
		Spd:       data["SPD"].(uint16),
		Alts:      altitude,
		Dir:       data["DIR"].(byte),
		Dirh:      data["DIRH"].(byte),
		Odm:       data["ODM"].(uint32),
		SDR:       *sdr,
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
