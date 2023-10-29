package service

import (
	"math"
	"time"

	"github.com/KuYaki/waffler_server/internal/models"
)

const (
	q            = 9                         // weight multiplier that decreases weights of older messages by 10 in 10 years
	p    float64 = 0.01                      // weight multiplier that decreases weights for low scores
	MAX  float64 = 100                       // maximum score possible
	year float64 = 1000 * 60 * 60 * 24 * 365 // 1 year in milliseconds
	T1           = 10 * year
)

// Weights are normalized from 1 to 0.1 with the coefficient 1/(1 + q*x).
// where x is min(t - t0, 10 years).
// This means that after 10 years there is no dropoff in weights.
// To amplify higher scores we multiply weights by p^(1 - (s / MAX)).
// Then we return the sum of weighted scores divided by the sum of the weights.
func averageRacismScore(records []models.RacismDTO) models.NullFloat64 {
	currentTime := time.Now().UnixMilli()
	var sumScores float64
	var sumWeights float64

	for _, r := range records {
		if !r.Score.Valid {
			continue
		}
		weight := 1 / (1 + q*math.Min(float64(currentTime)-float64(r.CreatedTs.UnixMilli()), T1)/float64(T1))
		weight = weight * math.Pow(p, 1-(float64(r.Score.Int64)/MAX))
		sumScores += weight * float64(r.Score.Int64)
		sumWeights += weight
	}

	if sumWeights == 0 {
		return models.NullFloat64{}
	}
	return models.NewNullFloat64(sumScores / sumWeights)
}

// Weights are normalized from 1 to 0.1 with the coefficients 1/(1 + q*x) and
// 1/(1 + q*y) where x = min(t - t_2, 10 years) and y = min(t2-t1, 10 years).
// This means that after 10 years there is no dropoff in weights.
// Then to re-normalize we take the geometric mean of the 2 weights.
// To amplify higher scores we multiply weights by p^(1 - (s / MAX)).
// Then we return the sum of weighted scores divided by the sum of the weights.
func averageWafflerScore(records []models.WafflerDTO) models.NullFloat64 {
	currentTime := time.Now().UnixMilli()
	var sumScores float64
	var sumWeights float64

	for _, r := range records {
		if !r.Score.Valid {
			continue
		}
		weight_1 := 1 / (1 + q*math.Min(float64(currentTime)-float64(r.CreatedTsAfter.UnixMilli()), T1)/float64(T1))
		weight_2 := 1 / (1 + q*math.Min(float64(r.CreatedTsAfter.UnixMilli())-float64(r.CreatedTsBefore.UnixMilli()), T1)/float64(T1))
		weight := math.Sqrt(weight_1 * weight_2)
		weight = weight * math.Pow(p, 1-(float64(r.Score.Int64)/MAX))
		sumScores += weight * float64(r.Score.Int64)
		sumWeights += weight
	}

	if sumWeights == 0 {
		return models.NullFloat64{}
	}
	return models.NewNullFloat64(sumScores / sumWeights)
}
