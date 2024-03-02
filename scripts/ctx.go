package scripts

import (
	"context"

	"github.com/Red-Sock/go_tg/model"
)

const ctxLangKey = "lang"

func EnrichCtx(in *model.MessageIn) context.Context {
	return context.WithValue(context.Background(), ctxLangKey, lang(in.From.LanguageCode))
}

func extractLang(ctx context.Context) lang {
	v := ctx.Value(ctxLangKey)
	l, ok := v.(lang)
	if ok {
		_, ok = scripts[l]
	}

	if ok {
		return l
	}

	return defaultLang
}

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
