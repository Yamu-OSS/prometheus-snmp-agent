package prometheus

import (
	"context"
	"time"

	"github.com/Yamu-OSS/snmp-agent/internal/pkg/log"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusModel "github.com/prometheus/common/model"
)

type Client struct {
	client v1.API
}

func New(addr string) *Client {
	c, err := api.NewClient(api.Config{
		Address: addr,
	})
	if err != nil {
		panic(err)
	}

	api := v1.NewAPI(c)
	return &Client{
		client: api,
	}
}

func (c *Client) Query(ctx context.Context, query string, ts time.Time) (prometheusModel.Value, error) {
	v, w, err := c.client.Query(ctx, query, ts)

	if len(w) > 0 {
		log.Info("Query", "warning", w)
	}

	return v, err
}
