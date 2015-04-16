package fhirterm

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	page := `
<!DOCTYPE html>
<html>
<head>
<title>Welcome to fhirterm!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to fhirterm!</h1>
<p>If you see this page, the fhirterm server is successfully installed and
working. Further configuration is required. Normally, this page should not be publicly accessible.</p>

<p><em>Thank you for using fhirterm.</em></p>
</body>
</html>`

	fmt.Fprint(w, page)
}

func ValueSetExpand(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := GetDb()
	row := db.QueryRow("SELECT long_common_name FROM loinc_loincs LIMIT 1")
	var display string
	row.Scan(&display)

	vs, _ := ExpandValueSet(ps.ByName("id"))

	fmt.Fprintf(w, "hello, %s!\n%s", ps.ByName("id"), vs.Identifier)
	fmt.Fprintf(w, "%s", display)
}

type HttpLogger struct {
	i int
}

func (l *HttpLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	log.Printf("%s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	log.Printf("Completed with %v %s in %v\n\n", res.Status(), http.StatusText(res.Status()), time.Since(start))
}

func setupCors(cfg *Config) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: cfg.HttpCorsAllowedOrigins,
	})
}

func StartServer(cfg *Config) {
	addr := fmt.Sprintf("%s:%d", cfg.HttpHost, cfg.HttpPort)
	router := httprouter.New()

	router.GET("/", Index)
	router.GET("/ValueSet/:id/$expand", ValueSetExpand)

	n := negroni.New()
	corsMw := setupCors(cfg)

	n.Use(&HttpLogger{42})
	n.Use(corsMw)
	n.UseHandler(router)

	log.Printf("Starting FHIRterm server on %s", addr)
	http.ListenAndServe(addr, n)
}
