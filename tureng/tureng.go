package tureng

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SearchResponse struct {
	Exception    string `json:"ExceptionMessage"`
	IsSuccessful bool   `json:"IsSuccessful"`

	Result struct {
		IsFound            int      `json:"IsFound"`
		IsEnglishToTurkish int      `json:"IsTRToEN"`
		Suggestions        []string `json:"Suggestions"`
		Results            []struct {
			Category string `json:"CategoryEN"`
			Term     string `json:"Term"`
			Type     string `json:"TypeEN"`
		} `json:"Results"`
	} `json:"MobileResult"`
}

type AutoCompleteResponse struct {
	Words []string
}

type SearchRequest struct {
	Term string `json:"Term"`
	Code string `json:"Code"`
}

func Search(word string) (*SearchResponse, error) {
	code := md5.Sum([]byte(fmt.Sprintf("%s%s", word, SECRET)))
	req := SearchRequest{word, hex.EncodeToString(code[:])}

	requestJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", SEARCH_URL, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", USER_AGENT)
	request.Header.Set("Content-Type", BODY_TYPE)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response := &SearchResponse{}

	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AutoComplete(term string) (*AutoCompleteResponse, error) {
	term = url.QueryEscape(term)
	requestUrl := fmt.Sprintf(AUTOCOMPLETE_URL, term)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var words []string
	err = json.Unmarshal(body, &words)

	if err != nil {
		return nil, err
	}

	return &AutoCompleteResponse{Words: words}, nil
}

const (
	SEARCH_URL       = "http://ws.tureng.com/TurengSearchServiceV4.svc/Search"
	AUTOCOMPLETE_URL = "https://ac.tureng.co?t=%s"
	SECRET           = "46E59BAC-E593-4F4F-A4DB-960857086F9C"
	BODY_TYPE        = "application/json"
	USER_AGENT       = "Mozilla/5.0 (Linux; Android 7.1.1; Android SDK built for x86 Build/NYC; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.68 Mobile Safari/537.36"
)
