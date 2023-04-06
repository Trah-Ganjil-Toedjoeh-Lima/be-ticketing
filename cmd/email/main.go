package main

import (
	"context"
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/app"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/frchandra/ticketing-gmcgo/injector"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"strconv"
)

func main() {
	fmt.Println("sending email")
	mailer := injector.InitializeEmail()
	logger := app.NewLogger(config.NewAppConfig())
	logutil := util.NewLogUtil(logger)

	ticketUtil := util.NewETicketUtil(config.NewAppConfig(), newMinio(), logutil)

	reciever := "nismara.chandra@gmail.com"
	data := map[string]any{
		"Name":  "Chandra Herd",
		"Seats": []string{"H31", "H32", "H33"},
	}

	var seatsName = []string{"H31", "H32", "H33"}
	var attachments map[string][]byte = make(map[string][]byte)

	for i := 31; i <= 33; i++ {
		ticket, _ := ticketUtil.GenerateETicket("H"+strconv.Itoa(i), "H"+strconv.Itoa(i))
		attachments["H"+strconv.Itoa(i)+".png"] = ticket
	}

	err := mailer.SendTicketEmail(data, reciever, attachments, seatsName)
	if err != nil {
		panic(err)
	}
}

func newMinio() *minio.Client {
	ctx := context.Background()
	appConfig := config.NewAppConfig()

	minioClient, errInit := minio.New(appConfig.MinioHost+":"+appConfig.MinioPort, &minio.Options{ // Initialize minio client object.
		Creds:  credentials.NewStaticV4(appConfig.MinioRootUser, appConfig.MinioRootPassword, ""),
		Secure: false,
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
		log.Printf("Successfully created bucket: %s\n", appConfig.MinioTicketsBucket)
	}
	return minioClient
}
