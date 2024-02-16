package usecase

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type ProjectUsecase interface {
	BaseUsecase[model.Project]
	FindPagination(param request.PaginationParam) ([]model.Project, response.Paging, error)
	PublishProject(id string) error
	DeactivateProject(id string) error
	// FindActiveProject() (*model.Project, error)
	// FindActiveProjectByCalibratorID(calibratorId string) (*response.ProjectCalibrationResponse, error)
	FindScoreDistributionByCalibratorID(businessUnitName string) (*model.Project, error)
	FindRatingQuotaByCalibratorID(calibratorId, prevCalibrator, businessUnitID, types string) (*response.RatingQuota, error)
	FindTotalActualScoreByCalibratorID(calibratorId, prevCalibrator, businessUnitName, types string) (*response.TotalActualScore, error)
	FindSummaryProjectByCalibratorID(calibratorId string) (*response.SummaryProject, error)
	FindCalibrationsByBusinessUnit(calibratorId, businessUnit string) (response.UserCalibration, error)
	FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) (response.UserCalibration, error)
	FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) (response.UserCalibration, error)
	FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnit string) (response.UserCalibration, error)
	FindCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorId, prevCalibrator, businessUnit, rating string) (response.UserCalibration, error)
	FindCalibrationsByBusinessUnitAndRating(calibratorId, prevCalibrator, rating string) (response.UserCalibration, error)
	FindCalibrationsByRating(calibratorId, rating string) (response.UserCalibration, error)
	FindCalibratorPhase(calibratorId string) (*model.ProjectPhase, error)
	FindActiveProjectPhase() ([]model.ProjectPhase, error)
	FindActiveManagerPhase() (model.ProjectPhase, error)
	FindActiveProject() (*model.Project, error)
	FindProjectRatingQuotaByBusinessUnit(businessUnitID string) (*model.Project, error)
	FindSummaryProjectTotalByCalibratorID(calibratorId string) (*response.SummaryTotal, error)
	ReportCalibrations(types, calibratorId, businessUnit, prevCalibrator string, c *gin.Context) (string, error)
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

func (r *projectUsecase) DeactivateProject(id string) error {
	return r.repo.NonactivateByID(id)
}

func (r *projectUsecase) FindActiveProject() (*model.Project, error) {
	return r.repo.GetActiveProject()
}

func (r *projectUsecase) FindScoreDistributionByCalibratorID(businessUnitName string) (*model.Project, error) {
	return r.repo.GetScoreDistributionByCalibratorID(businessUnitName)
}

func (r *projectUsecase) FindRatingQuotaByCalibratorID(calibratorId, prevCalibrator, businessUnitID, types string) (*response.RatingQuota, error) {
	var calibrations response.UserCalibration
	var err error
	if types == "numberOne" {
		calibrations, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitID)
	} else if types == "n-1" {
		calibrations, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnitID)
	} else if types == "default" {
		calibrations, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitID)
	} else {
		calibrations, err = r.FindCalibrationsByBusinessUnit(calibratorId, businessUnitID)
	}

	if err != nil {
		return nil, err
	}

	projects, err := r.repo.GetRatingQuotaByCalibratorID(businessUnitID)
	if err != nil {
		return nil, err
	}

	ratingQuota := projects.RatingQuotas[0]
	totalCalibrations := len(calibrations.UserData)
	// fmt.Println("TOTAL CALIBRATIONS = ", totalCalibrations)
	// Jangan di round down
	responses := response.RatingQuota{
		APlus: int(math.Round(((ratingQuota.APlusQuota) / float64(100)) * float64(totalCalibrations))),
		A:     int(math.Round(((ratingQuota.AQuota) / float64(100)) * float64(totalCalibrations))),
		BPlus: int(math.Round(((ratingQuota.BPlusQuota) / float64(100)) * float64(totalCalibrations))),
		B:     int(math.Round(((ratingQuota.BQuota) / float64(100)) * float64(totalCalibrations))),
		C:     int(math.Round(((ratingQuota.CQuota) / float64(100)) * float64(totalCalibrations))),
		D:     int(math.Round(((ratingQuota.DQuota) / float64(100)) * float64(totalCalibrations))),
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
			if responses.APlus-(total-totalCalibrations) > 0 {
				responses.APlus -= (total - totalCalibrations)
			} else {
				responses.BPlus -= (total - totalCalibrations)
			}
		} else if ratingQuota.Excess == "A" {
			if responses.A-(total-totalCalibrations) > 0 {
				responses.A -= (total - totalCalibrations)
			} else {
				responses.BPlus -= (total - totalCalibrations)
			}
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
	var calibrations response.UserCalibration
	var err error
	if types == "numberOne" {
		calibrations, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitName)
	} else if types == "n-1" {
		calibrations, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnitName)
	} else if types == "default" {
		calibrations, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnitName)
	} else {
		calibrations, err = r.FindCalibrationsByBusinessUnit(calibratorId, businessUnitName)
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

	for _, calibration := range calibrations.UserData {
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
		Summary:           []*response.BusinessUnitTotal{},
		APlusTotalScore:   0,
		ATotalScore:       0,
		BPlusTotalScore:   0,
		BTotalScore:       0,
		CTotalScore:       0,
		DTotalScore:       0,
		APlusGuidance:     0,
		AGuidance:         0,
		BPlusGuidance:     0,
		BGuidance:         0,
		CGuidance:         0,
		DGuidance:         0,
		AverageTotalScore: 0,
	}

	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return nil, err
	}

	prevCalibrator := map[string]string{}
	businessUnit := map[string]string{}
	picIDs := map[string]string{}
	resultSummary := map[string]*response.CalibratorBusinessUnit{}
	users, err := r.repo.GetAllUserCalibrationByCalibratorID(calibratorId, phase)
	if err != nil {
		return nil, err
	}

	totalUsers := 0
	countCalibratedScoresUsers := 0.0

	for _, user := range users {
		if user.ScoringMethod == "Score" {
			totalUsers += 1
			countCalibratedScoresUsers += user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationScore
		}

		// Grouping By Previous Calibrator and Business Unit
		pic := false
		picName := "N-1"
		picId := "N-1"
		calibrationLength := len(user.CalibrationScores)
		for _, calibration := range user.CalibrationScores {
			if calibration.ProjectPhase.Phase.Order == phase && calibration.CalibratorID == calibratorId {
				if _, isExist := prevCalibrator[user.Name]; calibrationLength == 1 && isExist {
					picName = user.Name
					picId = user.ID
				}
				// else if name, isExist := businessUnit[user.BusinessUnit.Name]; calibrationLength == 1 && isExist {
				// 	picName = name
				// 	picId = picIDs[user.BusinessUnit.Name]
				// }

				pic = true
				break
			} else if calibration.ProjectPhase.Phase.Order >= phase && calibration.CalibratorID != calibratorId {
				break
			}

			if calibration.ProjectPhase.Phase.Order < phase {
				prevCalibrator[calibration.Calibrator.Name] = calibration.Calibrator.Name
				picName = calibration.Calibrator.Name
				picId = calibration.CalibratorID
			}
		}

		bu := true
		if _, ok := resultSummary[picName+user.BusinessUnit.Name]; ok {
			bu = false
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
				TotalCalibratedScore:     0,
				UserCount:                0.0,
				AverageScore:             0,
			}

			if user.ScoringMethod == "Score" {
				resp.UserCount += 1
				resp.TotalCalibratedScore += user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationScore
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

			if user.CalibrationScores[calibrationLength-1].Status != "Complete" || user.CalibrationScores[calibrationLength-1].SpmoStatus == "Rejected" {
				resp.Status = "Calibrate"
			}
			if user.CalibrationScores[calibrationLength-1].Status == "Waiting" {
				resp.Status = "Waiting"
			}

			resultSummary[picName+user.BusinessUnit.Name] = resp
			// result.Summary = append(result.Summary, resp)
			// fmt.Println("SUMMARY 2A", result.Summary)
		} else {
			if summary, ok := resultSummary[picName+user.BusinessUnit.Name]; ok {
				if user.ScoringMethod == "Score" {
					summary.UserCount += 1
					summary.TotalCalibratedScore += user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationScore
				}

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

				if user.CalibrationScores[calibrationLength-1].Status == "Calibrate" || user.CalibrationScores[calibrationLength-1].SpmoStatus == "Rejected" {
					if summary.Status != "Waiting" {
						summary.Status = "Calibrate"
					}
				}

				if user.CalibrationScores[calibrationLength-1].Status == "Waiting" {
					summary.Status = "Waiting"
				}
			} else {
				if picName == "N-1" {
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
						TotalCalibratedScore:     0,
						UserCount:                0.0,
						AverageScore:             0,
					}

					if user.ScoringMethod == "Score" {
						resp.UserCount += 1
						resp.TotalCalibratedScore += user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationScore
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

					if user.CalibrationScores[calibrationLength-1].Status != "Complete" || user.CalibrationScores[calibrationLength-1].SpmoStatus == "Rejected" {
						resp.Status = "Calibrate"
					}
					if user.CalibrationScores[calibrationLength-1].Status == "Waiting" {
						resp.Status = "Waiting"
					}

					resultSummary[picName+user.BusinessUnit.Name] = resp
				}
			}
		}

		if _, isExist := businessUnit[user.BusinessUnit.Name]; !isExist {
			businessUnit[user.BusinessUnit.Name] = user.BusinessUnit.Name
			picIDs[user.BusinessUnit.Name] = picId
		}

		// if user.CalibrationScores[calibrationLength-1].CalibrationRating == "A+" {
		// 	businessUnit[user.BusinessUnit.Name].APlusCalibrated += 1
		// } else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "A" {
		// 	businessUnit[user.BusinessUnit.Name].ACalibrated += 1
		// } else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "B+" {
		// 	businessUnit[user.BusinessUnit.Name].BPlusCalibrated += 1
		// } else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "B" {
		// 	businessUnit[user.BusinessUnit.Name].BCalibrated += 1
		// } else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "C" {
		// 	businessUnit[user.BusinessUnit.Name].CCalibrated += 1
		// } else if user.CalibrationScores[calibrationLength-1].CalibrationRating == "D" {
		// 	businessUnit[user.BusinessUnit.Name].DCalibrated += 1
		// }
		// fmt.Println("Business Unit:= ", businessUnit)

	}

	buCheck := map[string]string{}
	finalData := map[string]*response.BusinessUnitTotal{}
	for _, summary := range resultSummary {
		types := "all"
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

		if _, isExist := buCheck[summary.CalibratorBusinessUnit]; !isExist {
			// types = "numberOne"
			buCheck[summary.CalibratorBusinessUnit] = summary.CalibratorBusinessUnit
			finalData[summary.CalibratorBusinessUnit] = &response.BusinessUnitTotal{
				CalibratorBusinessUnit:     []*response.CalibratorBusinessUnit{},
				CalibratorBusinessUnitName: summary.CalibratorBusinessUnit,
				CalibratorBusinessUnitID:   summary.CalibratorBusinessUnitID,
				APlusCalibrated:            0,
				ACalibrated:                0,
				BPlusCalibrated:            0,
				BCalibrated:                0,
				CCalibrated:                0,
				DCalibrated:                0,
				APlusGuidance:              0,
				AGuidance:                  0,
				BPlusGuidance:              0,
				BGuidance:                  0,
				CGuidance:                  0,
				DGuidance:                  0,
				TotalCalibratedScore:       0,
				UserCount:                  0,
				AverageScore:               0,
				Status:                     "Calibrate",
				Completed:                  true,
			}

			finalData[summary.CalibratorBusinessUnit].APlusGuidance = guidance.APlus
			finalData[summary.CalibratorBusinessUnit].AGuidance = guidance.A
			finalData[summary.CalibratorBusinessUnit].BPlusGuidance = guidance.BPlus
			finalData[summary.CalibratorBusinessUnit].BGuidance = guidance.B
			finalData[summary.CalibratorBusinessUnit].CGuidance = guidance.C
			finalData[summary.CalibratorBusinessUnit].DGuidance = guidance.D

			result.APlusGuidance += guidance.APlus
			result.AGuidance += guidance.A
			result.BPlusGuidance += guidance.BPlus
			result.BGuidance += guidance.B
			result.CGuidance += guidance.C
			result.DGuidance += guidance.D
		}

		if summary.UserCount > 0 {
			summary.AverageScore = summary.TotalCalibratedScore / float64(summary.UserCount)
		}

		finalData[summary.CalibratorBusinessUnit].APlusCalibrated += summary.APlus
		finalData[summary.CalibratorBusinessUnit].ACalibrated += summary.A
		finalData[summary.CalibratorBusinessUnit].BPlusCalibrated += summary.BPlus
		finalData[summary.CalibratorBusinessUnit].BCalibrated += summary.B
		finalData[summary.CalibratorBusinessUnit].CCalibrated += summary.C
		finalData[summary.CalibratorBusinessUnit].DCalibrated += summary.D
		finalData[summary.CalibratorBusinessUnit].TotalCalibratedScore += summary.TotalCalibratedScore
		finalData[summary.CalibratorBusinessUnit].UserCount += summary.UserCount
		finalData[summary.CalibratorBusinessUnit].CalibratorBusinessUnit = append(finalData[summary.CalibratorBusinessUnit].CalibratorBusinessUnit, summary)

		if finalData[summary.CalibratorBusinessUnit].UserCount > 0 {
			finalData[summary.CalibratorBusinessUnit].AverageScore = finalData[summary.CalibratorBusinessUnit].TotalCalibratedScore / float64(finalData[summary.CalibratorBusinessUnit].UserCount)
		}

		result.APlusTotalScore += summary.APlus
		result.ATotalScore += summary.A
		result.BPlusTotalScore += summary.BPlus
		result.BTotalScore += summary.B
		result.CTotalScore += summary.C
		result.DTotalScore += summary.D

		if summary.Status == "Waiting" {
			finalData[summary.CalibratorBusinessUnit].Status = "Waiting"
		}

		if summary.Status == "Complete" {
			finalData[summary.CalibratorBusinessUnit].Completed = finalData[summary.CalibratorBusinessUnit].Completed && true
		} else {
			finalData[summary.CalibratorBusinessUnit].Completed = finalData[summary.CalibratorBusinessUnit].Completed && false
		}
	}

	if totalUsers > 0 {
		result.AverageTotalScore = countCalibratedScoresUsers / float64(totalUsers)
	}

	for _, rSummary := range finalData {
		if rSummary.Completed == true {
			rSummary.Status = "Complete"
		}

		result.Summary = append(result.Summary, rSummary)
		// fmt.Println("ISI BU UNIT", key)
	}

	sort.Slice(result.Summary, func(i, j int) bool {
		if result.Summary[i].CalibratorBusinessUnitName != result.Summary[j].CalibratorBusinessUnitName {
			return result.Summary[i].CalibratorBusinessUnitName < result.Summary[j].CalibratorBusinessUnitName
		}
		return result.Summary[i].CalibratorBusinessUnitName < result.Summary[j].CalibratorBusinessUnitName
	})

	for _, businessUnit := range result.Summary {
		sort.Slice(businessUnit.CalibratorBusinessUnit, func(i, j int) bool {
			return businessUnit.CalibratorBusinessUnit[i].CalibratorName < businessUnit.CalibratorBusinessUnit[j].CalibratorName
		})
	}

	return result, nil
}

func (r *projectUsecase) FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibrationsByBusinessUnit(calibratorId, businessUnit string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByBusinessUnit(calibratorId, businessUnit, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	users, err := r.repo.GetNumberOneUserWhoCalibrator(calibratorId, businessUnit, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}

	results, err := r.repo.GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit, phase, users)
	if err != nil {
		return response.UserCalibration{}, err
	}

	return results, nil
}

func (r *projectUsecase) FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnit string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetNMinusOneCalibrationsByBusinessUnit(businessUnit, phase, calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorId, prevCalibrator, businessUnit, rating string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorId, prevCalibrator, businessUnit, rating, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibrationsByBusinessUnitAndRating(calibratorId, businessUnit, rating string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByBusinessUnitAndRating(calibratorId, businessUnit, rating, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibrationsByRating(calibratorId, rating string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByRating(calibratorId, rating, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindSummaryProjectTotalByCalibratorID(calibratorId string) (*response.SummaryTotal, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorId)
	if err != nil {
		return nil, err
	}

	allBusinessUnit, err := r.repo.GetAllBusinessUnitSummary(calibratorId, phase)
	if err != nil {
		return nil, err
	}

	results := &response.SummaryTotal{
		Data: []*response.RatingDataSummary{},
	}

	mapRating := map[string]*response.RatingDataSummary{}
	mapRating["A+"] = &response.RatingDataSummary{
		Rating: "A+",
	}
	mapRating["A"] = &response.RatingDataSummary{
		Rating: "A",
	}
	mapRating["B+"] = &response.RatingDataSummary{
		Rating: "B+",
	}
	mapRating["B"] = &response.RatingDataSummary{
		Rating: "B",
	}
	mapRating["C"] = &response.RatingDataSummary{
		Rating: "C",
	}
	mapRating["D"] = &response.RatingDataSummary{
		Rating: "D",
	}

	for _, businessUnit := range allBusinessUnit {
		ratingQuota, err := r.FindRatingQuotaByCalibratorID(calibratorId, "", businessUnit.ID, "all")
		if err != nil {
			return nil, err
		}

		mapRating["A+"].Guidance += ratingQuota.APlus
		mapRating["A"].Guidance += ratingQuota.A
		mapRating["B+"].Guidance += ratingQuota.BPlus
		mapRating["B"].Guidance += ratingQuota.B
		mapRating["C"].Guidance += ratingQuota.C
		mapRating["D"].Guidance += ratingQuota.D

		users, err := r.FindCalibrationsByBusinessUnit(calibratorId, businessUnit.ID)
		for _, user := range users.UserData {
			if user.ActualScores[0].ActualRating == "A+" {
				mapRating["A+"].ActualRating += 1
			} else if user.ActualScores[0].ActualRating == "A" {
				mapRating["A"].ActualRating += 1
			} else if user.ActualScores[0].ActualRating == "B+" {
				mapRating["B+"].ActualRating += 1
			} else if user.ActualScores[0].ActualRating == "B" {
				mapRating["B"].ActualRating += 1
			} else if user.ActualScores[0].ActualRating == "C" {
				mapRating["C"].ActualRating += 1
			} else {
				mapRating["D"].ActualRating += 1
			}

			if user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating == "A+" {
				mapRating["A+"].CalibratedRating += 1
			} else if user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating == "A" {
				mapRating["A"].CalibratedRating += 1
			} else if user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating == "B+" {
				mapRating["B+"].CalibratedRating += 1
			} else if user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating == "B" {
				mapRating["B"].CalibratedRating += 1
			} else if user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating == "C" {
				mapRating["C"].CalibratedRating += 1
			} else {
				mapRating["D"].CalibratedRating += 1
			}
		}

	}

	for _, mRating := range mapRating {
		mRating.Variance = mRating.CalibratedRating - mRating.Guidance
		results.Data = append(results.Data, mRating)
	}

	return results, nil
}

func (r *projectUsecase) FindCalibratorPhase(calibratorId string) (*model.ProjectPhase, error) {
	phase, err := r.repo.GetProjectPhase(calibratorId)
	if err != nil {
		return nil, err
	}

	return phase, nil
}

func (r *projectUsecase) FindActiveProjectPhase() ([]model.ProjectPhase, error) {
	projectPhase, err := r.repo.GetActiveProjectPhase()
	if err != nil {
		return nil, err
	}

	return projectPhase, nil
}

func (r *projectUsecase) FindActiveManagerPhase() (model.ProjectPhase, error) {
	projectPhase, err := r.repo.GetActiveManagerPhase()
	if err != nil {
		return model.ProjectPhase{}, err
	}

	return projectPhase, nil
}

func (r *projectUsecase) FindProjectRatingQuotaByBusinessUnit(businessUnitID string) (*model.Project, error) {
	projects, err := r.repo.GetRatingQuotaByCalibratorID(businessUnitID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *projectUsecase) ReportCalibrations(types, calibratorId, businessUnit, prevCalibrator string, c *gin.Context) (string, error) {
	var responseData response.UserCalibration
	var err error

	if types == "numberOne" {
		responseData, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit)
	} else if types == "n-1" {
		responseData, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, businessUnit)
	} else if types == "default" {
		responseData, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit)
	} else {
		responseData, err = r.FindCalibrationsByBusinessUnit(calibratorId, businessUnit)
	}

	projectPhase, err := r.FindCalibratorPhase(calibratorId)
	if err != nil {
		return "", err
	}

	file := excelize.NewFile()
	index := file.NewSheet(responseData.UserData[0].BusinessUnit.Name)

	headers := []string{
		"No",
		"Employee",
		"Grade",
		"BU",
		"OU",
		"Supervisor",
		"PrevRating Y-2",
		"PrevRating Y-1",
		"PTT Score",
		"PAT Score",
		"360 Score",
		"Actual Score",
		"Actual Rating",
	}

	colorPallete := []string{
		"#F2C4DE", "#71B1D9", "#AED8F2", "#ABD3DB", "#C2E6DF", "#D1EBD8", "#E5F5DC", "#F2DEA2", "#FFFFE1", "#F2CDC4",
	}

	for i := 0; i < projectPhase.Phase.Order; i++ {
		headers = append(headers, fmt.Sprintf("Calibration-%d Score %s", i+1, colorPallete[i]))
		headers = append(headers, fmt.Sprintf("Calibration-%d Rating %s", i+1, colorPallete[i]))
		headers = append(headers, fmt.Sprintf("JustificationType-%d %s", i+1, colorPallete[i]))
		headers = append(headers, fmt.Sprintf("Justification-%d %s", i+1, colorPallete[i]))
	}

	style1, err := file.NewStyle(`{"alignment":{"wrap_text":true, "vertical":"center"}, "font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	style2, err := file.NewStyle(`{"alignment":{"wrap_text":true, "vertical":"center"}}`)
	if err != nil {
		return "", err
	}

	firstAppear := true
	for col, header := range headers {
		colName := excelize.ToAlphaString(col)
		cellRef := fmt.Sprintf("%s%d", colName, 1)
		cellRefUnder := fmt.Sprintf("%s%d", colName, 2)
		nextColName := excelize.ToAlphaString(col + 1)
		cellNextCollRef := fmt.Sprintf("%s%d", nextColName, 1)
		file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, style1)

		if (strings.Contains(header, "Prev") || strings.Contains(header, "Actual") || strings.Contains(header, "Calibration")) && firstAppear {
			file.MergeCell(responseData.UserData[0].BusinessUnit.Name, cellRef, cellNextCollRef)
			words := strings.Fields(header)
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRef, words[0])
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, words[1])
			if strings.Contains(header, "Calibration") {
				backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"wrap_text":true, "vertical":"center"}, "fill":{"type":"pattern", "color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[2]))
				if err != nil {
					return "", err
				}
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, backgroundColorCalibration)
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, cellRefUnder, backgroundColorCalibration)
			} else {
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, style1)
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, cellRefUnder, style1)
			}
			firstAppear = false
		} else if (strings.Contains(header, "Prev") || strings.Contains(header, "Actual") || strings.Contains(header, "Calibration")) && !firstAppear {
			words := strings.Fields(header)
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, words[1])
			if strings.Contains(header, "Calibration") {
				backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"wrap_text":true, "vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[2]))
				if err != nil {
					return "", err
				}
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, backgroundColorCalibration)
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, cellRefUnder, backgroundColorCalibration)
			} else {
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, style1)
				file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, cellRefUnder, style1)
			}
			firstAppear = true
		} else if strings.Contains(header, "Justification") {
			file.MergeCell(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRefUnder)
			words := strings.Fields(header)
			backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"wrap_text":true, "vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[1]))
			if err != nil {
				return "", err
			}
			file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, backgroundColorCalibration)
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRef, words[0])
		} else {
			file.MergeCell(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRefUnder)
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRef, header)
		}
	}

	for i, user := range responseData.UserData {
		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("A%d", i+3), i+1)
		file.SetColWidth(user.BusinessUnit.Name, "A", "A", 4)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("A%d", i+3), fmt.Sprintf("A%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("B%d", i+3), user.Name)
		file.SetColWidth(user.BusinessUnit.Name, "B", "B", 20)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("B%d", i+3), fmt.Sprintf("B%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("C%d", i+3), user.Grade)
		file.SetColWidth(user.BusinessUnit.Name, "C", "C", 4)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("C%d", i+3), fmt.Sprintf("C%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("D%d", i+3), user.BusinessUnit.Name)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("D%d", i+3), fmt.Sprintf("D%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("E%d", i+3), user.OrganizationUnit)
		file.SetColWidth(user.BusinessUnit.Name, "E", "E", 15)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("E%d", i+3), fmt.Sprintf("E%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("F%d", i+3), user.SupervisorNames)
		file.SetColWidth(user.BusinessUnit.Name, "F", "F", 20)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("F%d", i+3), fmt.Sprintf("F%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("G%d", i+3), user.ActualScores[0].Y2Rating)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("G%d", i+3), fmt.Sprintf("G%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("H%d", i+3), user.ActualScores[0].Y1Rating)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("H%d", i+3), fmt.Sprintf("H%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("I%d", i+3), user.ActualScores[0].PTTScore)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("I%d", i+3), fmt.Sprintf("I%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+3), user.ActualScores[0].PATScore)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+3), fmt.Sprintf("J%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+3), user.ActualScores[0].Score360)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+3), fmt.Sprintf("K%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+3), user.ActualScores[0].ActualScore)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+3), fmt.Sprintf("L%d", i+3), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+3), user.ActualScores[0].ActualRating)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+3), fmt.Sprintf("M%d", i+3), style2)
		column := int('N')

		startCalibrationPhaseCol := 13
		for j := 1; j < projectPhase.Phase.Order+1; j++ {
			columnBefore := column
			words := strings.Fields(headers[startCalibrationPhaseCol])
			startCalibrationPhaseCol += 4
			backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"wrap_text":true, "vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}}`, words[2]))
			if err != nil {
				return "", err
			}

			for _, calibrationScore := range user.CalibrationScores {
				if j == calibrationScore.ProjectPhase.Phase.Order {

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), calibrationScore.CalibrationScore)
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), calibrationScore.CalibrationRating)
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), calibrationScore.JustificationType)
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
					column++

					if calibrationScore.JustificationType == "default" {
						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), calibrationScore.Comment)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
					} else if calibrationScore.JustificationType == "top" {
						var justifications string
						for _, topJustification := range calibrationScore.TopRemarks {
							justifications = justifications + fmt.Sprintf("Initiative: %s\nDescription: %s\nResult: %s\nComment: %s\nDate: %s - %s\nEvidence: %s",
								topJustification.Initiative,
								topJustification.Description,
								topJustification.Result,
								topJustification.Comment,
								topJustification.StartDate,
								topJustification.EndDate,
								topJustification.EvidenceLink,
							)
						}
						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), justifications)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
						file.SetColWidth(user.BusinessUnit.Name, fmt.Sprintf("%s", asciiToName(column)), fmt.Sprintf("%s", asciiToName(column)), 25)
					} else {
						justification := fmt.Sprintf("Attitude: %s\nIndisipliner: %s\nLow Performance: %s\nWarning Letter: %s",
							calibrationScore.BottomRemark.Attitude,
							calibrationScore.BottomRemark.Indisipliner,
							calibrationScore.BottomRemark.LowPerformance,
							calibrationScore.BottomRemark.WarningLetter,
						)
						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), justification)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
						file.SetColWidth(user.BusinessUnit.Name, fmt.Sprintf("%s", asciiToName(column)), fmt.Sprintf("%s", asciiToName(column)), 25)
					}
					column++
				}
			}
			if columnBefore == column {
				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
				column++

				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
				column++

				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
				column++

				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+3), fmt.Sprintf("%s%d", asciiToName(column), i+3), backgroundColorCalibration)
				column++
			}
		}
	}

	file.SetActiveSheet(index)

	err = file.SaveAs("report.xlsx")
	if err != nil {
		return "", err
	}

	return "report.xlsx", nil
}

func asciiToName(column int) string {
	columnName := fmt.Sprintf("%c", column)
	if column > int('Z') {
		offset := (column - int('A')) % 26
		firstLetterOffset := int((column - int('A')) / 26)
		firstLetter := int('A') - 1 + firstLetterOffset
		columnName = fmt.Sprintf("%c%c", firstLetter, offset+int('A'))
	}
	return columnName
}

func NewProjectUsecase(repo repository.ProjectRepo) ProjectUsecase {
	return &projectUsecase{
		repo: repo,
	}
}
