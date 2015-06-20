package main

import (
	//"html/template"
	//"io/ioutil"
	"log"
	//"regexp"

	"net/http"

	"flag"

	//"github.com/JackKnifed/blackfriday"
	"github.com/JackKnifed/gnosis"
	//"github.com/russross/blackfriday"
)

var indexes gnosisIndex[]

var configFile = flag.String("config", "config.json", "specify a configuration file")

func main() {
	flag.Parse()

	// ##TODO## check for false returnear- if null, the config could not be loaded
	gnosis.LoadConfig(*configFile)

	config := gnosis.GetConfig()

	// ##TODO## add global redirects and set them up in here

	gnosis.ParseTemplates(config.Global)

	for _, individualIndex := range config.Indexes {
		index, err := gnosis.openIndex(individualIndex)
		if err != nil {
			log.Println(err)
		} else {
			indexes = append(indexes, index...)
		}
	}

	for _, redirect := range config.Redirects {
		http.Handle(redirect.Requested, http.RedirectHandler(redirect.Target, redirect.Code))
	}
	for _, individualServer := range config.Server {
		http.HandleFunc(individualServer.Prefix, gnosis.MakeHandler(individualServer))
	}

	log.Println(http.ListenAndServe(":"+config.Global.Port, nil))

}
