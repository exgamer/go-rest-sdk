package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
)

func NewQueryBuilder(client *elasticsearch.Client, indexName string) *QueryBuilder {
	return &QueryBuilder{
		Client:    client,
		IndexName: indexName,
	}
}

type QueryBuilder struct {
	Client    *elasticsearch.Client
	IndexName string
	Query     map[string]interface{}
}

func (builder QueryBuilder) SearchByQuery(query map[string]interface{}) (*map[string]interface{}, error) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s\n", err.Error())

		return nil, err
	}

	res, err := builder.Client.Search(
		builder.Client.Search.WithIndex(builder.IndexName),
		builder.Client.Search.WithBody(&buf),
		builder.Client.Search.WithTrackTotalHits(true),
		builder.Client.Search.WithPretty(),
		builder.Client.Search.WithContext(context.Background()),
	)

	defer res.Body.Close()

	if err != nil {
		log.Printf("Error getting response: %s", err.Error())

		return nil, err
	}

	if res.IsError() {
		var e map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err.Error())

			return nil, err
		} else {
			// Print the response status and error information.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)

			return nil, errors.New("Elastic response error: " + res.Status())
		}
	}

	if err != nil {
		log.Printf("Error getting response: %s\n", err.Error())

		return nil, err
	}

	r := make(map[string]interface{})

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err.Error())
	}

	return &r, nil
}
