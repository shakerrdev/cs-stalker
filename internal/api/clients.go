package api

import (
	"net/http"
	"time"
)

var (
	steamClient   = &http.Client{Timeout: 10 * time.Second}
	leetifyClient = &http.Client{Timeout: 10 * time.Second}
)
