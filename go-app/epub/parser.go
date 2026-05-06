package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"strings"
	"golang.org/x/net/html"
	"path"
)

type Container struct {
	RootFiles []struct {
		FullPath string `xml:"full-path,attr"`
	} `xml:"rootfiles>rootfile"`
}

type Pack struct {
	Items []struct {
		HRef string `xml:"href,attr"`
		ID string `xml:"id,attr"`
	} `xml:"manifest>item"`

	ItemRefs []struct {
		IDRef string `xml:"idref,attr"`
	} `xml:"spine>itemref"`
}

func extractCleanText(reader *zip.ReadCloser, spine []string) []string {
	var chapters []string

	for _, chapterPath := range spine {
		for _, file := range reader.File {
			if file.Name == chapterPath {
				rc, _ := file.Open()

				tokenizer := html.NewTokenizer(rc)
				var sb strings.Builder
				skip := false

				for {
					tokenType := tokenizer.Next()

					if tokenType == html.ErrorToken {
						break // End of chapter
					}

					token := tokenizer.Token()

					switch tokenType {
					case html.StartTagToken, html.SelfClosingTagToken:
						if token.Data == "script" || token.Data == "style" {
							skip = true
						}
						sb.WriteString(" ")
					case html.EndTagToken:
						if token.Data == "script" || token.Data == "style" {
							skip = false
						}
						sb.WriteString(" ")

					case html.TextToken:
						if !skip {
							sb.WriteString(token.Data)
						}
					}
				}
				rc.Close()
				words := strings.Fields(sb.String())
				if len(words) > 0 {
					cleanText := strings.Join(words, " ")
					chapters = append(chapters, cleanText)
				}
			}
		}
	}
	return chapters
}

func parseOPF(reader *zip.ReadCloser, filename string) []string {
	for _, file := range reader.File {
		if file.Name == filename {
			rc, _ := file.Open()
			defer rc.Close()

			var pack Pack
			decoder := xml.NewDecoder(rc)
			_ = decoder.Decode(&pack)

			manifest := make(map[string]string)
			var spine []string
			for _, item := range pack.Items {
				manifest[item.ID] = item.HRef
			}

			baseDir := path.Dir(filename)
			for _, item := range pack.ItemRefs {
				href := manifest[item.IDRef]
				if baseDir != "." {
					href = baseDir + "/" + href
				}
				spine = append(spine, href)
			}
			return spine
		}
	}
	return nil
}

func getOPF(reader *zip.ReadCloser) string {
	for _, file := range reader.File {
		if file.Name == "META-INF/container.xml" {
			rc, _ := file.Open()
			defer rc.Close()

			var container Container
			decoder := xml.NewDecoder(rc)
			_ = decoder.Decode(&container)
			opfFile := container.RootFiles[0].FullPath
			return opfFile
		}
	}
	return ""
}

func GetChapters(filename string) []string {
	reader, _ := zip.OpenReader(filename)
	defer reader.Close()

	opf := getOPF(reader)
	spine := parseOPF(reader, opf)
	chapters := extractCleanText(reader, spine)

	return chapters
}
