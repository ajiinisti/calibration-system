package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"calibration-system.com/config"
)

func SendWhatsAppNotif(config config.WhatsAppConfig, phoneNumber, data1, data2, data3 string) error {
	// Ganti nomor telpon ke phoneNumber, ini hanya untuk UAT
	if phoneNumber != "" {
		formData := bytes.NewBufferString(fmt.Sprintf(`message={"11": "%s","22": "%s","33": "%s"}&template_id=%s&api_key=%s&shorten_url=%s&to_no=%s`,
			data1, data2, data3, config.TemplateID, config.ApiKey, config.ShortenUrl, phoneNumber))
		fmt.Println(formData)
		resp2, err := http.Post(config.URL, "application/x-www-form-urlencoded", formData)
		if err != nil {
			// handle error
			return err
		}

		rbody2, err := ioutil.ReadAll(resp2.Body)
		if err != nil {
			fmt.Println("Error pesan", err.Error())
		}
		fmt.Println(string(rbody2))
	}
	return nil
}
