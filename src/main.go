package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/gob"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"os"
	"time"
)

func main() {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	<-time.After(3 * time.Second)

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := whatsapp.TextMessage {
			Info: whatsapp.MessageInfo{
				RemoteJid: r.FormValue("receiver")+"@s.whatsapp.net",
			},
			Text: r.FormValue("message"),
		}

		msgId,err := wac.Send(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		} else {
			fmt.Println("Message Sent -> ID : "+msgId)
		}
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func restoreSession() {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "restoring failed: %v\n", err)
			return
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error saving session: %v\n", err)
	}

	fmt.Printf("login successful, session: %v\n", session)
}