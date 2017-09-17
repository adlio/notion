package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type API struct {
	client   *http.Client
	BaseURL  string
	APIToken string
}

type EchoResponse struct {
	Hello  string `json:"hello"`
	Status string `json:"status"`
}

func NewAPI(token string) *API {
	a := API{}
	a.client = http.DefaultClient
	a.APIToken = token
	a.BaseURL = "https://app.usenotion.com/api/v1"
	return &a
}

func (a *API) EchoTest() string {
	response := EchoResponse{}
	a.Get("/echo", &response)
	return response.Status
}

func (a *API) SendSingleIngredientReport(report *IngredientReport) (response IngredientReportResponse, err error) {
	err = a.Post("/report", report, &response)
	return
}

func (a *API) SendBatchIngredientReport(report *BatchIngredientReport) (response IngredientReportResponse, err error) {
	err = a.Post("/batch_report", report, &response)
	return
}

func (a *API) Get(path string, target interface{}) error {
	url := fmt.Sprintf("%s%s", a.BaseURL, path)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "github.com/adlio/notion")
	req.Header.Set("Authorization", a.APIToken)
	if err != nil {
		return errors.Wrapf(err, "Invalid GET request %s", url)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failure on %s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		return errors.Errorf("HTTP request failure on %s: %s %s", url, string(body), err)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(target)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Wrapf(err, "JSON decode failed on %s: %s", url, string(body))
	}

	return nil
}

func (a *API) Post(path string, postData interface{}, target interface{}) error {
	url := fmt.Sprintf("%s%s", a.BaseURL, path)

	buffer := new(bytes.Buffer)
	if postData != nil {
		json.NewEncoder(buffer).Encode(postData)
	}

	req, err := http.NewRequest("POST", url, buffer)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "github.com/adlio/notion")
	req.Header.Set("Authorization", a.APIToken)
	if err != nil {
		return errors.Wrapf(err, "Invalid POST request %s", url)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failure on %s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		return errors.Wrapf(err, "HTTP request failure on %s: %s %s", url, string(body), err)
	}

	return nil
}
