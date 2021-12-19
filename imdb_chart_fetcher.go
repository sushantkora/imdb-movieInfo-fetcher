package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gocolly/colly"
	"flag"
)

type movie struct {
	Title      			string `json:"title"`
	Movie_release_year  string `json:"movie_release_year"`
	Imdb_rating     	string `json:"imdb_rating"`
	Summary				string `json:"summary"`
	Duration 			string `json:"duration"`
	Genre 				string `json:"genre"`
}



func main() {

	c_url := flag.String("chart_url", "https://www.imdb.com/india/top-rated-telugu-movies", "")
	i_count := flag.Int("items_count", 1, "")
    flag.Parse()

    chart_url := *c_url
	
 	
	var m_url string =chart_url
	var itemsCount int=*i_count

	fmt.Println("Fetchig data...")
	chart_fetcher(itemsCount,m_url)
}

func chart_fetcher(n_moviesList int, movie_url string) {

	movie_list := []movie{}
	items_count := 0

	if n_moviesList == 0 {
		n_moviesList = 1
	}

	m_datacollector := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com", "imdb.com"),
	)
	colly.Async(true)

	per_movieCollector := m_datacollector.Clone()

	m_datacollector.OnHTML(".titleColumn", func(element *colly.HTMLElement) {
		
		per_movieUrl := element.ChildAttr("a", "href")
		per_movieUrl = element.Request.AbsoluteURL(per_movieUrl)
		
		if items_count < n_moviesList {			
			per_movieCollector.Visit(per_movieUrl)
		}

	})

	per_movieCollector.OnHTML(".fARFJI", func(element *colly.HTMLElement) {
		movie_info := movie{}
		
	
		movie_info.Title = element.ChildText(".dxSWFG")	
		movie_info.Movie_release_year = element.ChildText(".rgaOW")
		imdbRating := element.ChildText(".iTLWoV");
		movie_info.Imdb_rating = imdbRating[3:]
		movie_info.Summary = element.ChildText("span.dcFkRD")
		movie_duration := element.ChildText("div.hWHMKr > ul > li.ipc-inline-list__item")
		movie_info.Duration = movie_duration
		if len(movie_duration) >= 6 { 
			movie_info.Duration = movie_duration[len(movie_duration)-6:len(movie_duration)]
		}
		movie_info.Genre = element.ChildText("span.ipc-chip__text")		
		movie_list = append(movie_list, movie_info)
		items_count++
	
	})

	//m_datacollector.Visit("https://www.imdb.com/india/top-rated-indian-movies")
	m_datacollector.Visit(movie_url)

	movies_jsonString, err := json.MarshalIndent(movie_list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(movies_jsonString))
	fmt.Println("Success !")

}

