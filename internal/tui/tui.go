package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tsuen4/wepo/pkg/wepo"
)

// labels
var (
	quitLabel = "Quit"
	okLabel   = "OK"
)

// page names
var (
	inputPage = "input"
	errorPage = "error"
)

// tview containers
var (
	app        *tview.Application
	page       *tview.Pages
	inputField *tview.InputField
	errorModal *tview.Modal
)

var inputLabel = "Enter a text: "

func init() {
	// initialize tview
	app = tview.NewApplication()
	inputField = tview.NewInputField().SetLabel(inputLabel).
		SetFieldStyle(tcell.StyleDefault.Background(tview.Styles.PrimitiveBackgroundColor))
	errorModal = tview.NewModal()
	page = tview.NewPages().
		AddPage(inputPage, inputField, true, true).
		AddPage(errorPage, errorModal, true, false)
}

func handleError(err error, pageName string) {
	errorModal.SetFocus(0).SetText(err.Error())
	page.SwitchToPage(errorPage)
}

// Run starts tui view.
func Run(cfgDirPath string, args []string) error {
	// initialize
	client, err := wepo.New(cfgDirPath)
	if err != nil {
		return err
	}

	inputField.SetText(strings.Join(args, " ")).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEnter:
				contents, err := client.NewContents(inputField.GetText())
				if err != nil {
					handleError(err, errorPage)
				}

				if err := client.PostContents(contents); err != nil {
					handleError(err, errorPage)
				}

				inputField.SetText("")
				return nil
			}
			return event
		})

	errorModal.AddButtons([]string{quitLabel, okLabel}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == quitLabel {
				app.Stop()
			} else if buttonLabel == okLabel {
				errorModal.SetText("")
				page.SwitchToPage(inputPage)
			}
		})

	return app.SetRoot(page, true).EnableMouse(true).Run()
}
