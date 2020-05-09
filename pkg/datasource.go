package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/prometheus/client_golang/prometheus"
)

const metricNamespace = "departureboard_io_datasource"

// DepatureBoardIODataSource handler for departureboard.io API.
type DepatureBoardIODataSource struct {
	client http.Client
}

// NewDataSource creates the departureboard.io datasource
func NewDataSource(mux *http.ServeMux) *DepatureBoardIODataSource {
	queriesTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "data_query_total",
			Help:      "data query counter",
			Namespace: metricNamespace,
		},
		[]string{"scenario"},
	)

	ds := &DepatureBoardIODataSource{
		client: http.Client{},
	}

	prometheus.MustRegister(queriesTotal)
	return ds
}

// CheckHealth checks if the plugin is running properly
func (ds *DepatureBoardIODataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}

	// Just checking that the plugin exe is alive and running
	if req.PluginContext.DataSourceInstanceSettings == nil {
		res.Status = backend.HealthStatusOk
		res.Message = "Plugin is Running"
		return res, nil
	}

	settings, err := LoadSettings(req.PluginContext)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Invalid config"
		return res, nil
	}

	url := settings.APIEndpoint + "/getDeparturesByCRS/PAD?apikey=" + settings.APIKey
	httpReq, err := http.NewRequest("GET", url, nil)
	httpReq.Header.Add("X-API-Consumer", "DBIO-GRAFANA-PLUGIN")
	response, err := ds.client.Do(httpReq)
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Error making request to server"
		return res, nil
	}

	if response.StatusCode != http.StatusOK {
		backend.Logger.Error("Query failed", "url", url)
		res.Status = backend.HealthStatusError
		res.Message = "Invalid response from server"
		return res, nil
	}

	res.Status = backend.HealthStatusOk
	res.Message = "Success"
	return res, nil
}

// QueryData queries for data.
func (ds *DepatureBoardIODataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	settings, err := LoadSettings(req.PluginContext)
	if err != nil {
		return nil, err
	}

	for _, q := range req.Queries {
		dr := backend.DataResponse{}
		model := &DeparturesByCRSQueryModel{}
		err = json.Unmarshal(q.JSON, model)
		if err != nil {
			backend.Logger.Error("Query failed", "model", spew.Sdump(model))
		}
		url := settings.APIEndpoint + "/getDeparturesByCRS/" + model.StationCRS + "?apikey=" + settings.APIKey
		httpReq, err := http.NewRequest("GET", url, nil)
		httpReq.Header.Add("X-API-Consumer", "DBIO-GRAFANA-PLUGIN")
		response, err := ds.client.Do(httpReq)
		if err != nil {
			backend.Logger.Error("Query failed", "url", url)
		}

		if response.StatusCode != http.StatusOK {
			backend.Logger.Error("Query failed", "url", url)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			backend.Logger.Error("Query failed", "url", url)
		}

		board := Board{}
		json.Unmarshal(body, &board)
		var destinations, platforms, std, etd []string
		for _, service := range board.TrainServices {
			std = append(std, service.STD)
			etd = append(etd, service.ETD)
			destinations = append(destinations, service.Destination[0].LocationName)
			platforms = append(platforms, service.Platform)
		}
		frame := data.NewFrame(model.StationCRS,
			data.NewField("Scheduled", data.Labels{}, std),
			data.NewField("Estimated", data.Labels{}, etd),
			data.NewField("Destination", data.Labels{}, destinations),
			data.NewField("Platform", data.Labels{}, platforms),
		)
		dr.Frames = append(dr.Frames, frame)

		res.Responses[q.RefID] = dr
	}

	return res, nil
}
