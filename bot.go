package main

import (
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

	b.Handle("/listAll", func(m *tb.Message) {
		users := db.getAllRecords()

		_, err := b.Send(m.Chat, fmt.Sprintf("%+v", users), tb.ModeMarkdown)

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
			oldKarma := db.get(name).Karma
			amount := getKarmaChanges(m.Text)
			println("Amount: ", amount)

			if amount == 0 {
				return
			}

			if name == strings.Replace(m.Sender.Username, "@", "", 1) {
				b.Send(m.Chat, "TA CHUPANDO TEU PROPRIO CU AE PORRA", tb.ModeMarkdown)

				db.updateKarma(name, -1)
				newKarma := db.get(name).Karma

				formatted := fmt.Sprintf("*%s* karma has decreased from _%d_ to _%d_", name, oldKarma, newKarma)
				_, err := b.Send(m.Chat, formatted, tb.ModeMarkdown)

				if err != nil {
					fmt.Printf("%s", err)
				}

				return
			}

			db.updateKarma(name, amount)
			newKarma := db.get(name).Karma

			var verb string
			if amount > 0 {
				verb = "increased"
			} else {
				verb = "decreased"
			}

			formatted := fmt.Sprintf("*%s* karma has %s from _%d_ to _%d_", name, verb, oldKarma, newKarma)
			_, err := b.Send(m.Chat, formatted, tb.ModeMarkdown)

			if err != nil {
				fmt.Printf("%s", err)
			}
		}
	})

	b.Start()
}

func getKarmaChanges(message string) int {
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
