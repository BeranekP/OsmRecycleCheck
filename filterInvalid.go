package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BeranekP/OsmRecycleCheck/lib"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-U" {
			recycle.Recycle()
		}
	}

	data, err := os.ReadFile("containers.json")

	if err != nil {
		log.Fatal(err)
	}

	var containers recycle.ResponseData
	json_err := json.Unmarshal(data, &containers)

	if json_err != nil {
		log.Fatal(json_err)
	}

	var missingAmenity []recycle.Container
	var missingRecycling []recycle.Container
	var missingType []recycle.Container
	var invalidType []recycle.Container

	for _, container := range containers.Elements {

		if container.Tags["recycling_type"] == "" {
			missingType = append(missingType, container)
		}
		if container.Tags["amenity"] == "" {
			missingAmenity = append(missingAmenity, container)

		}

		if (container.Tags["recycling_type"] != "container") && (container.Tags["recycling_type"] != "centre") {
			invalidType = append(invalidType, container)
		}

		checkSubstance := 0
		for key, value := range container.Tags {
			if strings.Contains(key, "recycling:") && value == "yes" {
				checkSubstance += 1
			}
		}

		if checkSubstance == 0 {

			missingRecycling = append(missingRecycling, container)
		}

	}

	fmt.Println("missingRecycling: ", len(missingRecycling))
	fmt.Println("missingAmenity: ", len(missingAmenity))
	fmt.Println("missingType: ", len(missingType))
	fmt.Println("invalidType: ", len(invalidType))

	output := OutputData{Elements: missingRecycling}
	file, _ := json.MarshalIndent(output, "", " ")
	os.WriteFile("missingRecycling.json", file, 0644)

	if len(missingRecycling) < 50 {
		printArray(missingRecycling, "Missing recycling:*=yes")

	}

	if len(missingAmenity) > 0 {
		printArray(missingAmenity, "Missing amenity:recycling")

		output := OutputData{Elements: missingAmenity}
		file, _ := json.MarshalIndent(output, "", " ")
		os.WriteFile("missingAmenity.json", file, 0644)

	}

	if len(missingType) > 0 {
		printArray(missingType, "Missing recycling_type = container/centre")
	}
	if len(invalidType) > 0 {
		printArray(invalidType, "Recycling type not one of container/centre")
	}

}

type OutputData struct {
	Elements []recycle.Container `json:"elements"`
}

func printArray(array []recycle.Container, title string) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println("-----------------------------")
	for i, item := range array {
		url := fmt.Sprintf("http://openstreetmap.org/%s/%d", item.Type, item.Id)
		fmt.Println(i, ": ", item, url)
	}
	fmt.Println()
}
