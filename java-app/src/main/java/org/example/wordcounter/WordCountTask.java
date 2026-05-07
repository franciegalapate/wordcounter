package org.example.wordcounter;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Callable;

public class WordCountTask implements Callable<Map<String, Integer>> {

    private final String textChunk;

    public WordCountTask(String textChunk) {
        this.textChunk = textChunk;
    }

    @Override
    public Map<String, Integer> call() {
        Map<String, Integer> wordCount = new HashMap<>();

        if (textChunk == null || textChunk.isBlank()) {
            return wordCount;
        }
        String cleaned = textChunk.toLowerCase().replaceAll("[^a-z0-9\\s]", " ");

        String[] words = cleaned.trim().split("\\s+");

        for (String word : words) {
            if (!word.isEmpty()) {
                wordCount.merge(word, 1, Integer::sum);
            }
        }

        return wordCount;
    }
}