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

	"github.com/joho/godotenv"
)

func main() {
	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	apiKey := os.Getenv("API_KEY")

	records := readCSV("russia.csv")

	for idx, record := range records[1:2] {
		os.MkdirAll(fmt.Sprint(idx), 0755)
		generateImages(record[0], record[1], idx, apiKey)
		println(fmt.Sprint(idx), "done")
	}
}

func createSlice(start, end, step float64) []float64 {
	if step <= 0 || end < start {
		return []float64{}
	}
	s := make([]float64, 0, int(1+(end-start)/step))
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

func generateImages(lat string, long string, idx int, apiKey string) {
	latF64, _ := strconv.ParseFloat(lat, 64)
	longF64, _ := strconv.ParseFloat(long, 64)

	var lats = createSlice(latF64-0.01, latF64+0.01, 0.005)
	var longs = createSlice(longF64-0.01, longF64+0.01, 0.005)

	var tuples [][]float64

	for _, i := range lats {
		for _, j := range longs {
			tuples = append(tuples, []float64{i, j})
		}
	}

	for _, tuple := range tuples {
		resp, err := http.Get(fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%f,%f&zoom=17&size=1280x1280&scale=2&maptype=satellite&key=%s", tuple[0], tuple[1], apiKey))
		if err != nil {
			log.Fatalln(err)
		}
		img, _, err := image.Decode(resp.Body)

		if err != nil {
			log.Fatalln(err)
		}

		f, err := os.Create(fmt.Sprint(idx, "/", time.Now().UnixNano(), ".png"))
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
