package adapter

import (
	"app/model"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type StrapiAdapter struct {
	HttpClient *http.Client
	Host       string
}

func NewStrapiAdapter() *StrapiAdapter {
	return &StrapiAdapter{
		HttpClient: &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			Timeout:   60 * time.Second,
		},
		Host: os.Getenv("STRAPI_URL"),
	}
}

func (adapter *StrapiAdapter) CreateCategory(body model.StrapiCategoryCreateWrapper) model.StrapiCategoryCreateResponseWrapper {
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(adapter.Host, "api", "business-categories")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp model.StrapiCategoryCreateResponseWrapper
	response, _ := adapter.HttpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp
}

func (adapter *StrapiAdapter) Localize(id int, body model.StrapiCategoryCreate) model.StrapiCategoryCreateResponse {
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(adapter.Host, "api", "business-categories", fmt.Sprintf("%d", id), "localizations")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp model.StrapiCategoryCreateResponse
	response, _ := adapter.HttpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp
}

func (adapter *StrapiAdapter) Parent(id int, reqa model.StrapiCategoryParenting) {
	body := model.StrapiCategoryParentingWrapper{
		Data: reqa,
	}
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(adapter.Host, "api", "business-categories", fmt.Sprintf("%d", id))
	req, _ := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	response, _ := adapter.HttpClient.Do(req)
	defer response.Body.Close()
}

func (adapter *StrapiAdapter) Insert(body model.BusinessPageInsert) model.BusinessPageResponseWrapper {
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(adapter.Host, "api", "business-pages")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp model.BusinessPageResponseWrapper
	response, _ := adapter.HttpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp
}

func (adapter *StrapiAdapter) Localizations(body model.BusinessPage, pageId int) model.BusinessPageResponseWrapper {
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(adapter.Host, "api", "business-pages", fmt.Sprintf("%d", pageId), "localizations")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp model.BusinessPageResponseWrapper
	response, err := adapter.HttpClient.Do(req)
	if err != nil {
		return model.BusinessPageResponseWrapper{}
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return model.BusinessPageResponseWrapper{}
	}
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		fmt.Println(err)
		return model.BusinessPageResponseWrapper{}
	}
	return resp
}
