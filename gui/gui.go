package gui

import (
	"bytes"
	"fmt"
	"weird/db/engine/executor"

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

	exec executor.DbExecutor
}

func New(exec executor.DbExecutor) *GUI {
	a := app.New()
	return &GUI{
		app:  a,
		win:  a.NewWindow("Database Query UI"),
		exec: exec,
	}
}

func (g *GUI) Start() {
	g.query = widget.NewMultiLineEntry()
	g.query.SetPlaceHolder("Enter your SQL query here...")
	g.query.SetMinRowsVisible(5)

	g.status = widget.NewLabel("Ready")

	g.result = widget.NewLabel("Results will appear here...")
	g.result.Wrapping = fyne.TextWrapWord
	resultScroll := container.NewScroll(g.result)
	resultScroll.SetMinSize(fyne.NewSize(600, 300))

	executeBtn := widget.NewButton("Execute Query", func() {
		fmt.Println("EXECUTING")
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
	g.win.SetFixedSize(true)
	g.win.ShowAndRun()
}

func (g *GUI) executeQuery() {
	q := g.query.Text
	resps, err := g.exec.ExecuteQuery(q)
	if err != nil {
		g.result.SetText(err.Error())
		return
	}

	var out bytes.Buffer
	for _, resp := range resps {
		out.WriteString(resp.Message)
	}

	g.result.SetText(out.String())
}

func (g *GUI) Quit() {
	g.app.Quit()
}
