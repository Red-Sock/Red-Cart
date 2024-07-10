package scripts

import (
	"context"
)

type lang string

const (
	ruLang = "ru"
	enLang = "en"

	defaultLang = ruLang
)

var scripts = map[lang]map[PhraseKey]string{
	ruLang: ru,
	enLang: en,
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
	_, ok := scripts[lang(in)]
	if ok {
		return in
	}

	return ruLang

}
