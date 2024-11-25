package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"access_token":"12345"}`)
	})

	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		msg, _ := io.ReadAll(r.Body)
		fmt.Println(string(msg))
		w.WriteHeader(200)
	})

	http.HandleFunc("/visit", func(w http.ResponseWriter, r *http.Request) {
		msg, _ := io.ReadAll(r.Body)
		if _, err := os.Stat("data/visit.json"); err == nil {
			err = os.Remove("data/visit.json")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
		}
		err := os.WriteFile("data/visit.json", msg, 0644)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("data/visit.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var visits struct {
			Request []struct {
				Mid          string `json:"mid"`
				Tid          string `json:"tid"`
				Vendor       string `json:"vendor"`
				VisitTime    string `json:"visitTime"`
				TicketTime   string `json:"ticketTime"`
				VisitType    string `json:"visitType"`
				VisitReason  string `json:"visitReason"`
				VisitResult  string `json:"visitResult"`
				VisitStatus  string `json:"visitStatus"`
				LinkEvidence string `json:"linkEvidence"`
			} `json:"request"`
		}
		err = json.Unmarshal(data, &visits)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><table border=1>")
		for _, v := range visits.Request {
			fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td><a href=%q>%s</a></td></tr>",
				v.Mid, v.Tid, v.Vendor, v.VisitTime, v.TicketTime, v.LinkEvidence, v.VisitResult)
		}
		fmt.Fprintf(w, "</table></body></html>")
	})

	http.ListenAndServe(":8484", nil)
}
