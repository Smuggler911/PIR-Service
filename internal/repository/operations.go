package repository

import (
	"pir-serverSide/internal/repository/models"
	"strconv"
)

func convertStringUint(strValue string) uint64 {
	Value, err := strconv.ParseUint(strValue, 10, 64)
	if err != nil {
		return 0
	}
	return Value
}

func convertStringInt(strValue string) int32 {
	Value, err := strconv.ParseInt(strValue, 10, 32)
	if err != nil {
		return 0
	}
	return int32(Value)

}
func calcAveragePoints(reviews []models.Review) float64 {
	if len(reviews) != 0 {
		result := float64(0)
		for i := 0; i < len(reviews); i++ {
			rate := reviews[i].Rate
			result = result + float64(rate)
		}
		return round(result/float64(len(reviews)), 1)
	} else {
		return 0
	}
}

func likeCount(likes []models.Like) float64 {
	if len(likes) != 0 {
		result := float64(0)
		for i := 0; i < len(likes); i++ {
			like := float64(likes[i].Liked)
			result = result + like
		}
		return result

	} else {
		return 0
	}
}

func viewCount(views []models.Views) float64 {
	if len(views) != 0 {
		result := float64(0)
		for i := 0; i < len(views); i++ {
			view := float64(views[i].Viewed)
			result = result + view
		}
		return result
	} else {
		return 0
	}
}

func round(number float64, decimals int) float64 {
	output := strconv.FormatFloat(number, 'f', decimals, 64)
	result, _ := strconv.ParseFloat(output, 64)
	return result
}
