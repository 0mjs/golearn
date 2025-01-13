package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type URLData struct {
	Original string `json:"original"`
	Domain   string `json:"domain"`
}

type URLShortener struct {
	url   string
	data  map[string]URLData
	mutex sync.RWMutex
	alpha string
	nums  string
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		url:   "http://localhost:8080/",
		data:  make(map[string]URLData),
		alpha: "abcdefghijklmnopqrstuvwxyz",
		nums:  "0123456789",
	}
}

func (s *URLShortener) getURLs(
	w http.ResponseWriter,
	r *http.Request,
) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	json.NewEncoder(w).Encode(s.data)
}

func (s *URLShortener) shorten(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	chars := s.alpha + strings.ToUpper(s.alpha) + s.nums
	b := make([]byte, 12)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(b)

	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	code := string(b)
	shortURL := s.url + code

	s.mutex.Lock()
	s.data[code] = URLData{
		Original: req.URL,
		Domain:   shortURL,
	}
	s.mutex.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"url": shortURL,
	})
}

func (s *URLShortener) redirect(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	code := r.URL.Path[1:]
	urlData, exists := s.data[code]

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, urlData.Original, http.StatusFound)
}

func main() {
	s := NewURLShortener()

	http.HandleFunc("/", s.redirect)
	http.HandleFunc("/shorten", s.shorten)
	http.HandleFunc("/urls", s.getURLs)

	http.ListenAndServe(":8080", nil)
}
