package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"calibration-system.com/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL          string
	FirstName    string
	Subject      string
	PhaseOrder   int
	Comment      string
	Calibrator   string
	Deadline     string
	EmployeeName string
}

func SendMail(to []string, data *EmailData, templateDir string, templatePath string, cfg config.SMTPConfig) error {
	emailTemplate, err := ParseTemplateFile(templateDir, templatePath)
	if err != nil {
		return fmt.Errorf("Could not parse email template: " + err.Error())
	}

	var emailContent bytes.Buffer
	err = emailTemplate.Execute(&emailContent, &data)
	if err != nil {
		return fmt.Errorf("Error executing template: " + err.Error())
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", cfg.SMTPSenderName)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", data.Subject)
	fmt.Println("SCONTENT BODY", emailContent.String())
	mailer.SetBody("text/html", emailContent.String())
	// mailer.AddAlternative("text/plain", html2text.HTML2Text(emailContent.String()))

	port, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(
		cfg.SMTPHost,
		port,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}

func ParseTemplateFile(templateDir, templatePath string) (*template.Template, error) {
	// Combine templateDir and templatePath to get the absolute path to the template file
	templateFullPath := filepath.Join(templateDir, templatePath)

	// Parse the template file
	template, err := template.ParseFiles(templateFullPath)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
