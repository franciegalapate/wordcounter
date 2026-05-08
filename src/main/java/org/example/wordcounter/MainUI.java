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
import java.util.HashMap;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

public class MainUI extends VBox {

    private final TextArea reportArea = new TextArea();

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

        reportArea.setEditable(false);
        reportArea.setPrefHeight(300);
        reportArea.setPromptText("Analysis report will appear here...");
        reportArea.setStyle("-fx-font-family: 'Monospaced'; -fx-font-size: 12px;");

        this.getChildren().addAll(instructionLabel, browseButton, reportArea);

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
        reportArea.clear();

        try {
            EpubParser parser = new EpubParser();
            String cleanText = parser.extractText(file.getAbsolutePath());
            List<String> chunks = parser.splitText(cleanText, 3);

            int threadCount = chunks.size();

            long startTime = System.nanoTime();
            Map<String, Integer> sequentialResults = new HashMap<>();
            for (String chunk: chunks){
                WordCountTask task = new WordCountTask(chunk);
                Map<String, Integer> output = task.call();
                output.forEach((key, value) -> {
                    sequentialResults.merge(key, value, Integer::sum);
                });
            }
            long endTime = System.nanoTime();
            long sequentialDuration = (endTime - startTime);
            int seqUniqueWords = sequentialResults.size();
            int seqTotalWords = sequentialResults.values().stream().mapToInt(Integer::intValue).sum();


            ExecutorService executor = Executors.newFixedThreadPool(threadCount);

            List<Future<Map<String, Integer>>> futures = new ArrayList<>();
            for (String chunk : chunks) {
                WordCountTask task = new WordCountTask(chunk);
                Future<Map<String, Integer>> future = executor.submit(task);
                futures.add(future);
            }

            startTime = System.nanoTime();
            Map<String, Integer> results = new HashMap<>();
            for (Future<Map<String, Integer>> future: futures){
                Map<String, Integer> output = future.get();
                output.forEach((key, value) -> {
                    results.merge(key, value, Integer::sum);
                });
            }

            executor.shutdown();
            endTime = System.nanoTime();
            long parallelDuration = (endTime - startTime);
            int uniqueWords = results.size();
            int totalWords = results.values().stream().mapToInt(Integer::intValue).sum();

            Platform.runLater(() -> {
                appendReport("--- PARSER REPORT ---");
                appendReport("Total words extracted (approx): " + cleanText.split("\\s+").length);
                for (int i = 0; i < chunks.size(); i++) {
                    appendReport("Chunk " + (i + 1) + " length: " + chunks.get(i).length());
                }

                appendReport("\n--- DISPATCH REPORT ---");
                appendReport("Tasks submitted: " + futures.size());

                appendReport("\n--- PARALLEL REPORT ---");
                appendReport("Total Unique Words: " + results.size());
                appendReport("Total Number of Words: " + results.values().stream().mapToInt(Integer::intValue).sum());
                appendReport("Total Time Elapsed: " + parallelDuration + " ms");

                appendReport("\n--- SEQUENTIAL REPORT ---");
                appendReport("Total Unique Words: " + seqUniqueWords);
                appendReport("Total Number of Words: " + seqTotalWords);
                appendReport("Total Time Elapsed: " + sequentialDuration + " ms");
            });

            

        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }

    private void appendReport(String text) {
        reportArea.appendText(text + "\n");
    }
}