package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Monitor interface {
	Collect(ctx context.Context, url string) (int, error)
	Evaluate(ctx context.Context) ([]Message, error)
}

type GamestopHandler struct {
	Name                  string
	Products              []Products
	Message               []Message
	alreadyPingedProducts map[string]bool
}

func NewGamestopHandler(name string) (Monitor, error) {
	return &GamestopHandler{
		Name:                  name,
		Products:              []Products{},
		Message:               []Message{},
		alreadyPingedProducts: make(map[string]bool),
	}, nil
}

type Products struct {
	ID           string
	Title        string
	ImageUrl     string
	Link         string
	Series       string
	Price        string
	Availability bool
}

func getHtml(url string) *http.Response {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil
	}

	if res.StatusCode > 400 {
		fmt.Printf("Status code: %v \n", res.StatusCode)
	}

	return res
}

type DataProductAttribute struct {
	Id    string
	Name  string
	Price string
	Brand string
}

func (g *GamestopHandler) scrapePageData(doc *goquery.Document) (p []Products) {
	doc.Find("div.prodList>div.searchTileLayout").Each(func(index int, item *goquery.Selection) {
		dataProductAttribute, _ := item.Attr("data-product")

		var dataProductInstances []DataProductAttribute

		err := json.Unmarshal([]byte(dataProductAttribute), &dataProductInstances)

		if err != nil {
			fmt.Println(fmt.Sprintf("failed to unmarshall json into go struct: %s", err.Error()))
			return
		}

		dataProductInstance := dataProductInstances[0]

		link, _ := item.Find("div.searchTilePriceDesktop>h3.desktopSearchProductTitle>a").Attr("href")
		series := item.Find("div.searchTilePriceDesktop>h4.platLogo").Text()
		imageUrl, _ := item.Find("div.searchProductImage").First().Find("img").First().Attr("data-llsrc")
		availability := item.Find("button").HasClass("SPTenabled")

		p = append(p, Products{
			ID:           dataProductInstance.Id,
			Title:        dataProductInstance.Name,
			Link:         "https://www.gamestop.de" + link,
			ImageUrl:     imageUrl,
			Series:       series,
			Price:        dataProductInstance.Price,
			Availability: availability,
		})
	})
	return p
}

func (g *GamestopHandler) Collect(ctx context.Context, url string) (int, error) {
	// fetch data. return statuscode or error
	response := getHtml(url)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return 0, err
	}

	data := g.scrapePageData(doc)

	g.Products = data

	return response.StatusCode, nil
}

func (g *GamestopHandler) Evaluate(ctx context.Context) (m []Message, err error) {
	// evaluate data. return list of messages or error
	for _, p := range g.Products {
		isProductAvailable := p.Availability

		if !isProductAvailable {
			if g.alreadyPingedProducts[p.ID] {
				g.alreadyPingedProducts[p.ID] = false
			}

			continue
		}

		if g.alreadyPingedProducts[p.ID] == true {
			fmt.Println(fmt.Sprintf("skipping product %s - was already pinged", p.ID))
			continue
		}

		m = append(m, Message{
			ID:       p.ID,
			Name:     p.Title,
			Price:    p.Price,
			ImageUrl: p.Link,
		})

		g.alreadyPingedProducts[p.ID] = true
	}

	return m, nil
}
