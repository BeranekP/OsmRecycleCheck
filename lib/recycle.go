package recycle

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Container struct {
	Type  string  `json:"type"`
	Id    int      `json:"id"` 
	Nodes []int   `json:"nodes, omitempty"`
	Lat   float32 `json:"lat, omitempty"`
	Lon   float32 `json:"lon, omitempty"`
	Tags  map[string]string `json:"tags"`
}

type TimeStamps struct {
	TimestampOsmBase   string
	TimestampAreasBase string
	Copyright          string
}

type ResponseData struct {
	Version   float32
	Generator string
	Osm3s     TimeStamps
	Elements  []Container `json:"elements"`
}

func Recycle() {
	client := http.Client{}
	timeout := 60
    geocodes := getGeocodes()
    id := geocodes.CZ
	query := fmt.Sprintf(`[out:json][timeout:%d];
                area(id:%d)->.searchArea;
                nwr["recycling_type" = "container"](area.searchArea);
                out;`, timeout, id)

	form := url.Values{}
	form.Add("data", query)
	body := form.Encode()

	req, err := http.NewRequest("POST", "https://overpass-api.de/api/interpreter/", strings.NewReader(body))

	if err != nil {
		log.Fatal("Error creating request", err)
	}
    fmt.Println("Making request to overpass-api")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error making request", err)
	}

	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error in response", err)
	}

	//fmt.Println(string(payload))
	var containers ResponseData
	e := json.Unmarshal(payload, &containers)
	if e != nil {
		log.Fatal("Error parsing data:", e)
	}
	//fmt.Println(containers)



	file, _ := json.MarshalIndent(containers, "", " ")
	os.WriteFile("containers.json", file, 0644)

}
