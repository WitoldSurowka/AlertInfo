package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GoogleSearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

// Define a struct to parse the JSON
var data struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func scrapeWebpage(pageURL string) (string, error) {
	// Fetch the webpage content
	resp, err := http.Get(pageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch webpage: %w", err)
	}
	defer resp.Body.Close()

	// Parse the URL into *url.URL type
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract the main content using readability
	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", fmt.Errorf("failed to extract main content: %w", err)
	}

	return article.Content, nil
}

func askGPT(request string) string {
	// Set your OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY") // Make sure to set the OPENAI_API_KEY environment variable
	//fmt.Println("klucz:" + apiKey)
	// API endpoint
	url := "https://api.openai.com/v1/chat/completions"

	// JSON payload
	payload := `{
       "model": "gpt-4o-mini",
       "messages": [
          {
             "role": "system",
             "content": "You are a helpful assistant."
          },
          {
             "role": "user",
             "content": "` + request + `"
          }
       ]
    }`
	//fmt.Println(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making API request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Parse the JSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Error Unmarshalling: %v", err)
	}

	//fmt.Println(string(body), data)
	// Print the response
	//fmt.Println(data.Choices[0].Message.Content)
	output := data.Choices[0].Message.Content
	if len(output) > 250 {
		output = output[0:249]
	}
	return output
}

func stripHTMLTags(input string) (string, error) {
	// Create a new reader from the input string
	reader := strings.NewReader(input)

	// Parse the HTML content
	tokenizer := html.NewTokenizer(reader)

	var result strings.Builder

	// Iterate over the tokens (elements) in the HTML
	for {
		// Get the next token
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			// If we reach the end of the document, return the result
			return result.String(), nil
		case html.TextToken:
			// Append text content to the result if the token is text
			text := tokenizer.Text()
			result.WriteString(string(text))
		}
	}
}

func compactText(input string) string {
	// Split the string by newline characters
	lines := strings.Split(input, "\n")

	// Filter out empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	outcome := strings.Join(nonEmptyLines, "\n")
	outcome = strings.TrimSpace(outcome)
	outcome = strings.Join(strings.Fields(outcome), " ")
	outcome = strings.ReplaceAll(outcome, `"`, `\"`)
	outcome = strings.ReplaceAll(outcome, "\n", " ")
	outcome = strings.ReplaceAll(outcome, `\`, `\\`)

	// Join the non-empty lines back into a single string
	return outcome
}

func SearchGoogle(query string) ([]string, error) {
	apiKey := os.Getenv("GOOGLE_SEARCH_API_KEY")
	//fmt.Println("klucz:" + apiKey)
	cx := "123abc" // Your custom search engine ID
	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&cx=%s&key=%s", query, cx, apiKey)
	//url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&cx=%s&key=%s&sort=date", query, cx, apiKey)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var searchResponse GoogleSearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Collect the first three links
	var links []string
	for i, item := range searchResponse.Items {
		if i >= 6 {
			break
		}
		links = append(links, item.Link)
	}

	return links, nil
}

func normalizeQuery(query string) string {
	return url.QueryEscape(query)
}

func findInfoForAlert(usersQuery string) string {

	prefix1 := `jako odpowiedz podaj JEDYNIE link do strony internetowej, gdzie można zaleźć informację odpowiadjącą na pyatnie: `
	prefix2 := strings.ToUpper(usersQuery)
	queryGPTlinks := prefix1 + prefix2

	GPTwebsiteProposition := askGPT(queryGPTlinks)
	//fmt.Println(GPTwebsiteProposition)
	normalizedQuery := normalizeQuery(GPTwebsiteProposition + " " + usersQuery)
	//fmt.Println(normalizedQuery)
	googleLinks, err := SearchGoogle(normalizedQuery)
	if len(googleLinks) < 6 {
		normalizedQuery = normalizeQuery(usersQuery)
		googleLinks, err = SearchGoogle(normalizedQuery)
	}
	if len(googleLinks) < 6 {
		return "błąd podczas przeszukiwania internetu"
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "błąd podczas przeszukiwania internetu"
	}

	err = errors.New("something went wrong")
	context := "błąd"
	n := 0
	for err != nil {
		if n == 6 {
			log.Fatalf("zadna z 6 stron google nie została prawidłowo pobrana przez funckje scrapeWebpage. Ostatniłąd: %v", err)
			break
		}
		context, err = scrapeWebpage(googleLinks[n])
		n += 1
	}
	fmt.Println(googleLinks[n])
	n = 0

	strippedText, err := stripHTMLTags(context)
	if err != nil {
		fmt.Println("Error:", err)
		return "błąd podczas przeszukiwania internetu"
	}
	strippedComapctText := compactText(strippedText)

	prefix1 = `z podanego kontekstu wydobądź informację, będącą odpowiedzią na zapytanie: `
	prefix3 := `. przedstaw tą informację w sposób minimalistyczny; maksymalna liczba znaków odpowiedzi to 250: `
	queryGPTanswer := prefix1 + prefix2 + prefix3 + strippedComapctText

	return askGPT(queryGPTanswer)
}
