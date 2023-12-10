package clients

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var client = http.Client{Timeout: time.Second * 10}

const REQUEST_ERROR_LOG = "Could not send request to '%s'. Root cause:\n%s"
const RESPONSE_ERROR_LOG = "Could not parse body. Root cause:\n%s"

func PrepareRequest(method, url string) *http.Request {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf(REQUEST_ERROR_LOG, url, err)
		return nil
	}
	request.Header.Add("User-Agent", os.Getenv("EMAIl"))
	request.Header.Add("Accept-Encoding", "gzip, deflate")

	return request
}

func SendRequestAndGetBody(request *http.Request) []byte {
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf(REQUEST_ERROR_LOG, request.URL, err)
		return nil
	}
	defer response.Body.Close()

	var reader io.Reader
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			log.Fatalf(RESPONSE_ERROR_LOG, err)
			return nil
		}
		defer reader.(*gzip.Reader).Close()
	default:
		reader = response.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf(RESPONSE_ERROR_LOG, err)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Got a response with statuc code %d and body:\n%s", response.StatusCode, string(body))
		return nil
	}

	return body
}
