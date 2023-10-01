package service

import (
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/wrapper/data_source"
)

func updateScoreRecods(records []*models.RecordDTO, source *models.SourceDTO, dataTelegram *data_source.DataTelegram) {
	var typeScore models.ScoreType
	switch dataTelegram.Records[0].ScoreType {
	case models.Waffler:
		typeScore = models.Waffler
	case models.Racism:
		typeScore = models.Racism
	}

	source.WafflerScore = countingScore(records, typeScore) / len(records)

}

func countingScore(records []*models.RecordDTO, ScoreType models.ScoreType) int {
	score := 0
	for _, r := range records {
		if r.ScoreType == ScoreType {
			score += r.Score
		}
	}
	return score
}
