package scripts

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ctxLangKey = "lang"

func EnrichCtx(ctx context.Context, in tgbotapi.Message) context.Context {
	return context.WithValue(ctx, ctxLangKey, lang(in.From.LanguageCode))
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
