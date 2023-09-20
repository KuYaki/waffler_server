package service

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/service/telegram"
	"github.com/KuYaki/waffler_server/internal/models"
)

func updateScoreRecods(records []*models.RecordDTO, source *models.SourceDTO, dataTelegram *telegram.DataTelegram) {
	var typeScore models.ScoreType
	switch dataTelegram.Records[0].ScoreType {
	case models.Waffler:
		typeScore = models.Waffler
	case models.Racism:
		typeScore = models.Racism
	}

	source.WaffelScore = countingScore(records, typeScore) / len(records)

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
