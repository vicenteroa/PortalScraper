package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"portalscraper/internal/excel"
	"portalscraper/internal/models"
	"portalscraper/internal/ollama"
	"portalscraper/internal/scraper"
)

func main() {
	baseURL := "https://www.portalinmobiliario.com/venta/casa/propiedades-usadas/las-condes-metropolitana"
	results := make([]models.Property, 0)

	// Configuración de scraping
	var maxPages, maxProperties int

	// Solicitar cantidad de páginas
	for {
		fmt.Print("Ingrese cantidad de páginas a scrapear: ")
		_, err := fmt.Scanln(&maxPages)
		if err != nil || maxPages < 1 {
			fmt.Println("Error: Debe ingresar un número mayor a 0")
			continue
		}
		break
	}

	// Solicitar máximo de propiedades
	for {
		fmt.Print("Ingrese máximo de propiedades a obtener: ")
		_, err := fmt.Scanln(&maxProperties)
		if err != nil || maxProperties < 1 {
			fmt.Println("Error: Debe ingresar un número mayor a 0")
			continue
		}
		break
	}

	// Ejecutar scraping
	fmt.Printf("\n🔍 Iniciando scraping (%d páginas, máximo %d propiedades)...\n", maxPages, maxProperties)

	for page := 1; page <= maxPages; page++ {
		if len(results) >= maxProperties {
			break
		}

		url := fmt.Sprintf("%s?_PAGE=%d", baseURL, page)
		fmt.Printf("📃 Página %d\n", page)

		props, err := scraper.MainPage(url)
		if err != nil {
			log.Printf("⚠️ Error página %d: %v", page, err)
			break
		}

		// Controlar máximo de propiedades
		remaining := maxProperties - len(results)
		if len(props) > remaining {
			props = props[:remaining]
		}

		results = append(results, props...)
		time.Sleep(1 * time.Second)
	}

	// Asegurar no exceder el máximo
	if len(results) > maxProperties {
		results = results[:maxProperties]
	}

	mostrarMenu(results)
}

// Resto del código sin modificaciones
func mostrarMenu(props []models.Property) {
	for {
		clearScreen()
		mostrarEncabezado()

		var opcion int
		fmt.Print("\nSelecciona una opción: ")
		_, err := fmt.Scanln(&opcion)
		if err != nil {
			fmt.Println("Error: Ingresa un número válido")
			time.Sleep(1 * time.Second)
			continue
		}

		switch opcion {
		case 1:
			analizarPropiedades(props)
			return
		case 2:
			exportarExcel(props)
			return
		default:
			fmt.Println("Opción no válida")
			time.Sleep(1 * time.Second)
		}
	}
}

func mostrarEncabezado() {
	titulo := `
╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║██████╗  ██████╗ ██████╗ ████████╗ █████╗ ██╗         ███████╗ ██████╗██████╗  █████╗ ██████╗ ███████╗██████╗║
║██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔══██╗██║         ██╔════╝██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗║
║██████╔╝██║   ██║██████╔╝   ██║   ███████║██║         ███████╗██║     ██████╔╝███████║██████╔╝█████╗  ██████╔╝║
║██╔═══╝ ██║   ██║██╔══██╗   ██║   ██╔══██║██║         ╚════██║██║     ██╔══██╗██╔══██║██╔═══╝ ██╔══╝  ██╔══██╗║
║██║     ╚██████╔╝██║  ██║   ██║   ██║  ██║███████╗    ███████║╚██████╗██║  ██║██║  ██║██║     ███████╗██║  ██║║
║╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚══════╝    ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝║
║                                        v1.1.0 - Portal Scraper                                      
                                             @vicenteroa
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`

	fmt.Println(titulo)
	//fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
	fmt.Println("                                          ¿Qué deseas hacer?                ")
	fmt.Println("═════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
	fmt.Println("                                     1. Análisis Inteligente con IA")
	fmt.Println("                                     2. Exportar datos a Excel")
	fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════\n")
}

func analizarPropiedades(props []models.Property) {
	clearScreen()
	mostrarResultados(props)

	// Generar prompt compacto
	fmt.Println("\n🧠 Generando análisis...")
	client := ollama.NewClient()

	prompt := construirPromptCompacto(props)
	fmt.Printf("📤 Prompt enviado (%d caracteres):\n%.200s...\n", len(prompt), prompt)

	respuesta, err := client.Generate(prompt)
	if err != nil {
		log.Fatalf("💥 Error crítico: %v", err)
	}

	fmt.Println("\n📊 Análisis generado:")
	fmt.Println(strings.TrimSpace(respuesta))
}

func exportarExcel(props []models.Property) {
	clearScreen()
	fmt.Println("📤 Exportando a Excel...")

	filename, err := excel.ExportToExcel(props)
	if err != nil {
		log.Fatalf("💥 Error exportando: %v", err)
	}

	fmt.Printf("✅ Archivo creado: %s\n", filename)
	fmt.Println("📁 Revisa el archivo Excel en tu directorio actual")
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func construirPromptCompacto(props []models.Property) string {
	var sb strings.Builder
	sb.WriteString("Analiza las propiedades calculando UF/m² y detecta oportunidades:\n")
	sb.WriteString("Instrucciones:\n")
	sb.WriteString("1. Calcular UF/m² para cada propiedad (Precio_UF / m²)\n")
	sb.WriteString("2. Ordenar de menor a mayor UF/m²\n")
	sb.WriteString("3. Identificar las 3 con menor ratio como oportunidades\n\n")
	sb.WriteString("Datos en formato tabla:\n")
	sb.WriteString("| # | Propiedad | UF | m² | UF/m² |\n")
	sb.WriteString("|---|---|---|---|---|\n")

	for i, prop := range props {
		uf := extraerValorUF(prop.Price)
		m2 := extraerValorM2(prop.M2)

		sb.WriteString(fmt.Sprintf("| %d | %s | %.2f UF | %.2f m² | \n",
			i+1,
			prop.Title,
			uf,
			m2,
		))
	}

	sb.WriteString("\nFormato requerido para la respuesta:\n")
	sb.WriteString("Oportunidades detectadas:\n")
	sb.WriteString("1. [Nombre propiedad] - UF/m²: [valor] (Motivo)\n")
	sb.WriteString("2. [Nombre propiedad] - UF/m²: [valor] (Motivo)\n")
	sb.WriteString("3. [Nombre propiedad] - UF/m²: [valor] (Motivo)\n")
	sb.WriteString("Análisis comparativo: [Breve explicación]")

	return sb.String()
}

func extraerValorUF(precio string) float64 {
	re := regexp.MustCompile(`UF(\d+\.?\d*)`)
	match := re.FindStringSubmatch(precio)
	if len(match) > 1 {
		valor, _ := strconv.ParseFloat(match[1], 64)
		return valor
	}
	return 0
}

func extraerValorM2(m2 string) float64 {
	re := regexp.MustCompile(`(\d+\.?\d*)\s*m²`)
	match := re.FindStringSubmatch(m2)
	if len(match) > 1 {
		valor, _ := strconv.ParseFloat(match[1], 64)
		return valor
	}
	return 0
}

func mostrarResultados(props []models.Property) {
	fmt.Printf("\n🏠 %d propiedades encontradas:\n", len(props))
	for i, prop := range props {
		fmt.Printf("\n%d. %s\n", i+1, prop.Title)
		fmt.Printf("   💵 %s\n", prop.Price)
		fmt.Printf("   📏 %s\n", prop.M2)
		fmt.Printf("   📍 %s\n", prop.Location)
		fmt.Printf("  🛜 %s\n", prop.Link)
	}
}
