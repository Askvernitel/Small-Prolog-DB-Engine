package gui

import (
	"fmt"
	"weird/db/engine/client"
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
	result *fyne.Container
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

	// Create the result container
	g.result = container.NewVBox()
	g.result.Add(widget.NewLabel("Results will appear here..."))

	// Wrap result in a scroll container
	/*	resultScroll := container.NewScroll(g.result)
		resultScroll.SetMinSize(fyne.NewSize(600, 300))*/

	executeBtn := widget.NewButton("Execute Query", func() {
		fmt.Println("EXECUTING")
		g.executeQuery()
	})

	clearBtn := widget.NewButton("Clear", func() {
		g.query.SetText("")
		g.result.Objects = []fyne.CanvasObject{widget.NewLabel("Results will appear here...")}
		g.result.Refresh()
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
			g.result,
			//resultScroll, // Use the scroll container here
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
		g.status.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	// Clear previous results
	g.result.Objects = []fyne.CanvasObject{}

	// Add new results
	g.result.Add(g.outputResponse(resps))
	g.result.Refresh()
	g.status.SetText("Query completed")
}

// outputResponse converts database responses into formatted table output
func (g *GUI) outputResponse(resps []*client.Response) fyne.CanvasObject {
	out := container.NewVBox()

	for idx, resp := range resps {
		// Add response header info
		header := widget.NewLabel(fmt.Sprintf("Result %d: %s - %s (Table: %s, Count: %d)",
			idx+1, resp.Status, resp.Message, resp.Table, resp.Count))
		header.TextStyle = fyne.TextStyle{Bold: true}
		out.Add(header)

		// Check if there are any rows
		if len(resp.Rows) == 0 || len(resp.Columns) == 0 {
			out.Add(widget.NewLabel("No data to display"))
			out.Add(widget.NewSeparator())
			continue
		}

		// Create a grid for the table
		gridContainer := container.NewVBox()

		// Add header row
		headerRow := container.NewHBox()
		for _, col := range resp.Columns {
			headerLabel := widget.NewLabel(col)
			headerLabel.TextStyle = fyne.TextStyle{Bold: true}
			headerLabel.Resize(fyne.NewSize(150, 35))
			headerRow.Add(headerLabel)
		}
		gridContainer.Add(headerRow)
		gridContainer.Add(widget.NewSeparator())

		// Add data rows
		for _, row := range resp.Rows {
			dataRow := container.NewHBox()
			for _, cellData := range row.Data {
				cellLabel := widget.NewLabel(cellData)
				cellLabel.Resize(fyne.NewSize(150, 35))
				dataRow.Add(cellLabel)
			}
			gridContainer.Add(dataRow)
		}

		out.Add(gridContainer)
		out.Add(widget.NewSeparator())
	}

	return out
}

func (g *GUI) Quit() {
	g.app.Quit()
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
