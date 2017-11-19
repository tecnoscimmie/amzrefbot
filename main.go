package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	tgAPIKey          string
	botInstance       *tgbotapi.BotAPI
	certPath          string
	keyPath           string
	debug             bool
	domain            string
	port              string
	refs              Refs
	refsFilePath      string
	userAdministrator int64
)

func main() {
	setupParams()
	var err error

	// Load a Refs instance
	if refsFilePath != "" {
		refs = NewRefsFromPath(refsFilePath)
	} else {
		refs = NewRefs()
	}

	botInstance, err = setupTelegramBot(tgAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	// set telegram api lib debug level to the same one we use globally
	botInstance.Debug = debug

	log.Printf("Authorized on account %s", botInstance.Self.UserName)
	log.Println("Referral codes active:")
	for _, refcode := range refs.ReferralCodes {
		log.Printf("\t - %s, with code %s\n", refcode.AssociatedUser, refcode.Code)
	}

	updates := botInstance.ListenForWebhook("/" + tgAPIKey)
	go startServer()
	for update := range updates {
		// handle every update in a separate goroutine
		go handleUpdate(update)
	}
}

func setupTelegramBot(apiKey string) (bot *tgbotapi.BotAPI, err error) {
	bot, err = tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://" + domain + ":" + port + "/" + apiKey))
	if err != nil {
		log.Fatal(err)
	}

	return
}

func startServer() {
	go log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+port, certPath, keyPath, nil))
}

func setupParams() {
	flag.StringVar(&certPath, "cert", "", "required, TLS certificate path")
	flag.StringVar(&keyPath, "key", "", "required, TLS key path")
	flag.StringVar(&tgAPIKey, "apikey", "", "required, Telegram bot API key")
	flag.StringVar(&domain, "domain", "", "required, domain associated to the TLS cert+key and the server where this bot will be running")
	flag.StringVar(&port, "port", "88", "port to run on, must be 443, 80, 88, 8443")
	flag.BoolVar(&debug, "debug", false, "debug Telegram bot interactions")
	flag.StringVar(&refsFilePath, "refsfile", "", "file containing referral codes, written by this bot")
	flag.Int64Var(&userAdministrator, "admin", 0, "user capable of accepting or denying referral link add requests")
	flag.Parse()

	if certPath == "" || keyPath == "" || tgAPIKey == "" || domain == "" || userAdministrator == 0 {
		flag.Usage()
		os.Exit(1)
	}
}
