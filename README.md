# go-flickr-pip

Tools for reverse-geocoding geotagged Flickr photos using the Flickr API and the go-whosonfirst-spatial package.

## Example

### Step 1

Download [SFO Museum architecture data](https://github.com/sfomuseum-data/sfomuseum-data-architecture).

```
$> git clone git@github.com:sfomuseum-data/sfomuseum-data-architecture.git /path/to/sfomuseum-data-architecture
```

### Step 2

Create a spatially-enabled SQLite database for the SFO Museum architecture data using the tools in the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package.

```
$> cd go-whosonfirst-sqlite-features-index

$> make cli
go build -mod vendor -o bin/wof-sqlite-index-features cmd/wof-sqlite-index-features/main.go

$> ./bin/wof-sqlite-index-features -all -dsn /path/to/sfomuseum-architecture.db /path/to/sfomuseum-data-architecture
```

### Step 3

Expose the SQLite database with WOF data for Denmark as a simple HTTP service using the tools in the [go-whosonfirst-spatial-www-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite) package.

```
$> cd go-whosonfirst-sqlite-spatial-www-sqlite
go build -mod vendor -o bin/server cmd/server/main.go

$> make cli

$> ./bin/server -spatial-database-uri 'sqlite://?dsn=/path/to/sfomuseum-architecture.db'
```

### Step 4

In a different terminal window run the `pip` tool to fetch openly-licensed and geotagged photos from the [airports-sfo](https://www.flickr.com/groups/airports-sfo/pool/) Flickr group and to reverse-geocode them (using the `server` tool you've just started in the other terminal):

```
$> go run -mod vendor cmd/pip/main.go \
	-client-uri 'oauth1://?consumer_key=...&consumer_secret=...' \
	-param method=flickr.photos.search \
	-param group_id=95693046@N00 \
	-param has_geo=1
	-param extras=geo \
	-param license=1,2,3,4,5,6,7,8,9,10

photo_id,latitude,longitude,whosonfirst_id,whosonfirst_name,whosonfirst_placetype
51057165676,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50954440458,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50758638548,37.623545,-122.389712,1730008851,Taxiway Q,custom
50731164373,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50731892551,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50713143507,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50697680261,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50697680261,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50683567541,37.622049,-122.383017,1730008749,RUNWAY 10R/28L,custom
50683567541,37.622049,-122.383017,1730008749,RUNWAY 10R/28L,custom
50654967468,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50654967468,37.616339,-122.387223,1360665043,Central Parking Garage,wing
50654966103,37.616339,-122.387223,1360665043,Central Parking Garage,wing
... and so on
```	

The `pip` tool is designed to work with any Flickr API method that returns a [standard photo response](https://code.flickr.net/2008/08/19/standard-photos-response-apis-for-civilized-age/).

## See also

* https://github.com/aaronland/go-flickr-api
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite