package server

import (
	"fmt"
	"net/http"
)

func onlyGetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

var whoHasRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO"))
})

var whereToUploadRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO"))
})

var introduceMeRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO"))
})

func initHTTP(httpPort int) {

	apiMux := http.NewServeMux()
	apiMux.Handle("/who-has", onlyGetHandler(whoHasRequestHandler))
	apiMux.Handle("/where-to-upload", onlyGetHandler(whereToUploadRequestHandler))
	apiMux.Handle("/introduce-me", onlyGetHandler(introduceMeRequestHandler))

	http.Handle("/api/", http.StripPrefix("/api", apiMux))
	logger.Printf("Server is listening on port 0.0.0.0:%v\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil)

}
