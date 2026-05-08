package org.example.wordcounter;

import javafx.application.Platform;
import javafx.geometry.Pos;
import javafx.scene.control.Button;
import javafx.scene.control.Label;
import javafx.scene.control.TextArea;
import javafx.scene.input.TransferMode;
import javafx.scene.layout.VBox;
import javafx.stage.FileChooser;
import javafx.stage.Stage;
import java.io.File;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;

public class MainUI extends VBox {

    private final EpubParser parser = new EpubParser();
    private final TextArea resultArea = new TextArea();
    private final Label statusLabel = new Label("Ready");

    public MainUI(Stage stage) {
        this.setAlignment(Pos.CENTER);
        this.setSpacing(20);
        this.setStyle("-fx-background-color: #f9f9fb; " +
                "-fx-border-color: #4f46e5; " +
                "-fx-border-style: dashed; " +
                "-fx-border-width: 2; " +
                "-fx-border-radius: 15; " +
                "-fx-padding: 50;");

        Label instructionLabel = new Label("Drag your EPUB here to count words");
        instructionLabel.setStyle("-fx-font-size: 18px; -fx-font-weight: bold;");

        Button browseButton = new Button("Browse Files");
        browseButton.setStyle("-fx-background-color: #4f46e5; -fx-text-fill: white;");

        resultArea.setEditable(false);
        resultArea.setPrefHeight(200);
        resultArea.setWrapText(true);
        resultArea.setStyle("-fx-font-family: 'Monospaced';");
        statusLabel.setStyle("-fx-text-fill: #666666;");

        this.getChildren().addAll(instructionLabel, browseButton, statusLabel, resultArea);

        this.setOnDragOver(event -> {
            if (event.getDragboard().hasFiles()) {
                event.acceptTransferModes(TransferMode.COPY);
            }
            event.consume();
        });

        this.setOnDragDropped(event -> {
            var db = event.getDragboard();
            if (db.hasFiles()) {
                processFile(db.getFiles().get(0));
            }
            event.setDropCompleted(true);
            event.consume();
        });

        browseButton.setOnAction(e -> {
            FileChooser fileChooser = new FileChooser();
            fileChooser.getExtensionFilters().add(new FileChooser.ExtensionFilter("EPUB Files", "*.epub"));
            File file = fileChooser.showOpenDialog(stage);
            if (file != null) processFile(file);
        });
    }

    private void processFile(File file) {
        new Thread(() -> {
            try {
                long startTime = System.currentTimeMillis();

                Platform.runLater(() -> {
                    statusLabel.setText("Processing: " + file.getName() + "...");
                    resultArea.clear();
                });

                String cleanText = parser.extractText(file.getAbsolutePath());
                List<String> chunks = parser.splitText(cleanText, 3);

                long endTime = System.currentTimeMillis();
                double durationSeconds = (endTime - startTime) / 1000.0;
                String finishedAt = ZonedDateTime.now().format(DateTimeFormatter.ISO_OFFSET_DATE_TIME);

                StringBuilder report = new StringBuilder();
                report.append("--- PARSER REPORT ---\n");
                report.append("Total words extracted (approx): ").append(cleanText.split("\\s+").length).append("\n");
                for (int i = 0; i < chunks.size(); i++) {
                    report.append("Chunk ").append(i+1).append(" length: ").append(chunks.get(i).length()).append("\n");
                }

                report.append("\nTotal time: ").append(String.format("%.3f", durationSeconds)).append(" s\n");
                report.append("Finished at: ").append(finishedAt).append("\n");

                Platform.runLater(() -> {
                    resultArea.setText(report.toString());
                    statusLabel.setText("Success!");
                });

            } catch (Exception e) {
                Platform.runLater(() -> {
                    statusLabel.setText("Error occurred.");
                    resultArea.setText("Error: " + e.getMessage());
                });
            }
        }).start();
    }
}