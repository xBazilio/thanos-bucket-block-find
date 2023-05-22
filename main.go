package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/oklog/ulid"
	v1 "github.com/thanos-io/thanos/pkg/api/blocks"
	"github.com/thanos-io/thanos/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	mint     = kingpin.Flag("min-time", "Start of time range to search blocks").Default("0000-01-01T00:00:00Z").String()
	maxt     = kingpin.Flag("max-time", "End of time range to search blocks").Default("9999-12-31T23:59:59Z").String()
	endpoint = kingpin.Flag("endpoint", "Endpoint host of thanos bucket web to fetch data from").Required().String()
)

type bucketWebResponse struct {
	Status string
	Data   *v1.BlocksInfo
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	MinTime, err := time.Parse(time.RFC3339, *mint)
	if err != nil {
		panic(err)
	}
	MaxTime, err := time.Parse(time.RFC3339, *maxt)
	if err != nil {
		panic(err)
	}

	blocksData, err := getBlocksData(*endpoint)
	if err != nil {
		panic(err)
	}

	mintimeMilli := MinTime.UnixMilli()
	maxtimeMilli := MaxTime.UnixMilli()
	var ulids []ulid.ULID
	for _, m := range blocksData.Data.Blocks {
		if m.MaxTime < mintimeMilli || m.MinTime > maxtimeMilli {
			continue
		}
		ulids = append(ulids, m.ULID)
	}

	for _, u := range ulids {
		fmt.Println(u.String())
	}
}

func getBlocksData(endpoint string) (*bucketWebResponse, error) {
	resp, err := http.DefaultClient.Get("http://" + endpoint + "/api/v1/blocks")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading body")
	}

	if err := resp.Body.Close(); err != nil {
		return nil, errors.Wrapf(err, "error closing body")
	}

	var respData bucketWebResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling body")
	}

	if respData.Status != "success" {
		return nil, errors.Newf("status is not success, got %s", respData.Status)
	}

	return &respData, nil
}
