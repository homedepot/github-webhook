package holler

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

const Version = "holler 0.0.1"

var yellerHostnames = []string{
	"https://collector1.yellerapp.com",
	"https://collector2.yellerapp.com",
	"https://collector3.yellerapp.com",
	"https://collector4.yellerapp.com",
	"https://collector5.yellerapp.com",
}

func NewYeller(
	token string,
	env string,
	options ...func(*Yeller),
) *Yeller {
	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			timeout := time.Duration(1 * time.Second)
			return net.DialTimeout(network, addr, timeout)
		},
	}
	httpClient := &http.Client{Transport: &transport}
	hostname, _ := os.Hostname()

	client := &Yeller{
		token:    token,
		env:      env,
		hostname: hostname,

		httpClient: httpClient,
		hostnames:  yellerHostnames,
		random:     rand.New(rand.NewSource(time.Now().Unix())),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func UseCollectors(collectors ...string) func(*Yeller) {
	return func(client *Yeller) {
		client.hostnames = collectors
	}
}

func UseRandomSource(r rand.Source) func(*Yeller) {
	return func(client *Yeller) {
		client.random = rand.New(r)
	}
}

type Yeller struct {
	token    string
	env      string
	hostname string

	httpClient *http.Client
	hostnames  []string

	random *rand.Rand
}

func (y *Yeller) Notify(typ string, err error, data map[string]interface{}) {
	notification := errorNotification{
		Type:          typ,
		Message:       err.Error(),
		Host:          y.hostname,
		Environment:   y.env,
		ClientVersion: Version,
		CustomData:    data,
	}

	body, _ := json.Marshal(notification)

	numHostnames := len(y.hostnames)
	startingIndex := y.random.Int() % numHostnames

	for i := startingIndex; i < startingIndex+numHostnames; i++ {
		index := i % numHostnames
		tryAgain := y.tryNotify(y.hostnames[index], body)
		if !tryAgain {
			return
		}
	}

	log.Println("all yeller collectors rejected the error")
}

func (y *Yeller) tryNotify(hostname string, body []byte) bool {
	buffer := bytes.NewBuffer(body)
	response, err := y.httpClient.Post(hostname+"/"+y.token, "application/json", buffer)
	if err != nil {
		log.Printf("failed to emit error: %s\n", err.Error())
		return true
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		log.Println("yeller rejected authentication")
		return false
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		log.Printf("yeller responded with error: %d", response.StatusCode)
		return true
	}

	return false
}

type errorNotification struct {
	Type          string                 `json:"type"`
	Message       string                 `json:"message"`
	Host          string                 `json:"host"`
	Environment   string                 `json:"application-environment"`
	ClientVersion string                 `json:"client-version"`
	CustomData    map[string]interface{} `json:"custom-data"`

	// We don't populate these but the API requires that these keys are sent.
	StackTrace []stackFrame `json:"stacktrace"`
	URL        string       `json:"url"`
	Location   string       `json:"location"`
}

type stackFrame struct {
	Filename     string
	LineNumber   string
	FunctionName string
	Options      map[string]interface{}
}
