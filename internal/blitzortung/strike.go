package blitzortung

import "encoding/json"

type Strike struct {
	Time int64   `json:"time"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Alt  float64 `json:"alt"`
	Pol  int     `json:"pol"`
	Mds  int64   `json:"mds"`
	Scs  int     `json:"scs"`
}

func ParseStrike(payload []byte) (Strike, error) {
	var s Strike
	if err := json.Unmarshal(payload, &s); err != nil {
		return Strike{}, err
	}
	return s, nil
}
