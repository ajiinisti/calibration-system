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
	// FindActiveProjectByCalibratorID(calibratorID string) (*response.ProjectCalibrationResponse, error)
	FindScoreDistributionByCalibratorID(businessUnitName, projectID string) (*model.Project, error)
	FindRatingQuotaByCalibratorID(calibratorID, prevCalibrator, businessUnitID, types, projectID string, countCurrentUser int) (*response.RatingQuota, error)
	FindTotalActualScoreByCalibratorID(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (*response.TotalActualScore, error)
	FindTotalCalibratedByCalibratorID(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (*response.TotalCalibratedRating, error)
	FindAverageScoreByCalibratorID(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (float32, error)
	FindAllEmployeeName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error)
	FindAllSupervisorName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error)
	FindAllGrade(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error)
	FindSummaryProjectByCalibratorID(calibratorID, projectID string, prevCalibratorIDs []string) (*response.SummaryProject, error)
	FindCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error)
	FindCalibrationsByBusinessUnitPaginate(calibratorID, businessUnit, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error)
	FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string) (response.UserCalibration, error)
	FindCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, prevCalibrator, businessUnit, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error)
	FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string) (response.UserCalibration, error)
	FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error)
	FindNMinusOneCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, businessUnit, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error)
	FindCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error)
	FindCalibrationsByBusinessUnitAndRating(calibratorID, prevCalibrator, rating, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error)
	FindCalibrationsByRating(calibratorID, rating, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error)
	FindCalibratorPhase(calibratorID, projectID string) (*model.ProjectPhase, error)
	FindActiveProjectPhase(projectID string) ([]model.ProjectPhase, error)
	FindActiveManagerPhase() (model.ProjectPhase, error)
	FindActiveProject() ([]model.Project, error)
	FindProjectRatingQuotaByBusinessUnit(businessUnitID, projectID string) (*model.Project, error)
	FindSummaryProjectTotalByCalibratorID(calibratorID, projectID string) (*response.SummaryTotal, error)
	ReportCalibrations(types, calibratorID, businessUnit, prevCalibrator, projectID string, c *gin.Context) (string, error)
	SummaryReportCalibrations(calibratorID, projectID string, c *gin.Context) (string, error)
	FindRatingQuotaByCalibratorIDforSummaryHelper(calibratorID, prevCalibrator, businessUnitID, types, projectID string, countCurrentUser int) (*response.RatingQuota, error)
	FindActiveProjectByCalibratorID(calibratorID string) ([]model.Project, error)
	FindActiveProjectBySpmoID(spmoID string) ([]model.Project, error)
	FindReportCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error)
	FindReportNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error)
	FindReportCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string) (response.UserCalibration, error)
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
	return r.repo.ActivateByID(id)
}

func (r *projectUsecase) DeactivateProject(id string) error {
	return r.repo.NonactivateByID(id)
}

func (r *projectUsecase) FindActiveProject() ([]model.Project, error) {
	return r.repo.GetActiveProject()
}

func (r *projectUsecase) FindScoreDistributionByCalibratorID(businessUnitName, projectID string) (*model.Project, error) {
	return r.repo.GetScoreDistributionByCalibratorID(businessUnitName, projectID)
}

func (r *projectUsecase) FindRatingQuotaByCalibratorID(calibratorID, prevCalibrator, businessUnitID, types, projectID string, countCurrentUser int) (*response.RatingQuota, error) {
	var calibrations response.UserCalibration
	var err error
	if types == "numberOne" {
		calibrations, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnitID, projectID)
	} else if types == "n-1" {
		calibrations, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnitID, projectID)
	} else if types == "default" {
		calibrations, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnitID, projectID)
	} else {
		calibrations, err = r.FindCalibrationsByBusinessUnit(calibratorID, businessUnitID, projectID)
	}

	if err != nil {
		return nil, err
	}

	projects, err := r.repo.GetRatingQuotaByCalibratorID(businessUnitID, projectID)
	if err != nil {
		return nil, err
	}

	ratingQuota := projects.RatingQuotas[0]
	totalCalibrations := len(calibrations.UserData)
	if countCurrentUser > 0 {
		totalCalibrations = countCurrentUser
	}
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

func (r *projectUsecase) FindTotalActualScoreByCalibratorID(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (*response.TotalActualScore, error) {
	var calibrations response.UserCalibration
	var err error
	if types == "numberOne" {
		calibrations, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnitName, projectID)
	} else if types == "n-1" {
		calibrations, err = r.FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnitName, projectID)
	} else if types == "default" {
		calibrations, err = r.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnitName, projectID)
	} else {
		calibrations, err = r.FindCalibrationsByBusinessUnit(calibratorID, businessUnitName, projectID)
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

func (r *projectUsecase) FindTotalCalibratedByCalibratorID(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (*response.TotalCalibratedRating, error) {
	return r.repo.GetCalibratedRating(calibratorID, prevCalibrator, businessUnitName, types, projectID)
}

func (r *projectUsecase) FindAverageScoreByCalibratorID(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (float32, error) {
	return r.repo.GetAverageScore(calibratorID, prevCalibrator, businessUnitName, types, projectID)
}

func (r *projectUsecase) FindAllEmployeeName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error) {
	return r.repo.GetAllEmployeeName(calibratorID, prevCalibrator, businessUnitName, types, projectID)
}

func (r *projectUsecase) FindAllSupervisorName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error) {
	return r.repo.GetAllSupervisorName(calibratorID, prevCalibrator, businessUnitName, types, projectID)
}

func (r *projectUsecase) FindAllGrade(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error) {
	return r.repo.GetAllGrade(calibratorID, prevCalibrator, businessUnitName, types, projectID)
}

func (r *projectUsecase) FindSummaryProjectByCalibratorID(calibratorID, projectID string, prevCalibratorIDs []string) (*response.SummaryProject, error) {
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

	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return nil, err
	}

	prevCalibrator := map[string]string{}
	businessUnit := map[string]string{}
	picIDs := map[string]string{}
	resultSummary := map[string]*response.CalibratorBusinessUnit{}
	users, err := r.repo.GetAllUserCalibrationByCalibratorID(calibratorID, projectID, phase)
	if err != nil {
		return nil, err
	}

	totalUsers := 0
	countCalibratedScoresUsers := 0.0

	// Grouping By Previous Calibrator and Business Unit
	for _, user := range users {
		if user.ScoringMethod == "Score" {
			totalUsers += 1
			countCalibratedScoresUsers += user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationScore
		}

		pic := false
		picName := "N-1"
		picId := "N-1"
		calibrationLength := len(user.CalibrationScores)
		for _, calibration := range user.CalibrationScores {
			if calibration.ProjectPhase.Phase.Order == phase && calibration.CalibratorID == calibratorID {
				if _, isExist := prevCalibrator[user.Name+*user.BusinessUnitId]; calibrationLength == 1 && !isExist {
					// check if n-1 or not
					checkIfTrue, err := r.repo.FindIfCalibratorOnPhaseBefore(user.ID, projectID, calibration.ProjectPhase.Phase.Order)
					if err != nil {
						return nil, err
					}

					fmt.Println("==========================CHECK PERNAH KALIBRASI ATAU GA==============================", checkIfTrue)
					if checkIfTrue {
						prevCalibrator[user.Name+*user.BusinessUnitId] = user.Name
						picName = user.Name
						picId = user.ID
					}
				}
				pic = true
				// else if name, isExist := businessUnit[user.BusinessUnit.Name]; calibrationLength == 1 && isExist {
				// 	picName = name
				// 	picId = picIDs[user.BusinessUnit.Name]
				// }
				break
			} else if calibration.ProjectPhase.Phase.Order >= phase && calibration.CalibratorID != calibratorID {
				break
			}

			if calibration.ProjectPhase.Phase.Order < phase {
				prevCalibrator[calibration.Calibrator.Name+*user.BusinessUnitId] = calibration.Calibrator.Name
				picName = calibration.Calibrator.Name
				picId = calibration.CalibratorID
			}
		}

		fmt.Println("===================DATA USER=================", user.Name, "========", picName, calibrationLength)

		filterCheck := true
		// fmt.Println("DATA KITA", prevCalibratorIDs, fmt.Sprintf("%s-%s", picId, *user.BusinessUnitId), contains(prevCalibratorIDs, fmt.Sprintf("%s-%s", picId, *user.BusinessUnitId)))

		checkCalibrator := picId
		checkCalibrator = fmt.Sprintf("%s-%s", picId, *user.BusinessUnitId)
		if len(prevCalibratorIDs) > 0 && contains(prevCalibratorIDs, checkCalibrator) {
			filterCheck = true
		} else if len(prevCalibratorIDs) > 0 && !contains(prevCalibratorIDs, checkCalibrator) {
			filterCheck = false
		}

		if filterCheck {
			bu := true
			if _, ok := resultSummary[picName+user.BusinessUnit.Name]; ok {
				bu = false
			}

			if _, isExist := businessUnit[user.BusinessUnit.Name]; bu && pic && (picName != "N-1" || !isExist) {
				resp := &response.CalibratorBusinessUnit{
					Pillar:                   user.BusinessUnit.Pillar,
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
			} else {
				if summary, ok := resultSummary[picName+user.BusinessUnit.Name]; ok {
					fmt.Println("=============PICNAME N-1======================")
					if picName == "N-1" {
						fmt.Println("================user===================", user.Name, user.CalibrationScores[calibrationLength-1].CalibrationRating)
					}
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
						} else {
							summary.Status = "Waiting"
						}
					}

					if user.CalibrationScores[calibrationLength-1].Status == "Waiting" {
						summary.Status = "Waiting"
					}
				} else {
					if picName == "N-1" {
						fmt.Println("================user===================", user.Name, user.CalibrationScores[calibrationLength-1].CalibrationRating)
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
		}
	}

	var summaries []*response.CalibratorBusinessUnit
	for _, summary := range resultSummary {
		summaries = append(summaries, summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].CalibratorBusinessUnit < summaries[j].CalibratorBusinessUnit
	})

	buCheck := map[string]string{}
	finalData := map[string]*response.BusinessUnitTotal{}
	for _, summary := range summaries {
		types := "all"
		// if summary.CalibratorName == "N-1" {
		// 	types = "n-1"
		// }

		guidance, err := r.FindRatingQuotaByCalibratorIDforSummaryHelper(calibratorID, summary.CalibratorID, summary.CalibratorBusinessUnitID, types, projectID, 0)
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
				APlusGuidance:              guidance.APlus,
				AGuidance:                  guidance.A,
				BPlusGuidance:              guidance.BPlus,
				BGuidance:                  guidance.B,
				CGuidance:                  guidance.C,
				DGuidance:                  guidance.D,
				TotalCalibratedScore:       0,
				UserCount:                  0,
				AverageScore:               0,
				Status:                     "Calibrate",
				Completed:                  true,
			}

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
		if result.Summary[i].CalibratorBusinessUnitID != result.Summary[j].CalibratorBusinessUnitID {
			return result.Summary[i].CalibratorBusinessUnitID < result.Summary[j].CalibratorBusinessUnitID
		}
		return result.Summary[i].CalibratorBusinessUnitID < result.Summary[j].CalibratorBusinessUnitID
	})

	for _, businessUnit := range result.Summary {
		sort.Slice(businessUnit.CalibratorBusinessUnit, func(i, j int) bool {
			return businessUnit.CalibratorBusinessUnit[i].CalibratorName < businessUnit.CalibratorBusinessUnit[j].CalibratorName
		})
	}

	return result, nil
}

func (r *projectUsecase) FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, prevCalibrator, businessUnit, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.GetCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, prevCalibrator, businessUnit, projectID, phase, paginationQuery)
}

func (r *projectUsecase) FindCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindCalibrationsByBusinessUnitPaginate(calibratorID, businessUnit, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.GetCalibrationsByBusinessUnitPaginate(calibratorID, businessUnit, projectID, phase, paginationQuery)
}

func (r *projectUsecase) FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	users, err := r.repo.GetNumberOneUserWhoCalibrator(calibratorID, businessUnit, projectID, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}

	results, err := r.repo.GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, phase, users)
	if err != nil {
		return response.UserCalibration{}, err
	}

	return results, nil
}

func (r *projectUsecase) FindNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetNMinusOneCalibrationsByBusinessUnit(businessUnit, phase, calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindNMinusOneCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, businessUnit, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.GetNMinusOneCalibrationsByBusinessUnitPaginate(businessUnit, phase, calibratorID, projectID, paginationQuery)
}

func (r *projectUsecase) FindCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.GetCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID, phase, paginationQuery)
}

func (r *projectUsecase) FindCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.GetCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID, phase, paginationQuery)
}

func (r *projectUsecase) FindCalibrationsByRating(calibratorID, rating, projectID string, param request.PaginationParam) (response.UserCalibrationNew, response.Paging, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.GetCalibrationsByRating(calibratorID, rating, projectID, phase, paginationQuery)
}

func (r *projectUsecase) FindSummaryProjectTotalByCalibratorID(calibratorID, projectID string) (*response.SummaryTotal, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return nil, err
	}

	allBusinessUnit, err := r.repo.GetAllBusinessUnitSummary(calibratorID, projectID, phase)
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
		ratingQuota, err := r.FindRatingQuotaByCalibratorIDforSummaryHelper(calibratorID, "", businessUnit.ID, "all", projectID, 0)
		if err != nil {
			return nil, err
		}

		mapRating["A+"].Guidance += ratingQuota.APlus
		mapRating["A"].Guidance += ratingQuota.A
		mapRating["B+"].Guidance += ratingQuota.BPlus
		mapRating["B"].Guidance += ratingQuota.B
		mapRating["C"].Guidance += ratingQuota.C
		mapRating["D"].Guidance += ratingQuota.D

		users, err := r.FindCalibrationsByBusinessUnit(calibratorID, businessUnit.ID, projectID)
		if err != nil {
			return nil, err
		}
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

func (r *projectUsecase) FindCalibratorPhase(calibratorID, projectID string) (*model.ProjectPhase, error) {
	phase, err := r.repo.GetProjectPhase(calibratorID, projectID)
	if err != nil {
		return nil, err
	}

	return phase, nil
}

func (r *projectUsecase) FindActiveProjectPhase(projectID string) ([]model.ProjectPhase, error) {
	projectPhase, err := r.repo.GetActiveProjectPhase(projectID)
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

func (r *projectUsecase) FindProjectRatingQuotaByBusinessUnit(businessUnitID, projectID string) (*model.Project, error) {
	projects, err := r.repo.GetRatingQuotaByCalibratorID(businessUnitID, projectID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func truncateFloat(f float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	truncated := math.Trunc(f*factor) / factor
	return truncated
}

func (r *projectUsecase) FindReportCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetAllDataCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindReportCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetAllDataCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID, phase)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) FindReportNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnit, projectID string) (response.UserCalibration, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}

	calibration, err := r.repo.GetAllDataNMinusOneCalibrationsByBusinessUnit(businessUnit, phase, calibratorID, projectID)
	if err != nil {
		return response.UserCalibration{}, err
	}
	return calibration, nil
}

func (r *projectUsecase) ReportCalibrations(types, calibratorID, businessUnit, prevCalibrator, projectID string, c *gin.Context) (string, error) {
	var responseData response.UserCalibration
	var err error

	if types == "numberOne" {
		responseData, err = r.FindNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID)
	} else if types == "n-1" {
		responseData, err = r.FindReportNMinusOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, businessUnit, projectID)
	} else if types == "default" {
		responseData, err = r.FindReportCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID)
	} else {
		responseData, err = r.FindReportCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID)
	}
	if err != nil {
		return "", err
	}

	projectPhase, err := r.FindCalibratorPhase(calibratorID, projectID)
	if err != nil {
		return "", err
	}

	actualScore, err := r.FindTotalActualScoreByCalibratorID(calibratorID, prevCalibrator, businessUnit, types, projectID)
	if err != nil {
		return "", err
	}

	ratingQuota, err := r.FindRatingQuotaByCalibratorIDforSummaryHelper(calibratorID, prevCalibrator, businessUnit, types, projectID, 0)
	if err != nil {
		return "", err
	}

	file := excelize.NewFile()
	sheetName := responseData.UserData[0].BusinessUnit.Name
	index := file.NewSheet(sheetName)

	headers := []string{
		"No",
		"Nik",
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

	style1, err := file.NewStyle(`{"alignment":{"vertical":"center"}, "font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	style2, err := file.NewStyle(`{"alignment":{"vertical":"center"}}`)
	if err != nil {
		return "", err
	}

	countCalibrated := map[string]int{
		"APlus": 0,
		"A":     0,
		"BPlus": 0,
		"B":     0,
		"C":     0,
		"D":     0,
	}

	firstAppear := true
	for col, header := range headers {
		colName := excelize.ToAlphaString(col)
		cellRef := fmt.Sprintf("%s%d", colName, 13)
		cellRefUnder := fmt.Sprintf("%s%d", colName, 14)
		nextColName := excelize.ToAlphaString(col + 1)
		cellNextCollRef := fmt.Sprintf("%s%d", nextColName, 13)
		file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, style1)

		if (strings.Contains(header, "Prev") || strings.Contains(header, "Actual") || strings.Contains(header, "Calibration")) && firstAppear {
			file.MergeCell(responseData.UserData[0].BusinessUnit.Name, cellRef, cellNextCollRef)
			words := strings.Fields(header)
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRef, words[0])
			file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, words[1])
			if strings.Contains(header, "Calibration") {
				backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern", "color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[2]))
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
				backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[2]))
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
			backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[1]))
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

	fmt.Println("==========================DATA USER", responseData.UserData[0].Name, responseData.UserData[0].CalibrationScores)
	for i, user := range responseData.UserData {
		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("A%d", i+15), i+1)
		file.SetColWidth(user.BusinessUnit.Name, "A", "A", 4)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("A%d", i+15), fmt.Sprintf("A%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("B%d", i+15), user.Nik)
		file.SetColWidth(user.BusinessUnit.Name, "B", "B", 15)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("B%d", i+15), fmt.Sprintf("B%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("C%d", i+15), user.Name)
		file.SetColWidth(user.BusinessUnit.Name, "C", "C", 20)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("C%d", i+15), fmt.Sprintf("C%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("D%d", i+15), user.Grade)
		file.SetColWidth(user.BusinessUnit.Name, "D", "D", 4)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("D%d", i+15), fmt.Sprintf("D%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("E%d", i+15), user.BusinessUnit.Name)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("E%d", i+15), fmt.Sprintf("E%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("F%d", i+15), user.OrganizationUnit)
		file.SetColWidth(user.BusinessUnit.Name, "F", "F", 15)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("F%d", i+15), fmt.Sprintf("F%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("G%d", i+15), user.SupervisorNames)
		file.SetColWidth(user.BusinessUnit.Name, "G", "G", 20)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("G%d", i+15), fmt.Sprintf("G%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("H%d", i+15), user.ActualScores[0].Y2Rating)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("H%d", i+15), fmt.Sprintf("H%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("I%d", i+15), user.ActualScores[0].Y1Rating)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("I%d", i+15), fmt.Sprintf("I%d", i+15), style2)

		// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+15), truncateFloat(user.ActualScores[0].PTTScore, 2))
		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].PTTScore))
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+15), fmt.Sprintf("J%d", i+15), style2)

		// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+15), truncateFloat(user.ActualScores[0].PATScore, 2))
		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].PATScore))
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+15), fmt.Sprintf("K%d", i+15), style2)

		// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+15), truncateFloat(user.ActualScores[0].Score360, 2))
		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].Score360))
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+15), fmt.Sprintf("L%d", i+15), style2)

		// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+15), truncateFloat(user.ActualScores[0].ActualScore, 2))
		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].ActualScore))
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+15), fmt.Sprintf("M%d", i+15), style2)

		file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("N%d", i+15), user.ActualScores[0].ActualRating)
		file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("N%d", i+15), fmt.Sprintf("N%d", i+15), style2)
		column := int('O')

		rating := user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating
		switch rating {
		case "A+":
			countCalibrated["APlus"] += 1
		case "A":
			countCalibrated["A"] += 1
		case "B+":
			countCalibrated["BPlus"] += 1
		case "B":
			countCalibrated["B"] += 1
		case "C":
			countCalibrated["C"] += 1
		case "D":
			countCalibrated["D"] += 1
		default:
			//
		}

		startCalibrationPhaseCol := 14
		for j := 1; j < projectPhase.Phase.Order+1; j++ {
			columnBefore := column
			words := strings.Fields(headers[startCalibrationPhaseCol])
			startCalibrationPhaseCol += 4
			backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}}`, words[2]))
			if err != nil {
				return "", err
			}

			for _, calibrationScore := range user.CalibrationScores {
				if j == calibrationScore.ProjectPhase.Phase.Order {
					fmt.Println("sheet= ", user.BusinessUnit.Name, " axis=", fmt.Sprintf("%s%d", asciiToName(column), i+15), " value=", truncateFloat(calibrationScore.CalibrationScore, 2))
					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), truncateFloat(calibrationScore.CalibrationScore, 2))
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), calibrationScore.CalibrationRating)
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), calibrationScore.JustificationType)
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++

					if calibrationScore.JustificationType == "default" {
						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), calibrationScore.Comment)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
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
						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), justifications)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
						file.SetColWidth(user.BusinessUnit.Name, fmt.Sprintf("%s", asciiToName(column)), fmt.Sprintf("%s", asciiToName(column)), 25)
					} else {
						justification := fmt.Sprintf("Attitude: %s\nIndisipliner: %s\nLow Performance: %s\nWarning Letter: %s",
							calibrationScore.BottomRemark.Attitude,
							calibrationScore.BottomRemark.Indisipliner,
							calibrationScore.BottomRemark.LowPerformance,
							calibrationScore.BottomRemark.WarningLetter,
						)
						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), justification)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
						file.SetColWidth(user.BusinessUnit.Name, fmt.Sprintf("%s", asciiToName(column)), fmt.Sprintf("%s", asciiToName(column)), 25)
					}
					column++
				}
			}
			if columnBefore == column {
				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
				column++

				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
				column++

				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
				column++

				file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
				file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
				column++
			}
		}
	}

	file.SetActiveSheet(index)
	file.DeleteSheet("Sheet1")

	deviationRatingAP := countCalibrated["APlus"] - ratingQuota.APlus
	deviationRatingA := countCalibrated["A"] - ratingQuota.A
	deviationRatingBP := countCalibrated["BPlus"] - ratingQuota.BPlus
	deviationRatingB := countCalibrated["B"] - ratingQuota.B
	deviationRatingC := countCalibrated["C"] - ratingQuota.C
	deviationRatingD := countCalibrated["D"] - ratingQuota.D

	// Add data for the chart
	data := map[string]interface{}{
		"H3": "Category", "I3": "Rating Scale", "J3": "Guidance", "K3": "Actual", "L3": "Calibrated", "M3": "Deviation Rating",
		"H4": "A+", "I4": "5-5", "J4": ratingQuota.APlus, "K4": actualScore.APlus, "L4": countCalibrated["APlus"], "M4": deviationRatingAP,
		"H5": "A", "I5": "4.5 - 4.99", "J5": ratingQuota.A, "K5": actualScore.A, "L5": countCalibrated["A"], "M5": deviationRatingA,
		"H6": "B+", "I6": "3.5 - 4.49", "J6": ratingQuota.BPlus, "K6": actualScore.BPlus, "L6": countCalibrated["BPlus"], "M6": deviationRatingBP,
		"H7": "B", "I7": "3 - 3.49", "J7": ratingQuota.B, "K7": actualScore.B, "L7": countCalibrated["B"], "M7": deviationRatingB,
		"H8": "C", "I8": "2 - 2.99", "J8": ratingQuota.C, "K8": actualScore.C, "L8": countCalibrated["C"], "M8": deviationRatingC,
		"H9": "D", "I9": "0 - 1.99", "J9": ratingQuota.D, "K9": actualScore.D, "L9": countCalibrated["D"], "M9": deviationRatingD,
		"H10": "Total", "J10": ratingQuota.APlus + ratingQuota.A + ratingQuota.BPlus + ratingQuota.B + ratingQuota.C + ratingQuota.D,
		"K10": actualScore.APlus + actualScore.A + actualScore.BPlus + actualScore.B + actualScore.C + actualScore.D,
		"L10": countCalibrated["APlus"] + countCalibrated["A"] + countCalibrated["BPlus"] + countCalibrated["B"] + countCalibrated["C"] + countCalibrated["D"],
		"M10": deviationRatingAP + deviationRatingA + deviationRatingBP + deviationRatingB + deviationRatingC + deviationRatingD,
	}
	for k, v := range data {
		file.SetCellValue(sheetName, k, v)
	}

	chartJSON := fmt.Sprintf(`{
		"type": "line",
		"dimension": {
			"width": 600,
			"height": 225
		},
		"series": [
			{
				"name": "'%s'!$J$3",
				"categories": "'%s'!$H$4:$H$9",
				"values": "'%s'!$J$4:$J$9",
				"line": {
					"color": "#BEBEBE"
				}
			},
			{
				"name": "'%s'!$K$3",
				"categories": "'%s'!$H$4:$H$9",
				"values": "'%s'!$K$4:$K$9",
				"line": {
					"color": "#02B4CC"
				}
			},
			{
				"name": "'%s'!$L$3",
				"categories": "'%s'!$H$4:$H$9",
				"values": "'%s'!$L$4:$L$9",
				"line": {
					"color": "#FF7300"
				}
			}
		],
		"legend": {
			"position": "bottom",
			"show_legend_key": true
		},
		"title": {
			"name": "Chart"
		},
		"plotarea": {
			"show_val": false
		},
		"x_axis": {
			"reverse_order": false
		},
		"y_axis": {
			"maximum": %d,
			"minimum": 0
		}
	}`, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, len(responseData.UserData))
	// fmt.Println("=====================,", len(responseData.UserData), chartJSON)

	// Add the chart at a specific location
	err = file.AddChart(sheetName, "A1", chartJSON)
	if err != nil {
		return "", err
	}

	err = file.SaveAs("report.xlsx")
	if err != nil {
		return "", err
	}

	return "report.xlsx", nil
}

func (r *projectUsecase) SummaryReportCalibrations(calibratorID, projectID string, c *gin.Context) (string, error) {
	file := excelize.NewFile()
	summary, err := r.FindSummaryProjectByCalibratorID(calibratorID, projectID, []string{})
	if err != nil {
		return "", err
	}

	style1, err := file.NewStyle(`{"alignment":{"vertical":"center"}, "font":{"bold":true}}`)
	if err != nil {
		return "", err
	}

	style2, err := file.NewStyle(`{"alignment":{"vertical":"center"}}`)
	if err != nil {
		return "", err
	}

	summaryHeaders := []string{
		"Business Unit",
		"Previous Calibrator",
		"Indicator",
		"A+",
		"A",
		"B+",
		"B",
		"C",
		"D",
		"Total",
		"Average",
		"Status",
	}

	for col, header := range summaryHeaders {
		colName := excelize.ToAlphaString(col)
		cellRef := fmt.Sprintf("%s%d", colName, 1)
		file.SetCellStyle("Sheet1", cellRef, cellRef, style2)
		file.SetCellValue("Sheet1", cellRef, header)
	}

	firstRow := 2
	for _, summaryData := range summary.Summary {
		startRow := firstRow
		for _, prevCalibratorData := range summaryData.CalibratorBusinessUnit {
			file.SetCellValue("Sheet1", fmt.Sprintf("A%d", firstRow), summaryData.CalibratorBusinessUnitName)
			file.SetCellStyle("Sheet1", fmt.Sprintf("A%d", firstRow), fmt.Sprintf("A%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("B%d", firstRow), prevCalibratorData.CalibratorName)
			file.SetCellStyle("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("B%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("C%d", firstRow), "Calibrated")
			file.SetCellStyle("Sheet1", fmt.Sprintf("C%d", firstRow), fmt.Sprintf("C%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("D%d", firstRow), prevCalibratorData.APlus)
			file.SetCellStyle("Sheet1", fmt.Sprintf("D%d", firstRow), fmt.Sprintf("D%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("E%d", firstRow), prevCalibratorData.A)
			file.SetCellStyle("Sheet1", fmt.Sprintf("E%d", firstRow), fmt.Sprintf("E%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("F%d", firstRow), prevCalibratorData.BPlus)
			file.SetCellStyle("Sheet1", fmt.Sprintf("F%d", firstRow), fmt.Sprintf("F%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("G%d", firstRow), prevCalibratorData.B)
			file.SetCellStyle("Sheet1", fmt.Sprintf("G%d", firstRow), fmt.Sprintf("G%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("H%d", firstRow), prevCalibratorData.C)
			file.SetCellStyle("Sheet1", fmt.Sprintf("H%d", firstRow), fmt.Sprintf("H%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("I%d", firstRow), prevCalibratorData.D)
			file.SetCellStyle("Sheet1", fmt.Sprintf("I%d", firstRow), fmt.Sprintf("I%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("J%d", firstRow), prevCalibratorData.TotalCalibratedScore)
			file.SetCellStyle("Sheet1", fmt.Sprintf("J%d", firstRow), fmt.Sprintf("J%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("K%d", firstRow), prevCalibratorData.AverageScore)
			file.SetCellStyle("Sheet1", fmt.Sprintf("K%d", firstRow), fmt.Sprintf("K%d", firstRow), style2)

			file.SetCellValue("Sheet1", fmt.Sprintf("L%d", firstRow), prevCalibratorData.Status)
			file.SetCellStyle("Sheet1", fmt.Sprintf("L%d", firstRow), fmt.Sprintf("L%d", firstRow), style2)
			firstRow += 1
		}
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", firstRow), summaryData.CalibratorBusinessUnitName)
		file.SetCellStyle("Sheet1", fmt.Sprintf("A%d", firstRow), fmt.Sprintf("A%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", firstRow), "Rating")
		file.SetCellStyle("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("B%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", firstRow), "Guidance")
		file.SetCellStyle("Sheet1", fmt.Sprintf("C%d", firstRow), fmt.Sprintf("C%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", firstRow), summaryData.APlusGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("D%d", firstRow), fmt.Sprintf("D%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", firstRow), summaryData.AGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("E%d", firstRow), fmt.Sprintf("E%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", firstRow), summaryData.BPlusGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("F%d", firstRow), fmt.Sprintf("F%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", firstRow), summaryData.BGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("G%d", firstRow), fmt.Sprintf("G%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("H%d", firstRow), summaryData.CGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("H%d", firstRow), fmt.Sprintf("H%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("I%d", firstRow), summaryData.DGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("I%d", firstRow), fmt.Sprintf("I%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("J%d", firstRow), summaryData.APlusGuidance+summaryData.AGuidance+summaryData.BPlusGuidance+summaryData.BGuidance+summaryData.CGuidance+summaryData.DGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("J%d", firstRow), fmt.Sprintf("J%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("K%d", firstRow), summaryData.AverageScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("K%d", firstRow), fmt.Sprintf("K%d", firstRow), style2)
		firstRow += 1

		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", firstRow), summaryData.CalibratorBusinessUnitName)
		file.SetCellStyle("Sheet1", fmt.Sprintf("A%d", firstRow), fmt.Sprintf("A%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", firstRow), "Total")
		file.SetCellStyle("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("B%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", firstRow), "Calibrated")
		file.SetCellStyle("Sheet1", fmt.Sprintf("C%d", firstRow), fmt.Sprintf("C%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", firstRow), summaryData.APlusCalibrated)
		file.SetCellStyle("Sheet1", fmt.Sprintf("D%d", firstRow), fmt.Sprintf("D%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", firstRow), summaryData.ACalibrated)
		file.SetCellStyle("Sheet1", fmt.Sprintf("E%d", firstRow), fmt.Sprintf("E%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", firstRow), summaryData.BPlusCalibrated)
		file.SetCellStyle("Sheet1", fmt.Sprintf("F%d", firstRow), fmt.Sprintf("F%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", firstRow), summaryData.BCalibrated)
		file.SetCellStyle("Sheet1", fmt.Sprintf("G%d", firstRow), fmt.Sprintf("G%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("H%d", firstRow), summaryData.CCalibrated)
		file.SetCellStyle("Sheet1", fmt.Sprintf("H%d", firstRow), fmt.Sprintf("H%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("I%d", firstRow), summaryData.DCalibrated)
		file.SetCellStyle("Sheet1", fmt.Sprintf("I%d", firstRow), fmt.Sprintf("I%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("J%d", firstRow), summaryData.TotalCalibratedScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("J%d", firstRow), fmt.Sprintf("J%d", firstRow), style2)
		firstRow += 1

		endRow := firstRow - 1  // Adjust to get the last row of the group
		if startRow != endRow { // Only merge if there's more than one row
			file.MergeCell("Sheet1", fmt.Sprintf("A%d", startRow), fmt.Sprintf("A%d", endRow))
		}
	}

	if len(summary.Summary) > 1 {
		total := firstRow
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", firstRow), "Grand Total")
		file.SetCellStyle("Sheet1", fmt.Sprintf("A%d", firstRow), fmt.Sprintf("A%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", firstRow), "Guidance")
		file.SetCellStyle("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("B%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", firstRow), summary.APlusGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("D%d", firstRow), fmt.Sprintf("D%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", firstRow), summary.AGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("E%d", firstRow), fmt.Sprintf("E%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", firstRow), summary.BPlusGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("F%d", firstRow), fmt.Sprintf("F%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", firstRow), summary.BGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("G%d", firstRow), fmt.Sprintf("G%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("H%d", firstRow), summary.CGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("H%d", firstRow), fmt.Sprintf("H%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("I%d", firstRow), summary.DGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("I%d", firstRow), fmt.Sprintf("I%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("J%d", firstRow), summary.APlusGuidance+summary.AGuidance+summary.BPlusGuidance+summary.BGuidance+summary.CGuidance+summary.DGuidance)
		file.SetCellStyle("Sheet1", fmt.Sprintf("J%d", firstRow), fmt.Sprintf("J%d", firstRow), style2)
		file.MergeCell("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("C%d", firstRow))

		firstRow += 1

		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", firstRow), "Grand Total")
		file.SetCellStyle("Sheet1", fmt.Sprintf("A%d", firstRow), fmt.Sprintf("A%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", firstRow), "Calibrated")
		file.SetCellStyle("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("B%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", firstRow), summary.APlusTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("D%d", firstRow), fmt.Sprintf("D%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", firstRow), summary.ATotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("E%d", firstRow), fmt.Sprintf("E%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", firstRow), summary.BPlusTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("F%d", firstRow), fmt.Sprintf("F%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", firstRow), summary.BTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("G%d", firstRow), fmt.Sprintf("G%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("H%d", firstRow), summary.CTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("H%d", firstRow), fmt.Sprintf("H%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("I%d", firstRow), summary.DTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("I%d", firstRow), fmt.Sprintf("I%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("J%d", firstRow), summary.APlusTotalScore+summary.ATotalScore+summary.BPlusTotalScore+summary.BTotalScore+summary.CTotalScore+summary.DTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("J%d", firstRow), fmt.Sprintf("J%d", firstRow), style2)

		file.SetCellValue("Sheet1", fmt.Sprintf("K%d", firstRow), summary.AverageTotalScore)
		file.SetCellStyle("Sheet1", fmt.Sprintf("K%d", firstRow), fmt.Sprintf("K%d", firstRow), style2)
		file.MergeCell("Sheet1", fmt.Sprintf("B%d", firstRow), fmt.Sprintf("C%d", firstRow))

		firstRow += 1
		file.MergeCell("Sheet1", fmt.Sprintf("A%d", total), fmt.Sprintf("A%d", firstRow-1))
	}
	file.SetSheetName("Sheet1", "Summary")

	for _, summaryBusinessUnit := range summary.Summary {
		responseData, err := r.FindReportCalibrationsByBusinessUnit(calibratorID, summaryBusinessUnit.CalibratorBusinessUnitID, projectID)
		if err != nil {
			return "", err
		}

		actualScore, err := r.FindTotalActualScoreByCalibratorID(calibratorID, "", summaryBusinessUnit.CalibratorBusinessUnitID, "all", projectID)
		if err != nil {
			return "", err
		}

		ratingQuota, err := r.FindRatingQuotaByCalibratorIDforSummaryHelper(calibratorID, "", summaryBusinessUnit.CalibratorBusinessUnitID, "all", projectID, 0)
		if err != nil {
			return "", err
		}

		projectPhase, err := r.FindCalibratorPhase(calibratorID, projectID)
		if err != nil {
			return "", err
		}

		sheetName := responseData.UserData[0].BusinessUnit.Name
		index := file.NewSheet(sheetName)

		headers := []string{
			"No",
			"Nik",
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

		countCalibrated := map[string]int{
			"APlus": 0,
			"A":     0,
			"BPlus": 0,
			"B":     0,
			"C":     0,
			"D":     0,
		}

		firstAppear := true
		for col, header := range headers {
			colName := excelize.ToAlphaString(col)
			cellRef := fmt.Sprintf("%s%d", colName, 13)
			cellRefUnder := fmt.Sprintf("%s%d", colName, 14)
			nextColName := excelize.ToAlphaString(col + 1)
			cellNextCollRef := fmt.Sprintf("%s%d", nextColName, 13)
			file.SetCellStyle(responseData.UserData[0].BusinessUnit.Name, cellRef, cellRef, style1)

			if (strings.Contains(header, "Prev") || strings.Contains(header, "Actual") || strings.Contains(header, "Calibration")) && firstAppear {
				file.MergeCell(responseData.UserData[0].BusinessUnit.Name, cellRef, cellNextCollRef)
				words := strings.Fields(header)
				file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRef, words[0])
				file.SetCellValue(responseData.UserData[0].BusinessUnit.Name, cellRefUnder, words[1])
				if strings.Contains(header, "Calibration") {
					backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern", "color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[2]))
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
					backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[2]))
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
				backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}, "font":{"bold":true}}`, words[1]))
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
			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("A%d", i+15), i+1)
			file.SetColWidth(user.BusinessUnit.Name, "A", "A", 4)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("A%d", i+15), fmt.Sprintf("A%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("B%d", i+15), user.Nik)
			file.SetColWidth(user.BusinessUnit.Name, "B", "B", 15)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("B%d", i+15), fmt.Sprintf("B%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("C%d", i+15), user.Name)
			file.SetColWidth(user.BusinessUnit.Name, "C", "C", 20)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("C%d", i+15), fmt.Sprintf("C%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("D%d", i+15), user.Grade)
			file.SetColWidth(user.BusinessUnit.Name, "D", "D", 4)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("D%d", i+15), fmt.Sprintf("D%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("E%d", i+15), user.BusinessUnit.Name)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("E%d", i+15), fmt.Sprintf("E%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("F%d", i+15), user.OrganizationUnit)
			file.SetColWidth(user.BusinessUnit.Name, "F", "F", 15)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("F%d", i+15), fmt.Sprintf("F%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("G%d", i+15), user.SupervisorNames)
			file.SetColWidth(user.BusinessUnit.Name, "G", "G", 20)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("G%d", i+15), fmt.Sprintf("G%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("H%d", i+15), user.ActualScores[0].Y2Rating)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("H%d", i+15), fmt.Sprintf("H%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("I%d", i+15), user.ActualScores[0].Y1Rating)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("I%d", i+15), fmt.Sprintf("I%d", i+15), style2)

			// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+3), truncateFloat(user.ActualScores[0].PTTScore, 2))
			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].PTTScore))
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("J%d", i+15), fmt.Sprintf("J%d", i+15), style2)

			// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+15), truncateFloat(user.ActualScores[0].PATScore, 2))
			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].PATScore))
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("K%d", i+15), fmt.Sprintf("K%d", i+15), style2)

			// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+15), truncateFloat(user.ActualScores[0].Score360, 2))
			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].Score360))
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("L%d", i+15), fmt.Sprintf("L%d", i+15), style2)

			// file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+15), truncateFloat(user.ActualScores[0].ActualScore, 2))
			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+15), fmt.Sprintf("%.2f", user.ActualScores[0].ActualScore))
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("M%d", i+15), fmt.Sprintf("M%d", i+15), style2)

			file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("N%d", i+15), user.ActualScores[0].ActualRating)
			file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("N%d", i+15), fmt.Sprintf("N%d", i+15), style2)
			column := int('O')

			rating := user.CalibrationScores[len(user.CalibrationScores)-1].CalibrationRating
			switch rating {
			case "A+":
				countCalibrated["APlus"] += 1
			case "A":
				countCalibrated["A"] += 1
			case "B+":
				countCalibrated["BPlus"] += 1
			case "B":
				countCalibrated["B"] += 1
			case "C":
				countCalibrated["C"] += 1
			case "D":
				countCalibrated["D"] += 1
			default:
				//
			}

			startCalibrationPhaseCol := 14
			for j := 1; j < projectPhase.Phase.Order+1; j++ {
				columnBefore := column
				words := strings.Fields(headers[startCalibrationPhaseCol])
				startCalibrationPhaseCol += 4
				backgroundColorCalibration, err := file.NewStyle(fmt.Sprintf(`{"alignment":{"vertical":"center"}, "fill":{"type":"pattern","color":["%s"],"pattern":1}}`, words[2]))
				if err != nil {
					return "", err
				}

				for _, calibrationScore := range user.CalibrationScores {
					if j == calibrationScore.ProjectPhase.Phase.Order {

						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), truncateFloat(calibrationScore.CalibrationScore, 2))
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
						column++

						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), calibrationScore.CalibrationRating)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
						column++

						file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), calibrationScore.JustificationType)
						file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
						column++

						if calibrationScore.JustificationType == "default" {
							file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), calibrationScore.Comment)
							file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
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
							file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), justifications)
							file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
							file.SetColWidth(user.BusinessUnit.Name, fmt.Sprintf("%s", asciiToName(column)), fmt.Sprintf("%s", asciiToName(column)), 25)
						} else {
							justification := fmt.Sprintf("Attitude: %s\nIndisipliner: %s\nLow Performance: %s\nWarning Letter: %s",
								calibrationScore.BottomRemark.Attitude,
								calibrationScore.BottomRemark.Indisipliner,
								calibrationScore.BottomRemark.LowPerformance,
								calibrationScore.BottomRemark.WarningLetter,
							)
							file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), justification)
							file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
							file.SetColWidth(user.BusinessUnit.Name, fmt.Sprintf("%s", asciiToName(column)), fmt.Sprintf("%s", asciiToName(column)), 25)
						}
						column++
					}
				}
				if columnBefore == column {
					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++

					file.SetCellValue(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), "-")
					file.SetCellStyle(user.BusinessUnit.Name, fmt.Sprintf("%s%d", asciiToName(column), i+15), fmt.Sprintf("%s%d", asciiToName(column), i+15), backgroundColorCalibration)
					column++
				}
			}
		}
		file.SetActiveSheet(index)

		deviationRatingAP := countCalibrated["APlus"] - ratingQuota.APlus
		deviationRatingA := countCalibrated["A"] - ratingQuota.A
		deviationRatingBP := countCalibrated["BPlus"] - ratingQuota.BPlus
		deviationRatingB := countCalibrated["B"] - ratingQuota.B
		deviationRatingC := countCalibrated["C"] - ratingQuota.C
		deviationRatingD := countCalibrated["D"] - ratingQuota.D

		// Add data for the chart
		data := map[string]interface{}{
			"H3": "Category", "I3": "Rating Scale", "J3": "Guidance", "K3": "Actual", "L3": "Calibrated", "M3": "Deviation Rating",
			"H4": "A+", "I4": "5-5", "J4": ratingQuota.APlus, "K4": actualScore.APlus, "L4": countCalibrated["APlus"], "M4": deviationRatingAP,
			"H5": "A", "I5": "4.5 - 4.99", "J5": ratingQuota.A, "K5": actualScore.A, "L5": countCalibrated["A"], "M5": deviationRatingA,
			"H6": "B+", "I6": "3.5 - 4.49", "J6": ratingQuota.BPlus, "K6": actualScore.BPlus, "L6": countCalibrated["BPlus"], "M6": deviationRatingBP,
			"H7": "B", "I7": "3 - 3.49", "J7": ratingQuota.B, "K7": actualScore.B, "L7": countCalibrated["B"], "M7": deviationRatingB,
			"H8": "C", "I8": "2 - 2.99", "J8": ratingQuota.C, "K8": actualScore.C, "L8": countCalibrated["C"], "M8": deviationRatingC,
			"H9": "D", "I9": "0 - 1.99", "J9": ratingQuota.D, "K9": actualScore.D, "L9": countCalibrated["D"], "M9": deviationRatingD,
			"H10": "Total", "J10": ratingQuota.APlus + ratingQuota.A + ratingQuota.BPlus + ratingQuota.B + ratingQuota.C + ratingQuota.D,
			"K10": actualScore.APlus + actualScore.A + actualScore.BPlus + actualScore.B + actualScore.C + actualScore.D,
			"L10": countCalibrated["APlus"] + countCalibrated["A"] + countCalibrated["BPlus"] + countCalibrated["B"] + countCalibrated["C"] + countCalibrated["D"],
			"M10": deviationRatingAP + deviationRatingA + deviationRatingBP + deviationRatingB + deviationRatingC + deviationRatingD,
		}
		for k, v := range data {
			file.SetCellValue(sheetName, k, v)
		}

		chartJSON := fmt.Sprintf(`{
			"type": "line",
			"dimension": {
				"width": 600,
				"height": 225
			},
			"series": [
				{
					"name": "'%s'!$J$3",
					"categories": "'%s'!$H$4:$H$9",
					"values": "'%s'!$J$4:$J$9",
					"line": {
						"color": "#BEBEBE"
					}
				},
				{
					"name": "'%s'!$K$3",
					"categories": "'%s'!$H$4:$H$9",
					"values": "'%s'!$K$4:$K$9",
					"line": {
						"color": "#02B4CC"
					}
				},
				{
					"name": "'%s'!$L$3",
					"categories": "'%s'!$H$4:$H$9",
					"values": "'%s'!$L$4:$L$9",
					"line": {
						"color": "#FF7300"
					}
				}
			],
			"legend": {
				"position": "bottom",
				"show_legend_key": true
			},
			"title": {
				"name": "Chart"
			},
			"plotarea": {
				"show_val": false
			},
			"x_axis": {
				"reverse_order": false
			},
			"y_axis": {
				"maximum": %d,
				"minimum": 0
			}
		}`, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, sheetName, len(responseData.UserData))

		// Add the chart at a specific location
		err = file.AddChart(sheetName, "A1", chartJSON)
		if err != nil {
			return "", err
		}
	}

	sheetIndex := file.GetSheetIndex("Summary")
	file.SetActiveSheet(sheetIndex)
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

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if strings.HasPrefix(v, item) {
			return true
		}
	}
	return false
}

func (r *projectUsecase) FindRatingQuotaByCalibratorIDforSummaryHelper(calibratorID, prevCalibrator, businessUnitID, types, projectID string, countCurrentUser int) (*response.RatingQuota, error) {
	phase, err := r.repo.GetProjectPhaseOrder(calibratorID, projectID)
	if err != nil {
		return nil, err
	}

	calibrationsCount, err := r.repo.GetCalibrationsForSummaryHelper(types, calibratorID, prevCalibrator, businessUnitID, projectID, phase)
	if err != nil {
		return nil, err
	}

	projects, err := r.repo.GetRatingQuotaByCalibratorID(businessUnitID, projectID)
	if err != nil {
		return nil, err
	}

	ratingQuota := projects.RatingQuotas[0]
	totalCalibrations := calibrationsCount
	if countCurrentUser > 0 {
		totalCalibrations = countCurrentUser
	}
	// fmt.Println("TOTAL CALIBRATIONS =========================================================== ", totalCalibrations)
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

func (r *projectUsecase) FindActiveProjectByCalibratorID(calibratorID string) ([]model.Project, error) {
	return r.repo.GetAllActiveProjectByCalibratorID(calibratorID)
}

func (r *projectUsecase) FindActiveProjectBySpmoID(spmoID string) ([]model.Project, error) {
	return r.repo.GetAllActiveProjectBySpmoID(spmoID)
}

func NewProjectUsecase(repo repository.ProjectRepo) ProjectUsecase {
	return &projectUsecase{
		repo: repo,
	}
}
