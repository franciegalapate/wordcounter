package org.example.wordcounter;

import javafx.geometry.Pos;
import javafx.scene.control.Button;
import javafx.scene.control.Label;
import javafx.scene.input.TransferMode;
import javafx.scene.layout.VBox;
import javafx.stage.FileChooser;
import javafx.stage.Stage;
import java.io.File;
import java.util.List;

public class MainUI extends VBox {

    private final EpubParser parser = new EpubParser();

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

        this.getChildren().addAll(instructionLabel, browseButton);

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
        try {
            EpubParser parser = new EpubParser();
            String cleanText = parser.extractText(file.getAbsolutePath());

            List<String> chunks = parser.splitText(cleanText, 3);

            System.out.println("--- PARSER REPORT ---");
            System.out.println("Total words extracted (approx): " + cleanText.split("\\s+").length);
            for (int i = 0; i < chunks.size(); i++) {
                System.out.println("Chunk " + (i+1) + " length: " + chunks.get(i).length());
            }

        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
        }
    }
}