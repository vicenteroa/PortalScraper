package excel

import (
	"fmt"
	"portalscraper/internal/models"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportToExcel(props []models.Property) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Crear nueva hoja
	index, err := f.NewSheet("Propiedades")
	if err != nil {
		return "", err
	}

	// Encabezados
	headers := []string{"Título", "Precio", "Ubicación", "m²", "Dormitorios", "Baños", "Enlace"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Propiedades", cell, h)
	}

	// Datos
	for rowIdx, prop := range props {
		row := rowIdx + 2 // Empieza en fila 2
		f.SetCellValue("Propiedades", fmt.Sprintf("A%d", row), prop.Title)
		f.SetCellValue("Propiedades", fmt.Sprintf("B%d", row), prop.Price)
		f.SetCellValue("Propiedades", fmt.Sprintf("C%d", row), prop.Location)
		f.SetCellValue("Propiedades", fmt.Sprintf("D%d", row), prop.M2)
		f.SetCellValue("Propiedades", fmt.Sprintf("E%d", row), prop.Bedrooms)
		f.SetCellValue("Propiedades", fmt.Sprintf("F%d", row), prop.Bathrooms)
		f.SetCellValue("Propiedades", fmt.Sprintf("G%d", row), prop.Link)
	}

	f.SetActiveSheet(index)

	filename := fmt.Sprintf("propiedades_%s.xlsx", time.Now().Format("20060102_150405"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}

	return filename, nil
}
