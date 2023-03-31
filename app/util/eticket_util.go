package util

import (
	"bytes"
	"context"
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/kumparan/bimg"
	"github.com/minio/minio-go/v7"
	"github.com/skip2/go-qrcode"
	"log"
)

type ETicketUtil struct {
	config *config.AppConfig
	minio  *minio.Client
}

func NewETicketUtil(config *config.AppConfig, minio *minio.Client) *ETicketUtil {
	return &ETicketUtil{config: config, minio: minio}
}

func (e *ETicketUtil) GenerateETicket(seatName, seatLink string) ([]byte, error) {

	url := e.config.AppUrl + ":" + e.config.AppPort + "/api/v1/seat/" + seatLink
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	title := bimg.Watermark{
		Top:         2,
		Left:        38,
		Text:        "GRAND CONCERT GMCO 2023",
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 9",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	seat := bimg.Watermark{
		Top:         12,
		Left:        90,
		Text:        "SEAT - " + seatName,
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 11",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	newQrCode, err := bimg.NewImage(qr).Watermark(title)
	newQrCode, err = bimg.NewImage(newQrCode).Watermark(seat)

	content := bimg.WatermarkImage{
		Left:    0,
		Top:     0,
		Buf:     newQrCode,
		Opacity: 1,
	}

	frame, err := bimg.Read("./storage/picture/polite_cat.png")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	ticket, err := bimg.NewImage(frame).WatermarkImage(content)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Upload the file with to minio
	bucketName := e.config.MinioTicketsBucket
	objectName := seatName + ".png"
	fileBuffer := bytes.NewReader(ticket)
	fileSize := fileBuffer.Size()
	contentType := "picture"

	_, err = e.minio.PutObject(context.Background(), bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("Successfully uploaded %s\n", objectName)

	return ticket, nil
}
