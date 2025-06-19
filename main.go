package main

import (
	"fmt"
	"time"
	"os"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("settings.env")
	if err != nil {
		fmt.Println("Error")
	}


	DOCUMENT_ID := os.Getenv("DOCUMENT_ID")
	SERVICE_ACCOUNT_FILE:= os.Getenv("SERVICE_ACCOUNT_FILE")

	
	USERNAME := os.Getenv("USERNAME")
	COLOR := os.Getenv("COLOR")

	PASSWORD := os.Getenv("PASSWORD")

	app := tview.NewApplication()
	/*
	messages := []Message{
		{time.Now(), "Alice", "red", "Ciao!"},
		{time.Now(), "Bob", "green", "Salve!"},
	}
	*/
	messages, _ := ReadDocMessages(DOCUMENT_ID, SERVICE_ACCOUNT_FILE)

	var chatView *tview.TextView
	chatView = tview.NewTextView().
	SetDynamicColors(true).
	SetScrollable(true).
	SetChangedFunc(func() {
		app.Draw()
		chatView.ScrollToEnd()
	})

	updateChatView := func() {
		chatView.Clear()
		for _, msg := range messages {
			sender := msg.Sender
			if msg.Sender == USERNAME {
				sender = "You"
			}
			decMessage, _ := decrypt(msg.Message, PASSWORD)
			timeStr := msg.Timestamp.Format("15:04:05")
			fmt.Fprintf(chatView, "[gray]%s [gray]<@[white][%s]%s[white][gray]> [white]%s\n",
			timeStr, msg.Color, sender, decMessage)
		}
	}


	var input *tview.InputField

	input = tview.NewInputField().
	    SetLabel("message >> ").
	    SetFieldWidth(0).
	    SetFieldBackgroundColor(tcell.ColorDefault).
	    SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
		    encText, _:= encrypt(input.GetText(), PASSWORD)
		    input.SetText("") 
		    go func(txt string) {  
			_ = AppendMessageToDoc(DOCUMENT_ID, SERVICE_ACCOUNT_FILE, USERNAME, COLOR, txt)
		    }(encText)
		}
	    })




	updateChatView()


	header := tview.NewTextView()
	header.SetBackgroundColor(tcell.ColorPurple)
	header.SetText(" docTalk v0.0.1")

	style := tcell.StyleDefault.
	    Foreground(tcell.ColorBlack).
	    Background(tcell.NewHexColor(0x00ff11)).
	    Bold(true).
	    Blink(true)


	footer := tview.NewTextView()
	footer.SetBackgroundColor(tcell.ColorGreen)
	footer.SetTextStyle(style)
	footer.SetText(time.Now().Format("15:04:05") + " Writing as: " + USERNAME)
	
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			newMessages, err := ReadDocMessages(DOCUMENT_ID, SERVICE_ACCOUNT_FILE)
			if err != nil {
				continue
			}
			app.QueueUpdateDraw(func() {
				messages = newMessages
				updateChatView()
			})
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			app.QueueUpdateDraw(func() {
				footer.SetText(time.Now().Format("15:04:05") + " Writing as: " + USERNAME)
			})
		}
	}()


	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(chatView, 0, 1, false).
		AddItem(footer, 1, 0, false).
		AddItem(input, 1, 0, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	    if event.Key() == tcell.KeyTAB {
		focused := app.GetFocus()
		if focused == input {
		    app.SetFocus(chatView)
		} else {
		    app.SetFocus(input)
		}
		return nil
	    }
	    return event
	})


	app.SetRoot(layout, true)
	app.SetFocus(input) 

	if err := app.Run(); err != nil {
		panic(err)
	}
}
