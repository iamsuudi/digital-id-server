package storage

import (
	"context"
	"digital-id-server/shared/utils"
	"log"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinIO() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	// Initialize MinIO client
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Use true for HTTPS in production
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to check bucket existence: %v", err)
	}
	if !exists {
		err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}
}

/*
 * UploadToMinio uploads a file to MinIO.
 *
 */
func UploadToMinio(file *multipart.FileHeader) (string, error) {
	filename := utils.MakeFileName(file.Filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload to MinIO
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	_, err = MinioClient.PutObject(context.Background(), bucketName, filename, src, file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})

	return filename, err
}
