package manager

import "calibration-system.com/repository"

type RepoManager interface {
	RoleRepo() repository.RoleRepo
	UserRepo() repository.UserRepo
	BusinessUnitRepo() repository.BusinessUnitRepo
	GroupBusinessUnitRepo() repository.GroupBusinessUnitRepo
	ActualScoreRepo() repository.ActualScoreRepo
	CalibrationRepo() repository.CalibrationRepo
	PhaseRepo() repository.PhaseRepo
	ProjectRepo() repository.ProjectRepo
	ProjectPhaseRepo() repository.ProjectPhaseRepo
	RatingQuotaRepo() repository.RatingQuotaRepo
	ScoreDistributionRepo() repository.ScoreDistributionRepo
	RemarkSettingRepo() repository.RemarkSettingRepo
	TopRemarkRepo() repository.TopRemarkRepo
	BottomRemarkRepo() repository.BottomRemarkRepo
	NotificationRepo() repository.NotificationRepo
}

type repoManager struct {
	infra InfraManager
}

func (r *repoManager) RoleRepo() repository.RoleRepo {
	return repository.NewRoleRepo(r.infra.Conn())
}

func (r *repoManager) UserRepo() repository.UserRepo {
	return repository.NewUserRepo(r.infra.Conn())
}

func (r *repoManager) BusinessUnitRepo() repository.BusinessUnitRepo {
	return repository.NewBusinessUnitRepo(r.infra.Conn())
}

func (r *repoManager) GroupBusinessUnitRepo() repository.GroupBusinessUnitRepo {
	return repository.NewGroupBusinessUnitRepo(r.infra.Conn())
}

func (r *repoManager) ActualScoreRepo() repository.ActualScoreRepo {
	return repository.NewActualScoreRepo(r.infra.Conn())
}

func (r *repoManager) CalibrationRepo() repository.CalibrationRepo {
	return repository.NewCalibrationRepo(r.infra.Conn())
}

func (r *repoManager) PhaseRepo() repository.PhaseRepo {
	return repository.NewPhaseRepo(r.infra.Conn())
}

func (r *repoManager) ProjectRepo() repository.ProjectRepo {
	return repository.NewProjectRepo(r.infra.Conn())
}

func (r *repoManager) ProjectPhaseRepo() repository.ProjectPhaseRepo {
	return repository.NewProjectPhaseRepo(r.infra.Conn())
}

func (r *repoManager) RatingQuotaRepo() repository.RatingQuotaRepo {
	return repository.NewRatingQuotaRepo(r.infra.Conn())
}

func (r *repoManager) ScoreDistributionRepo() repository.ScoreDistributionRepo {
	return repository.NewScoreDistributionRepo(r.infra.Conn())
}

func (r *repoManager) RemarkSettingRepo() repository.RemarkSettingRepo {
	return repository.NewRemarkSettingRepo(r.infra.Conn())
}

func (r *repoManager) TopRemarkRepo() repository.TopRemarkRepo {
	return repository.NewTopRemarkRepo(r.infra.Conn())
}

func (r *repoManager) BottomRemarkRepo() repository.BottomRemarkRepo {
	return repository.NewBottomRemarkRepo(r.infra.Conn())
}

func (r *repoManager) NotificationRepo() repository.NotificationRepo {
	return repository.NewNotificationRepo(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{
		infra: infra,
	}
}
