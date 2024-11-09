package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type POIS []POI

type Tags []string

type POI struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image"`
	Tags        Tags   `json:"tags"`
}

type Result struct {
	Items ItemList `xml:"items"`
}

type ItemList struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID    string          `xml:"id"`
	Name  string          `xml:"title"`
	Texts TextList        `xml:"texts"`
	Media MediaObjectList `xml:"media_objects"`
}

type TextList struct {
	Texts []Text `xml:"text"`
}

type Text struct {
	Text string `xml:",chardata"`
	Type string `xml:"rel,attr"`
}

type MediaObjectList struct {
	MediaObjects []MediaObject `xml:"media_object"`
}

type MediaObject struct {
	Url         string `xml:"url,attr"`
	Description string `xml:",chardata"`
	Type        string `xml:"rel,attr"`
}

func main() {
	http.HandleFunc("/api/moods", moods)
	http.HandleFunc("/api/pois", pois)
	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

func pois(w http.ResponseWriter, req *http.Request) {
	fn := fmt.Sprintf("%s/pois.xml", os.Getenv("KO_DATA_PATH"))
	r, err := os.Open(fn)

	if err != nil {
		log.Printf("Could not open file: %s", fn)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	doc := buf.Bytes()

	res := Result{}

	xml.Unmarshal(doc, &res)

	p := make(POIS, 0)

	for _, item := range res.Items.Items {
		poi := POI{}
		poi.Name = item.Name
		for _, t := range item.Texts.Texts {
			if t.Type == "details" {
				poi.Description = t.Text
			}
		}
		poi.ImageUrl = getImageForUid(item.ID, item.Media.MediaObjects)
		poi.Tags = getTagsForUid(item.ID)
		p = append(p, poi)
	}

	enc := json.NewEncoder(w)
	enc.Encode(p)
}

func getImageForUid(uid string, mediaObjects []MediaObject) string {
	imageMapping := map[string]int{
		"91a2dcec-9914-4714-8c4c-bc133aa198a9": 1,
		"1202ae41-3413-4225-870c-32586270a52a": 2,
		"aa68c862-00ac-4c25-a59b-00fd4adc8e72": 1,
		"d8d971ed-f8ca-4115-959e-71a03604e68f": 1,
		"dab14386-5116-47ea-83ee-841d6d17ed4b": 2,
	}

	return mediaObjects[imageMapping[uid]].Url
}

func getTagsForUid(uid string) []string {
	tagMapping := map[string]Tags{
		"91a2dcec-9914-4714-8c4c-bc133aa198a9": {"culture", "museum"},
		"1202ae41-3413-4225-870c-32586270a52a": {"nature"},
		"aa68c862-00ac-4c25-a59b-00fd4adc8e72": {"spiritual", "culture"},
		"d8d971ed-f8ca-4115-959e-71a03604e68f": {"sports"},
		"dab14386-5116-47ea-83ee-841d6d17ed4b": {"culture", "music"},
	}

	return tagMapping[uid]
}

func moods(w http.ResponseWriter, req *http.Request) {
	fn := fmt.Sprintf("%s/moods.json", os.Getenv("KO_DATA_PATH"))
	r, err := os.Open(fn)

	if err != nil {
		log.Printf("Could not open file: %s", fn)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(w, r)
}
