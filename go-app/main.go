package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"go-app/gui"
)

func main() {
	// Initialize the Fyne app
	myApp := app.New()
	
	// Create the main window
	myWindow := myApp.NewWindow("Parallel Word Counter")

	// Set the UI content from our gui package
	myWindow.SetContent(gui.SetupUI(myWindow))

	// Set a default window size and launch the app
	myWindow.Resize(fyne.NewSize(400, 250))
	myWindow.ShowAndRun()
}