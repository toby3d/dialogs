package alice

import (
	"bytes"
	"log"

	"github.com/kirillDanshin/dlog"
	json "github.com/pquerna/ffjson/ffjson"
	http "github.com/valyala/fasthttp"
)

// Questions является вебхук-каналом входящих запросов пользователя.
type Questions <-chan Question

// Listen создаёт роутер для прослушивания входящих запросов с дальнейшей их
// пересылкой в канал.
func Listen(addr, path string) Questions {
	var err error
	channel := make(chan Question)

	handleFunc := func(ctx *http.RequestCtx) {
		dlog.Ln("Тело входящего запроса:")
		dlog.D(ctx.Request)

		if !bytes.HasPrefix(ctx.Path(), []byte(path)) {
			dlog.Ln("Получен неподдерживаемый запрос")
			return
		}
		dlog.Ln("Получен поддерживаемый запрос")

		var question Question
		dlog.Ln("Декодируем запрос...")
		if err = json.Unmarshal(ctx.Request.Body(), &question); err != nil {
			log.Println("Ошибка:", err.Error())
			return
		}

		dlog.Ln("Отправляем запрос в канал...")
		channel <- question
	}

	go func() {
		dlog.Ln("Создаём простой роутер...")
		if err = http.ListenAndServe(addr, handleFunc); err != nil {
			log.Fatalln("Ошибка:", err.Error())
		}
	}()

	return channel
}
