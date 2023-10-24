package usecase

import (
	"fmt"
	"math"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
)

type ProjectUsecase interface {
	BaseUsecase[model.Project]
	FindPagination(param request.PaginationParam) ([]model.Project, response.Paging, error)
	PublishProject(id string) error
	// FindActiveProject() (*model.Project, error)
	// FindActiveProjectByCalibratorID(calibratorId string) (*response.ProjectCalibrationResponse, error)
	FindScoreDistributionByCalibratorID(businessUnitName string) (*model.Project, error)
	FindRatingQuotaByCalibratorID(calibratorId, prevCalibrator, businessUnitName, types string) (*response.RatingQuota, error)
	FindTotalActualScoreByCalibratorID(calibratorId, prevCalibrator, businessUnitName, types string) (*response.TotalActualScore, error)
	FindSummaryProjectByCalibratorID(calibratorId string) (*response.SummaryProject, error)
	FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) ([]response.UserResponse, error)
	FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) ([]response.UserResponse, error)
	FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnit string) ([]response.UserResponse, error)
	FindCalibratorPhase(calibratorId string) (*model.ProjectPhase, error)
	FindActiveProject() (*model.Project, error)
}

type projectUsecase struct {
	repo repository.ProjectRepo
}

func (r *projectUsecase) FindAll() ([]model.Project, error) {
	return r.repo.List()
}

func (r *projectUsecase) FindPagination(param request.PaginationParam) ([]model.Project, response.Paging, error) {
	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.PaginateList(paginationQuery)
}

func (r *projectUsecase) FindById(id string) (*model.Project, error) {
	return r.repo.Get(id)
}

func (r *projectUsecase) SaveData(payload *model.Project) error {
	return r.repo.Save(payload)
}

func (r *projectUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func (r *projectUsecase) PublishProject(id string) error {
	err := r.repo.ActivateByID(id)
	if err != nil {
		return err
	}

	return r.repo.DeactivateAllExceptID(id)
}

func (r *projectUsecase) FindActiveProject() (*model.Project, error) {
	return r.repo.GetActiveProject()
}

func (r *projectUsecase) FindScoreDistributionByCalibratorID(businessUnitName string) (*model.Project, error) {
	return r.repo.GetScoreDistributionByCalibratorID(businessUnitName)
}

func (r *projectUsecase) FindRatingQuotaByCalibratorID(calibratorId, prevCalibrator, businessUnitName, types string) (*response.RatingQuota, error) {
	var calibrations []response.UserResponse
	var err error
	if types == "numberOne" {
		calibrations, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitName)
	} else if types == "n-1" {
		calibrations, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnitName)
	} else {
		calibrations, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitName)
	}

	if err != nil {
		return nil, err
	}

	projects, err := r.repo.GetRatingQuotaByCalibratorID(businessUnitName)
	if err != nil {
		return nil, err
	}

	ratingQuota := projects.RatingQuotas[0]
	totalCalibrations := len(calibrations)
	// fmt.Println("TOTAL CALIBRATIONS = ", totalCalibrations)
	responses := response.RatingQuota{
		APlus: int(math.Floor(((ratingQuota.APlusQuota) / float64(100)) * float64(totalCalibrations))),
		A:     int(math.Floor(((ratingQuota.AQuota) / float64(100)) * float64(totalCalibrations))),
		BPlus: int(math.Floor(((ratingQuota.BPlusQuota) / float64(100)) * float64(totalCalibrations))),
		B:     int(math.Floor(((ratingQuota.BQuota) / float64(100)) * float64(totalCalibrations))),
		C:     int(math.Floor(((ratingQuota.CQuota) / float64(100)) * float64(totalCalibrations))),
		D:     int(math.Floor(((ratingQuota.DQuota) / float64(100)) * float64(totalCalibrations))),
	}

	var total = responses.APlus + responses.A +
		responses.BPlus + responses.B +
		responses.C + responses.D

	if total < totalCalibrations {
		if ratingQuota.Remaining == "A+" {
			responses.APlus += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "A" {
			responses.A += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "B+" {
			responses.BPlus += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "B" {
			responses.B += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "C" {
			responses.C += (totalCalibrations - total)
		} else {
			responses.D += (totalCalibrations - total)
		}
		total += (totalCalibrations - total)
	}

	if total > totalCalibrations {
		if ratingQuota.Excess == "A+" {
			responses.APlus -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "A" {
			responses.A -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "B+" {
			responses.BPlus -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "B" {
			responses.B -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "C" {
			responses.C -= (total - totalCalibrations)
		} else {
			responses.D -= (total - totalCalibrations)
		}
		total -= (total - totalCalibrations)
	}
	responses.Total = total

	return &responses, nil
}

func (r *projectUsecase) FindTotalActualScoreByCalibratorID(calibratorId, prevCalibrator, businessUnitName, types string) (*response.TotalActualScore, error) {
	var calibrations []response.UserResponse
	var err error
	if types == "numberOne" {
		calibrations, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitName)
	} else if types == "n-1" {
		calibrations, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnitName)
	} else {
		calibrations, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitName)
	}

	if err != nil {
		return nil, err
	}

	totalActualScore := response.TotalActualScore{
		APlus: 0,
		A:     0,
		BPlus: 0,
		B:     0,
		C:     0,
		D:     0,
		Total: 0,
	}

	for _, calibration := range calibrations {
		if calibration.ActualScores[0].ActualRating == "A+" {
			totalActualScore.APlus += 1
		} else if calibration.ActualScores[0].ActualRating == "A" {
			totalActualScore.A += 1
		} else if calibration.ActualScores[0].ActualRating == "B+" {
			totalActualScore.BPlus += 1
		} else if calibration.ActualScores[0].ActualRating == "B" {
			totalActualScore.B += 1
		} else if calibration.ActualScores[0].ActualRating == "C" {
			totalActualScore.C += 1
		} else {
			totalActualScore.D += 1
		}
		totalActualScore.Total += 1

	}

	return &totalActualScore, nil
}

func (r *projectUsecase) FindSummaryProjectByCalibratorID(calibratorId string) (*response.SummaryProject, error) {
	result := &response.SummaryProject{
		Summary: []*response.CalibratorBusinessUnit{},
	}

	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return nil, err
	}

	prevCalibrator := map[string]string{}
	businessUnit := map[string]string{}
	picIDs := map[string]string{}
	users, err := r.repo.GetAllCalibrationByCalibratorID(calibratorId, phase)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		pic := false
		picName := "N-1"
		var picId string
		calibrationLength := len(user.CalibrationScores)
		for _, calibration := range user.CalibrationScores {
			if calibration.ProjectPhase.Phase.Order == phase && calibration.CalibratorID == calibratorId {
				if _, isExist := prevCalibrator[user.Name]; calibrationLength == 1 && isExist {
					picName = user.Name
					picId = user.ID
				} else if name, isExist := businessUnit[user.BusinessUnit.Name]; calibrationLength == 1 && isExist {
					picName = name
					picId = picIDs[user.BusinessUnit.Name]
				}

				pic = true
				break
			} else if calibration.ProjectPhase.Phase.Order == phase && calibration.CalibratorID != calibratorId {
				break
			}

			if calibration.ProjectPhase.Phase.Order < phase {
				prevCalibrator[calibration.Calibrator.Name] = calibration.Calibrator.Name
				picName = calibration.Calibrator.Name
				picId = calibration.CalibratorID
			}
		}

		bu := true
		for _, summary := range result.Summary {
			if summary.CalibratorName == picName && summary.CalibratorBusinessUnit == user.BusinessUnit.Name {
				bu = false
			}
		}

		if _, isExist := businessUnit[user.BusinessUnit.Name]; bu && pic && (picName != "N-1" || !isExist) {
			resp := &response.CalibratorBusinessUnit{
				CalibratorName:           picName,
				CalibratorID:             picId,
				CalibratorBusinessUnit:   user.BusinessUnit.Name,
				CalibratorBusinessUnitID: *user.BusinessUnitId,
				APlus:                    0,
				A:                        0,
				BPlus:                    0,
				B:                        0,
				C:                        0,
				D:                        0,
				APlusGuidance:            0,
				AGuidance:                0,
				BPlusGuidance:            0,
				BGuidance:                0,
				CGuidance:                0,
				DGuidance:                0,
				Status:                   "Complete",
			}

			if user.CalibrationScores[calibrationLength-1].CalibrationRating == "A+" {
				resp.APlus += 1
			} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "A" {
				resp.A += 1
			} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "B+" {
				resp.BPlus += 1
			} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "B" {
				resp.B += 1
			} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "C" {
				resp.C += 1
			} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "D" {
				resp.D += 1
			}

			if user.CalibrationScores[calibrationLength-1].Status != "Complete" {
				resp.Status = "Pending"
			}
			result.Summary = append(result.Summary, resp)
			// fmt.Println("SUMMARY 2A", result.Summary)
		} else {
			for _, summary := range result.Summary {
				if summary.CalibratorName == picName && summary.CalibratorBusinessUnit == user.BusinessUnit.Name {
					if user.CalibrationScores[calibrationLength-1].CalibrationRating == "A+" {
						summary.APlus += 1
					} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "A" {
						summary.A += 1
					} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "B+" {
						summary.BPlus += 1
					} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "B" {
						summary.B += 1
					} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "C" {
						summary.C += 1
					} else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "D" {
						summary.D += 1
					}

					if user.CalibrationScores[calibrationLength-1].Status != "Complete" {
						summary.Status = "Pending"
					}
				}
			}
		}

		if _, isExist := businessUnit[user.BusinessUnit.Name]; !isExist && picName != "N-1" {
			businessUnit[user.BusinessUnit.Name] = picName
			picIDs[user.BusinessUnit.Name] = picId
		}
		// fmt.Println("Business Unit:= ", businessUnit)

		buCheck := map[string]string{}
		for _, summary := range result.Summary {
			types := "default"
			if _, isExist := buCheck[summary.CalibratorBusinessUnit]; !isExist {
				types = "numberOne"
				buCheck[summary.CalibratorBusinessUnit] = summary.CalibratorBusinessUnit
			}

			if summary.CalibratorName == "N-1" {
				types = "n-1"
			}
			guidance, err := r.FindRatingQuotaByCalibratorID(calibratorId, summary.CalibratorID, summary.CalibratorBusinessUnitID, types)
			if err != nil {
				return nil, err
			}

			summary.APlusGuidance = guidance.APlus
			summary.AGuidance = guidance.A
			summary.BPlusGuidance = guidance.BPlus
			summary.BGuidance = guidance.B
			summary.CGuidance = guidance.C
			summary.DGuidance = guidance.D
		}
	}
	return result, nil
}

func (r *projectUsecase) FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) ([]response.UserResponse, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return nil, err
	}

	calibration, err := r.repo.GetCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit, phase)
	if err != nil {
		return nil, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) ([]response.UserResponse, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return nil, err
	}

	users, err := r.repo.GetNumberOneUserWhoCalibrator(calibratorId, businessUnit, phase)
	if err != nil {
		return nil, err
	}

	results, err := r.repo.GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit, phase, users)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *projectUsecase) FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnit string) ([]response.UserResponse, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return nil, err
	}

	calibration, err := r.repo.GetNMinusOneCalibrationsByBusinessUnit(businessUnit, phase)
	if err != nil {
		return nil, err
	}

	for _, data := range calibration {
		fmt.Println(data.Name)
		fmt.Println(data.CalibrationScores)
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibratorPhase(calibratorId string) (*model.ProjectPhase, error) {
	phase, err := r.repo.GetProjectPhase(calibratorId)
	if err != nil {
		return nil, err
	}

	return phase, nil
}

func NewProjectUsecase(repo repository.ProjectRepo) ProjectUsecase {
	return &projectUsecase{
		repo: repo,
	}
}
