package message

import "github.com/KuYaki/waffler_server/internal/models"

func ValidateParser(local int) bool {
	var parser = []models.ParserType{
		models.GPT3_5TURBO,
		models.GPT4,
		models.YakiModel_GPT3_5TURBO,
	}
	var successLocale bool
	for _, l := range parser {
		if models.ParserType(local) == l {
			successLocale = true
			break
		}
	}

	return successLocale
}

func ValidateLocale(lang string) bool {
	var locale = []string{
		"RU", "EN",
	}
	var successLocale bool
	for _, l := range locale {
		if lang == l {
			successLocale = true
			break
		}
	}

	return successLocale
}
