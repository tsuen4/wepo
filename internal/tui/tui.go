package tui

import (
	"os"

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
		SetFieldStyle(tcell.StyleDefault.Background(tcell.ColorDefault))
	inputField.SetBackgroundColor(tcell.ColorDefault)

	errorModal = tview.NewModal()
	errorModal.SetBackgroundColor(tcell.ColorDefault)

	page = tview.NewPages().
		AddPage(inputPage, inputField, true, true).
		AddPage(errorPage, errorModal, true, false)
	page.SetBackgroundColor(tcell.ColorDefault)
}

func handleError(err error, pageName string) {
	errorModal.SetFocus(0).SetText(err.Error())
	page.SwitchToPage(errorPage)
}

// Run starts tui view.
func Run(iniPath, section string, args []string) error {
	// initialize
	client, err := wepo.New(iniPath, section)
	if err != nil {
		return err
	}

	input, err := wepo.Input(args, int(os.Stdin.Fd()))
	if err != nil {
		if err != wepo.ErrEmptyValue {
			return err
		}
	}

	inputField.SetText(input).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyEnter:
				causedErr := false
				contents, err := client.NewContents(inputField.GetText())
				if err != nil {
					causedErr = true
					// FIXME: add page
					handleError(err, errorPage)
				}

				if err := client.PostContents(contents); err != nil {
					causedErr = true
					// FIXME: add page
					handleError(err, errorPage)
				}

				if !causedErr {
					inputField.SetText("")
				}
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
