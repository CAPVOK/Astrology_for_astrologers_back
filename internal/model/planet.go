// Package model ...
package model

// Planet представляет информацию о багаже.
type Planet struct {
	PlanetID     uint   `gorm:"type:serial;primarykey" json:"planet_id"`
	Name         string `json:"name" example:"Планета"`
	Discovered   string `json:"discovered" example:"Неизвестно"`
	Mass         string `json:"mass" example:"Неизвестно"`
	Distance     string `json:"distance" example:"Неизвестно"`
	Info         string `json:"info" example:"Неизвестно"`
	Color1       string `json:"color1" example:"#ababab"`
	Color2       string `json:"color2" example:"#8a8a8a"`
	PlanetStatus string `json:"status" example:"активна"`
	ImageName    string `json:"image_name" example:"http://example.com/mars.jpg"`
}

// PlanetRequest представляет запрос на создание планеты.
type PlanetRequest struct {
	Name       string `json:"name" example:"Планета"`
	Discovered string `json:"discovered" example:"Неизвестно"`
	Mass       string `json:"mass" example:"Неизвестно"`
	Distance   string `json:"distance" example:"Неизвестно"`
	Info       string `json:"info" example:"Неизвестно"`
	Color1     string `json:"color1" example:"#ababab"`
	Color2     string `json:"color2" example:"#8a8a8a"`
}

// PlanetsGetResponse представляет ответ с информацией о багажах и идентификаторе созвездия.
type PlanetsGetResponse struct {
	Planets         []Planet `json:"Planets"`
	ConstellationID uint     `json:"constellation_id" example:"1"`
}
