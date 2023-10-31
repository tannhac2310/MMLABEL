package upload

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mmlabel.gitlab.com/mm-printing-backend/internal/hydra/configs"
)

type Service interface {
	Upload(ctx context.Context, uploadType, originName string, file *bytes.Buffer) (string, error)
}

type uploadService struct {
	cfg *configs.Config

	uploader *s3manager.Uploader
}

func NewService(
	cfg *configs.Config,
	uploader *s3manager.Uploader,
) Service {
	return &uploadService{
		cfg:      cfg,
		uploader: uploader,
	}
}
