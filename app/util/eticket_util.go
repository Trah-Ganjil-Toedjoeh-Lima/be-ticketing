package util

import (
	"bytes"
	"context"
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

func (e *ETicketUtil) GenerateETicket(seatName, seatLink string) ([]byte, error) {

	url := e.config.AppUrl + ":" + e.config.AppPort + "/api/v1/seat/" + seatLink //creating basic qr code
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		e.log.BasicLog(err, "when generating e-ticket qr code")
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

	newQrCode, err := bimg.NewImage(qr).Watermark(title) //pasting additional information to the qrcode as a watermark
	newQrCode, err = bimg.NewImage(newQrCode).Watermark(seat)

	content := bimg.WatermarkImage{
		Left:    0,
		Top:     0,
		Buf:     newQrCode,
		Opacity: 1,
	}

	frame, _ := bimg.Read("./storage/picture/polite_cat.png")

	ticket, err := bimg.NewImage(frame).WatermarkImage(content)
	if err != nil {
		e.log.BasicLog(err, "when generating e-ticket")
		return nil, err
	}

	bucketName := e.config.MinioTicketsBucket // Upload the file with to minio
	objectName := seatName + ".png"
	fileBuffer := bytes.NewReader(ticket)
	fileSize := fileBuffer.Size()
	contentType := "png"

	_, err = e.minio.PutObject(context.Background(), bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		e.log.BasicLog(err, "when storing e-ticket to minio")
		return nil, err
	}

	return ticket, nil
}
