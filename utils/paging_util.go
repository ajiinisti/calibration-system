package utils

import (
	"math"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
)

func GetPaginationParams(params request.PaginationParam) model.PaginationQuery {
	var page int
	var take int
	var skip int

	if params.Page > 0 {
		page = params.Page
	} else {
		page = 1
	}
	if params.Limit > 0 {
		take = params.Limit
	} else {
		take = 10
	}
	skip = (page-1)*take + params.Offset
	return model.PaginationQuery{
		Page:                          page,
		Take:                          take,
		Skip:                          skip,
		Name:                          params.Name,
		SupervisorName:                params.SupervisorName,
		Grade:                         params.Grade,
		EmployeeName:                  params.EmployeeName,
		OrderGrade:                    params.OrderGrade,
		OrderEmployeeName:             params.OrderEmployeeName,
		OrderCalibrationScore:         params.OrderCalibrationScore,
		OrderCalibrationRating:        params.OrderCalibrationRating,
		FilterCalibrationRating:       params.FilterCalibrationRating,
		CalibrationPhaseBefore:        params.CalibrationPhaseBefore,
		OrderCalibrationScoreBefore:   params.OrderCalibrationScoreBefore,
		OrderCalibrationRatingBefore:  params.OrderCalibrationRatingBefore,
		FilterCalibrationRatingBefore: params.FilterCalibrationRatingBefore,
		RatingChangedStatus:           params.RatingChangedStatus,
		RatingChanged:                 params.RatingChanged,
	}
}

func Paginate(page, limit, totalRows int) response.Paging {
	return response.Paging{
		Page:        page,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(limit))),
		TotalRows:   totalRows,
		RowsPerPage: limit,
	}
}
