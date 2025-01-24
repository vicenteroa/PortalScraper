package main

import (
	"fmt"
	"log"
	"time"

	"portalscraper/internal/models"
	"portalscraper/internal/scraper"
)

func main() {
	baseURL := "https://www.portalinmobiliario.com/venta/casa/propiedades-usadas/las-condes-metropolitana"
	maxPages := 3
	results := make([]models.Property, 0)

	for page := 1; page <= maxPages; page++ {
		url := fmt.Sprintf("%s?_PAGE=%d", baseURL, page)
		fmt.Printf("Scrapeando página %d: %s\n", page, url)

		props, err := scraper.MainPage(url)
		if err != nil {
			log.Printf("Error en página %d: %v", page, err)
			break
		}

		results = append(results, props...)
		time.Sleep(2 * time.Second)
	}

	printResults(results)
}

func printResults(results []models.Property) {
	fmt.Printf("\nTotal propiedades: %d\n", len(results))
	for i, prop := range results {
		fmt.Printf("\nPropiedad #%d:\n", i+1)
		fmt.Printf("Título: %s\n", prop.Title)
		fmt.Printf("Precio: %s\n", prop.Price)
		fmt.Printf("Ubicación: %s\n", prop.Location)
		fmt.Printf("m²: %s\n", prop.M2)
		fmt.Printf("Dormitorios: %s\n", prop.Bedrooms)
		fmt.Printf("Baños: %s\n", prop.Bathrooms)
		fmt.Printf("Enlace: %s\n", prop.Link)
	}
}
