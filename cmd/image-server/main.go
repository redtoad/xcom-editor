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
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"sort"
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

	var paths []string
	for pth := range resources.Images {
		paths = append(paths, pth)
	}
	sort.Strings(paths)

	w.Header().Set("Content-Type", "text/html")

	files := []string{
		"templates/base.html",
		"templates/index.html",
	}
	ts, err := template.ParseFS(templates, files...)
	if err != nil {
		log.Printf("could not load templates: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", struct {
		Images []string
	}{
		paths,
	})
	if err != nil {
		log.Printf("error: %v", err)
	}

}

func ResourceDetails(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	files := []string{
		"templates/base.html",
		"templates/resource.html",
	}
	ts, err := template.ParseFS(templates, files...)
	if err != nil {
		log.Printf("could not load templates: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", struct {
		Path string
		Meta resources.ImageEntry
	}{
		r.URL.Path[9:],
		resources.Images[r.URL.Path[9:]],
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

// OpenURL opens the specified URL in the default browser of the user.
// Shamelessly stolen from https://stackoverflow.com/a/39324149.
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
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
	r.PathPrefix("/details").HandlerFunc(ResourceDetails)
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

	log.Println("Opening browser...")
	_ = OpenURL("http://localhost:8080")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	log.Println("Press Ctrl+C to stop.")

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
