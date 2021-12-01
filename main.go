package main

import (
	"app/infrastructure/config"
	"app/infrastructure/debug"
	"app/infrastructure/log"
	"app/service/dict"
	"app/service/home"
	"app/service/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

type HttpMethod string

const (
	GET  HttpMethod = "GET"  // read
	POST HttpMethod = "POST" // update
	PUT  HttpMethod = "PUT"  // create
)

type api struct {
	method  HttpMethod
	path    string
	handler gin.HandlerFunc
}

var apis = []api{
	{GET, "/ping", home.HandlePagePing},
	{GET, "/", home.HandlePageIndex},

	// user page
	{GET, "/login", user.HandlePageLogin},

	// user api
	{POST, "/api/login", user.HandleAPILogin},

	// dict page
	{GET, "/dicts", dict.HandlePageQueryDict}, // ?tags=xxx
	{GET, "/dicts/:id", dict.HandlePageDict},
	{GET, "/dicts/:id/:categoryLinkId", dict.HandlePageDict},

	{GET, "/words", dict.HandlePageQueryWord},         // word page ?dictId=xxx&catalogLinkId=xxx
	{GET, "/words/:writing/:id", dict.HandlePageWord}, // word page
	{GET, "/words/:writing", dict.HandlePageWord},     // word page

	{GET, "/editor/dict", dict.HandlePageEditDict}, // ?dictId=
	{GET, "/editor/word", dict.HandlePageEditWord}, // ?dictId=&wordId=

	// dict api
	{GET, "/api/dicts/:id", dict.HandleAPIGetDict},
	{POST, "/api/dicts", dict.HandleAPICreateDict},
	{PUT, "/api/dicts/:id", dict.HandleAPIUpdateDict},
	{GET, "/api/words", dict.HandleAPIGetWord},
	{GET, "/api/words/:id", dict.HandleAPIGetWord},
	{POST, "/api/words", dict.HandleAPICreateWord},
	{PUT, "/api/words/:id", dict.HandleAPIUpdateWord},

	// search
	{GET, "/search/dicts", dict.HandlePageSearchDicts},
	{GET, "/search/words", dict.HandlePageSearchWords}, // ?kw=
}

func main() {
	debug.PprofRouter()
	router := gin.Default()
	for _, item := range apis {
		router.Handle(string(item.method), item.path, item.handler)
	}
	router.Static("css", "./static/css")
	router.Static("js", "./static/js")
	router.Static("img", "./static/img")
	// router.Static("font", "./static/font")
	router.StaticFile("favicon.ico", "./static/img/favicon.png")
	router.StaticFile("robots.txt", "./static/robots.txt")
	router.StaticFile("sitemap.xml", "./static/sitemap.xml")

	server := &http.Server{
		Addr:           ":80",
		Handler:        router,
		IdleTimeout:    3 * time.Minute,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if config.Get().Domain != "localhost" {
		secureMiddleware := secure.New(secure.Options{SSLRedirect: true})
		secureRouter := secureMiddleware.Handler(router)
		server.Handler = secureRouter
		secureServer := &http.Server{
			Addr:           ":443",
			Handler:        router,
			IdleTimeout:    3 * time.Minute,
			ReadTimeout:    20 * time.Second,
			WriteTimeout:   20 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			if err := secureServer.ListenAndServeTLS("server.pem", "server.key"); err != nil {
				log.Error("run server error", log.String("error", err.Error()))
			}
		}()
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("run server error", log.String("error", err.Error()))
	}

}
