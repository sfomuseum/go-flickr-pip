package main

import (
	"context"
	"flag"
	"log"
	"github.com/sfomuseum/go-flickr-pip"
	"github.com/aaronland/go-flickr-api/client"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"io"
	"net/url"
	_ "os"
	_ "fmt"
)

func main(){

	client_uri := flag.String("client-uri", "", "")

	var params multi.KeyValueString	
	flag.Var(&params, "param", "Zero or more {KEY}={VALUE} Flickr API parameters to include with your uploads.")
	
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

	cb := func(ctx context.Context, fh io.ReadSeekCloser, err error) error {

		if err != nil {
			return err
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return err
		}
		
		photos_rsp := gjson.GetBytes(body, "photos.photo")

		for _, ph := range photos_rsp.Array(){
			
			lat_rsp := ph.Get("latitude")
			lon_rsp := ph.Get("longitude")

			if !lat_rsp.Exists() || !lon_rsp.Exists(){
				continue
			}

			id_rsp := ph.Get("id")
			ph_id := id_rsp.Int()
			
			lat := lat_rsp.Float()
			lon := lon_rsp.Float()

			rsp, err := pip_cl.Query(ctx, lat, lon)

			if err != nil {
				log.Println(lat, lon, err)
				continue
			}

			for _, pl := range rsp.Places {
				log.Println(ph_id, lat, lon, pl.Id, pl.Placetype, pl.Name)
			}
		}
		
		return nil
	}

	args := &url.Values{}

	for _, kv := range params {
		args.Set(kv.Key(), kv.Value().(string))
	}
	
	err := client.ExecuteMethodPaginatedWithClient(ctx, api_cl, args, cb)
	
	if err != nil {
		log.Fatalf("Failed to write method results, %v", err)
	}
	
}

/*

> go run -mod vendor cmd/pip/main.go -client-uri 'oauth1://?consumer_key=&consumer_secret=' -param method=flickr.photos.search -param user_id=161215698@N03 -param has_geo=1 -param extras=geo
2021/04/21 18:25:38 51130478394 37.615555 -122.388889 1729792579 wing International Terminal
2021/04/21 18:25:38 51130478394 37.615555 -122.388889 1729792387 building SFO Terminal Complex
2021/04/21 18:25:38 51130478394 37.615555 -122.388889 1729792387 building SFO Terminal Complex
2021/04/21 18:25:38 51130478394 37.615555 -122.388889 1729792679 concourse International Terminal Main Hall
2021/04/21 18:25:38 51130478394 37.615555 -122.388889 1729792681 concourse International Terminal Connector
2021/04/21 18:25:38 51131288670 37.612777 -122.361112 1730008749 custom RUNWAY 10R/28L

*/
