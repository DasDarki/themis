package main

import (
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"log"
)

func main() {
	if err := clipboard.Init(); err != nil {
		log.Println("Error initializing clipboard: ", err)
		return
	}

	storage := LoadStorage()

	app := tview.NewApplication()

	renderMainMenu(storage, app)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func renderMainMenu(storage Storage, app *tview.Application) {
	list := tview.NewList()

	list.AddItem("Aktenzeichen generieren", "Gibt ein neues Aktenzeichen aus", 'a', func() {
		renderGenerateCaseID(storage, app)
	})

	list.AddItem("Aktenzeichen hinzufügen", "Fügt ein neues Aktenzeichen-Format hinzu", 'h', func() {
		renderAddCaseID(storage, app)
	})

	list.AddItem("Aktenzeichen entfernen", "Löscht ein Aktenzeichen-Format", 'e', func() {
		renderRemoveCaseID(storage, app)
	})

	list.AddItem("Beenden", "Beendet Themis", 'q', func() {
		app.Stop()
	})

	app.SetRoot(list, true)
	app.SetFocus(list)
}

func renderGenerateCaseID(storage Storage, app *tview.Application) {
	list := tview.NewList()

	i := '1'
	for _, c := range storage {
		c := c
		list.AddItem(c.Name, "Generiert ein neues Aktenzeichen", i, func() {
			entry := c.Next()
			modal := tview.NewModal().
				SetText(entry).
				AddButtons([]string{"Kopieren", "OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Kopieren" {
						clipboard.Write(clipboard.FmtText, []byte(entry))
					}

					app.SetRoot(list, true)
					app.SetFocus(list)
				})

			app.SetRoot(modal, true)
		})
		i++
	}

	list.AddItem("Zurück", "Geht zurück zum Hauptmenü", 'z', func() {
		renderMainMenu(storage, app)
	})

	list.AddItem("Beenden", "Beendet Themis", 'q', func() {
		app.Stop()
	})

	app.SetRoot(list, true)
	app.SetFocus(list)
}

func renderAddCaseID(storage Storage, app *tview.Application) {
	form := tview.NewForm().
		AddInputField("Name", "", 0, nil, nil).
		AddInputField("Format", "", 0, nil, nil)

	form.AddButton("Hinzufügen", func() {
		name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
		format := form.GetFormItemByLabel("Format").(*tview.InputField).GetText()

		storage.Create(name, format)

		renderMainMenu(storage, app)
	})

	form.AddButton("Abbrechen", func() {
		renderMainMenu(storage, app)
	})

	app.SetRoot(form, true)
	app.SetFocus(form)
}

func renderRemoveCaseID(storage Storage, app *tview.Application) {
	list := tview.NewList()

	i := '1'
	for _, c := range storage {
		c := c
		list.AddItem(c.Name, "Löscht ein Aktenzeichen-Format", i, func() {
			confirm := tview.NewModal().
				SetText("Soll das Aktenzeichen-Format wirklich gelöscht werden?").
				AddButtons([]string{"Ja", "Nein"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Ja" {
						storage.Remove(c.Name)
					}

					app.SetRoot(list, true)
					app.SetFocus(list)
				})

			app.SetRoot(confirm, true)
			app.SetFocus(confirm)
		})
		i++
	}

	list.AddItem("Zurück", "Geht zurück zum Hauptmenü", 'z', func() {
		renderMainMenu(storage, app)
	})

	list.AddItem("Beenden", "Beendet Themis", 'q', func() {
		app.Stop()
	})

	app.SetRoot(list, true)
	app.SetFocus(list)
}
