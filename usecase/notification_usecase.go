package usecase

import (
	"fmt"
	"time"

	"calibration-system.com/config"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
)

type NotificationUsecase interface {
	NotifyCalibrator() error
	NotifyCalibrators(ids []string, deadline time.Time) error
	NotifyApprovedCalibrationToCalibrator(ids []string) error
	NotifyRejectedCalibrationToCalibrator(id, comment string) error
	NotifyCalibrationToSpmo(calibrator *model.User, listOfSpmo []*model.User) error
}

type notificationUsecase struct {
	repo     repository.NotificationRepo
	employee UserUsecase
	project  ProjectUsecase
	cfg      config.Config
}

func (n *notificationUsecase) NotifyCalibrator() error {
	year, month, day := time.Now().Date()
	fmt.Println("Tanggal sekarang", year, month, day)

	projectPhases, err := n.project.FindActiveProjectPhase()
	if err != nil {
		return err
	}

	for _, projectPhase := range projectPhases {
		year_pp, month_pp, day_pp := projectPhase.StartDate.Date()
		if year == year_pp && month == month_pp && day == day_pp {

		}
		fmt.Println(year_pp, month_pp, day_pp)

		email, err := n.repo.GetCalibratorEmailOnProjectPhase(projectPhase.ID)
		if err != nil {
			return err
		}

		// Disini kirim email dan wanya
		fmt.Println("Email CALIBRATOR", email)
	}

	emailData := utils.EmailData{
		URL:        fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"),
		FirstName:  "Aji",
		Subject:    "Calibration Assignment",
		PhaseOrder: 1,
		Deadline:   "12 November 2023",
	}

	err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

	data2 := fmt.Sprintf("As a calibrator for Calibration System, you are requested to complete the phase %d of the calibration process. Please log in to Calibration System and complete before %s.", emailData.PhaseOrder, "11 November 2023")
	err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, "6285210971537", emailData.FirstName, data2, fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"))
	if err != nil {
		return err
	}
	return nil
}

func (n *notificationUsecase) NotifyCalibrators(ids []string, deadline time.Time) error {
	for _, calibratorID := range ids {
		employee, err := n.employee.FindById(calibratorID)
		if err != nil {
			return err
		}

		emailData := utils.EmailData{
			URL:        fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"),
			FirstName:  employee.Name,
			Subject:    "Calibration Assignment",
			PhaseOrder: 1,
			Deadline:   deadline.Format("2006-01-02"),
		}

		err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := fmt.Sprintf("As a calibrator for Calibration System, you are requested to complete the phase %d of the calibration process. Please log in to Calibration System and complete before %s.", emailData.PhaseOrder, "11 November 2023")
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, employee.PhoneNumber, emailData.FirstName, data2, fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *notificationUsecase) NotifyApprovedCalibrationToCalibrator(ids []string) error {
	for _, id := range ids {
		user, err := n.employee.FindById(id)
		if err != nil {
			return err
		}

		emailData := utils.EmailData{
			URL:       "http://localhost:3000/",
			FirstName: user.Name,
			Subject:   "Approved Calibration",
		}

		err = utils.SendMail([]string{user.Email}, &emailData, "./utils/templates", "approvedCalibrationEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := fmt.Sprintf("SPMO has approved your calibration worksheet, and it will now be forwarded to the next phase's calibrator. We would greatly appreciate it if you do not disclose these interim results to anyone. Thank you for your attention and cooperation.")
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, user.PhoneNumber, emailData.FirstName, data2, fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"))
		if err != nil {
			return err
		}
	}

	// emailData := utils.EmailData{
	// 	URL:       "http://localhost:3000/",
	// 	FirstName: "Aji Inisti Udma Wijaya",
	// 	Subject:   "Approved Calibration",
	// }

	// err := utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "approvedCalibrationEmail.html", n.cfg.SMTPConfig)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (n *notificationUsecase) NotifyRejectedCalibrationToCalibrator(id, comment string) error {
	user, err := n.employee.FindById(id)
	if err != nil {
		return err
	}

	emailData := utils.EmailData{
		URL:       "http://localhost:3000/",
		FirstName: user.Name,
		Subject:   "Rejected Calibration",
		Comment:   comment,
	}

	err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "rejectedCalibrationEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

	data2 := fmt.Sprintf("SPMO has rejected your calibration worksheet. Please re-do and re-submit your calibration worksheet.")
	err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, user.PhoneNumber, emailData.FirstName, data2, fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"))
	if err != nil {
		return err
	}

	return nil
}

func (n *notificationUsecase) NotifyCalibrationToSpmo(calibrator *model.User, listOfSpmo []*model.User) error {
	for _, spmo := range listOfSpmo {
		emailData := utils.EmailData{
			URL:        "http://localhost:3000/",
			FirstName:  spmo.Name,
			Subject:    "Submitted Worksheet",
			Calibrator: calibrator.Name,
		}

		err := utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "spmoEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := fmt.Sprintf("%s on Calibration System has submitted the calibration worksheet. Please review and approve as soon as possible to proceed to the next phase.", emailData.FirstName)
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, spmo.PhoneNumber, emailData.FirstName, data2, fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"))
		if err != nil {
			return err
		}
	}
	return nil
}

func NewNotificationUsecase(repo repository.NotificationRepo, employee UserUsecase, project ProjectUsecase, cfg config.Config) NotificationUsecase {
	return &notificationUsecase{
		repo:     repo,
		employee: employee,
		project:  project,
		cfg:      cfg,
	}
}
