package main

import (
    "context"
    "flag"
    "fmt"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    "os"
    "os/signal"
    "syscall"
    waLog "go.mau.fi/whatsmeow/util/log"
    _ "github.com/ncruces/go-sqlite3/driver"
    _ "github.com/ncruces/go-sqlite3/embed"
)

type Ev struct {
    ChatPresence *events.ChatPresence
    Message *events.Message
    Receipt *events.Receipt
    More interface {}
}
func Connect(nomor string, cb func(conn *whatsmeow.Client, evt Ev)) {
    dbLog := waLog.Stdout("Database", "ERROR", true)
    container, err := sqlstore.New("sqlite3", "file:"+nomor+".db?nolock=1", dbLog)
    if err != nil { fmt.Println("GoError:",err); return }
    deviceStore, err := container.GetFirstDevice()
    if err != nil { fmt.Println("GoError:",err); return }
    clientLog := waLog.Stdout("Client", "ERROR", true)
    client := whatsmeow.NewClient(deviceStore, clientLog)
    client.AddEventHandler(func(evt interface {}) {
      switch evt.(type) {
        case *events.Message:
          cb(client, Ev{Message:evt.(*events.Message)})
        case *events.Receipt:
          cb(client, Ev{Receipt:evt.(*events.Receipt)})
        case *events.ChatPresence:
          cb(client, Ev{ChatPresence:evt.(*events.ChatPresence)})
        default:
          cb(client, Ev{More:evt})
      }
    })
    if client.Store.ID == nil {
        err = client.Connect()
        if err != nil { fmt.Println("GoError:",err); return }
        linkingCode, gagal := client.PairPhone(nomor, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
        if gagal != nil { fmt.Println("GoError:",gagal); return }
        fmt.Println(nomor,">",linkingCode)
    } else {
        err = client.Connect()
        if err != nil { fmt.Println("GoError:",err); return }
        fmt.Println(nomor,">","Connected")
        //client.SendPresence(types.PresenceAvailable)
        client.SendPresence(types.PresenceUnavailable)
    }
}

func mes(client *whatsmeow.Client, evt Ev) {
    if evt.Message != nil {
        v := evt.Message
        if v.Info.Chat.String() == "status@broadcast" {
            if v.Info.Type != "reaction" {
                reaction := client.BuildReaction(v.Info.Chat, v.Info.Sender, v.Info.ID, "🗿")
                extras := []whatsmeow.SendRequestExtra{}
                client.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender)
                client.SendMessage(context.Background(), v.Info.Chat, reaction, extras...)
            }
        }
    }
}
func main() {
    nomer := flag.String("n", "6283169480682", "n")
    flag.Parse()
    go Connect(*nomer, mes)
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    os.Exit(0)
}