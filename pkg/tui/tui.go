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

func handleError(err error, buttonIndex int, pageName string) {
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

	// initialize tview
	app = tview.NewApplication()
	inputField = tview.NewInputField().SetLabel("Enter a text: ").SetText(strings.Join(args, " "))
	errorModal = tview.NewModal()
	page = tview.NewPages().
		AddPage(inputPage, inputField, true, true).
		AddPage(errorPage, errorModal, true, false)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			contents, err := client.NewContents(inputField.GetText())
			if err != nil {
				handleError(err, 0, errorPage)
			}

			if err := client.PostContents(contents); err != nil {
				handleError(err, 0, errorPage)
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
