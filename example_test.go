package dialogs_test

import (
	"log"
	"strings"

	"gitlab.com/toby3d/dialogs"
)

var (
	questions dialogs.Questions
	answers   dialogs.Answers

	answer   dialogs.Answer
	question = dialogs.Question{
		Meta: dialogs.Meta{
			ClientID: "Developer Console",
			Locale:   "ru-RU",
			TimeZone: "UTC",
		},
		Request: dialogs.Request{
			Command:           "привет",
			OriginalUtterance: "привет",
			Type:              dialogs.TypeSimpleUtterance,
		},
		Session: dialogs.Session{
			MessageID: 42,
			New:       false,
			SessionID: "1ab234cd-56e7890f-gh1j23k-45l6",
			SkillID:   "ab1c2d34-5e67-8f90-g12h-3456jkl78901",
			UserID:    "1A2BC3456789D0E12F345GHJ67890K1LMN23OP4567QR8901234567ST89UV0W12",
		},
		Version: "1.0",
	}

	textButton, urlButton, dataButton dialogs.Button
	buttons                           []dialogs.Button
)

func errCheck(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func Example_fastStart() {
	log.Println("Стартуем!..")
	questions, answers = dialogs.New("127.0.0.1:2368", "/alice", "", "")

	for question := range questions {
		switch {
		case strings.EqualFold(question.Request.Command, "привет"):
			// Это команда приветствия. Надо ответить взаимностью!

			// Готовим ответ на реплику и приветствуем пользователя.
			answer = dialogs.NewAnswer(question, "Привет!")

			// Корректно озвучиваем реплику
			answer.Response.TTS = "прив+ет!"

			// Результат отправляем в канал
			answers <- answer
		case question.Request.Command != "":
			// Это какая-то команда, которую мы не знаем. Нужно извиниться.
			answer = dialogs.NewAnswer(question, "Простите, я не поняла.")
			answer.Response.TTS = "Прост+ите, я не понял+а."
			answers <- answer
		default:
			continue // Это что-то совсем иное - ничего не делаем.
		}
	}
}

func ExampleNew() {
	// New принимает аргументы в следующем порядке:
	// * Локальный адрес и порт.
	// * Роут по которому нужно слушать входящий трафик; желательно использовать
	//   уникальный и секретный путь, по которому можно однозначно
	//   идентифицировать трафик как запрос от Яндекса.
	// * Файл сертификата (если необходим).
	// * Файл ключа сертификата (если необходим).
	//
	// В случае ошибки возникнет паника. В случае успеха будет создан роутер по
	// указанному адресу:порту/пути который будет слушать входящий трафик.
	//
	// В ответ будут возвращены два канала: для чтения запросов и отправки
	// ответов соответственно.
	questions, answers = dialogs.New("127.0.0.1:2368", "/alice", "", "")
}

func ExampleNewAnswer() {
	// Привязываем новый ответ к идентификаторам входящей реплики
	answer = dialogs.NewAnswer(question, "Прощай, жестокий мир!")

	// Можно дополнить реплику дополнительными возможностями вроде кнопок,
	// озвучки и/или параметром, обозначающим конец разговора и выхода из Навыка.
	answer.Response.TTS = "Прощ+ай, жест+окий м+ир!"
	answer.Response.Buttons = buttons
	answer.Response.EndSession = true

	// Оформленный ответ нужно отправить не позднее 1,5 секунд после получения
	// реплики пользователя.
	answers <- answer
}

func ExampleNewButton() {
	textButton = dialogs.NewButton("я просто кнопка")

	dataButton = dialogs.NewButton("я кнопка с данными")
	var payload dialogs.Payload
	// Произвольные данные должны быть в формате JSON
	payload = struct {
		Count int    `json:"count"`
		Word  string `json:"word"`
	}{
		Count: 42,
		Word:  "Алиса",
	}
	dataButton.Payload = payload

	urlButton = dialogs.NewButton("я ссылка")
	urlButton.URL = "https://toby3d.gitlab.io"
}

func ExampleNewButtons() {
	answer.Response.Buttons = dialogs.NewButtons(textButton, dataButton, urlButton)
}
