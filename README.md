# go-json-rest-prometheus

Prometheus middleware for [go-json-rest](https://github.com/ant0ine/go-json-rest)

This is heavily influenced by [negroni-prometheus](https://github.com/zbindenren/negroni-prometheus). I wanted to have the same functionality on go-json-rest.

###A note on the Handler. 
go-json-rest rest.HandlerFunc it not compatible with http.HandlerFunc because the ResponseWriter is different (see this issue [https://github.com/ant0ine/go-json-rest/issues/192](https://github.com/ant0ine/go-json-rest/issues/192) so the Handler definition on the endpoint is awkward.

### Example
This example can be found in the [example](https://github.com/Telmo/go-json-rest-prometheus/tree/master/example) directory.

```go
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
	api.Use(&restprometheus.PromMiddleware{})
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/metrics", func(w rest.ResponseWriter, r *rest.Request) {
			prometheus.InstrumentHandler("prometheus", prometheus.Handler())(w.(http.ResponseWriter), r.Request)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
```

### TODO

- Add custom bucket definitions
- Learn more go
