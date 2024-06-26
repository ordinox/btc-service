package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetFastestFeeRate() (uint, error) {
	fee := struct {
		FastestFee uint `json:"fastestFee"`
		AvgFee     uint `json:"halfHourFee"`
		SlowFee    uint `json:"hourFee"`
		MinFee     uint `json:"minimumFee"`
	}{}

	res, err := http.Get("https://mempool.space/api/v1/fees/recommended")
	if err != nil {
		fmt.Println("Http Error getting fee rate", err)
		return 1, err
	}
	defer res.Body.Close()
	o, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading fee rate response", err)
		return 1, err
	}
	err = json.Unmarshal(o, &fee)
	if err != nil {
		fmt.Println("Error unmarshalling json", err)
		return 1, err
	}
	return fee.FastestFee, nil
}
