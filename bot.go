package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const maxChangesPerMessage = 10
const incrementChar = "+"
const decrementChar = "-"

func initBot(db Connection) {
	b, err := tb.NewBot(tb.Settings{
		Token:  "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Sender, "FUCK YOU BITCH")
	})

	b.Handle("/list", func(m *tb.Message) {
		users := db.getAllRecords()

		var buffer bytes.Buffer

		for index := 0; index < len(users); index++ {
			user := users[index]
			formatted := fmt.Sprintf(" - *%s*: _%d_\n", user.ID, user.Karma)
			buffer.WriteString(formatted)
		}

		_, err := b.Send(m.Chat, buffer.String(), tb.ModeMarkdown)

		if err != nil {
			fmt.Printf("%s", err)
		}
	})

	b.Handle("/self", func(m *tb.Message) {
		user := db.get(m.Sender.Username)
		formatted := fmt.Sprintf("@%s: *%d*", user.ID, user.Karma)
		_, err := b.Send(m.Chat, formatted, tb.ModeMarkdown)

		if err != nil {
			fmt.Printf("%s", err)
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if len(m.Entities) != 1 {
			println("You have to do it one at a time")
			return
		}

		if len(m.Text) > 50 {
			m.Text = m.Text[0:50]
		}

		if m.Entities[0].Type == tb.EntityMention {
			name := getUserName(*m)
			amount := calcKarmaChanges(m.Text)

			if amount == 0 {
				return
			}

			if name == strings.Replace(m.Sender.Username, "@", "", 1) {
				printSelfMessage(m.Chat, b)
				setKarma(name, -1, m.Chat, db, b)
				return
			}

			setKarma(name, amount, m.Chat, db, b)
		}
	})

	b.Start()
}

func printSelfMessage(chat *tb.Chat, b *tb.Bot) {
	selfMessage := []byte{84, 65, 32, 67, 72, 85, 80, 65, 78, 68, 79, 32, 84, 69, 85, 32, 80, 82, 79, 80, 82, 73, 79, 32, 67, 85, 32, 65, 69, 32, 80, 79, 82, 82, 65}
	b.Send(chat, string(selfMessage), tb.ModeMarkdown)
}

func setKarma(name string, amount int, chat *tb.Chat, db Connection, b *tb.Bot) {
	var verb string

	if amount > 0 {
		verb = "increased"
	} else {
		verb = "decreased"
	}

	db.updateKarma(name, amount)
	newKarma := db.get(name).Karma

	formatted := fmt.Sprintf("*%s* karma has %s to _%d_ (%d)", verb, name, newKarma, amount)
	_, err := b.Send(chat, formatted, tb.ModeMarkdown)

	if err != nil {
		fmt.Printf("%s", err)
	}
}

func calcKarmaChanges(message string) int {
	change := 0
	msg := strings.Split(message, "")

	for _, letter := range msg {
		if letter == incrementChar {
			change++
		}

		if letter == decrementChar {
			change--
		}
	}

	if change <= maxChangesPerMessage && change >= -maxChangesPerMessage {
		return change
	}

	if change <= -maxChangesPerMessage {
		return -maxChangesPerMessage
	}

	return maxChangesPerMessage
}

func getUserName(m tb.Message) string {
	var name string

	for index := 0; index < len(m.Entities); index++ {
		entity := m.Entities[index]
		name = strings.Split(m.Text, " ")[entity.Offset]
	}

	// TODO: Fix this shit
	name = strings.Replace(name, decrementChar, "", 1000)
	name = strings.Replace(name, incrementChar, "", 1000)
	return strings.Replace(name, "@", "", 10)
}
