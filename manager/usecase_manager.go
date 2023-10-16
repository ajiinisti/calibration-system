package manager

import (
	"calibration-system.com/config"
	"calibration-system.com/usecase"
)

type UsecaseManager interface {
	UserUc() usecase.UserUsecase
	RoleUc() usecase.RoleUsecase
	AuthUc() usecase.AuthUsecase
	BusinessUnitUc() usecase.BusinessUnitUsecase
	GroupBusinessUnitUc() usecase.GroupBusinessUnitUsecase
	ActualScoreUc() usecase.ActualScoreUsecase
	CalibrationUc() usecase.CalibrationUsecase
	PhaseUc() usecase.PhaseUsecase
	ProjectUc() usecase.ProjectUsecase
	ProjectPhaseUc() usecase.ProjectPhaseUsecase
	RatingQuotaUc() usecase.RatingQuotaUsecase
	ScoreDistributionUc() usecase.ScoreDistributionUsecase
	RemarkSettingUc() usecase.RemarkSettingUsecase
}

type usecaseManager struct {
	repo RepoManager
	cfg  *config.Config
}

func (u *usecaseManager) RoleUc() usecase.RoleUsecase {
	return usecase.NewRoleUsecase(u.repo.RoleRepo())
}

func (u *usecaseManager) UserUc() usecase.UserUsecase {
	return usecase.NewUserUseCase(u.repo.UserRepo(), u.RoleUc(), u.BusinessUnitUc(), u.cfg)
}

func (u *usecaseManager) AuthUc() usecase.AuthUsecase {
	return usecase.NewAuthUsecase(u.UserUc())
}

func (u *usecaseManager) GroupBusinessUnitUc() usecase.GroupBusinessUnitUsecase {
	return usecase.NewGroupBusinessUnitUsecase(u.repo.GroupBusinessUnitRepo())
}

func (u *usecaseManager) BusinessUnitUc() usecase.BusinessUnitUsecase {
	return usecase.NewBusinessUnitUsecase(u.repo.BusinessUnitRepo(), u.GroupBusinessUnitUc())
}

func (u *usecaseManager) PhaseUc() usecase.PhaseUsecase {
	return usecase.NewPhaseUsecase(u.repo.PhaseRepo())
}

func (u *usecaseManager) ProjectUc() usecase.ProjectUsecase {
	return usecase.NewProjectUsecase(u.repo.ProjectRepo())
}

func (u *usecaseManager) ProjectPhaseUc() usecase.ProjectPhaseUsecase {
	return usecase.NewProjectPhaseUsecase(u.repo.ProjectPhaseRepo(), u.PhaseUc(), u.ProjectUc())
}

func (u *usecaseManager) ActualScoreUc() usecase.ActualScoreUsecase {
	return usecase.NewActualScoreUsecase(u.repo.ActualScoreRepo(), u.UserUc(), u.ProjectUc())
}

func (u *usecaseManager) CalibrationUc() usecase.CalibrationUsecase {
	return usecase.NewCalibrationUsecase(u.repo.CalibrationRepo(), u.UserUc(), u.ProjectUc(), u.ProjectPhaseUc())
}

func (u *usecaseManager) RatingQuotaUc() usecase.RatingQuotaUsecase {
	return usecase.NewRatingQuotaUsecase(u.repo.RatingQuotaRepo(), u.BusinessUnitUc(), u.ProjectUc())
}

func (u *usecaseManager) ScoreDistributionUc() usecase.ScoreDistributionUsecase {
	return usecase.NewScoreDistributionUsecase(u.repo.ScoreDistributionRepo(), u.GroupBusinessUnitUc(), u.ProjectUc())
}

func (u *usecaseManager) RemarkSettingUc() usecase.RemarkSettingUsecase {
	return usecase.NewRemarkSettingUsecase(u.repo.RemarkSettingRepo(), u.ProjectUc())
}

func NewUsecaseManager(repo RepoManager, cfg *config.Config) UsecaseManager {
	return &usecaseManager{
		repo: repo,
		cfg:  cfg,
	}
}
