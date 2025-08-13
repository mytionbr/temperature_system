package model

type WeatherApiResp struct {
	Current struct {
		TempC float32 `json:"temp_c"`
	} `json:"current"`
}
