package gjson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func Test_Size(t *testing.T) {
	file, _ := os.Open("../../testdata/featureCollection.geojson")
	geojson := &GeoJSON{}
	geojson.Unmarshal(file)
	defer file.Close()

	if geojson.Size() != 896 {
		t.Errorf("Invalid size: expected %d, but got %d", 896, geojson.Size())
	}
}

func Test_ToURLEncoded(t *testing.T) {
	file, _ := os.Open("../../testdata/featureCollection.geojson")
	geojson := &GeoJSON{}
	geojson.Unmarshal(file)
	defer file.Close()

	data, _ := json.Marshal(geojson.JSON)

	if geojson.ToURLEncoded() != url.QueryEscape(string(data)) {
		t.Errorf("Invalid encoded URL: expected %v, but got %v",
			geojson.ToURLEncoded(),
			url.QueryEscape(string(data)),
		)
	}
}

func Test_Unmarshal(t *testing.T) {
	files, _ := filepath.Glob("../../testdata/*")

	for _, file := range files {
		f, _ := os.Open(file)
		geojson := &GeoJSON{}
		geojson.Unmarshal(f)
		f.Close()

		rawJSON, _ := ioutil.ReadFile(file)
		marshalledJSON, _ := json.Marshal(geojson.JSON)

		compactedBuffer := new(bytes.Buffer)
		json.Compact(compactedBuffer, rawJSON)

		if !bytes.Equal(compactedBuffer.Bytes(), marshalledJSON) {
			t.Errorf("Incorrect unmarshal: expected %v, but got %v",
				string(compactedBuffer.Bytes()),
				string(marshalledJSON),
			)
		}
	}
}
