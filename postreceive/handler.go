package postreceive

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"encoding/json"
	"crypto/sha1"
	"crypto/hmac"
	"strings"
	"crypto/subtle"
	"github.com/homedepot/github-webhook/events"
	"github.com/homedepot/github-webhook/executor"
	"expvar"
)

//https://developer.github.com/webhooks/

const POST = "POST"
const APPLICATION_JSON = "application/json"
const USER_AGENT = "User-Agent"
const X_GITHUB_DELIVERY = "X-Github-Delivery"
const X_GITHUB_EVENT = "X-GitHub-Event"

func HttpHandler(w http.ResponseWriter, req *http.Request) {

	metrics.Add("requests", 1)

	if req.Body == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		metrics.Add("error", 1)
		return
	}

	defer req.Body.Close()

	if !strings.EqualFold(APPLICATION_JSON, req.Header.Get("Content-Type")) {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		metrics.Add("error", 1)
		return
	}

	if !strings.EqualFold(POST,req.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		metrics.Add("error", 1)
		return
	}

	if userAgent := req.Header.Get(USER_AGENT); !strings.HasPrefix(userAgent,"GitHub") {
		http.Error(w, "Invalid User Agent Header", http.StatusBadRequest)
		metrics.Add("error", 1)
		return
	}

	delivery := req.Header.Get(X_GITHUB_DELIVERY)
	if len(delivery) == 0 {
		http.Error(w, "Missing X-Github-Delivery Header", http.StatusBadRequest)
		metrics.Add("error", 1)
		return
	}

	eventType := req.Header.Get(X_GITHUB_EVENT)
	if len(eventType) == 0 {
		http.Error(w, "Missing X-GitHub-Event Header", http.StatusBadRequest)
		metrics.Add("error", 1)
		return
	}

	eventTypeId, found := events.EventTypes[eventType]
	if !found {
		http.Error(w, "Unsupported Github-Event Type", http.StatusBadRequest)
		metrics.Add("error", 1)
		return
	}

	metrics.Add(eventType, 1)

	if eventTypeId == events.PING {
		w.WriteHeader(http.StatusOK)
		return
	}

	payload, readError := ioutil.ReadAll(req.Body)
	if readError != nil {
		http.Error(w, fmt.Sprint(readError), http.StatusInternalServerError)
		metrics.Add("error", 1)
		return
	}

	if !validateSignature(payload, secret, req) {
		http.Error(w, "Signature was not valid", http.StatusUnauthorized)
		metrics.Add("error", 1)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		http.Error(w, fmt.Sprint(readError), http.StatusInternalServerError)
		metrics.Add("error", 1)
		return
	}

	event := events.Event{
		ID: delivery,
		Type: eventTypeId,
		Name: eventType,
		Data: data,
	}

	err := executor.ProcessEvent(event)
	if err != nil {
		http.Error(w, fmt.Sprint(readError), http.StatusInternalServerError)
		metrics.Add("error", 1)
		return
	}

	w.WriteHeader(http.StatusOK)

}

var secret string

func SetSecret(value string) {
	secret = value
}

func validateSignature(payload []byte, secret string, req *http.Request) bool {

	if len(secret) == 0 {
		return true
	}

	requestSignature := req.Header.Get("X-Hub-Signature")
	if len(requestSignature) == 0 {
		return false
	}
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Reset()
	mac.Write([]byte(payload))
	calculatedSignature := fmt.Sprintf("sha1=%x", mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(requestSignature), []byte(calculatedSignature)) == 1
}

var metrics *expvar.Map

func SetMetrics(m *expvar.Map) {
	metrics = m
}

