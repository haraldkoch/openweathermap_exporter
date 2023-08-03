package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/caarlos0/env"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config stores the parameters used to fetch the data
type Config struct {
	pollingInterval time.Duration
	requestTimeout  time.Duration
	APIKey          string `env:"OWM_API_KEY"`
	Location        string `env:"OWM_LOCATION" envDefault:"Amsterdam,NL"`
	ServerPort      int    `env:"SERVER_PORT" envDefault:"2112"`
}

func loadMetrics(ctx context.Context, location string) <-chan error {
	errC := make(chan error)
	go func() {
		c := time.Tick(cfg.pollingInterval)
		for {
			select {
			case <-ctx.Done():
				return // returning not to leak the goroutine
			case <-c:
				client := &http.Client{
					Timeout: cfg.requestTimeout,
				}

				w, err := owm.NewCurrent("C", "en", cfg.APIKey, owm.WithHttpClient(client))
				if err != nil {
					errC <- err
					continue
				}

				err = w.CurrentByName(location)
				if err != nil {
					errC <- err
					continue
				}

				temp.WithLabelValues(location).Set(w.Main.Temp)

				dewpoint.WithLabelValues(location).Set(w.Main.Temp - (float64(100 - w.Main.Humidity)/5.0))

				feelslike.WithLabelValues(location).Set(w.Main.FeelsLike)

				pressure.WithLabelValues(location).Set(w.Main.Pressure)

				humidity.WithLabelValues(location).Set(float64(w.Main.Humidity))

				wind.WithLabelValues(location).Set(w.Wind.Speed)

				clouds.WithLabelValues(location).Set(float64(w.Clouds.All))

				rain.WithLabelValues(location).Set(w.Rain.ThreeH)
			}
		}
	}()
	return errC
}

var (
	cfg = Config{
		pollingInterval: 5 * time.Second,
		requestTimeout:  1 * time.Second,
	}
	temp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "temperature_celsius",
		Help:      "Temperature in °C",
	}, []string{"location"})

	feelslike = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "feelslike_temperature_celsius",
		Help:      "Apparent Temperature in °C",
	}, []string{"location"})

	pressure = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "pressure_hpa",
		Help:      "Atmospheric pressure in hPa",
	}, []string{"location"})

	humidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "humidity_percent",
		Help:      "Humidity in Percent",
	}, []string{"location"})

	dewpoint = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "dewpoint_temperature_celsius",
		Help:      "Dew Point Temperature in °C",
	}, []string{"location"})
       
	wind = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "wind_mps",
		Help:      "Wind speed in m/s",
	}, []string{"location"})

	clouds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "cloudiness_percent",
		Help:      "Cloudiness in Percent",
	}, []string{"location"})

	rain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "openweathermap",
		Name:      "rain",
		Help:      "Rain contents 3h",
	}, []string{"location"})
)

func main() {
	env.Parse(&cfg)
	if cfg.APIKey == "" {
		log.Fatal("Please provide openWeatherMap API key by setting env var OWM_API_KEY")
	}
	prometheus.Register(temp)
	prometheus.Register(feelslike)
	prometheus.Register(pressure)
	prometheus.Register(humidity)
	prometheus.Register(dewpoint)
	prometheus.Register(wind)
	prometheus.Register(clouds)
	prometheus.Register(rain)

	errC := loadMetrics(context.TODO(), cfg.Location)
	go func() {
		for err := range errC {
			log.Println(err)
		}
	}()

	log.Printf("Starting openWeatherMap exporter on port %d", cfg.ServerPort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), nil)
}
