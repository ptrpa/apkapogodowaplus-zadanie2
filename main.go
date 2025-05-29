package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

var token = os.Getenv("WEATHER_API_KEY")
var tmpl = template.Must(template.ParseFiles("index.html"))

func main() {
	// tryb healthcheck
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	log.Printf("Aplikacja uruchomiona %s\n", os.Getenv("TZ"))
	log.Println("Autor: Piotr TwojeNazwisko")
	log.Println("Port: 8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, map[string]string{
			"Miasto": "",
			"Wynik":  "",
		})
	})

	http.HandleFunc("/pogoda", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Błąd formularza", http.StatusBadRequest)
			return
		}

		miasto := r.FormValue("miasto")
		akcja := r.FormValue("akcja")

		dane := map[string]string{
			"Miasto": miasto,
			"Wynik":  "",
		}

		if akcja == "odśwież" || miasto == "" {
			// reset do pustej opcji
			dane["Miasto"] = ""
			tmpl.Execute(w, dane)
			return
		}

		pogodatekst, err := pobierzPogode(miasto)
		if err != nil {
			dane["Wynik"] = "Błąd pobierania pogody."
		} else {
			dane["Wynik"] = pogodatekst
		}

		tmpl.Execute(w, dane)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pobierzPogode(miasto string) (string, error) {
	miasto = usunPolskieZnaki(miasto)

	apiUrl := fmt.Sprintf(
		"https://dobrapogoda24.pl/api/v1/weather/simple?city=%s&day=1&token=%s",
		url.QueryEscape(miasto), token)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var dane struct {
		Day struct {
			TempMax       float64 `json:"temp_max"`
			TempMin       float64 `json:"temp_min"`
			WindVelocity  float64 `json:"wind_velocity"`
			Precipitation string  `json:"precipitation"`
			Humidity      int     `json:"humidity"`
			Pressure      int     `json:"pressure"`
		} `json:"day"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&dane); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Temperatura: %.1f°C do %.1f°C\nWiatr: %.0f km/h\nOpady: %s mm\nWilgotność: %d%%\nCiśnienie: %d hPa",
		dane.Day.TempMin, dane.Day.TempMax, dane.Day.WindVelocity, dane.Day.Precipitation, dane.Day.Humidity, dane.Day.Pressure), nil
}

func usunPolskieZnaki(s string) string {
	zamiany := map[rune]rune{
		'ą': 'a', 'ć': 'c', 'ę': 'e', 'ł': 'l',
		'ń': 'n', 'ó': 'o', 'ś': 's', 'ź': 'z', 'ż': 'z',
		'Ą': 'A', 'Ć': 'C', 'Ę': 'E', 'Ł': 'L',
		'Ń': 'N', 'Ó': 'O', 'Ś': 'S', 'Ź': 'Z', 'Ż': 'Z',
	}

	out := []rune{}
	for _, r := range s {
		if zamieniony, ok := zamiany[r]; ok {
			out = append(out, zamieniony)
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}

