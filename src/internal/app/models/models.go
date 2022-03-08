package models

import (
	"gorm.io/gorm"
	"time"
)

type Tracker struct {
	ID              uint16         `json:"id" gorm:"primarykey"`
	Title           string         `json:"title"`
	Imei            string         `json:"imei"`
	TransportNumber string         `json:"transport_number"`
	Description     string         `json:"description"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ServiceDataRecord struct {
	ID               uint64         `json:"id" gorm:"primarykey"`
	PacketId         uint16         `json:"packet_id"`
	TrackerId        uint16         `json:"tracker_id"`
	Tracker          Tracker        `gorm:"foreignKey:TrackerId"`
	RecordNumber     uint16         `json:"record_number"`
	ObjectIdentifier uint32         `json:"object_identifier"`
	UpdatedAt        time.Time      `json:"updated_at"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type SrPosData struct {
	ID                  uint64            `json:"id" gorm:"primarykey"`
	ServiceDataRecordId uint64            `json:"-"`
	SDR                 ServiceDataRecord `gorm:"foreignKey:ServiceDataRecordId" json:"-"`
	Ntm                 time.Time         `json:"ntm"`
	Latitude            float64           `json:"latitude"`
	Longitude           float64           `json:"longitude"`
	Mv                  byte              `json:"mv"`
	Bb                  byte              `json:"bb"`
	Vld                 byte              `json:"vld"`
	Spd                 uint16            `json:"spd"`
	Alts                int32             `json:"alts"`
	Dir                 byte              `json:"dir"`
	Dirh                byte              `json:"dirh"`
	Odm                 uint32            `json:"odm"`
	Satellites          uint              `json:"satellites"`
	DisplayName         string            `json:"display_name"`
	RecordNumber        int16             `json:"record_number"`
	UpdatedAt           time.Time         `json:"-"`
	CreatedAt           time.Time         `json:"-"`
	DeletedAt           gorm.DeletedAt    `json:"-" gorm:"index"`
}
