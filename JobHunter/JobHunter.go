package JobHunter

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gocolly/colly/v2"
)

type APIResponse struct {
	SearchMetadata   SearchMetadata   `json:"search_metadata"`
	SearchParameters SearchParameters `json:"search_parameters"`
	JobsResults      []JobsResult     `json:"jobs_results"`
	Chips            []Chip           `json:"chips"`
}

type SearchMetadata struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	JsonEndpoint   string  `json:"json_endpoint"`
	CreatedAt      string  `json:"created_at"`
	ProcessedAt    string  `json:"processed_at"`
	GoogleJobsURL  string  `json:"google_jobs_url"`
	RawHtmlFile    string  `json:"raw_html_file"`
	TotalTimeTaken float64 `json:"total_time_taken"`
}

type SearchParameters struct {
	Q                 string `json:"q"`
	Engine            string `json:"engine"`
	LocationRequested string `json:"location_requested"`
	LocationUsed      string `json:"location_used"`
	GoogleDomain      string `json:"google_domain"`
	HL                string `json:"hl"`
	GL                string `json:"gl"`
	Chips             string `json:"chips"`
	Ltype             string `json:"ltype"`
}

type JobsResult struct {
	Title              string             `json:"title"`
	CompanyName        string             `json:"company_name"`
	Location           string             `json:"location"`
	Via                string             `json:"via"`
	Description        string             `json:"description"`
	JobHighlights      []JobHighlight     `json:"job_highlights"`
	RelatedLinks       []RelatedLink      `json:"related_links"`
	Thumbnail          string             `json:"thumbnail"`
	Extensions         []string           `json:"extensions"`
	DetectedExtensions DetectedExtensions `json:"detected_extensions"`
	JobID              string             `json:"job_id"`
}

type JobHighlight struct {
	Title string   `json:"title"`
	Items []string `json:"items"`
}

type RelatedLink struct {
	Link string `json:"link"`
	Text string `json:"text"`
}

type DetectedExtensions struct {
	PostedAt     string `json:"posted_at"`
	ScheduleType string `json:"schedule_type"`
	Salary       string `json:"salary"`
	WorkFromHome bool   `json:"work_from_home"`
}

type Chip struct {
	Type    string       `json:"type"`
	Param   string       `json:"param"`
	Options []ChipOption `json:"options"`
}

type ChipOption struct {
	Text  string `json:"text"`
	Value string `json:"value,omitempty"` // Use omitempty to allow for the absence of this field in some options
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func Scrape(url string) {
	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(url)
}

func GetSerp() []JobsResult {

	SerpKey, serpExists := os.LookupEnv("SERP_API")

	if !serpExists {
		log.Fatal("No Serp Key!")
	}

	log.Println(SerpKey)

	baseURL := "https://serpapi.com/search"

	params := url.Values{}
	params.Add("api_key", SerpKey)
	params.Add("engine", "google_jobs")
	params.Add("google_domain", "google.com")
	params.Add("q", "Software+Engineer")
	params.Add("hl", "en")
	params.Add("gl", "us")
	params.Add("location", "United States")
	params.Add("ltype", "1")
	params.Add("chips", "date_posted:today")

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	checkNilErr(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	checkNilErr(err)

	var response APIResponse

	jsonErr := json.Unmarshal([]byte(body), &response)

	checkNilErr(jsonErr)

	return response.JobsResults
}
