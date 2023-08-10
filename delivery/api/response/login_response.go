package response

import "calibration-system.com/model"

type LoginResponse struct {
	AccessToken string
	TokenModel  model.TokenModel
}
