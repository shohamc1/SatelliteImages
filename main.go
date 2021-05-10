package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	records := readCSV("russia.csv")

	for idx, record := range records[1:] {
		os.MkdirAll(fmt.Sprint(idx), 0755)
		generateImages(record[0], record[1], idx)
		println(fmt.Sprint(idx), "done")
	}

}

func generateImages(lat string, long string, idx int) {
	latF64, _ := strconv.ParseFloat(lat, 64)
	longF64, _ := strconv.ParseFloat(long, 64)

	lats := [3]float64{latF64 - 0.005, latF64, latF64 + 0.005}
	longs := [3]float64{longF64 - 0.005, longF64, longF64 + 0.005}

	tuples := [9][2]float64{{lats[0], longs[0]}, {lats[1], longs[0]}, {lats[2], longs[0]}, {lats[0], longs[1]}, {lats[1], longs[1]}, {lats[2], longs[1]}, {lats[0], longs[2]}, {lats[1], longs[2]}, {lats[2], longs[2]}}

	for _, tuple := range tuples {
		resp, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%f,%f&zoom=17&size=1280x1280&scale=2&maptype=satellite&key=", tuple[0], tuple[1]))
		if err != nil {
			log.Fatalln(err)
		}
		img, _, err := image.Decode(resp.Body)

		if err != nil {
			log.Fatalln(err)
		}

		f, err := os.Create(fmt.Sprint(idx, "\\", time.Now().UnixNano(), ".png"))
		if err != nil {
			log.Fatalln(err)
		}

		defer f.Close()

		png.Encode(f, img)
	}
}

func readCSV(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalln(err)
	}

	return records
}
