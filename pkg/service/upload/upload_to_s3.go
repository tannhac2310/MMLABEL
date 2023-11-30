package upload

import (
	"bytes"
	"context"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *uploadService) Upload(ctx context.Context, uploadType, originName string, file *bytes.Buffer) (string, error) {
	dateStr := time.Now().Format("2006-01-02")

	fileName := fmt.Sprintf("%s/%s/origin/%sss1%sss1%s", uploadType, dateStr, idutil.ULIDNow(), time.Now().Nanosecond(), originName)
	//fmt.Println(fileName)
	url, err := u.uploadToS3(ctx, u.cfg.S3Storage.Bucket, file, fileName)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return url, nil
}

func (u *uploadService) uploadToS3(ctx context.Context, bucket string, file *bytes.Buffer, path string) (string, error) {
	contentType := http.DetectContentType(file.Bytes())

	// Upload the file to S3.
	r, err := u.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(path),
		Body:        file,
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return r.Location, nil
}
