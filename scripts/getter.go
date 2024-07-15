package scripts

import (
	"context"
)

const (
	ruLang lang = "ru"
	enLang lang = "en"

	defaultLang = ruLang
)

var (
	scripts = map[lang]map[PhraseKey]string{
		ruLang: ru,
		enLang: en,
	}
)

type lang string

func Get(ctx context.Context, key PhraseKey) string {
	return extractScript(ctx)[key]
}

func GetInstructions() map[string]map[string]PhraseKey {
	m := map[string]map[string]PhraseKey{}

	for lng, src := range scripts {
		m[string(lng)] = map[string]PhraseKey{
			src[CreateCartAction]: CreateCartAction,
			src[OpenClearMenu]:    OpenClearMenu,
		}
	}

	return m
}

func GetLang(in string) string {
	_, ok := scripts[lang(in)]
	if ok {
		return in
	}

	return string(ruLang)
}
