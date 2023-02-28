package plugin

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestQueryData(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("id", "secretKey", ""),
		),
	)
	if err != nil {
		t.Error(err)
	}

	s3PresignClient := s3.NewPresignClient(s3.NewFromConfig(cfg))
	ds := Datasource{bucket: "test", presigner: s3PresignClient}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: []byte("{\"datasource\":{\"type\":\"babylonmf-farmpics-datasource\",\"uid\":\"PvTOH7x4z\"},\"datasourceId\":1,\"image_keys\":\"33884862/packs/c0de1446-0000-feed-f00d-5a1ad2c0ffee/zone1_20230223-175257_DEF.png, 33884862/packs/c0de1446-0000-feed-f00d-5a1ad2c0ffee/zone2_20230223-175348_DEF.png\",\"intervalMs\":2000,\"key\":\"Q-d6afae93-a253-47c4-a14b-9f6606edbddc-0\",\"maxDataPoints\":1863,\"refId\":\"A\"}")},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}
