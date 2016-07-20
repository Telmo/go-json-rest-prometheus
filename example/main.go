package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/telmo/go-json-rest-prometheus"
)

func main() {
	api := rest.NewApi()
	//	statusMw := statsd.StatsdMiddleware{}
	api.Use(&restprometheus.PromMiddleware{})
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/metrics", rest.HandlerFunc(prometheus.InstrumentHandler("prometheus", prometheus.Handler()))),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
