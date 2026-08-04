package api

import (
	"net/http"
)

const layout = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}
