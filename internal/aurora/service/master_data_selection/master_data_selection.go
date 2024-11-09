package masterdataselection

import (
	"context"
	"io"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type MasterDataSelectionService interface {
	ImportExcel(ctx context.Context, r io.Reader, createdBy string) error
	ExampleFile(ctx context.Context, w io.Writer) error
}

type masterDataSelectionService struct {
	masterDataSelectionRepo repository.MasterDataSelectionRepo
}
