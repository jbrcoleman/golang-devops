package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type WordsPage struct {
	Page
	Words
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	out := []string{}
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurrence))
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func (a api) DoGetRequest(requestURL string) (Response, error) {
	response, err := a.Client.Get(requestURL)

	if err != nil {
		return nil, fmt.Errorf("http Get error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output (HTTP code %d): %s\n", response.StatusCode, string(body))
	}

	var page Page

	if !json.Valid(body) {
		return nil, RequestError{
			Err:      fmt.Sprintf("Response is not a json"),
			HTTPCode: response.StatusCode,
			Body:     string(body),
		}
	}

	err = json.Unmarshal(body, &page)

	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page Unmarshal error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words

		err = json.Unmarshal(body, &words)

		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words Unmarshal error: %s", err),
			}
		}

		return words, nil
	case "occurrence":
		var occurrence Occurrence

		err = json.Unmarshal(body, &occurrence)

		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurrences Unmarshal error: %s", err),
			}
		}

		return occurrence, nil
	}

	return nil, nil

}
