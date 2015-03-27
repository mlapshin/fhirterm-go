package fhirterm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type RestStorage struct {
	BaseUrl *url.URL
	Client  http.Client
}

type HttpParams map[string]string

func (this RestStorage) request(method string, u string, params HttpParams) (*ValueSet, error) {
	urlParsed, _ := url.Parse(u)

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	// force JSON output
	values.Add("_format", "application/json+fhir")

	urlParsed.RawQuery = values.Encode()
	absUrl := this.BaseUrl.ResolveReference(urlParsed).String()

	request, _ := http.NewRequest(method, absUrl, nil)
	resp, err := this.Client.Do(request)
	if err != nil {
		log.Printf("[FHIR REST] %s", err)
		return nil, err
	}

	log.Printf("[FHIR REST] %s %s => %s", method, absUrl, resp.Status)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vs ValueSet
	err = json.Unmarshal(body, &vs)
	if err != nil {
		return nil, err
	}

	return &vs, nil
}

func (this RestStorage) FindValueSetById(id string) (*ValueSet, error) {
	vs, err := this.request(
		"GET",
		"/ValueSet/"+url.QueryEscape(id),
		HttpParams{})

	if err != nil {
		log.Printf("Error: %s", err)
	}

	return vs, nil
}

func MakeRestStorage(cfg JsonObject) (Storage, error) {
	baseUrl, ok := cfg["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'base_url' attribute in config for 'rest' Storage")
	}

	baseUrlParsed, _ := url.Parse(baseUrl)
	return RestStorage{
		BaseUrl: baseUrlParsed,
		Client:  http.Client{},
	}, nil
}
