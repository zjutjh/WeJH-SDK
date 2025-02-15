package excel

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// File excel各工作表数据结构
type File struct {
	Sheets []Sheet `json:"sheets"`
}

// Sheet excel工作表数据结构
type Sheet struct {
	Name    string   `json:"name"`    // 工作表名称
	Headers []string `json:"headers"` // 列名
	Rows    [][]any  `json:"rows"`    // 行数据（支持不同类型）
}

// CreateExcelFile 生成 Excel 文件
func CreateExcelFile(data File, fileName, filePath, host string) (string, error) {
	if err := validateFileData(data); err != nil {
		return "", err
	}
	f := excelize.NewFile()

	// 遍历 sheets
	for i, sheet := range data.Sheets {
		sheetName := sheet.Name
		if err := createSheet(f, sheet, sheetName, i); err != nil {
			return "", err
		}
	}

	// 设置默认显示的 Sheet
	f.SetActiveSheet(0)

	// 保存文件
	fullPath := filepath.Join(filePath, fileName)
	if err := ensureDirExists(filePath); err != nil {
		return "", err
	}

	if err := removeOldFile(fullPath); err != nil {
		return "", err
	}

	if err := f.SaveAs(fullPath); err != nil {
		return "", err
	}

	return filepath.Join(host, fullPath), nil
}

// createSheet 创建一个工作表，并填充数据
func createSheet(f *excelize.File, sheet Sheet, sheetName string, index int) error {
	var err error
	if index == 0 {
		// 默认创建的 Sheet1 需要重命名
		err = f.SetSheetName("Sheet1", sheetName)
	} else {
		_, err = f.NewSheet(sheetName)
	}
	if err != nil {
		return err
	}

	// 写入列名
	if err = writeHeaders(f, sheetName, sheet.Headers); err != nil {
		return err
	}

	// 设置列宽
	autoAdjustColumnWidth(f, sheetName, sheet.Headers, sheet.Rows)

	// 写入数据
	return writeRows(f, sheetName, sheet.Rows)
}

// writeHeaders 写入列名
func writeHeaders(f *excelize.File, sheetName string, headers []string) error {
	cellRef, _ := excelize.CoordinatesToCellName(1, 1)
	// 批量写入列名
	return f.SetSheetRow(sheetName, cellRef, &headers)
}

// writeRows 批量写入数据行
func writeRows(f *excelize.File, sheetName string, rows [][]any) error {
	for rowIndex, row := range rows {
		cellRef, _ := excelize.CoordinatesToCellName(1, rowIndex+2) // +2 是因为第一行是列名
		// 批量写入每一行数据
		if err := f.SetSheetRow(sheetName, cellRef, &row); err != nil {
			return err
		}
	}
	return nil
}

// ensureDirExists 确保目录存在
func ensureDirExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.MkdirAll(filePath, 0750)
	}
	return nil
}

// removeOldFile 删除旧文件
func removeOldFile(fullPath string) error {
	if _, err := os.Stat(fullPath); err == nil {
		return os.Remove(fullPath)
	}
	return nil
}

// autoAdjustColumnWidth 自动调整列宽
func autoAdjustColumnWidth(f *excelize.File, sheetName string, headers []string, rows [][]any) {
	for colIndex, header := range headers {
		maxWidth := len(header)
		for _, row := range rows {
			cellValue := fmt.Sprintf("%v", row[colIndex])
			if len(cellValue) > maxWidth {
				maxWidth = len(cellValue)
			}
		}
		colName, _ := excelize.ColumnNumberToName(colIndex + 1)
		if err := f.SetColWidth(sheetName, colName, colName, float64(maxWidth+1)); err != nil {
			log.Printf("Failed to set column width for %s: %v\n", colName, err)
		}
	}
}

// validateFileData 验证文件数据
func validateFileData(data File) error {
	if len(data.Sheets) == 0 {
		return errors.New("no sheets provided")
	}
	for _, sheet := range data.Sheets {
		if len(sheet.Headers) == 0 {
			return fmt.Errorf("sheet '%s' has no headers", sheet.Name)
		}
		if len(sheet.Rows) == 0 {
			return fmt.Errorf("sheet '%s' has no rows", sheet.Name)
		}
	}
	return nil
}
