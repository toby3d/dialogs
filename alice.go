package alice

import (
	"bytes"
	"log"
	"strings"
	"time"

	"github.com/kirillDanshin/dlog"
	json "github.com/pquerna/ffjson/ffjson"
	http "github.com/valyala/fasthttp"
)

type (
	// Questions является вебхук-каналом входящих запросов от пользователя.
	Questions <-chan Question

	// Answers является вебхук-каналом исходящих ответов к пользователям.
	Answers chan Answer
)

// Listen создаёт простой роутер для прослушивания входящих данных по вебхуку и
// возвращает два канала: для чтения запросов и отправки ответов соответственно.
func Listen(addr, path string) (Questions, Answers) {
	var err error
	questions := make(chan Question)
	answers := make(chan Answer)

	handleFunc := func(ctx *http.RequestCtx) {
		dlog.Ln("Тело входящего запроса:")
		dlog.D(&ctx.Request)

		if !bytes.HasPrefix(ctx.Path(), []byte(path)) {
			dlog.Ln("Получен неподдерживаемый запрос")
			return
		}
		dlog.Ln("Получен поддерживаемый запрос")

		dlog.Ln("Декодируем запрос...")
		var question Question
		if err = json.Unmarshal(ctx.Request.Body(), &question); err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		dlog.Ln("Отправляем запрос в канал...")
		questions <- question

		var answer Answer
		for answer = range answers {
			if !strings.EqualFold(answer.Session.SessionID, question.Session.SessionID) ||
				!strings.EqualFold(answer.Session.UserID, question.Session.UserID) ||
				answer.Session.MessageID != question.Session.MessageID {
				dlog.Ln("Это не тот ответ...")
				continue
			}

			dlog.Ln("Обнаружен подходящий запрос! Отвечаем...")
			dlog.D(answer)
			break
		}

		dlog.Ln("Дождались нужный ответ! Отправляем его...")
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.SetStatusCode(http.StatusOK)

		dlog.Ln("Кодируем ответ...")
		if err = json.NewEncoder(ctx).Encode(answer); err != nil {
			dlog.Ln("Ошибка:", err.Error())
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}

		dlog.Ln("Готово, ответ доставлен!")
	}
	handleFunc = http.TimeoutHandler(handleFunc, 1500*time.Millisecond, "oh no")

	go func() {
		dlog.Ln("Создаём простой роутер...")
		if err = http.ListenAndServe(addr, handleFunc); err != nil {
			log.Fatalln("Ошибка:", err.Error())
		}
	}()

	return questions, answers
}
