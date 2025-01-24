package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Property struct {
	Title     string
	Price     string
	Location  string
	M2        string // Campo exportado (mayúscula)
	Bedrooms  string
	Bathrooms string
	Link      string
}

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}

func main() {
	baseURL := "https://www.portalinmobiliario.com/venta/casa/propiedades-usadas/las-condes-metropolitana"
	maxPages := 3
	results := make([]Property, 0)

	for page := 1; page <= maxPages; page++ {
		url := fmt.Sprintf("%s?_PAGE=%d", baseURL, page)
		fmt.Printf("Scrapeando página %d: %s\n", page, url)

		props, err := scrapeMainPage(url)
		if err != nil {
			log.Printf("Error en página %d: %v", page, err)
			break
		}

		results = append(results, props...)
		time.Sleep(2 * time.Second)
	}

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

func scrapeMainPage(url string) ([]Property, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var properties []Property

	doc.Find("li.ui-search-layout__item").Each(func(i int, s *goquery.Selection) {
		priceSymbol := cleanText(s.Find(".andes-money-amount__currency-symbol").First().Text())
		priceAmount := cleanText(s.Find(".andes-money-amount__fraction").First().Text())

		prop := Property{
			Title:    cleanText(s.Find(".poly-component__title-wrapper").Text()),
			Price:    priceSymbol + priceAmount,
			Location: cleanText(s.Find(".poly-component__location").Text()),
			Link:     extractLink(s),
		}

		// Extraer atributos de manera consistente
		attributes := s.Find(".poly-attributes-list__item.poly-attributes-list__separator")
		if attributes.Length() >= 3 {
			prop.Bedrooms = cleanText(attributes.Eq(0).Text())
			prop.Bathrooms = cleanText(attributes.Eq(1).Text())
			prop.M2 = cleanText(attributes.Eq(2).Text())
		}

		properties = append(properties, prop)
	})

	return properties, nil
}

func cleanText(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}

func extractLink(s *goquery.Selection) string {
	link, exists := s.Find("a.poly-component__title").Attr("href")
	if !exists {
		return ""
	}
	return link
}

// Función temporalmente deshabilitada
func scrapeDetailPage(_ string) (map[string]string, error) {
	return nil, nil
}
