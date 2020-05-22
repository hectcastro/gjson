package gjson

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-github/github"
	geojson "github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// GeoJSON ...
type GeoJSON struct {
	JSON   interface{}
	Length int
}

// BigGeoJSONLength is an arbitrary length threshold for a big
// blob of GeoJSON.
const BigGeoJSONLength = 5000000

// MaxURLLength is an arbitrary length threshold for the maximum
// length of a URL.
const MaxURLLength = 150000

// Unmarshal attempts to unmarshal GeoJSON FeatureCollections and
// standalone Geometries.
func (g *GeoJSON) Unmarshal(file *os.File) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	g.Length = len(data)

	featureCollection := &geojson.FeatureCollection{}
	unmarshalErr := json.Unmarshal(data, featureCollection)
	if unmarshalErr == nil && featureCollection.Type == "FeatureCollection" {
		g.JSON = featureCollection
		return nil
	}

	geometry := &geojson.Geometry{}
	unmarshalErr = json.Unmarshal(data, geometry)
	if unmarshalErr == nil && geometry.Type == geojson.GeometryType(geometry.Type) {
		g.JSON = geometry
		return nil
	}

	if unmarshalErr != nil {
		return unmarshalErr
	}

	return errors.New("Unable to unmarshal GeoJSON")
}

// Size returns the number of bytes of the raw, marshalled
// GeoJSON.
func (g *GeoJSON) Size() int {
	return g.Length
}

// ToURLEncoded URL encodes the GeoJSON data structure.
func (g *GeoJSON) ToURLEncoded() string {
	data, err := json.Marshal(g.JSON)
	if err != nil {
		log.Fatal(err)
	}

	return url.QueryEscape(string(data))
}

// ToGist persists the GeoJSON data structure to a public
// GitHub Gist.
func (g *GeoJSON) ToGist() string {
	data, err := json.MarshalIndent(g.JSON, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	client := github.NewClient(tokenHTTPClient())
	gist, _, err := client.Gists.Create(context.Background(), &github.Gist{
		Description: github.String(""),
		Public:      github.Bool(true),
		Files: map[github.GistFilename]github.GistFile{
			"map.geojson": {
				Content: github.String(string(data)),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return gist.GetID()
}

func tokenHTTPClient() *http.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: os.Getenv("GITHUB_TOKEN"),
	})
	tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	return tokenClient
}
