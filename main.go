package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func addRestHeader(w http.ResponseWriter) {
	w.Header().Set("Server", "Mux Web Server")
	w.Header().Set("DevAdmin", "info@webserver.com")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Robots-Tag", "noindex")
}

var ok bool

func main() {
	internalConfig, ok = InitializeConfiguration()
	if !ok {
		log.Println("InitializeConfiguration")
	}

	Debug("Connecting DB")
	if !ConnectDB() {
		Error("No DB Connected")
		//os.Exit(-1)
	} else {
		Config = make(map[string]string)
		if !LoadConfiguration() {
			os.Exit(-1)
		}
		Debug("DB Connected")
		defer DisconnectDB()
		params := make([]interface{}, 0)
		Row, ok := GetSingleRow("SELECT NOW() AS 'Time'", params, "default")
		if ok {
			Info("DB Time " + Row["Time"])
		}
	}

	routes := UpdateRoute()
	time.Sleep(1 * time.Second)

	Log("******************************************")
	Log("*              Ready to Serve            *")
	Log("******************************************")
	StartServer(routes)

}

func UpdateRoute() *mux.Router {
	Info("HTTP server is Running ")
	Router := mux.NewRouter().StrictSlash(true)
	Router.StrictSlash(true)
	Router.Handle(internalConfig.AssetsFileAbsolute, http.StripPrefix(internalConfig.AssetsFileAbsolute, http.FileServer(http.Dir(internalConfig.AssetsFileRelative))))
	for _, r := range Routes {
		Router.Methods(r.Method).Path(r.Path).Name(r.Name).Handler(r.Handler)
	}
	return Router
}

func StartServer(f *mux.Router) {
	for {
		if internalConfig.IsSsl == "0" {
			e := http.ListenAndServe(":"+internalConfig.Port, f)
			if e != nil {
				Error("HTTP Server not started : " + e.Error())
				time.Sleep(5 * time.Second)
				continue
			}
		} else {
			e := http.ListenAndServeTLS(":"+internalConfig.SslPort, internalConfig.SslCertificate, internalConfig.SslCertificateKey, f)
			if e != nil {
				Error("HTTP Server not started : " + e.Error())
				time.Sleep(5 * time.Second)
				continue
			}
		}
		break
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/404.html")
}
