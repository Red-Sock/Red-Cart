package scripts

import (
	"context"
)

const (
	ctxLangKey ctxLangKeyType = "lang"
)

type ctxLangKeyType string

func extractScript(ctx context.Context) (resp map[PhraseKey]string) {
	v := ctx.Value(ctxLangKey)
	l, ok := v.(lang)
	if ok {
		resp, ok = scripts[l]
	}

	if ok {
		return resp
	}

	return scripts[defaultLang]
}
