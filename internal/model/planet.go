// Package model ...
package model

// Planet представляет информацию о багаже.
type Planet struct {
	PlanetID     uint   `gorm:"type:serial;primarykey" json:"planetId"`
	Name         string `json:"name" example:"Планета"`
	Discovered   string `json:"discovered" example:"Неизвестно"`
	Mass         string `json:"mass" example:"Неизвестно"`
	Distance     string `json:"distance" example:"Неизвестно"`
	Info         string `json:"info" example:"Неизвестно"`
	Color1       string `json:"color1" example:"#ababab"`
	Color2       string `json:"color2" example:"#8a8a8a"`
	PlanetStatus string `json:"status" example:"активна"`
	ImageName    string `json:"imageName" example:"http://example.com/mars.jpg"`
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
	Planets         []Planet `json:"planets"`
	ConstellationID uint     `json:"constellationId" example:"1"`
}

type PlanetInConstellation struct {
	PlanetID  string `json:"id"`
	Name      string `json:"name" example:"Планета"`
	Color1    string `json:"color1" example:"#ababab"`
	Color2    string `json:"color2" example:"#8a8a8a"`
	ImageName string `json:"imageName" example:"http://example.com/mars.jpg"`
}
