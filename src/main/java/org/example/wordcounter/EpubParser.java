package org.example.wordcounter;

import nl.siegmann.epublib.domain.Book;
import nl.siegmann.epublib.domain.Resource;
import nl.siegmann.epublib.epub.EpubReader;
import org.jsoup.Jsoup;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class EpubParser {

    public String extractText(String filePath) throws IOException {
        StringBuilder sb = new StringBuilder();

        try (FileInputStream fis = new FileInputStream(filePath)) {
            Book book = new EpubReader().readEpub(fis);

            for (Resource resource : book.getContents()) {
                String mediaType = resource.getMediaType().toString();

                if (mediaType.contains("xhtml") || mediaType.contains("html")) {
                    String htmlContent = new String(resource.getData(), "UTF-8");

                    String cleanText = Jsoup.parse(htmlContent).text();
                    sb.append(cleanText).append(" ");
                }
            }
        }

        return sb.toString().toLowerCase().replaceAll("[^a-z0-9\\s]", " ");
    }

    public List<String> splitText(String text, int parts) {
        List<String> chunks = new ArrayList<>();
        int totalLength = text.length();

        if (totalLength == 0) return chunks;

        int targetChunkSize = totalLength / parts;
        int currentPos = 0;

        for (int i = 0; i < parts; i++) {
            if (i == parts - 1) {
                chunks.add(text.substring(currentPos).trim());
                break;
            }

            int endPos = currentPos + targetChunkSize;

            int nextSpace = text.indexOf(" ", endPos);

            if (nextSpace == -1 || nextSpace >= totalLength) {
                chunks.add(text.substring(currentPos).trim());
                break;
            }

            chunks.add(text.substring(currentPos, nextSpace).trim());
            currentPos = nextSpace + 1;
        }

        return chunks;
    }
}