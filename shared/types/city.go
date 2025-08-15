package types

type CityInput struct {
	Name string  `json:"name" binding:"required"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}
