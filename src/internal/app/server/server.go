package server

import (
	"fmt"
	"github.com/kuznetsovin/egts-protocol/libs/egts"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	db2 "tracker/internal/app/db"
	"tracker/internal/app/debug"
	"tracker/internal/app/logger"
	"tracker/internal/app/models"
	"tracker/internal/app/parser"
)

var db *gorm.DB
var om *ObjectManager

func HandleConnection(c net.Conn, timeout int) {

	var track *models.Tracker

	defer c.Close()
	//i := 19

	prefix := c.RemoteAddr().String()

	log.Println(prefix, "New Connection")
	for {

		bytes := make([]byte, 65535)
		log.Println(prefix, "Try to read...")
		//c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		n, err := c.Read(bytes)

		log.Println(prefix, "Read...")

		if err != nil {

			imei := "-"
			if track != nil {
				imei = track.Imei
			}

			log.Println(prefix, "Ошибка чтения данных TCP-соединения. Завершение работы  соединения от устройства ", imei, err)
			c.Close()
			return
		}

		logger.LogETGSConnectionData(bytes[:n], true, "")

		pack, servType, err := parser.ParseMessage(bytes[:n], prefix)

		if err != nil {
			log.Println(err)
			continue
		}

		if servType == parser.Auth {

			imei := parser.GetIMEI(pack)

			if imei == "" {
				log.Fatalln("Не опознан IMEI")
			}

			//logger.LogETGSConnectionData(bytes[:n], true, imei)

			imei = strings.Trim(imei, "\x00")
			if imei == "" {
				log.Fatalln("Не опознан IMEI")
			}

			log.Println(prefix, "IMEI:", imei)

			track, _ = om.GetTracker(imei)
			if track == nil {
				track, err = om.SaveTracker(imei)
				if err != nil {

					log.Println(prefix, "Ошибка создания трекера. Завершение работы  соединения от устройства", track.Imei, err)
					c.Close()
					return
				}
			}

			authResponse := parser.CreateAuthResponse(pack)

			authResponseBytes, err := authResponse.Encode()
			if err != nil {
				log.Println(prefix, "Ошибка формирования ответа:", err)
				c.Close()
				return
			}

			n, err = c.Write(authResponseBytes)
			if err != nil {
				log.Println(prefix, "Ошибка записи пакета ответа:", err)
				c.Close()
				return
			}

			logger.LogETGSConnectionData(authResponseBytes[:n], false, "")
			//logger.LogETGSConnectionData(authResponseBytes[:n], false, imei)

			//log.Println(prefix, "Send Sr Result Code")

			srResultCodeResp := parser.CreateSrResultCodeResponse(pack)

			srResultCodeRespBytes, err := srResultCodeResp.Encode()
			if err != nil {
				log.Println(prefix, "Ошибка формирования ответа:", err)
				c.Close()
				return
			}

			//log.Println(srResultCodeRespBytes)
			n2, err := c.Write(srResultCodeRespBytes)
			if err != nil {
				log.Println(prefix, "Ошибка записи пакета SrResCode:", err)
				c.Close()
				return
			}

			logger.LogETGSConnectionData(srResultCodeRespBytes[:n2], false, "")
			//logger.LogETGSConnectionData(srResultCodeRespBytes[:n2], false, imei)

		} else {
			if track == nil {
				log.Fatalln("Не определен трекер")
			}

			//logger.LogETGSConnectionData(bytes[:n], true, track.Imei)

			if servType == parser.Tele {

				dataSet := pack.ServicesFrameData.(*egts.ServiceDataSet)
				for _, sfd := range *dataSet {
					sdr, err := om.SaveSDR(pack.PacketIdentifier, sfd.RecordNumber, sfd.ObjectIdentifier, track)
					if err != nil {
						log.Fatalln(err)
					}

					if sdr != nil {
						locDef := false
						satDef := false

						var data map[string]interface{}

						for _, rd := range sfd.RecordDataSet {
							if rd.SubrecordType == 16 {
								rdd := rd.SubrecordData.(*egts.SrPosData)

								data = map[string]interface{}{
									"Latitude":  rdd.Latitude,
									"Longitude": rdd.Longitude,
									"NTM":       rdd.NavigationTime,
									"MV":        rdd.MV,
									"BB":        rdd.BB,
									"SPD":       rdd.Speed,
									"VLD":       rdd.VLD,
									"ALT":       rdd.Altitude,
									"DIR":       rdd.Direction,
									"DIRH":      rdd.DirectionHighestBit,
									"ODM":       rdd.Odometer,
									"LOHS":      rdd.LOHS,
									"LAHS":      rdd.LAHS,
									"ALTS":      rdd.AltitudeSign,
								}

								locDef = true
							}

							if rd.SubrecordType == 17 {
								rdsd := rd.SubrecordData.(*egts.SrExtPosData)
								if locDef {
									data["Satellites"] = rdsd.Satellites
									satDef = true
								}
							}
						}

						if locDef && satDef {
							srd, err := om.SaveSrPosData(sdr, data)
							if err != nil {
								log.Fatalln(err)
							}

							fmt.Sprint(srd)
						}
					}
				}

				resp := parser.CreatePtResponse(pack, 0)

				respBytes, err := resp.Encode()
				if err != nil {
					log.Println(prefix, "Ошибка формирования ответа:", err)
					c.Close()
					return
				}

				//log.Println(srResultCodeRespBytes)
				n2, err := c.Write(respBytes)
				if err != nil {
					log.Println(prefix, "Ошибка записи пакета PtResponse:", err)
					c.Close()
					return
				}

				logger.LogETGSConnectionData(respBytes[:n2], false, "")
				//logger.LogETGSConnectionData(respBytes[:n2], false, track.Imei)
			}
		}
	}
}

func MemoryLogger() {
	debug.PrintMemUsage()

	for range time.Tick(time.Minute * 20) {
		debug.PrintMemUsage()
	}
}

func Start(port int) {

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		return
	}

	defer l.Close()

	timeout, ok := os.LookupEnv("read_timeout")
	if !ok {
		log.Fatalln("Не задан таймаут для соединения")
	}

	timeoutInt, _ := strconv.Atoi(timeout)
	db = db2.DBConnect()
	om = &ObjectManager{
		db: db,
	}

	go MemoryLogger()

	log.Println("Start TCP server on port", port)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go HandleConnection(c, timeoutInt)
	}
}
