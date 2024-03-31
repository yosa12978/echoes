package data

import (
	"sync"

	"github.com/elastic/go-elasticsearch"
)

var (
	es          *elasticsearch.Client
	elasticOnce sync.Once
)

func Elastic() *elasticsearch.Client {
	elasticOnce.Do(func() {
		client, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{"localhost:9200"},
		})
		if err != nil {
			panic(err)
		}
		es = client
	})
	return es
}
