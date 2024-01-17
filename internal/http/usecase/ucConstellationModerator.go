package usecase

import (
	"errors"
	"strings"

	"space/internal/model"
)

func (uc *UseCase) GetConstellationsModerator(name, startFormationDate, endFormationDate, constellationStatus string, moderatorID uint) ([]model.ConstellationRequest, error) {
	name = strings.ToUpper(name + "%")
	constellationStatus = strings.ToLower(constellationStatus + "%")
	if moderatorID <= 0 {
		return nil, errors.New("недопустимый ИД модератора")
	}
	constellations, err := uc.Repository.GetConstellationsModerator(name, startFormationDate, endFormationDate, constellationStatus, moderatorID)
	if err != nil {
		return nil, err
	}
	return constellations, nil
}

func (uc *UseCase) GetConstellationByIDModerator(constellationID, moderatorID uint) (model.ConstellationGetResponse, error) {
	if constellationID <= 0 {
		return model.ConstellationGetResponse{}, errors.New("недопустимый ИД созвездия")
	}
	if moderatorID <= 0 {
		return model.ConstellationGetResponse{}, errors.New("недопустимый ИД модератора")
	}
	constellation, err := uc.Repository.GetConstellationByIDModerator(constellationID, moderatorID)
	if err != nil {
		return model.ConstellationGetResponse{}, err
	}
	return constellation, nil
}

/* func (uc *UseCase) UpdateFlightNumberModerator(constellationID, moderatorID uint, constellation model.ConstellationUpdateRequest) error {
	if constellationID <= 0 {
		return errors.New("недопустимый ИД созвездия")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if len(constellation.Name) != 6 {
		return errors.New("недопустимый номер рейса")
	}
	err := uc.Repository.UpdateConstellationModerator(constellationID, moderatorID, constellation)
	if err != nil {
		return err
	}
	return nil
} */

func (uc *UseCase) UpdateConstellationStatusModerator(constellationID, moderatorID uint, constellationStatus model.ConstellationUpdateStatusRequest) error {
	if constellationID <= 0 {
		return errors.New("недопустимый ИД созвездия")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if constellationStatus.ConstellationStatus == model.CONSTELLATION_STATUS_WORK {
		err := uc.Repository.UpdateConstellationStatusUser(constellationID, moderatorID)
		if err != nil {
			return err
		}
		return nil
	} else if constellationStatus.ConstellationStatus == model.CONSTELLATION_STATUS_COMPLETED || constellationStatus.ConstellationStatus == model.CONSTELLATION_STATUS_REJECTED {
		constellation, err := uc.Repository.GetConstellationByIDModerator(constellationID, moderatorID)
		if err != nil {
			return err
		}
		if constellation.ConstellationStatus != model.CONSTELLATION_STATUS_WORK {
			return errors.New("текущее cозвездие не черновое. Ошибка 234549")
		}
		err = uc.Repository.UpdateConstellationStatusModerator(constellationID, moderatorID, constellationStatus)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("текущий статус созвездия уже завершен или отклонен. Ошибка 309234")
}
