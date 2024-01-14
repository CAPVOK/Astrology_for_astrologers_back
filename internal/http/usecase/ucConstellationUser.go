package usecase

import (
	"errors"
	"strings"

	"space/internal/model"
)

func (uc *UseCase) GetConstellationsUser(name, startFormationDate, endFormationDate, constellationStatus string, userID uint) ([]model.ConstellationRequest, error) {
	name = strings.ToUpper(name + "%")
	constellationStatus = strings.ToLower(constellationStatus + "%")
	if userID <= 0 {
		return nil, errors.New("недопустимый ИД пользователя")
	}
	constellations, err := uc.Repository.GetConstellationsUser(name, startFormationDate, endFormationDate, constellationStatus, userID)
	if err != nil {
		return nil, err
	}
	return constellations, nil
}

func (uc *UseCase) GetConstellationByIDUser(constellationID, userID uint) (model.ConstellationGetResponse, error) {
	if constellationID <= 0 {
		return model.ConstellationGetResponse{}, errors.New("недопустимый ИД созвездия")
	}
	if userID <= 0 {
		return model.ConstellationGetResponse{}, errors.New("недопустимый ИД пользователя")
	}
	constellations, err := uc.Repository.GetConstellationByIDUser(constellationID, userID)
	if err != nil {
		return model.ConstellationGetResponse{}, err
	}
	return constellations, nil
}

func (uc *UseCase) DeleteConstellationUser(constellationID, userID uint) error {
	if constellationID <= 0 {
		return errors.New("недопустимый ИД созвездия")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	err := uc.Repository.DeleteConstellationUser(constellationID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) UpdateConstellationUser(constellationID, userID uint, constellation model.ConstellationUpdateRequest) error {
	if constellationID <= 0 {
		return errors.New("недопустимый ИД созвездия")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	err := uc.Repository.UpdateConstellationUser(constellationID, userID, constellation)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) UpdateConstellationStatusUser(userID uint) (uint, error) {
	if userID <= 0 {
		return 0, errors.New("недопустимый ИД пользователя")
	}
	data, err := uc.Repository.GetPlanets("", userID)
	if err != nil {
		return 0, err
	}
	constellation, err := uc.Repository.GetConstellationByIDUser(data.ConstellationID, userID)
	if err != nil {
		return 0, err
	}
	if constellation.ConstellationID == 0 {
		return 0, errors.New("чернового созведия не существует")
	}
	if len(constellation.Planets) == 0 {
		return 0, errors.New("пустое создездие")
	}
	err = uc.Repository.UpdateConstellationStatusUser(data.ConstellationID, userID)
	if err != nil {
		return 0, err
	}
	return constellation.ConstellationID, nil
}
