package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/stretchr/testify/assert"
)

const accuracy = 0.001

func TestAverageRacismScore(t *testing.T) {
	testData := []struct {
		records []models.RacismDTO
		score   float64
	}{
		{score: 0, records: []models.RacismDTO{{}}},
		{score: 10, records: []models.RacismDTO{{CreatedTs: time.Now(), Score: 10}}},
		{score: 10, records: []models.RacismDTO{{CreatedTs: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 10}}},
		{score: 1.368, records: []models.RacismDTO{
			{CreatedTs: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 10},
			{CreatedTs: time.Now(), Score: 0},
		}},
		{score: 1.368, records: []models.RacismDTO{
			{CreatedTs: time.Now().Add(-20 * time.Millisecond * time.Duration(year)), Score: 10},
			{CreatedTs: time.Now(), Score: 0},
		}},
		{score: 94.786, records: []models.RacismDTO{
			{CreatedTs: time.Now().Add(-5 * time.Millisecond * time.Duration(year)), Score: 100},
			{CreatedTs: time.Now(), Score: 0},
		}},
		{score: 67.738, records: []models.RacismDTO{
			{CreatedTs: time.Now().Add(-5 * time.Millisecond * time.Duration(year)), Score: 78},
			{CreatedTs: time.Now(), Score: 0},
		}},
		{score: 61.050, records: []models.RacismDTO{
			{CreatedTs: time.Now().Add(-5 * time.Millisecond * time.Duration(year)), Score: 78},
			{CreatedTs: time.Now(), Score: 24},
		}},
	}
	for count, tt := range testData {
		t.Run(fmt.Sprint("test ", count), func(t *testing.T) {
			score := averageRacismScore(tt.records)
			assert.InDelta(t, tt.score, score, accuracy)
		})
	}
}
func TestAverageWafflerScore(t *testing.T) {
	testData := []struct {
		records []models.WafflerDTO
		score   float64
	}{
		{score: 0, records: []models.WafflerDTO{{}}},
		{score: 10, records: []models.WafflerDTO{{CreatedTsAfter: time.Now(), Score: 10}}},
		{score: 10, records: []models.WafflerDTO{{CreatedTsAfter: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 10}}},
		{score: 6.131, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), Score: 10},
			{CreatedTsAfter: time.Now(), Score: 0},
		}},
		{score: 6.131, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 10},
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 0},
		}},
		{score: 6.131, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-20 * time.Millisecond * time.Duration(year)), Score: 10},
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 0},
		}},
		{score: 5.284, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 10},
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-time.Millisecond * time.Duration(40*year/9)), Score: 0},
		}},
		{score: 98.605, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 100},
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-time.Millisecond * time.Duration(40*year/9)), Score: 0},
		}},
		{score: 75.075, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 78},
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-time.Millisecond * time.Duration(40*year/9)), Score: 0},
		}},
		{score: 72.316, records: []models.WafflerDTO{
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-10 * time.Millisecond * time.Duration(year)), Score: 78},
			{CreatedTsAfter: time.Now(), CreatedTsBefore: time.Now().Add(-time.Millisecond * time.Duration(40*year/9)), Score: 24},
		}},
	}
	for count, tt := range testData {
		t.Run(fmt.Sprint("test ", count), func(t *testing.T) {
			score := averageWafflerScore(tt.records)
			assert.InDelta(t, tt.score, score, accuracy)
		})
	}
}
