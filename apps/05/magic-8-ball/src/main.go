// A Spin component written in Go that returns "Hello, Fermyon!"
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/fermyon/spin/sdk/go/key_value"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		// We need to seed the RNG here for each request because no program
		// state is persisted between calls to the API
		rand.Seed(time.Now().UnixNano())
		w.Header().Set("Content-Type", "text/plain")

		question, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		answer, err := getOrSetAnswer(string(question))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "{\"answer\": \"%s\"}", answer)
	})
}

func getOrSetAnswer(q string) (string, error) {
	store, err := key_value.Open("default")
	defer key_value.Close(store)

	if err != nil {
		return "", err
	}
	response := ""

	y, err := key_value.Exists(store, q)
	if err != nil {
		return "", err
	}
	if y {
		v, err := key_value.Get(store, q)
		if err != nil {
			return "", err
		}
		response = string(v)
	}

	if response == "Ask again later." || response == "" {
		response = answer()
		key_value.Set(store, q, []byte(response))
	}

	return response, nil
}

func answer() string {
	answers := []string{"Ask again later.", "Absolutely!", "Unlikely", "Simply put, no"}

	idx := rand.Intn(len(answers))

	return answers[idx]
}

func main() {}
