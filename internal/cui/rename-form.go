package cui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RenameForm(oldName string) string {
	app := tview.NewApplication()

	newName := oldName

	// Create a TextView to display the static instruction
	instruction := tview.NewTextView().
		SetText(fmt.Sprintf("Rename the context '%s' to:", oldName)).
		SetTextAlign(tview.AlignLeft)

	// Create an input field for renaming
	inputField := tview.NewInputField().
		//SetLabel("Rename to: ").
		//SetLabelColor(tcell.ColorWhite).
		SetText(oldName).
		SetFieldWidth(40).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetChangedFunc(func(text string) {
			newName = text
		})

	// Create a Flex layout to hold the instruction and the input field
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(instruction, 1, 1, false). // Add instruction text
		AddItem(inputField, 0, 1, true)    // Add input field

	// Capture key events for ESC and Enter
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// Rename (Enter key) confirms the new name and exits the app
			app.Stop()
		case tcell.KeyEscape:
			// Cancel (ESC key) sets newName back to oldName and exits the app
			newName = oldName
			app.Stop()
		}
		return event
	})

	// Run the application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return newName
}
