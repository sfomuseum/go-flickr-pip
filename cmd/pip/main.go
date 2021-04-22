package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/go-flickr-pip"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {

	client_uri := flag.String("client-uri", "", "A valid aaronland/go-flickr-api client URI.")

	var params multi.KeyValueString
	flag.Var(&params, "param", "One or more {KEY}={VALUE} Flickr API parameters.")

	flag.Parse()

	ctx := context.Background()

	pip_cl, err := pip.NewClient()

	if err != nil {
		log.Fatalf("Failed to create new PIP client, %v", err)
	}

	api_cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create new API client, %v", err)
	}

	writers := []io.Writer{
		os.Stdout,
	}

	wr := io.MultiWriter(writers...)
	csv_wr := csv.NewWriter(wr)

	count := 0

	// This is where we handle Flickr results

	cb := func(ctx context.Context, fh io.ReadSeekCloser, err error) error {

		if err != nil {
			return err
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return err
		}

		photos_rsp := gjson.GetBytes(body, "photos.photo")

		// This is where we reverse-geocode each Flickr photo

		for _, ph := range photos_rsp.Array() {

			lat_rsp := ph.Get("latitude")
			lon_rsp := ph.Get("longitude")

			if !lat_rsp.Exists() || !lon_rsp.Exists() {
				continue
			}

			id_rsp := ph.Get("id")
			ph_id := id_rsp.Int()

			lat := lat_rsp.Float()
			lon := lon_rsp.Float()

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// This is where we call the go-whosonfirst-spatial-www-* server

			rsp, err := pip_cl.Query(ctx, lat, lon)

			if err != nil {
				log.Printf("Unable to determine location for photo %d (at %f,%f), %v\n", ph_id, lat, lon, err)
				continue
			}

			for _, pl := range rsp.Places {

				if count == 0 {

					out := []string{
						"photo_id",
						"latitude",
						"longitude",
						"whosonfirst_id",
						"whosonfirst_name",
						"whosonfirst_placetype",
					}

					err := csv_wr.Write(out)

					if err != nil {
						return fmt.Errorf("Failed to write output, %v", err)
					}
				}

				out := []string{
					strconv.FormatInt(ph_id, 10),
					strconv.FormatFloat(lat, 'f', -1, 64),
					strconv.FormatFloat(lon, 'f', -1, 64),
					pl.Id,
					pl.Name,
					pl.Placetype,
				}

				err := csv_wr.Write(out)

				if err != nil {
					return fmt.Errorf("Failed to write output, %v", err)
				}

				count += 1
			}

			csv_wr.Flush()
		}

		return nil
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}

	err = client.ExecuteMethodPaginatedWithClient(ctx, api_cl, args, cb)

	if err != nil {
		log.Fatalf("Failed to write method results, %v", err)
	}

	err = csv_wr.Error()

	if err != nil {
		log.Fatalf("Failed to write results, %v", err)
	}

}
