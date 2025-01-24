package scraper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"portalscraper/internal/models"
)

func MainPage(url string) ([]models.Property, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return parseProperties(res)
}

func parseProperties(res *http.Response) ([]models.Property, error) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var properties []models.Property

	doc.Find("li.ui-search-layout__item").Each(func(i int, s *goquery.Selection) {
		priceSymbol := cleanText(s.Find(".andes-money-amount__currency-symbol").First().Text())
		priceAmount := cleanText(s.Find(".andes-money-amount__fraction").First().Text())

		prop := models.Property{
			Title:    cleanText(s.Find(".poly-component__title-wrapper").Text()),
			Price:    priceSymbol + priceAmount,
			Location: cleanText(s.Find(".poly-component__location").Text()),
			Link:     extractLink(s),
		}

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
