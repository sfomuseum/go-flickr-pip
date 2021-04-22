package pip

import (
	"bytes"
	"context"
	"net/http"
	"encoding/json"
)

type SPRResults struct {
	Places []*SPRResult `json:"places"`
}

type SPRResult struct {
	Id string `json:"wof:id"`
	ParentId string `json:"wof:parent_id"`
	Name string `json:"wof:name"`
	Placetype string `json:"wof:placetype"`
}

type PointInPolygonRequest struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsCurrent []int64 `json:"is_current,omitempty"`
}

type Client struct {
	http_client *http.Client
}

func NewClient() (*Client, error) {

	http_client := &http.Client{}

	cl := &Client{
		http_client: http_client,
	}

	return cl, nil
}

func (cl *Client) Query(ctx context.Context, lat float64, lon float64) (*SPRResults, error) {

	req := PointInPolygonRequest{
		Latitude: lat,
		Longitude: lon,
		IsCurrent: []int64{ 1 },
	}

	body, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)
	
	http_req, err := http.NewRequest("POST", "http://localhost:8080/api/point-in-polygon", br)

	if err != nil {
		return nil, err
	}

	rsp, err := cl.http_client.Do(http_req)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	var spr *SPRResults

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&spr)

	if err != nil {
		return nil, err
	}

	return spr, nil
}

