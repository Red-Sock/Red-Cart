package scripts

import (
	"context"
)

type lang string

const (
	ruLang = "ru"

	defaultLang = ruLang
)

var scripts = map[lang]map[PhraseKey]string{
	ruLang: ru,
}

func Get(ctx context.Context, key PhraseKey) string {
	return extractScript(ctx)[key]
}

func GetInstructions() map[string]map[string]PhraseKey {
	m := map[string]map[string]PhraseKey{}

	for lng, src := range scripts {
		m[string(lng)] = map[string]PhraseKey{
			src[CreateCart]:  CreateCart,
			src[OpenSetting]: OpenSetting,
			src[Clear]:       Clear,
		}
	}

	return m
}

func GetLang(in string) string {
	switch in {
	case ruLang:
		return ruLang
	default:
		return ruLang
	}
}
