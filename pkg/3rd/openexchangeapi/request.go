package openexchangeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const urlFormat = "https://openexchangerates.org/api/latest.json?app_id=%s&symbols=%s"

var token = "2885605918874255931300ad0c789f1c"
var cacheDuration = time.Duration(1) * time.Hour

var cache Cache[string, float32]
var cacheInit sync.Once

type OpenExchangeResponse struct {
	Rates map[string]float32 `json:"rates"`
}

func RequestUSDPair(pair string) (rate float32, err error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf(urlFormat, token, pair))
	if err != nil {
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	jsonStruct := OpenExchangeResponse{}
	json.Unmarshal(data, &jsonStruct)
	rate, ok := jsonStruct.Rates[pair]
	if !ok {
		err = fmt.Errorf("value was not found in response")
	}
	return
}

func RequestUSDPairCached(pair string) (rate float32, err error) {
	cacheInit.Do(func() {
		cache = NewCache[string, float32](cacheDuration)
	})
	return cache.Get(pair, RequestUSDPair)
}
