package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/b75/yummy-wookie/util"
)

type Config struct {
	StaticDir string
	TplDir    string
}

var conf *Config = &Config{}

var configFile string

var tpls *template.Template

func init() {
	flag.StringVar(&configFile, "c", "config.json", "Main config json file")
}

func main() {
	var err error
	flag.Parse()
	if err = util.LoadConfig(conf, configFile); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if conf.StaticDir == "" {
		log.Fatal("StaticDir not set")
	}
	if conf.TplDir == "" {
		log.Fatal("TplDir not set")
	}

	/*
		fmap := template.FuncMap{
			"aidontime": func(t time.Time) string {
				return t.Format(AIDON_TIME_FORMAT)
			},
			"utcnow": func() time.Time {
				return time.Now().In(time.UTC)
			},
			"uuid": func() string {
				return util.Uuid()
			},
			"minuteadd": func(t time.Time, minutes int) time.Time {
				return t.Add(time.Duration(minutes) * time.Minute)
			},
		}
	*/

	fmap := template.FuncMap{}

	tpls, err = util.LoadTemplates(conf.TplDir, fmap)
	if err != nil {
		log.Fatalf("error loading templates: %v", err)
	}

	fs := http.FileServer(http.Dir(conf.StaticDir))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Print("listening on 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func executeTemplate(w io.Writer, name string, data interface{}) {
	if err := tpls.ExecuteTemplate(w, name, data); err != nil {
		panic(err)
	}
}
