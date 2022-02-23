package parser

import (
	"github.com/kuznetsovin/egts-protocol/libs/egts"
	"log"
)

type ServiceType int

const (
	Auth ServiceType = 1
	Tele ServiceType = 2
)

func CreateSrResultCodeResponse(pack *egts.Package) *egts.Package {

	return &egts.Package{
		ProtocolVersion:  1,
		SecurityKeyID:    0,
		Prefix:           "00",
		EncryptionAlg:    "00",
		Compression:      "00",
		Priority:         "00",
		HeaderEncoding:   0,
		PacketIdentifier: 2,
		TimeToLive:       10,
		PacketType:       egts.PtAppdataPacket,
		ServicesFrameData: &egts.ServiceDataSet{
			egts.ServiceDataRecord{
				RecordNumber:             0,
				SourceServiceOnDevice:    "0",
				RecipientServiceOnDevice: "1",
				Group:                    "0",
				RecordProcessingPriority: "00",
				SourceServiceType:        egts.AuthService,
				RecipientServiceType:     egts.AuthService,
				RecordDataSet: egts.RecordDataSet{
					egts.RecordData{
						SubrecordType: egts.SrResultCodeType,
						SubrecordData: &egts.SrResultCode{
							ResultCode: 0,
						},
					},
				},
			},
		},
	}

}

func CreateAuthResponse(pack *egts.Package) *egts.Package {

	return &egts.Package{
		ProtocolVersion: 1,
		SecurityKeyID:   0,
		Prefix:          "00",
		Route:           "0",
		EncryptionAlg:   "00",
		Compression:     "0",
		Priority:        "00",
		//HeaderLength:    11,
		HeaderEncoding: 0,
		//FrameDataLength:  3,
		PacketIdentifier: 1,
		TimeToLive:       10,
		PacketType:       egts.PtResponsePacket,
		//HeaderCheckSum:   74,
		ServicesFrameData: &egts.PtResponse{
			ResponsePacketID: pack.PacketIdentifier,
			ProcessingResult: 0,
		},
		//ServicesFrameDataCheckSum: 59443,
	}
}

func ParseMessage(bytes []byte, addr string) (*egts.Package, ServiceType, error) {

	//fmt.Println("Str: ", bytes)

	result := egts.Package{}
	var servType ServiceType

	_, err := result.Decode(bytes)

	if err != nil {
		return nil, 0, err
	}

	if result.PacketType == 1 {
		dataSet := result.ServicesFrameData.(*egts.ServiceDataSet)
		for _, ds := range *dataSet {
			if ds.SourceServiceType == 1 {
				servType = Auth
			}
			if ds.SourceServiceType == 2 {
				servType = Tele
			}
		}
	} else {
		//dataSet := result.ServicesFrameData.(*egts.PtResponse)
	}

	if result.PacketType == 1 {
		//fmt.Println(addr, "DDD")
	}

	return &result, servType, nil
}

func GetIMEI(p *egts.Package) string {

	if p.PacketType != 1 {
		return ""
	}

	dataSet := p.ServicesFrameData.(*egts.ServiceDataSet)
	for _, ds := range *dataSet {
		if ds.SourceServiceType == 1 {
			for _, rds := range ds.RecordDataSet {
				termId := rds.SubrecordData.(*egts.SrTermIdentity)
				return termId.IMEI
			}
		}
	}

	return ""
}

func DebugPack(p *egts.Package) {

	var dataSet *egts.ServiceDataSet

	log.Println("Package: ")
	log.Println("\tProtocolVersion: ", p.ProtocolVersion)
	log.Println("\tSecurityKeyID: ", p.SecurityKeyID)

	log.Println("\tPrefix: ", p.Prefix)
	log.Println("\tRoute: ", p.Route)
	log.Println("\tEncryptionAlg: ", p.EncryptionAlg)
	log.Println("\tCompression: ", p.Compression)
	log.Println("\tPriority: ", p.Priority)

	log.Println("\tHeaderLength: ", p.HeaderLength)
	log.Println("\tHeaderEncoding: ", p.HeaderEncoding)
	log.Println("\tFrameDataLength: ", p.FrameDataLength)
	log.Println("\tPacketIdentifier: ", p.PacketIdentifier)
	log.Println("\tPacketType: ", p.PacketType)
	log.Println("\tPeerAddress: ", p.PeerAddress)
	log.Println("\tRecipientAddress: ", p.RecipientAddress)
	log.Println("\tTimeToLive: ", p.TimeToLive)
	log.Println("\tHeaderCheckSum: ", p.HeaderCheckSum)
	log.Println("\tServicesFrameDataCheckSum: ", p.ServicesFrameDataCheckSum)

	log.Println("\tFrameDataLength: ", p.FrameDataLength)
	log.Println("\tServicesFrameData: ", p.ServicesFrameData.Length())

	if p.PacketType == 1 {
		dataSet := p.ServicesFrameData.(*egts.ServiceDataSet)
		for _, ds := range *dataSet {
			if ds.SourceServiceType == 1 {
				//servType = Auth
			}
		}
	}

	log.Println("==============ServicesFrameData==============")
	log.Println("DataSetLength", dataSet.Length())
	for _, ds := range *dataSet {
		log.Println("ds:EventIdentifier", ds.EventIdentifier)
		log.Println("ds:RecordDataSet", ds.RecordDataSet)
		log.Println("ds:RecordLength", ds.RecordLength)
		log.Println("ds:RecordNumber", ds.RecordNumber)
		log.Println("ds:SourceServiceOnDevice", ds.SourceServiceOnDevice)
		log.Println("ds:RecipientServiceOnDevice", ds.RecipientServiceOnDevice)
		log.Println("ds:Group", ds.Group)
		log.Println("ds:RecordProcessingPriority", ds.RecordProcessingPriority)
		log.Println("ds:TimeFieldExists", ds.TimeFieldExists)
		log.Println("ds:EventIDFieldExists", ds.EventIDFieldExists)
		log.Println("ds:ObjectIDFieldExists", ds.ObjectIDFieldExists)
		log.Println("ds:ObjectIdentifier", ds.ObjectIdentifier)
		log.Println("ds:EventIdentifier", ds.EventIdentifier)
		log.Println("ds:Time", ds.Time)
		log.Println("ds:SourceServiceType", ds.SourceServiceType)
		if ds.SourceServiceType == 1 {
			//servType = Auth
		}
		log.Println("ds:RecipientServiceType", ds.RecipientServiceType)
		//483
		log.Println("========ds:RecordDataSet========")
		for _, rds := range ds.RecordDataSet {
			log.Println("\tds:RecordDataSet:SubrecordType", rds.SubrecordType)
			log.Println("\tds:RecordDataSet:SubrecordLength", rds.SubrecordLength)
			termId := rds.SubrecordData.(*egts.SrTermIdentity)
			log.Println("\tds:RecordDataSet:SubrecordData:IMEI", termId.IMEI)
			log.Println("\tds:RecordDataSet:SubrecordData:BSE", termId.BSE)
			log.Println("\tds:RecordDataSet:SubrecordData:HDIDE", termId.HDIDE)
			log.Println("\tds:RecordDataSet:SubrecordData:BufferSize", termId.BufferSize)
			log.Println("\tds:RecordDataSet:SubrecordData:HomeDispatcherIdentifier", termId.HomeDispatcherIdentifier)
			log.Println("\tds:RecordDataSet:SubrecordData:IMEIE", termId.IMEIE)
			log.Println("\tds:RecordDataSet:SubrecordData:IMSI", termId.IMSI)
			log.Println("\tds:RecordDataSet:SubrecordData:IMSIE", termId.IMSIE)
			log.Println("\tds:RecordDataSet:SubrecordData:LanguageCode", termId.LanguageCode)
			log.Println("\tds:RecordDataSet:SubrecordData:LNGCE", termId.LNGCE)
			log.Println("\tds:RecordDataSet:SubrecordData:MNE", termId.MNE)
			log.Println("\tds:RecordDataSet:SubrecordData:NetworkIdentifier", termId.NetworkIdentifier)
			log.Println("\tds:RecordDataSet:SubrecordData:NIDE", termId.NIDE)
			log.Println("\tds:RecordDataSet:SubrecordData:SSRA", termId.SSRA)
			log.Println("\tds:RecordDataSet:SubrecordData:TerminalIdentifier", termId.TerminalIdentifier)
		}

	}
}
