package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var port string

// Entry Point into the application
func main() {
	// Weakly seed the random number generator
	rand.Seed(time.Now().Unix())

	// Check if we have a port env variable to follow
	if os.Getenv("PORT") == "" {
		port = "8000"
	} else {
		port = os.Getenv("PORT")
	}

	// Print a startup message
	fmt.Println("Starting Reflector on port " + port)

	// Start the server
	http.HandleFunc("/", reflect)
	http.ListenAndServe(":"+port, nil)
}

// This generates a weak random number
// It is only used for selecting a server
func random(min, max int) int {
	// If we give two numbers that are the same we should just return it
	if max-min == 0 {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}

func reflect(w http.ResponseWriter, r *http.Request) {
	// Retrieve the stream proxy database from the Github
	// Todo: Add caching to this
	resp, _ := http.Get("https://openstreamproject.github.io/StreamDatabase/stream_proxies.csv")
	body, _ := ioutil.ReadAll(resp.Body)
	csv := string(body[:])

	// Rather than loading a csv lib, just split the csv based on linebreaks
	lines := strings.Split(csv, "\n")

	// Create a slice for all usable proxy urls
	var usable_lines []string

	// Loop through each line and only add the non commented out lines
	for _, line := range lines {
		if len(line) != 0 && string(line[0]) != "#" {
			usable_lines = append(usable_lines, line)
		}
	}

	// Get the length of the usable lines
	length := len(usable_lines) - 1

	// Randomly select a stream proxy
	// Todo: Make this selection weighted based on throughput
	line := usable_lines[random(0, length)]

	// Remove the throughput, we'll use it eventually
	line = strings.Split(line, ",")[0]

	// Append the channel id
	line = line + "/channels/" + r.URL.Query().Get("channel") + "/playlist.m3u8"

	// Redirect to server
	http.Redirect(w, r, line, 301)
}
