package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Teddy-Schmitz/particleExporter/particle"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Metrics = map[string]*prometheus.GaugeVec{
	"temperature": prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "particle_event_temperature",
		Help: "Temperature events from Particle cloud",
	}, []string{"device_name"}),
	"humidity": prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "particle_event_humidity",
		Help: "Humidity events from Particle cloud",
	}, []string{"device_name"}),
	"sound": prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "particle_event_sound",
		Help: "Sound sensor events from Particle cloud",
	}, []string{"device_name"}),
	"dust": prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "particle_event_dust",
		Help: "Dust sensor events from Particle cloud",
	}, []string{"device_name"}),
}

var nameCache = map[string]string{}

var ParticleClient = particle.Client{
	AccessToken: os.Getenv("PARTICLE_ACCESS_TOKEN"),
}

func main() {

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	s := r.PathPrefix("/event").Subrouter()
	s.Use(authKeyMiddleware)
	s.HandleFunc("", receiveEvents).Methods("POST")

	log.Fatalln(http.ListenAndServe(":4987", r), "unable to start http server")
}

func authKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != os.Getenv("API_KEY") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func receiveEvents(w http.ResponseWriter, r *http.Request) {

	event := &particle.Event{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		log.Println(err)
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	devName, ok := nameCache[event.DeviceID]

	if !ok {
		devName = event.DeviceID
		if os.Getenv("PARTICLE_ACCESS_TOKEN") != "" {
			res, err := ParticleClient.GetDeviceInfo(event.DeviceID)
			if err != nil {
				log.Println(err)
			} else {
				nameCache[event.DeviceID] = res.Name
				devName = res.Name
			}
		}
	}

	m, ok := Metrics[event.Name]
	if !ok {
		log.Println(event)
		http.Error(w, "No metric", http.StatusInternalServerError)
		return
	}

	val, err := strconv.ParseFloat(event.Data, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad data", http.StatusBadRequest)
		return
	}

	m.WithLabelValues(devName).Set(val)
}

func init() {
	for _, m := range Metrics {
		prometheus.MustRegister(m)
	}
}
