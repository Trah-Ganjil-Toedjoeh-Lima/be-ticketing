package util

import (
	"bytes"
	"context"
	"fmt"

	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/kumparan/bimg"
	"github.com/minio/minio-go/v7"
	"github.com/skip2/go-qrcode"
)

type ETicketUtil struct {
	config *config.AppConfig
	minio  *minio.Client
	log    *LogUtil
}

func NewETicketUtil(config *config.AppConfig, minio *minio.Client, log *LogUtil) *ETicketUtil {
	return &ETicketUtil{config: config, minio: minio, log: log}
}

func (e *ETicketUtil) GenerateETicket(seatName, seatLink, seatCategory string) ([]byte, error) {

	url := "https://gmco-event.com/ticket/" + seatLink //creating basic qr code
	generatedQrImage, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err.Error())
	}

	eticketTemplate, err := bimg.Read("./storage/picture/eticket_template.jpg") //get e-ticket default eticketTemplate
	if err != nil {
		fmt.Println(err.Error())
	}

	generatedQrWatermark := bimg.WatermarkImage{ //watermark placement
		Left:    1020,
		Top:     100,
		Buf:     generatedQrImage,
		Opacity: 1,
	}
	eticket, err := bimg.NewImage(eticketTemplate).WatermarkImage(generatedQrWatermark) //pasting qrcode with watermark to the e-ticket
	if err != nil {
		fmt.Println(err.Error())
	}

	categoryText := bimg.Watermark{
		Top:         417,
		Left:        220,
		Text:        seatCategory,
		Opacity:     1,
		Width:       100,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 20",
		Background:  bimg.Color{246, 247, 241},
		NoReplicate: true,
	}
	eticket, err = bimg.NewImage(eticket).Watermark(categoryText)
	if err != nil {
		fmt.Println(err.Error())
	}

	seatNameText := bimg.Watermark{
		Top:         417,
		Left:        123,
		Text:        seatName,
		Opacity:     1,
		Width:       100,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 24",
		Background:  bimg.Color{246, 247, 241},
		NoReplicate: true,
	}
	eticket, err = bimg.NewImage(eticket).Watermark(seatNameText)
	if err != nil {
		fmt.Println(err.Error())
	}

	bucketName := e.config.MinioTicketsBucket // Upload the file with to minio
	objectName := seatName + "_" + seatLink + ".png"
	fileBuffer := bytes.NewReader(eticket)
	fileSize := fileBuffer.Size()
	contentType := "image/png"

	_, err = e.minio.PutObject(context.Background(), bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		e.log.BasicLog(err, "when storing e-ticket to minio")
		return nil, err
	}

	return eticket, nil
}
