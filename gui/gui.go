package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	app    fyne.App
	win    fyne.Window
	query  *widget.Entry
	result *widget.Label
	status *widget.Label
}

func New() *GUI {
	a := app.New()
	return &GUI{
		app: a,
		win: a.NewWindow("Database Query UI"),
	}
}

func (g *GUI) Start() {
	g.win.Resize(fyne.NewSize(800, 600))
	g.win.SetFixedSize(true)
	g.query = widget.NewMultiLineEntry()
	g.query.SetPlaceHolder("Enter your SQL query here...")
	g.query.SetMinRowsVisible(5)

	g.status = widget.NewLabel("Ready")

	g.result = widget.NewLabel("Results will appear here...")
	g.result.Wrapping = fyne.TextWrapWord
	resultScroll := container.NewScroll(g.result)
	resultScroll.SetMinSize(fyne.NewSize(600, 300))

	executeBtn := widget.NewButton("Execute Query", func() {
		g.executeQuery()
	})

	clearBtn := widget.NewButton("Clear", func() {
		g.query.SetText("")
		g.result.SetText("Results will appear here...")
		g.status.SetText("Ready")
	})

	g.win.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabel("SQL Query:"),
			g.query,
			container.NewHBox(
				executeBtn,
				clearBtn,
				layout.NewSpacer(),
				g.status,
			),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewVBox(
			widget.NewLabel("Results:"),
			resultScroll,
		),
	))

	g.win.Resize(fyne.NewSize(800, 600))
	g.win.ShowAndRun()
}

func (g *GUI) executeQuery() {
	query := g.query.Text

	if query == "" {
		g.status.SetText("Error: Empty query")
		return
	}

	g.status.SetText("Executing...")

	result := g.mockDatabaseQuery(query)

	g.result.SetText(result)
	g.status.SetText("Query executed successfully")
}

func (g *GUI) mockDatabaseQuery(query string) string {
	return fmt.Sprintf("Query executed:\n%s\n\nMock Results:\n"+
		"Row 1: Data A, Data B, Data C\n"+
		"Row 2: Data D, Data E, Data F\n"+
		"Row 3: Data G, Data H, Data I\n"+
		"\n3 rows returned", query)
}

func (g *GUI) Quit() {
	g.app.Quit()
}
