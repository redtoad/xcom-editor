package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/lib/resources"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

var port string // server port
var root string // root path of game (containing all images and save games)

var loader *resources.ResourceLoader

//go:embed templates/*
var templates embed.FS

func Index(w http.ResponseWriter, r *http.Request) {

	log.Println("Index")

	var paths []string
	for pth, _ := range resources.Images {
		println(pth)
		paths = append(paths, pth)
	}

	w.Header().Set("Content-Type", "text/html")

	tpl := template.Must(template.New("index.html").ParseFS(templates, "templates/index.html"))
	err := tpl.Execute(w, struct {
		Images []string
	}{
		paths,
	})
	if err != nil {
		log.Printf("error: %v", err)
	}

}

func ServeImage(w http.ResponseWriter, r *http.Request) {

	pth := r.URL.Path[10:]
	basename := path.Base(pth)
	var img image.Image
	var err error
	paletteNr := r.URL.Query().Get("palette")

	if paletteNr != "" {
		val, _ := strconv.Atoi(paletteNr)
		img, err = loader.LoadImageWithPalette(pth, val)
	} else {
		img, err = loader.LoadImage(pth)
	}
	if err != nil {
		if err == resources.ErrImageNotFound {
			log.Printf("Error: File %s not found!\n", pth)
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		log.Printf("Error: Could not load %s: %s\n", pth, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Serving file %s\n", pth)

	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", basename))
	_, err = w.Write(buf.Bytes())

}

func main() {
	flag.StringVar(&port, "port", "8080", "port server run on (default: 8080)")
	flag.Parse()
	if root = flag.Arg(0); root == "" {
		// use current dir as default
		root, _ = os.Getwd()
	}

	log.Printf("Starting server...\n")
	log.Printf("Version %s-%s %s %s\n", version, commit, date, builtBy)
	log.Printf("Game root: %s\n", root)
	log.Printf("Try opening http://localhost:%s/resource/UNITS/ZOMBIE.PCK\n", port)

	var err error
	loader, err = resources.NewResourceLoader(root)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/resource").HandlerFunc(ServeImage)
	r.PathPrefix("/").HandlerFunc(Index)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
