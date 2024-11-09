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

func getDescForUid(uid string) string {
	descMapping := map[string]string{
		"91a2dcec-9914-4714-8c4c-bc133aa198a9": "From 1773 to 1787, the Mozart family lived in the ‘Tanzmeisterhaus’ on today's Makartplatz. The flat on the first floor now houses a museum. Mozartplatz is located in the centre of Salzburg's old town and is surrounded by Residenzplatz and Waagplatz. The Mozart statue forms the centrepiece of the square.",
		"1202ae41-3413-4225-870c-32586270a52a": "Museum of Nature and Technology. A visit to the Haus der Natur is as diverse as life itself. From prehistoric dinosaurs to the great moments of space travel, from colourful underwater worlds to rare reptiles from all over the world, from legendary crystal treasures to the inner workings of our own bodies. Another highlight: experimenting in Austria's most versatile science centre.",
		"aa68c862-00ac-4c25-a59b-00fd4adc8e72": "A special feature of the Kajetanerkirche is the Holy Staircase, which was built in 1712 as an imitation of the Scala Santa in Rome. Like its model in Rome, the staircase in Salzburg has a cross on the 2nd, 11th and 28th steps to symbolise the drops of Jesus Christ's blood. History: In 1591, Archbishop Wolf Dietrich acquired a hospital and church to establish a seminary, which was to be under the direction of the Theatine Order (Cajetans). It was only under Archbishop Johann Ernst von Thun that the church was completed in 1696 and consecrated in 1700. The interior of the church looks very festive, elegant and clear due to the stucco. The light-giving dome dominates the room. Built into the gallery balustrade above the vestibule is the oldest preserved organ in Salzburg, which was built around 1700 by Christoph Egedacher.",
		"d8d971ed-f8ca-4115-959e-71a03604e68f": "Dorfloipe Eugendorf: cross-country skiing in classic style, three classic trails, plus wonderful views over fields glistening in the snow, across romantic forests to the Gaisberg and even as far as the Dachstein massif. All three trails in Eugendorf are groomed as soon as the snow and weather conditions allow. Use of the cross-country ski trails in Eugendorf is free of charge! The 4.9-kilometre Eugendorf village trail leads from the public playground in the centre of Eugendorf towards Eugendorf train station, along the forest back via Einleiten and slightly above the Ischlerbahnstraße back to the starting point. This trail is an easy classic trail, which is also suitable for beginners thanks to its length and only 47 metres in altitude. Depending on your cross-country skiing ability and speed, you can complete the trail in just over an hour.",
		"dab14386-5116-47ea-83ee-841d6d17ed4b": "The Mozarteum Orchestra is Salzburg's symphony, festival and opera orchestra. It was founded with the support of Mozart's widow Constanze and their sons Franz Xaver and Karl Thomas. Its origins go back to the Cathedral Music Society and Mozarteum, founded in 1841. Today, the Mozarteum Orchestra is one of Austria's leading symphony orchestras. Viennese classical music in a new light: The Mozarteum Orchestra's great speciality is a new, exciting interpretation of Viennese classical music. Of course, Salzburg's ‘local hero’, Wolfgang Amadeus Mozart, plays a special role in the repertoire. However, the orchestra also performs Romantic and contemporary works.",
	}

	return descMapping[uid]
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
