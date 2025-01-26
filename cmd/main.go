package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"portalscraper/internal/models"
	"portalscraper/internal/ollama"
	"portalscraper/internal/scraper"
)

func main() {
	baseURL := "https://www.portalinmobiliario.com/venta/casa/propiedades-usadas/las-condes-metropolitana"
	maxPages := 1
	results := make([]models.Property, 0)

	// Scraping rápido con límite
	fmt.Println("🔍 Iniciando scraping rápido...")
	for page := 1; page <= maxPages; page++ {
		url := fmt.Sprintf("%s?_PAGE=%d", baseURL, page)
		fmt.Printf("📃 Página %d\n", page)

		props, err := scraper.MainPage(url)
		if err != nil {
			log.Printf("⚠️ Error página %d: %v", page, err)
			break
		}

		results = append(results, props...)
		time.Sleep(1 * time.Second)
	}

	// Limitar resultados para prueba
	if len(results) > 20 {
		results = results[:20]
	}

	mostrarResultados(results)

	// Generar prompt compacto
	fmt.Println("\n🧠 Generando análisis...")
	client := ollama.NewClient()

	prompt := construirPromptCompacto(results)
	fmt.Printf("📤 Prompt enviado (%d caracteres):\n%.200s...\n", len(prompt), prompt)

	respuesta, err := client.Generate(prompt)
	if err != nil {
		log.Fatalf("💥 Error crítico: %v", err)
	}

	fmt.Println("\n📊 Análisis generado:")
	fmt.Println(strings.TrimSpace(respuesta))
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
			prop.Title, // Usar el título real de la propiedad
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

// Funciones de soporte

func extraerValorUF(precio string) float64 {
	re := regexp.MustCompile(`UF(\d+\.?\d*)`)
	match := re.FindStringSubmatch(precio)
	if len(match) > 1 {
		valor, _ := strconv.ParseFloat(match[1], 64)
		return valor // UF16.400 → 16.4 (sin multiplicar por 1000)
	}
	return 0
}
func extraerValorM2(m2 string) float64 {
	re := regexp.MustCompile(`(\d+\.?\d*)\s*m²`)
	match := re.FindStringSubmatch(m2)
	if len(match) > 1 {
		valor, _ := strconv.ParseFloat(match[1], 64)
		return valor // 178 m² útiles → 178.0
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
