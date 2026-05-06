package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
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

func extractCleanText(reader *zip.ReadCloser, spine []string) ([]string, error) {
	var chapters []string
	
	for _, chapterPath := range spine {
		for _, file := range reader.File {
			if file.Name == chapterPath {
				rc, err := file.Open()
				if err != nil {
					return nil, err
				}

				tokenizer := html.NewTokenizer(rc)
				var sb strings.Builder
				skip := false

				for {
					tokenType := tokenizer.Next()

					if tokenType == html.ErrorToken {
						if tokenizer.Err() == io.EOF {
							break // End of chapter
						}
						rc.Close()
						return nil, tokenizer.Err()
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
	return chapters, nil
}


func parseOPF(reader *zip.ReadCloser, filename string) ([]string, error) {
	for _, file := range reader.File {
		if file.Name == filename {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			var pack Pack
			decoder := xml.NewDecoder(rc)
			err = decoder.Decode(&pack)
			if err != nil {
				return nil, err
			}

			manifest := make(map[string]string)
			var spine []string
			for _, item := range pack.Items {
				manifest[item.ID] = item.HRef
			}

			baseDir := path.Dir(filename)
			for _, item := range pack.ItemRefs {
				href := manifest[item.IDRef]
				if baseDir != "." && baseDir != "" {
					href = path.Join(baseDir, href)
				}
				spine = append(spine, href)
			}
			return spine, nil
		}
	}
	return nil, fmt.Errorf("OPF file %s not found", filename)
}

func getOPF(reader *zip.ReadCloser) (string, error) {
	for _, file := range reader.File {
		if file.Name == "META-INF/container.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			var container Container
			decoder := xml.NewDecoder(rc)
			err = decoder.Decode(&container)
			if err != nil {
				return "", err
			}
			opfFile := container.RootFiles[0].FullPath
			return opfFile, nil
		}
	}
	return "", fmt.Errorf("META-INF/container.xml not found")
}

func GetChapters(filename string) ([]string, error) {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	opf, err := getOPF(reader)
	if err != nil {
		return nil, err
	}

	spine, err := parseOPF(reader, opf)
	if err != nil {
		return nil, err
	}

	chapters, err := extractCleanText(reader, spine)
	if err != nil {
		return nil, err
	}

	return chapters, nil
}

