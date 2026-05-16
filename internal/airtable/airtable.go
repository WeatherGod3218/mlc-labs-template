package airtable

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"net/http"
	"net/url"
	"os"
)

type AirtableRecord struct {
	ID          string       `json:"id"`
	CreatedTime time.Time    `json:"createdTime"`
	Fields      AirtableData `json:"fields"`
}

type AirtableAPIResponse struct {
	Records []AirtableRecord `json:"records"`
	Offset  string           `json:"offset"`
}

type AirtableData struct {
	Block           int    `json:"block"`
	Item            string `json:"item"`
	Stimulus        string `json:"stimulus"`
	CorrectKey      string `json:"correct_key"`
	StimulusType    string `json:"stimulus_type"`
	Trial           int    `json:"trial"`
	Category        string `json:"category"`
	Order           int    `json:"order"`
	TrialType       string `json:"trial_type"`
	CategoryDisplay string `json:"category_display"`
}

type AirtableClientResponse struct {
	Block           int    `json:"block"`
	Item            string `json:"item"`
	Stimulus        string `json:"stimulus"`
	CorrectKey      string `json:"correct_key"`
	StimulusType    string `json:"stimulus_type"`
	Trial           int    `json:"trial"`
	Category        string `json:"category"`
	Order           int    `json:"order"`
	TrialType       string `json:"trial_type"`
	CategoryDisplay string `json:"category_display"`
	Association     string `json:"association"`
}

func GetAirtableURI(table string) ([]AirtableRecord, error) {
	base := os.Getenv("AIRTABLE_BASE")
	baseURL := fmt.Sprintf(
		"https://api.airtable.com/v0/%s/%s",
		base,
		url.PathEscape(table),
	)

	var allRecords []AirtableRecord
	var offset string

	client := &http.Client{}

	for {
		fetchURL := baseURL

		if offset != "" {
			fetchURL += "?offset=" + url.QueryEscape(offset)
		}

		req, err := http.NewRequest("GET", fetchURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+os.Getenv("AIRTABLE_API_KEY"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(
				"airtable error: %d %s",
				resp.StatusCode,
				string(body),
			)
		}

		var result AirtableAPIResponse

		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		allRecords = append(allRecords, result.Records...)

		if result.Offset == "" {
			break
		}

		offset = result.Offset
	}

	return allRecords, nil
}

func GetResponse(field AirtableData, association string) AirtableClientResponse {
	return AirtableClientResponse{
		Block:           field.Block,
		Item:            field.Item,
		Stimulus:        field.Stimulus,
		CorrectKey:      field.CorrectKey,
		StimulusType:    field.StimulusType,
		Trial:           field.Block,
		Category:        field.Category,
		Order:           field.Block,
		TrialType:       field.TrialType,
		CategoryDisplay: field.CategoryDisplay,
		Association:     association,
	}
}
