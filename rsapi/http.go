package rsapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

// Log request, dump its content if required then make request and log response and dump it.
func (a *Api) PerformRequest(req *http.Request) (*http.Response, error) {
	var id string
	var startedAt time.Time
	if a.Logger != nil {
		startedAt = time.Now()
		b := make([]byte, 6)
		io.ReadFull(rand.Reader, b)
		id = base64.StdEncoding.EncodeToString(b)
		a.Logger.Printf("[%s] %s %s", id, req.Method, req.URL.String())
	}
	if a.DumpRequestResponse {
		b, err := httputil.DumpRequest(req, true)
		if err == nil {
			fmt.Printf("REQUEST\n-------\n%s\n\n", b)
		}
	}
	// Sign last so auth headers don't get printed or logged
	if err := a.Auth.Sign(req, a.Host); err != nil {
		return nil, err
	}
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if a.Logger != nil {
		d := time.Since(startedAt)
		a.Logger.Printf("[%s] %s in %s", id, resp.Status, d.String())
	}
	if a.DumpRequestResponse {
		b, err := httputil.DumpResponse(resp, false)
		if err == nil {
			fmt.Printf("RESPONSE\n--------\n%s", b)
		}
	}

	return resp, err
}

// Deserialize JSON response into generic object.
// If the response has a "Location" header then the returned object is a map with one key "Location"
// containing the value of the header.
func (a *Api) LoadResponse(resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()
	var respBody interface{}
	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response (%s)", err.Error())
	}
	if len(jsonResp) > 0 {
		err = json.Unmarshal(jsonResp, &respBody)
		if err != nil {
			return nil, fmt.Errorf("Failed to load response (%s)", err.Error())
		}
	}
	// Special case for "Location" header, assume that if there is a location there is no body
	loc := resp.Header.Get("Location")
	if len(loc) > 0 {
		var bodyMap = make(map[string]interface{})
		bodyMap["Location"] = loc
		respBody = interface{}(bodyMap)
	}
	return respBody, err
}
