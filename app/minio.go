package app

import (
	"context"

	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func NewMinio(appConfig *config.AppConfig, log *logrus.Logger) *minio.Client {
	ctx := context.Background()

	minioClient, errInit := minio.New(appConfig.MinioHost+":"+appConfig.MinioPort, &minio.Options{ // Initialize minio client object.
		Creds:  credentials.NewStaticV4(appConfig.MinioAccessKey, appConfig.MinioSecretKey, ""),
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
	return minioClient
}
