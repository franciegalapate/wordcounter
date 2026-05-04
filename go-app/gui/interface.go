package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SetupUI builds and returns the main layout for the application
func SetupUI(win fyne.Window) fyne.CanvasObject {
	// Display file selection status
	statusLabel := widget.NewLabel("Status: No file selected")
	
	// Show word count results in GUI
	resultLabel := widget.NewLabel("Word Count: N/A")

	// Add progress indicator during processing
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	progressBar.Hide() // Hidden until processing starts

	// Pre-declare filePicker so it can be called inside the button logic
	var filePicker dialog.Dialog

	// Create main window with file picker button
	openBtn := widget.NewButton("Select File to Count", func() {
		filePicker.Show()
	})

	// Initialize the file picker dialog
	filePicker = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		// Handle user interactions (errors and cancellations)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if reader == nil {
			statusLabel.SetText("Status: Selection canceled")
			return
		}
		
		defer reader.Close()

		filename := reader.URI().Name()
		statusLabel.SetText(fmt.Sprintf("Status: Selected '%s'", filename))

		// Launch processing in a separate goroutine to keep the GUI responsive
		go processFileSimulation(filename, progressBar, resultLabel, openBtn)

	}, win)

	// Combine all elements into a Vertical Box layout
	content := container.NewVBox(
		widget.NewLabelWithStyle("Parallel Programming Demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		statusLabel,
		openBtn,
		progressBar,
		resultLabel,
	)

	return content
}

// processFileSimulation mimics the parallel backend processing
func processFileSimulation(filename string, progress *widget.ProgressBar, result *widget.Label, btn *widget.Button) {
	// Prepare UI for processing
	btn.Disable()       // Prevent opening multiple files at once
	progress.Show()     // Reveal progress bar
	progress.SetValue(0)
	result.SetText("Processing in parallel...")

	// Simulate parallel processing time (e.g., chunking the file and counting)
	for i := 0.0; i <= 1.0; i += 0.1 {
		time.Sleep(150 * time.Millisecond) // Mock processing delay
		progress.SetValue(i)               // Update progress bar safely
	}

	// Mock result - this is where to call the actual parallel word counter
	mockWordCount := 142058 

	// Reset UI after processing
	progress.Hide()
	result.SetText(fmt.Sprintf("Word Count: %d words found in %s", mockWordCount, filename))
	btn.Enable()
}