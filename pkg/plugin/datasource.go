package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces- only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(dataSourceSettings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {

	log.DefaultLogger.Debug("settings from front end", "json data", dataSourceSettings.JSONData, "secure json data", dataSourceSettings.DecryptedSecureJSONData)
	var settings map[string]string
	err := json.Unmarshal(dataSourceSettings.JSONData, &settings)
	if err != nil {
		return nil, err
	}
	log.DefaultLogger.Debug("parsed json data", settings)
	bucket := settings["bucket"]
	id := dataSourceSettings.DecryptedSecureJSONData["aws_access_key_id"]
	secretKey := dataSourceSettings.DecryptedSecureJSONData["aws_secret_access_key"]

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(id, secretKey, ""),
		),
	)

	if err != nil {
		return nil, err
	}

	s3PresignClient := s3.NewPresignClient(s3.NewFromConfig(cfg))
	return &Datasource{bucket: bucket, presigner: s3PresignClient}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	bucket    string
	presigner *s3.PresignClient
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// when logging at a non-Debug level, make sure you don't include sensitive information in the message
	// (like the *backend.QueryDataRequest)
	log.DefaultLogger.Debug("QueryData called", "numQueries", len(req.Queries))

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	ImageKeys string `json:"image_keys"`
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	//get signed URLs from comma separated object keys
	imageKeys := strings.Split(strings.ReplaceAll(qm.ImageKeys, " ", ""), ",")
	log.DefaultLogger.Debug("split image keys", "keys", imageKeys)

	var signedUrls []string
	for _, key := range imageKeys {
		signedRequest, err := d.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &d.bucket,
			Key:    &key,
		})
		if err != nil {
			log.DefaultLogger.Error("Error signing request", "err", err.Error())
		}

		signedUrls = append(signedUrls, signedRequest.URL)
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/
	frame := data.NewFrame("response")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("URL", nil, signedUrls),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// when logging at a non-Debug level, make sure you don't include sensitive information in the message
	// (like the *backend.QueryDataRequest)
	log.DefaultLogger.Debug("CheckHealth called")

	signedRequest, err := d.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &d.bucket,
		Key:    aws.String("some/important/object"),
	})
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, err
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("ready to generate urls like: %s", signedRequest.URL),
	}, nil
}
