package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/kumparan/bimg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/skip2/go-qrcode"
)

func main() {
	seatName := "H44" //creating base qrcode
	url := config.NewAppConfig().AppUrl + config.NewAppConfig().AppPort
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err.Error())
	}

	title := bimg.Watermark{
		Top:         3,
		Left:        38,
		Text:        "DIAMOND",
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 9",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	seat := bimg.Watermark{
		Top:         15,
		Left:        90,
		Text:        "SEAT - " + seatName,
		Opacity:     1,
		Width:       200,
		DPI:         100,
		Margin:      0,
		Font:        "sans bold 12",
		Background:  bimg.Color{0, 0, 0},
		NoReplicate: true,
	}

	newQrCode, err := bimg.NewImage(qr).Watermark(title) //pasting watermark to base qrcode
	newQrCode, err = bimg.NewImage(newQrCode).Watermark(seat)

	content := bimg.WatermarkImage{ //watermark placement
		Left:    0,
		Top:     0,
		Buf:     newQrCode,
		Opacity: 1,
	}

	frame, err := bimg.Read("./storage/picture/eticket.png") //get e-ticket default frame
	if err != nil {
		fmt.Println(err.Error())

	}

	ticket, err := bimg.NewImage(frame).WatermarkImage(content) //pasting qrcode with watermark to the e-ticket
	if err != nil {
		fmt.Println(err.Error())

	}

	minioClient, err := newMinio()
	if err != nil {
		fmt.Println(err.Error())
	}

	bucketName := config.NewAppConfig().MinioTicketsBucket
	objectName := seatName + ".png"
	fileBuffer := bytes.NewReader(ticket)
	fileSize := fileBuffer.Size()
	contentType := "picture"

	// Upload the file with to minio
	info, err := minioClient.PutObject(context.Background(), bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	// Download the file from minio
	err = minioClient.FGetObject(context.Background(), bucketName, objectName, "./storage/ticket/"+objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func newMinio() (*minio.Client, error) {
	ctx := context.Background()
	appConfig := config.NewAppConfig()

	minioClient, errInit := minio.New(appConfig.MinioHost+":"+appConfig.MinioPort, &minio.Options{ // Initialize minio client object.
		Creds:  credentials.NewStaticV4(appConfig.MinioRootUser, appConfig.MinioRootPassword, ""),
		Secure: appConfig.MinioSecure,
	})
	if errInit != nil {
		log.Fatalln(errInit)
	}

	err := minioClient.MakeBucket(ctx, appConfig.MinioTicketsBucket, minio.MakeBucketOptions{Region: appConfig.MinioLocation})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, appConfig.MinioTicketsBucket) // Check to see if we already own this bucket (which happens if you run this twice)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", appConfig.MinioTicketsBucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", appConfig.MinioTicketsBucket)
	}
	return minioClient, errInit
}
