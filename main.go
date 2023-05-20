package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"log"
	"os"
	"os/signal"
	"syscall"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "github.com/mattn/go-sqlite3"
)

func startHttp() {
	os.Setenv("TZ", "Asia/Jakarta")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	app := fiber.New()
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendString("Mau ngapain hayoo")
	})
	log.Fatal(app.Listen(":" + PORT))
}

func registerHandler(client *whatsmeow.Client) func(evt interface{}) {
return func(evt interface{}) {
	if msg := evt.(*events.Message) != nil {
		if msg.Info.Chat.String() == "status@broadcast" {
			client.MarkRead([]types.MessageID{msg.Info.ID}, msg.Info.Timestamp, msg.Info.Chat, msg.Info.Sender)
		}
	}
}
}

func startClient(nama string) {
	if nama == "" { log.Fatal("HALAH KOSONG") }
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:"+nama+".db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	eventHandle := eventHandler(client)
	client.AddEventHandler(eventHandle)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil { panic(err) }
		fmt.Println("Login Success")
		client.SendPresence(types.PresenceAvailable)
		client.SendPresence(types.PresenceUnavailable)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
}

func main() {
	go startHttp()
	startClient("jshhshsj")
}