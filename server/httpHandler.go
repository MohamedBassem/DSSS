package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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
			close(responseChans[i])
			if resp.err != nil || resp.text != "1" {
				continue MainLoop
			}
			ret = append(ret, agents[i].id)
		default:
		}
	}

	str, _ := json.Marshal(struct {
		Addresses []string
	}{Addresses: ret})
	w.Write(str)

})

var whereToUploadRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// TODO: Consider the size
	agents := connectedAgents.getAllAgents()

	if len(agents) < replicationFactor {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Number of online agents is smaller than the replication factor"))
		return
	}

	var rep = replicationFactor

	ret := []string{}
	for ; rep > 0; rep-- {
		idx := rand.Intn(len(agents))
		ret = append(ret, agents[idx].id)
		agents = append(agents[:idx], agents[idx+1:]...)
	}

	str, _ := json.Marshal(struct {
		Addresses []string
	}{Addresses: ret})
	w.Write(str)

})

var introduceMeRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	to := r.FormValue("to")
	sizeStr := r.FormValue("size")
	hash := r.FormValue("hash")

	if to == "" || sizeStr == "" || hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	agent := connectedAgents.get(to)
	if agent == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	introductionRequest := IntroductionRequest{
		Address: r.RemoteAddr,
		Size:    size,
		Hash:    hash,
	}

	responseChan := make(chan response, 1)
	agent.queries <- query{
		text:     introductionRequest.String(),
		response: responseChan,
	}

	select {
	case resp := <-responseChan:
		if resp.err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		str, _ := json.Marshal(struct {
			IntroductionKey string `json:"introduction-key"`
		}{IntroductionKey: resp.text})
		w.Write(str)
	case <-time.After(time.Second * 3):
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
