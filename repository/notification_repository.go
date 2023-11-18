package repository

import (
	"calibration-system.com/config"
	"calibration-system.com/model"
	"gorm.io/gorm"
)

type NotificationRepo interface {
	GetCalibratorEmailOnProjectPhase(id string) ([]string, error)
}

type notificationRepo struct {
	db  *gorm.DB
	cfg config.Config
}

func (n *notificationRepo) GetCalibratorEmailOnProjectPhase(id string) ([]string, error) {
	var calibrations []model.Calibration
	var calibratorEmails []string

	err := n.db.
		Preload("Calibrator").
		Where("project_phase_id = ?", id).Find(&calibrations).Error
	if err != nil {
		return []string{}, err
	}

	for _, calibration := range calibrations {
		if calibration.Status == "Calibrate" && calibration.SpmoStatus == "-" {
			calibratorEmails = append(calibratorEmails, calibration.Calibrator.Email)
		}
	}

	uniqueEmail := removeDuplicates(calibratorEmails)
	return uniqueEmail, nil
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{
		db: db,
	}
}
