package goverly

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ServeOverlay(configFile string, content embed.FS) {
	startTime := time.Now()
	http.Handle("/", http.FileServer(http.FS(content)))

	http.HandleFunc("/last-update", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.FormatInt(startTime.Unix(), 10)))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Add("content-type", "application/json")

		configBytes, err := os.ReadFile(configFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "config file not found"}`))
			return
		}

		w.Write(configBytes)
	})

	fmt.Println("Serving...")
	http.ListenAndServe(":8080", nil)
}
