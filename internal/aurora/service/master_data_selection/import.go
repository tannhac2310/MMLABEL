package masterdataselection

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/xuri/excelize/v2"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (s *masterDataSelectionService) ImportExcel(ctx context.Context, r io.Reader, createdBy string) error {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	defer f.Close()

	allSelections := make([]*model.MasterDataSelection, 0)

	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		selections, err := importSheet(f, sheet, createdBy)
		if err != nil {
			return err
		}
		allSelections = append(allSelections, selections...)
	}

	return cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		for _, selection := range allSelections {
			if err = s.masterDataSelectionRepo.Insert(ctx, selection); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *masterDataSelectionService) ExampleFile(ctx context.Context, w io.Writer) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Công đoạn (tên ví dụ)"

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	data := [][]string{
		{"Tên dùng để hiển thị trên OA/UI", "Tên dùng để lưu trữ trong phần mềm", "Dữ liệu cho phép chọn nhiều", "Mô tả chi tiết"},
		{"In lụa", "IL", "Y/N", "In lụa"},
	}

	for rowNumber, columns := range data {
		for colNumber, column := range columns {
			cell := string(rune(64+colNumber+1)) + string(rowNumber+1)
			if err := f.SetCellValue(sheetName, cell, column); err != nil {
				return err
			}
		}
	}

	f.SetActiveSheet(index)

	return f.Write(w)
}

func importSheet(f *excelize.File, sheetName string, createdBy string) ([]*model.MasterDataSelection, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	} else if len(rows) == 0 {
		return nil, fmt.Errorf("empty sheet %s", sheetName)
	}

	selections := make([]*model.MasterDataSelection, 0)
	now := time.Now()

	for i, columns := range rows {
		if i == 0 {
			// skip header
			continue
		}

		selection := &model.MasterDataSelection{
			ID:             idutil.ULIDNow(),
			SelectionGroup: sheetName,
			SortOrder:      int16(i),
			Status:         1,
			CreatedAt:      now,
			UpdatedAt:      now,
			CreatedBy:      sql.NullString{String: createdBy, Valid: createdBy != ""},
			UpdatedBy:      sql.NullString{String: createdBy, Valid: createdBy != ""},
		}
		for j, column := range columns {
			value, err := f.GetCellValue(sheetName, column)
			if err != nil {
				return nil, err
			}
			if j == 0 {
				selection.DisplayValue = value
				selection.Value = value
			}
			if j == 1 && value != "" {
				selection.Value = value
			}
			if j == 2 && value == "Y" {
				selection.MultipleChoices = 1
			}
			if j == 3 {
				selection.Description = sql.NullString{String: value, Valid: true}
			}
		}
		selections = append(selections, selection)
	}

	return selections, nil
}
