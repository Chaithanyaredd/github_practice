package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

func shortenURLHandler(c echo.Context) error {
	var request struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	originalURL := request.URL

	urlDatabase.Lock()
	defer urlDatabase.Unlock()
	// Check if the URL has already been shortened
	if shortURL, found := urlDatabase.reverseLookup[originalURL]; found {
		return c.JSON(http.StatusOK, map[string]string{"short_url": shortURL})
	}

	// Generate a new short URL
	shortURL := generateShortURL()
	urlDatabase.urls[shortURL] = originalURL
	urlDatabase.reverseLookup[originalURL] = shortURL
	domain := getDomain(originalURL)
	urlDatabase.metrics[domain]++

	return c.JSON(http.StatusOK, map[string]string{"short_url": shortURL})
}

func redirectURLHandler(c echo.Context) error {
	shortURL := c.Param("shortURL")

	urlDatabase.RLock()
	defer urlDatabase.RUnlock()

	if originalURL, ok := urlDatabase.urls[shortURL]; ok {
		return c.Redirect(http.StatusMovedPermanently, originalURL)
	}
	return c.JSON(http.StatusNotFound, map[string]string{"error": "URL not found"})
}

func getMetricsHandler(c echo.Context) error {
	urlDatabase.RLock()
	defer urlDatabase.RUnlock()

	type metric struct {
		Domain string `json:"domain"`
		Count  int    `json:"count"`
	}
	var domainMetrics []metric
	for domain, count := range urlDatabase.metrics {
		domainMetrics = append(domainMetrics, metric{Domain: domain, Count: count})
	}
	sort.Slice(domainMetrics, func(i, j int) bool {
		return domainMetrics[i].Count > domainMetrics[j].Count
	})
	if len(domainMetrics) > 3 {
		domainMetrics = domainMetrics[:3]
	}

	return c.JSON(http.StatusOK, domainMetrics)
}

func getDomain(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is healthy"))
}
