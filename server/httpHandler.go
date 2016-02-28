package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func onlyGetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func loggingMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Got %v from %v\n", r.URL, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

var whoHasRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	hash := r.FormValue("q")
	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	agents := connectedAgents.getAllAgents()

	whoHashRequest := WhoHasRequest{Hash: hash}

	responseChans := make([]chan response, len(agents))
	for i, agent := range agents {
		responseChans[i] = make(chan response, 1)
		whoHasQuery := query{text: whoHashRequest.String(), response: responseChans[i]}
		select {
		case agent.queries <- whoHasQuery:
		default:
		}
	}

	time.Sleep(time.Second * 3)

	ret := []string{}
MainLoop:
	for i := range agents {
		select {
		case resp := <-responseChans[i]:
			if resp.err != nil || resp.text != "1" {
				continue MainLoop
			}
			ret = append(ret, agents[i].id)
			close(responseChans[i])
		default:
		}
	}

	str, _ := json.Marshal(struct {
		Addresses []string
	}{Addresses: ret})
	w.Write(str)

})

var whereToUploadRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO"))
})

var introduceMeRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO"))
})

func initHTTP(httpPort int) {

	apiMux := http.NewServeMux()
	apiMux.Handle("/who-has", onlyGetMiddleware(whoHasRequestHandler))
	apiMux.Handle("/where-to-upload", onlyGetMiddleware(whereToUploadRequestHandler))
	apiMux.Handle("/introduce-me", onlyGetMiddleware(introduceMeRequestHandler))

	http.Handle("/api/", loggingMiddelware(http.StripPrefix("/api", apiMux)))
	logger.Printf("Server is serving http on port 0.0.0.0:%v\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil)

}
