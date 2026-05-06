package org.example.wordcounter;

import javafx.application.Application;
import javafx.scene.Scene;
import javafx.stage.Stage;

public class App extends Application {
    @Override
    public void start(Stage primaryStage) {
        MainUI root = new MainUI(primaryStage);
        Scene scene = new Scene(root, 600, 450);

        primaryStage.setTitle("EPUB Word Counter");
        primaryStage.setScene(scene);
        primaryStage.show();
    }
}