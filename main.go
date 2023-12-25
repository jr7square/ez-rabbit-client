package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jr7square/ezrabbitclient/fileserver"
	"github.com/jr7square/ezrabbitclient/rabbit"
)

func main() {
	doPrompt := flag.Bool("noPrompt", true, "no prompt fo ez development")
	flag.Parse()
	dialUrl := "amqp://guest:guest@localhost:5672/"
	if *doPrompt {
		dialUrl = rabbit.PromptForDialURL()
	}
	rabbit.Producer = rabbit.InitRabbit(dialUrl)
	defer rabbit.Producer.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fileserver.InitFileServer(r)

	r.Get("/testRabbit", func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("Connecting to rabbit %s", dialUrl)
		w.Write([]byte(message))
	})
	r.Get("/", getSendMessagePage)
	r.Post("/sendMessage", postSendMessage)
	http.ListenAndServe(":3000", r)

}

func getSendMessagePage(w http.ResponseWriter, r *http.Request) {
	page(sendMessage("", "", "", "")).Render(r.Context(), w)
}

func postSendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	message := r.Form.Get("queueMessage")
	exchange := r.Form.Get("exchange")
	routingKey := r.Form.Get("routingKey")
	fmt.Println("textarea value", r.Form.Get("queueMessage"))
	if err := rabbit.Producer.Send(exchange, routingKey, message); err != nil {
		log.Println("Failed to send message", err.Error())
		page(sendMessage(message, exchange, routingKey, err.Error())).Render(r.Context(), w)
		return
	}
	page(sendMessage(message, exchange, routingKey, "")).Render(r.Context(), w)
}
