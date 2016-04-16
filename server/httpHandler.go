package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/MohamedBassem/DSSS/internal/structs"
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

func onlyPostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
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

	whoHashRequest := structs.WhoHasRequest{Hash: hash}

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

	str, _ := json.Marshal(structs.WhoHasResponseJSON{Addresses: ret})
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

var RelayRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var request structs.UploadRequestJSON

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request without a body"))
		return
	}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse request body"))
		return
	}

	agent := connectedAgents.get(request.To)
	if agent == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Agent not found .."))
		return
	}

	responseChan := make(chan response, 1)
	uploadRequest := structs.UploadRequest{Hash: request.Hash, Content: request.Content}
	uploadQuery := query{
		text:     uploadRequest.String(),
		response: responseChan,
	}

	agent.queries <- uploadQuery

	response := <-responseChan

	if response.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response.err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
})

var DownloadRequestHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var request structs.DownloadRequestJSON

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request without a body"))
		return
	}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse request body"))
		return
	}

	agent := connectedAgents.get(request.From)
	if agent == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Agent not found .."))
		return
	}

	responseChan := make(chan response, 1)
	downloadRequest := structs.DownloadRequest{Hash: request.Hash}
	downloadQuery := query{
		text:     downloadRequest.String(),
		response: responseChan,
	}

	agent.queries <- downloadQuery

	response := <-responseChan

	if response.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(response.err.Error()))
		return
	}

	if strings.HasPrefix(response.text, "ERROR") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(strings.TrimPrefix(response.text, "ERROR ")))
		return
	}

	w.WriteHeader(http.StatusOK)
	ret, _ := json.Marshal(structs.DownloadResponseJSON{
		Hash:    request.Hash,
		Content: response.text,
	})
	w.Write(ret)
})

func initHTTP(httpPort int) {

	apiMux := http.NewServeMux()
	apiMux.Handle("/who-has", onlyGetMiddleware(whoHasRequestHandler))
	apiMux.Handle("/where-to-upload", onlyGetMiddleware(whereToUploadRequestHandler))
	apiMux.Handle("/relay", onlyPostMiddleware(RelayRequestHandler))
	apiMux.Handle("/download", onlyPostMiddleware(DownloadRequestHandler))

	http.Handle("/api/", loggingMiddelware(http.StripPrefix("/api", apiMux)))
	logger.Printf("Server is serving http on port 0.0.0.0:%v\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil)

}
