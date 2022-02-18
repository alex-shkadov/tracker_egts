package app

import (
	"bufio"
	"fmt"
	"github.com/kuznetsovin/egts-protocol/libs/egts"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	db2 "tracker/internal/app/db"
	"tracker/internal/app/parser"
	"tracker/internal/app/server"
)

func TestParser(files []string) {
	db := db2.DBConnect()
	om := server.NewObjectManager(db)

	for _, filePath := range files {
		f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			row := scanner.Text()
			if row == "" {
				continue
			}

			if row[0:2] == "<<" {
				continue
			}

			re, _ := regexp.Compile(`\[.*\]`)
			res := re.FindString(row)
			if res != "" {
				bytes := []byte{}
				res = res[1 : len(res)-1]
				for _, integerStr := range strings.Split(res, " ") {
					integer, _ := strconv.Atoi(integerStr)

					bytes = append(bytes, byte(integer))
				}

				pkt, st := parser.ParseMessage(bytes, "")

				if len(bytes) == 45 {
					imei := parser.GetIMEI(pkt)

					imei = strings.Trim(imei, "\x00")
					if imei == "" {
						log.Fatalln("Не опознан IMEI")
					}

					fmt.Println("IMEI:", imei)

					track, _ := om.GetTracker(imei)
					if track == nil {
						track, err = om.SaveTracker(imei)
						if err != nil {

							fmt.Println("Ошибка создания трекера. Завершение работы  соединения.", err)
							return
						}
					}
				}
				if len(bytes) == 73 {

					txt, _ := pkt.Encode()
					fmt.Println(len(bytes), st, txt)

					dataSet := pkt.ServicesFrameData.(*egts.ServiceDataSet)

					for _, sfd := range *dataSet {
						track, _ := om.GetTracker("111")
						sdr, err := om.SaveSDR(sfd.RecordNumber, sfd.ObjectIdentifier, track)
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

									bb := false
									mv := false

									if rdd.MV == "1" {
										mv = true
									}
									if rdd.BB == "1" {
										bb = true
									}

									data = map[string]interface{}{
										"Latitude":  rdd.Latitude,
										"Longitude": rdd.Longitude,
										"NTM":       rdd.NavigationTime,
										"MV":        mv,
										"BB":        bb,
										"SPD":       rdd.Speed,
										"ALT":       rdd.Altitude,
										"DIR":       rdd.Direction,
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

								fmt.Println(srd)
							}
						}
					}

					//fmt.Println(sdr)
				}

				//if tp == Auth {
				//	pkt1 := CreateAuthResponse(pkt)
				//
				//	txt, _ := pkt1.Encode()
				//	fmt.Println(txt)
				//	dec, _ := ParseMessage(txt, "")
				//	_txt, _ := dec.Encode()
				//	fmt.Println(_txt)
				//
				//	pkt2 := CreateSrResultCodeResponse(pkt)
				//	txt, _ = pkt2.Encode()
				//	fmt.Println(txt)
				//	dec, _ = ParseMessage(txt, "")
				//
				//	_txt, _ = dec.Encode()
				//	fmt.Println(_txt)
				//
				//}

			}
		}
	}
}
