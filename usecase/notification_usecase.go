package usecase

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"calibration-system.com/config"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
)

type NotificationUsecase interface {
	NotifyCalibrator() error
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
		URL:        "http://localhost:3000/",
		FirstName:  "Aji Inisti Udma Wijaya",
		Subject:    "Calibration Assignment",
		PhaseOrder: 1,
		Deadline:   "12 November 2023",
	}

	err = utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "calibratorEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

	// form := url.Values{}
	// message := fmt.Sprintf(`{
	// 	"11": "%s",
	// 	"22": "As a calibrator for Calibration System, you are requested to complete the phase %d of the calibration process. Please log in to Calibration System and
	// 	complete before %s.",
	// 	"33": "%s"}`,
	// 	emailData.FirstName,
	// 	emailData.PhaseOrder,
	// 	"11 November 2023",
	// 	fmt.Sprintf("%s:%s", n.cfg.ApiHost, "3000"),
	// )
	// // objectString := fmt.Sprintf("%+v", myObject)
	// form.Add("message", message)
	// form.Add("template_id", n.cfg.WhatsAppConfig.TemplateID)
	// form.Add("api_key", n.cfg.WhatsAppConfig.ApiKey)
	// form.Add("shorten_url", n.cfg.WhatsAppConfig.ShortenUrl)
	// form.Add("to_no", "6285210971537")

	// resp, err := http.PostForm(n.cfg.WhatsAppConfig.URL, form)
	// if err != nil {
	// 	fmt.Println("Error pesan", err.Error())
	// 	return err
	// }
	// defer resp.Body.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	message := fmt.Sprintf(`{
		"11": "%s",
		"22": "As a calibrator for Calibration System, you are requested to complete the phase %d of the calibration process. Please log in to Calibration System and complete before %s.",
		"33": "%s"}`,
		emailData.FirstName,
		emailData.PhaseOrder,
		"11 November 2023",
		fmt.Sprintf("%s:%s", n.cfg.ApiHost, "3000"),
	)
	writer.WriteField("message", message)
	writer.WriteField("template_id", n.cfg.WhatsAppConfig.TemplateID)
	writer.WriteField("api_key", n.cfg.WhatsAppConfig.ApiKey)
	writer.WriteField("shorten_url", n.cfg.WhatsAppConfig.ShortenUrl)
	writer.WriteField("to_no", "6285210971537")

	writer.Close()
	contentType := writer.FormDataContentType()

	resp, err := http.Post(n.cfg.WhatsAppConfig.URL, contentType, &buf)
	if err != nil {
		fmt.Println("Error pesan", err.Error())
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error pesan", err.Error())
	}

	requestBody, err := ioutil.ReadAll(resp.Request.Body)
	if err != nil {
		fmt.Println("Error pesan", err.Error())
	}

	fmt.Println(string(requestBody))
	fmt.Println(string(responseBody))

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
	}

	emailData := utils.EmailData{
		URL:       "http://localhost:3000/",
		FirstName: "Aji Inisti Udma Wijaya",
		Subject:   "Approved Calibration",
	}

	err := utils.SendMail([]string{"aji.wijaya@techconnect.co.id"}, &emailData, "./utils/templates", "approvedCalibrationEmail.html", n.cfg.SMTPConfig)
	if err != nil {
		return err
	}

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
