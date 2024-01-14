package usecase

import (
	"errors"
	"strings"

	"space/internal/model"
)

type PlanetUseCase interface {
}

func (uc *UseCase) GetPlanets(searchCode string, userID uint) (model.PlanetsGetResponse, error) {
	searchCode = strings.ToUpper(searchCode + "%")
	planets, err := uc.Repository.GetPlanets(searchCode, userID)
	if err != nil {
		return model.PlanetsGetResponse{}, err
	}

	return planets, nil
}

func (uc *UseCase) GetPlanetByID(planetID, userID uint) (model.Planet, error) {
	if planetID <= 0 {
		return model.Planet{}, errors.New("недопустимый ИД планеты")
	}
	planet, err := uc.Repository.GetPlanetByID(planetID, userID)
	if err != nil {
		return model.Planet{}, err
	}
	return planet, nil
}

func (uc *UseCase) CreatePlanet(userID uint, requestPlanet model.PlanetRequest) error {
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if requestPlanet.Name == "" {
		return errors.New("название планеты должно быть заполнено")
	}
	planet := model.Planet{
		Name:         requestPlanet.Name,
		Discovered:   requestPlanet.Discovered,
		Info:         requestPlanet.Info,
		Mass:         requestPlanet.Mass,
		Distance:     requestPlanet.Distance,
		Color1:       requestPlanet.Color1,
		Color2:       requestPlanet.Color2,
		PlanetStatus: model.PLANET_STATUS_ACTIVE,
	}
	err := uc.Repository.CreatePlanet(userID, planet)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) DeletePlanet(planetID, userID uint) error {
	if planetID <= 0 {
		return errors.New("недопустимый ИД планеты")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	err := uc.Repository.DeletePlanet(planetID, userID)
	if err != nil {
		return err
	}
	err = uc.Repository.RemoveServiceImage(planetID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) UpdatePlanet(planetID, userID uint, requestPlanet model.PlanetRequest) error {
	if planetID <= 0 {
		return errors.New("недопустимый ИД планеты")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	planet := model.Planet{
		Name:       requestPlanet.Name,
		Discovered: requestPlanet.Discovered,
		Info:       requestPlanet.Info,
		Mass:       requestPlanet.Mass,
		Distance:   requestPlanet.Distance,
		Color1:     requestPlanet.Color1,
		Color2:     requestPlanet.Color2,
	}
	err := uc.Repository.UpdatePlanet(planetID, userID, planet)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) AddPlanetToConstellation(planetID, userID uint) error {
	if planetID <= 0 {
		return errors.New("недопустимый ИД планеты")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	err := uc.Repository.AddPlanetToConstellation(planetID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) RemovePlanetFromConstellation(planetID, userID uint) error {
	if planetID <= 0 {
		return errors.New("недопустимый ИД планеты")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	err := uc.Repository.RemovePlanetFromConstellation(planetID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) AddPlanetImage(planetID, userID uint, imageBytes []byte, ContentType string) error {
	if planetID <= 0 {
		return errors.New("недопустимый ИД планеты")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if imageBytes == nil {
		return errors.New("недопустимый imageBytes изображения")
	}
	if ContentType == "" {
		return errors.New("недопустимый ContentType изображения")
	}
	imageURL, err := uc.Repository.UploadServiceImage(planetID, userID, imageBytes, ContentType)
	if err != nil {
		return err
	}
	err = uc.Repository.AddPlanetImage(planetID, userID, imageURL)
	if err != nil {
		return err
	}
	return nil
}
