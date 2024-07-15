package scripts

import (
	"context"
)

const (
	ctxLangKey ctxLangKeyType = "lang"
)

type ctxLangKeyType string

func extractScript(ctx context.Context) (resp map[PhraseKey]string) {
	return scripts[GetLangFromCtx(ctx)]
}

func GetLangFromCtx(ctx context.Context) lang {
	v := ctx.Value(ctxLangKey)
	l, ok := v.(lang)
	if ok {
		return l
	}

	return defaultLang
}
