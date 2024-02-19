package usecase

import (
	"fmt"
	"time"

	"calibration-system.com/config"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
)

type NotificationUsecase interface {
	NotifyCalibrator() error
	NotifyManager(ids []string, deadline time.Time) error                                               // Send to Manager
	NotifyFirstCurrentCalibrators(data []response.NotificationModel) error                              // First Send Calibrator on Click in Project Active
	NotifyNextCalibrators(data []response.NotificationModel) error                                      // From Previous Phase
	NotifyApprovedCalibrationToCalibrators(data []response.NotificationModel) error                     // Spmo Submit
	NotifySubmittedCalibrationToNextCalibratorsWithoutReview(data response.NotificationModel) error     // Without Spmo Review to Prev Calibrator
	NotifyRejectedCalibrationToCalibrator(id, employee, comment string) error                           // Spmo Reject
	NotifySubmittedCalibrationToSpmo(calibrator *model.User, listOfSpmo []*model.User, phase int) error // Spmo When Submit
	NotifySendBackCalibrators(data []response.NotificationModel) error                                  // Send Back Calibration
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
		URL: fmt.Sprintf("%s/#/login", n.cfg.FrontEndApi),
		// URL:        fmt.Sprintf("%s/#/autologin/%s", n.cfg.FrontEndApi, "token"),
		FirstName:  "Aji",
		Subject:    "Calibration Assignment",
		PhaseOrder: 1,
		Deadline:   "12 November 2023",
	}

	err = utils.SendMail([]string{}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

	data2 := fmt.Sprintf("As a calibrator for Performance Calibration Process, you are requested to complete the phase %d of the calibration process. Please log in to Performance Calibration Process and complete before %s.", emailData.PhaseOrder, "11 November 2023")
	err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, "6285210971537", emailData.FirstName, data2, fmt.Sprintf("http://%s:%s", n.cfg.ApiHost, "3000"))
	if err != nil {
		return err
	}
	return nil
}

func (n *notificationUsecase) NotifyManager(ids []string, deadline time.Time) error {
	for _, calibratorID := range ids {
		employee, err := n.employee.FindById(calibratorID)
		if err != nil {
			return err
		}
		emailData := utils.EmailData{
			// URL: fmt.Sprintf("%s/#/login", n.cfg.FrontEndApi),
			URL:        fmt.Sprintf("%s/#/autologin/%s", n.cfg.FrontEndApi, employee.AccessTokenGenerate),
			FirstName:  employee.Name,
			Subject:    "Calibration Assignment",
			PhaseOrder: 1,
			Deadline:   deadline.Format("02-January-2006"),
		}

		err = utils.SendMail([]string{employee.Email}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
		// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := fmt.Sprintf("As a calibrator for Performance Calibration Process, you are requested to complete the phase %d of the calibration process. Please log in to Performance Calibration Process and complete before %s.", emailData.PhaseOrder, deadline)
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, employee.PhoneNumber, emailData.FirstName, data2, emailData.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *notificationUsecase) NotifyApprovedCalibrationToCalibrators(data []response.NotificationModel) error {
	for _, dataX := range data {
		user, err := n.employee.FindById(dataX.CalibratorID)
		if err != nil {
			return err
		}

		emailData := utils.EmailData{
			FirstName: user.Name,
			Subject:   "Approved Calibration",
		}

		err = utils.SendMail([]string{user.Email}, &emailData, "./utils/templates", "approvedCalibrationEmail.html", n.cfg.SMTPConfig)
		// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "approvedCalibrationEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := "SPMO has approved your calibration worksheet, and it will now be forwarded to the next phase's calibrator. We would greatly appreciate it if you do not disclose these interim results to anyone. Thank you for your attention and cooperation."
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, user.PhoneNumber, emailData.FirstName, data2, fmt.Sprintf("%s/#/login", n.cfg.FrontEndApi))
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *notificationUsecase) NotifySubmittedCalibrationToNextCalibratorsWithoutReview(data response.NotificationModel) error {
	user, err := n.employee.FindById(data.CalibratorID)
	if err != nil {
		return err
	}

	emailData := utils.EmailData{
		FirstName: user.Name,
		Subject:   "Submitted Calibration",
	}

	err = utils.SendMail([]string{user.Email}, &emailData, "./utils/templates", "submitCalibrationWithoutSpmoEmail.html", n.cfg.SMTPConfig)
	// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "approvedCalibrationEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

	data2 := "Your calibration has been submitted, and it will now be forwarded to the next phase's calibrator. We would greatly appreciate it if you do not disclose these interim results to anyone. Thank you for your attention and cooperation."
	err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, user.PhoneNumber, emailData.FirstName, data2, fmt.Sprintf("%s/#/login", n.cfg.FrontEndApi))
	if err != nil {
		return err
	}
	return nil
}

func (n *notificationUsecase) NotifyNextCalibrators(data []response.NotificationModel) error {
	for _, calibratorData := range data {
		employee, err := n.employee.FindById(calibratorData.CalibratorID)
		if err != nil {
			return err
		}

		emailData := utils.EmailData{
			// URL: fmt.Sprintf("%s/#/login", n.cfg.FrontEndApi),
			URL:        fmt.Sprintf("%s/#/autologin/%s/%s/%s/%s", n.cfg.FrontEndApi, employee.AccessTokenGenerate, calibratorData.PreviousBusinessUnitID, calibratorData.PreviousCalibratorID, calibratorData.PreviousCalibrator),
			FirstName:  employee.Name,
			Subject:    "Calibration Assignment",
			PhaseOrder: calibratorData.ProjectPhase,
			Deadline:   calibratorData.Deadline.Format("02-January-2006"),
			Calibrator: calibratorData.PreviousCalibrator,
		}

		err = utils.SendMail([]string{employee.Email}, &emailData, "./utils/templates", "calibratorEmailFromPrevious.html", n.cfg.SMTPConfig)
		// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "calibratorEmailFromPrevious.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		// fmt.Println("ACCESS TOKEN KEY", employee.AccessTokenGenerate)
		data2 := fmt.Sprintf("In regards to Performance Calibration Process, you will act as the next Calibrator previously from %s.", calibratorData.PreviousCalibrator)
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, employee.PhoneNumber, emailData.FirstName, data2, emailData.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *notificationUsecase) NotifySendBackCalibrators(data []response.NotificationModel) error {
	for _, calibratorData := range data {
		employee, err := n.employee.FindById(calibratorData.CalibratorID)
		if err != nil {
			return err
		}

		var url string
		if calibratorData.PreviousBusinessUnitID != "" && calibratorData.PreviousCalibratorID != "" && calibratorData.PreviousCalibrator != "" {
			url = fmt.Sprintf("%s/#/autologin/%s/%s/%s/%s", n.cfg.FrontEndApi, employee.AccessTokenGenerate, calibratorData.PreviousBusinessUnitID, calibratorData.PreviousCalibratorID, calibratorData.PreviousCalibrator)
		} else {
			url = fmt.Sprintf("%s/#/autologin/%s", n.cfg.FrontEndApi, employee.AccessTokenGenerate)

		}
		emailData := utils.EmailData{
			URL:        url,
			FirstName:  employee.Name,
			Subject:    "Calibration Assignment",
			PhaseOrder: calibratorData.ProjectPhase,
			Deadline:   calibratorData.Deadline.Format("02-January-2006"),
			Calibrator: calibratorData.NextCalibrator,
		}

		err = utils.SendMail([]string{employee.Email}, &emailData, "./utils/templates", "sendBackEmail.html", n.cfg.SMTPConfig)
		// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "sendBackEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		// fmt.Println("ACCESS TOKEN KEY", employee.AccessTokenGenerate)
		data2 := fmt.Sprintf("In regards to the Performance Calibration process, %s has sendback their calibrations to you.", calibratorData.PreviousCalibrator)
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, employee.PhoneNumber, emailData.FirstName, data2, emailData.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *notificationUsecase) NotifyFirstCurrentCalibrators(data []response.NotificationModel) error {
	for _, calibratorData := range data {
		employee, err := n.employee.FindById(calibratorData.CalibratorID)
		if err != nil {
			return err
		}

		emailData := utils.EmailData{
			URL:        fmt.Sprintf("%s/#/autologin/%s", n.cfg.FrontEndApi, employee.AccessTokenGenerate),
			FirstName:  employee.Name,
			Subject:    "Calibration Assignment",
			PhaseOrder: calibratorData.ProjectPhase,
			Deadline:   calibratorData.Deadline.Format("02-January-2006"),
		}

		err = utils.SendMail([]string{employee.Email}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
		// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := "We would like to inform you that Calibration 2023 has started. You will act as one of the Calibrators for your team member."
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, employee.PhoneNumber, emailData.FirstName, data2, emailData.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *notificationUsecase) NotifyRejectedCalibrationToCalibrator(id, employee, comment string) error {
	user, err := n.employee.FindById(id)
	if err != nil {
		return err
	}

	emailData := utils.EmailData{
		URL:          fmt.Sprintf("%s/#/autologin/%s", n.cfg.FrontEndApi, user.AccessTokenGenerate),
		FirstName:    user.Name,
		Subject:      "Rejected Calibration",
		Comment:      comment,
		EmployeeName: employee,
	}

	err = utils.SendMail([]string{user.Email}, &emailData, "./utils/templates", "rejectedCalibrationEmail.html", n.cfg.SMTPConfig)
	// err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "rejectedCalibrationEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

	data2 := ("SPMO has rejected your calibration worksheet. Please re-do and re-submit your calibration worksheet.")
	err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, user.PhoneNumber, emailData.FirstName, data2, emailData.URL)
	if err != nil {
		return err
	}

	return nil
}

func (n *notificationUsecase) NotifySubmittedCalibrationToSpmo(calibrator *model.User, listOfSpmo []*model.User, phase int) error {
	fmt.Println("DATA SPMO:", &listOfSpmo[0])
	for _, spmo := range listOfSpmo {
		emailData := utils.EmailData{
			// URL: fmt.Sprintf("%s/#/login", n.cfg.FrontEndApi),
			URL:        fmt.Sprintf("%s/#/autologin-spmo/%s/%s/%s/%d", n.cfg.FrontEndApi, spmo.AccessTokenGenerate, calibrator.ID, *calibrator.BusinessUnitId, phase),
			FirstName:  spmo.Name,
			Subject:    "Submitted Worksheet",
			Calibrator: calibrator.Name,
		}

		err := utils.SendMail([]string{spmo.Email}, &emailData, "./utils/templates", "spmoEmail.html", n.cfg.SMTPConfig)
		// err := utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "spmoEmail.html", n.cfg.SMTPConfig)
		if err != nil {
			return err
		}

		data2 := fmt.Sprintf("%s on Performance Calibration Process has submitted the calibration worksheet. Please review and approve as soon as possible to proceed to the next phase.", emailData.Calibrator)
		err = utils.SendWhatsAppNotif(n.cfg.WhatsAppConfig, spmo.PhoneNumber, emailData.FirstName, data2, emailData.URL)
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
