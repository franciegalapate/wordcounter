package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"go-app/counter"
	"go-app/epub"
)

func SetupUI(win fyne.Window) fyne.CanvasObject {
	statusLabel := widget.NewLabel("Status: No file selected")
	resultLabel := widget.NewLabel("Word Count: N/A")

	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	progressBar.Hide()

	var filePicker dialog.Dialog

	openBtn := widget.NewButton("Select File to Count", func() {
		filePicker.Show()
	})

	filePicker = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if reader == nil {
			statusLabel.SetText("Status: Selection canceled")
			return
		}
		
		filePath := reader.URI().Path()
		filename := reader.URI().Name()
		
		// Close reader since epub.GetChapters will open it again by path
		reader.Close()

		statusLabel.SetText(fmt.Sprintf("Status: Selected '%s'", filename))

		go processFile(filePath, filename, progressBar, resultLabel, openBtn)

	}, win)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Parallel Programming Demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		statusLabel,
		openBtn,
		progressBar,
		resultLabel,
	)

	return content
}

func processFile(filePath string, filename string, progress *widget.ProgressBar, result *widget.Label, btn *widget.Button) {
	btn.Disable()
	progress.Show()
	
	progress.SetValue(0.1)
	result.SetText("Extracting chapters...")

	chapters := epub.GetChapters(filePath)
	
	progress.SetValue(0.5)
	result.SetText("Counting words in parallel...")

	// 4 workers for parallel counting
	wordCountResult := counter.ParallelWordCount(chapters, 4)

	progress.SetValue(1.0)
	progress.Hide()
	
	resultText := fmt.Sprintf("Word Count: %d total words, %d unique words found in %s\nProcessing Time: %v", 
		wordCountResult.TotalWords, wordCountResult.UniqueWords, filename, wordCountResult.ProcessingTime)
	result.SetText(resultText)
	
	btn.Enable()
}