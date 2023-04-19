package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/skip2/go-qrcode"
	"log"

	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/kumparan/bimg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	seatName := "H44" //creating base qrcode
	url := config.NewAppConfig().AppUrl + config.NewAppConfig().AppPort
	generatedQrImage, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		fmt.Println(err.Error())
	}

	generatedQrWatermark := bimg.WatermarkImage{ //watermark placement
		Left:    1020,
		Top:     100,
		Buf:     generatedQrImage,
		Opacity: 1,
	}

	eticketTemplate, err := bimg.Read("./storage/picture/eticket_template.jpg") //get e-ticket default eticketTemplate
	if err != nil {
		fmt.Println(err.Error())

	}

	eticket, err := bimg.NewImage(eticketTemplate).WatermarkImage(generatedQrWatermark) //pasting qrcode with watermark to the e-ticket
	if err != nil {
		fmt.Println(err.Error())

	}

	categoryText := bimg.Watermark{
		Top:         417,
		Left:        220,
		Text:        "DIAMOND",
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
		Text:        "H44",
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

	minioClient, err := newMinio()
	if err != nil {
		fmt.Println(err.Error())
	}

	bucketName := config.NewAppConfig().MinioTicketsBucket
	objectName := seatName + ".png"
	fileBuffer := bytes.NewReader(eticket)
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
