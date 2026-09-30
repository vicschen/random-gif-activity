package main

import (
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	http.HandleFunc("/", rootHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Open http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	count := parsedCount(r.URL.Query().Get("count"))
	urls := pickRandomGIFs(count)

	var body strings.Builder
	body.WriteString(`<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>Random GIF Activity</title>
	<style>
		body {
			margin: 0;
			font-family: Arial, sans-serif;
			background: linear-gradient(135deg, #f5d0fe, #dbeafe, #dcfce7);
			color: #1f2937;
		}
		main {
			max-width: 900px;
			margin: 0 auto;
			padding: 40px 20px 60px;
			text-align: center;
		}
		h1 {
			margin-bottom: 20px;
		}
		.gallery {
			display: flex;
			flex-wrap: wrap;
			justify-content: center;
			gap: 16px;
			margin: 24px 0;
		}
		img {
			max-width: 320px;
			width: 100%;
			height: auto;
			border-radius: 12px;
			box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
		}
		form {
			display: inline-flex;
			align-items: center;
			gap: 10px;
			padding: 12px 18px;
			background: rgba(255, 255, 255, 0.7);
			border-radius: 999px;
		}
		input {
			padding: 8px 10px;
			border-radius: 8px;
			border: 1px solid #cbd5e1;
		}
		button {
			padding: 8px 14px;
			border: none;
			border-radius: 8px;
			background: #2563eb;
			color: white;
			cursor: pointer;
		}
	</style>
</head>
<body>
	<main>
		<h1>Random Animal GIFs</h1>
		<div class="gallery">
`)

	for _, url := range urls {
		body.WriteString(fmt.Sprintf("<img src=\"%s\" alt=\"Cute animal GIF\">\n", html.EscapeString(url)))
	}

	body.WriteString(`
		</div>
		<form method="get" action="/">
			<label for="count">How many GIFs?</label>
			<input type="number" id="count" name="count" min="1" max="3" value="`)
	body.WriteString(strconv.Itoa(count))
	body.WriteString(`" required>
			<button type="submit">Show me</button>
		</form>
	</main>
</body>
</html>`)

	if _, err := w.Write([]byte(body.String())); err != nil {
		log.Printf("write response: %v", err)
	}
}

func parsedCount(raw string) int {
	count, err := strconv.Atoi(raw)
	if err != nil || count < 1 || count > 3 {
		return 1
	}
	return count
}

func pickRandomGIFs(count int) []string {
	gifURLs := readGIFs("gifs.txt")
	if len(gifURLs) == 0 {
		return []string{"https://example.com/placeholder.gif"}
	}

	selected := make([]string, 0, count)
	for i := 0; i < count; i++ {
		selected = append(selected, gifURLs[rng.Intn(len(gifURLs))])
	}
	return selected
}

func readGIFs(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("read GIF list: %v", err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			urls = append(urls, trimmed)
		}
	}
	return urls
}
