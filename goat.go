package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type QueryProduct struct {
	Results []struct {
		Hits []struct {
			ProductTemplateID int       `json:"product_template_id"`
			ShoeCondition     string    `json:"shoe_condition"`
			Name              string    `json:"Name"`
			Category          []string  `json:"category"`
			Color             string    `json:"color"`
			SizeRange         []float64 `json:"size_range"`
			Sku               string    `json:"sku"`
		} `json:"hits"`
	} `json:"results"`
}

type QueryPrices []struct {
	Size             float32 `json:"size"`
	Lowestpricecents struct {
		Currency string `json:"currency"`
		Amount   int    `json:"amount"`
	} `json:"lowestPriceCents"`
}

func goatSearch(query, size string) (int, string, string) {
	url := "https://2fwotdvm2o-dsn.algolia.net/1/indexes/*/queries?x-algolia-agent=Algolia%2520for%2520JavaScript%2520(3.35.1)%253B%2520Browser%2520(lite)%253B%2520JS%2520Helper%2520(3.2.2)%253B%2520react%2520(16.13.1)%253B%2520react-instantsearch%2520(6.8.2)&x-algolia-application-id=2FWOTDVM2O&x-algolia-api-key=ac96de6fef0e02bb95d433d8d5c7038a"
	method := "POST"

	payload := strings.NewReader(`{"requests":[{"indexName":"product_variants_v2","params":"highlightPreTag=%3Cais-highlight-0000000000%3E&highlightPostTag=%3C%2Fais-highlight-0000000000%3E&distinct=true&query=` + query + `&maxValuesPerFacet=30&page=0&facets=%5B%22product_category%22%2C%22instant_ship_lowest_price_cents%22%2C%22single_gender%22%2C%22presentation_size%22%2C%22shoe_condition%22%2C%22brand_name%22%2C%22color%22%2C%22silhouette%22%2C%22designer%22%2C%22upper_material%22%2C%22midsole%22%2C%22category%22%2C%22release_date_name%22%5D&tagFilters=&facetFilters=%5B%5B%22product_category%3Ashoes%22%5D%5D"},{"indexName":"product_variants_v2","params":"highlightPreTag=%3Cais-highlight-0000000000%3E&highlightPostTag=%3C%2Fais-highlight-0000000000%3E&distinct=true&query=` + query + `&maxValuesPerFacet=30&page=0&hitsPerPage=1&attributesToRetrieve=%5B%5D&attributesToHighlight=%5B%5D&attributesToSnippet=%5B%5D&tagFilters=&analytics=false&clickAnalytics=false&facets=product_category"}]}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)

	}
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Content-Length", "953")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}

	var cont QueryProduct
	if err := json.Unmarshal(body, &cont); err != nil {
		fmt.Println("error:", err)

	}
	pid := cont.Results[0].Hits[0].ProductTemplateID

	shoeTitle := cont.Results[0].Hits[0].Name
	return pid, shoeTitle, size
}

func goatPrices(pid int, title string, size string) {

	url := ("https://www.goat.com/web-api/v1/product_variants?productTemplateId=" + strconv.Itoa(pid))
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Cookie", "_csrf=vfc0aHK-yzf2Bcz0ueZMN_ti")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var sizing QueryPrices
	if err := json.Unmarshal(body, &sizing); err != nil {
		// panic(err)
		fmt.Println("error:", err)

	}

	validateSize := size
	var sizeIndex int
	sizeIndex = 0

	for i, v := range sizing {
		if fmt.Sprintf("%.1f", v.Size) == validateSize {
			sizeIndex = i
		}
	}

	lowestPrice := (sizing[sizeIndex].Lowestpricecents.Amount / 100)
	atSize := fmt.Sprintf("%.1f", sizing[sizeIndex].Size)
	fmt.Fprintf(color.Output, "For A Size %s of the %s, The Current Lowest Price is: $%s\n", color.MagentaString(atSize), color.HiRedString(title), color.GreenString(strconv.Itoa(lowestPrice)))
}

func main() {

	size := os.Args[1]
	prodname := strings.Join(os.Args[2:], " ")
	goatPrices(goatSearch(prodname, size))

}
