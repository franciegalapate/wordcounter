package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
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
				rc, _ := file.Open() // Must open the file first
				defer rc.Close()

				tokenizer := html.NewTokenizer(rc)
				var sb strings.Builder
				skip := false

				for {
					tokenType := tokenizer.Next()

					if tokenType == html.ErrorToken {
						if tokenizer.Err() == io.EOF {
							break // End of chapter
						}
						fmt.Print(tokenizer.Err())
						break
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
				// Use strings.Fields to automatically strip all extra spaces and newlines
				// Then join them back together with a single space.
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
			rc, err := file.Open()
			if err != nil {
				fmt.Print(err)
				return nil
			}
			defer rc.Close()

			var pack Pack
			decoder := xml.NewDecoder(rc)
			err = decoder.Decode(&pack)
			if err != nil {
				fmt.Print(err)
				return nil
			}

			manifest := make(map[string]string)
			var spine []string
			for _, item := range pack.Items {
				manifest[item.ID] = item.HRef
			}

			for _, item := range pack.ItemRefs {
				spine = append(spine, manifest[item.IDRef])
			}
			return spine
		}
	}
	return nil
}

func getOPF(reader *zip.ReadCloser) string {
	for _, file := range reader.File {
		if file.Name == "META-INF/container.xml" {
			rc, err := file.Open()
			if err != nil {
				fmt.Print(err)
				return ""
			}
			defer rc.Close()

			var container Container
			decoder := xml.NewDecoder(rc)
			err = decoder.Decode(&container)
			if err != nil {
				fmt.Print(err)
				return ""
			}
			opfFile := container.RootFiles[0].FullPath
			return opfFile
		}
	}
	fmt.Println("META-INF/container.xml not found")
	return ""
}

func GetChapters() {
	filename := "../data/Around the World in 28 Languages.epub"
	reader, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer reader.Close()

	opf := getOPF(reader)
	spine := parseOPF(reader, opf)
	chapters := extractCleanText(reader, spine)
	fmt.Print(chapters[1])
}

