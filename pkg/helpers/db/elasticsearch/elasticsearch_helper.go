package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"log"
	"os"
)

func InitElasticNoSecurityClient(config *structures.ElasticNoSecurityConfig) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.ElasticHost,
		},
	}

	if config.ElasticShowLogs {
		cfg.Logger = &estransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true}
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	_, iErr := es.Info()

	if iErr != nil {
		log.Fatalf("Cannot connect to Elastic. Err: %s", iErr)
	}

	return es, nil
}
