package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Result struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

func main() {
	start := time.Now()

	// list of urls we are scraping
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// slice to store results
	var results []Result

	//print the url we are requesting
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// return entire content
	c.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
		result := Result{
			URL:     e.Request.URL.String(),
			Content: e.Text,
		}
		results = append(results, result)
	})

	// loop to visit each website in slice
	for _, url := range urls {
		err := c.Visit(url)
		if err != nil {
			log.Println("Failed to visit:", url, "Error:", err)
		}
	}

	// create jsonl file
	file, err := os.Create("output.jsonl")
	if err != nil {
		log.Fatal("Could not create file:", err)
	}
	defer file.Close()

	// write results to jsonl file
	for _, result := range results {
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			continue
		}
		_, err = file.Write(append(jsonData, '\n'))
		if err != nil {
			log.Println("Error writing to file:", err)
		}
	}
	duration := time.Since(start)
	fmt.Println("Scraping finished and data saved to output.jsonl")
	fmt.Printf("Run time: %s\n", duration)
}
