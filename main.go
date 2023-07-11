package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type spritesheet struct {
	name    string
	content []byte
}

func main() {
	srcPath := flag.String("src", "", "Source path")
	spritesPath := flag.String("sprites", "", "Spritesheets path")
	componentsPath := flag.String("components", "", "Components path")

	flag.Parse()

	if *srcPath == "" {
		log.Fatal("Source path is required")
	}

	if *spritesPath == "" {
		log.Fatal("Spritesheets path is required")
	}

	if *componentsPath == "" {
		log.Fatal("Components path is required")
	}

	*srcPath = strings.TrimPrefix(*srcPath, "."+string(os.PathSeparator))

	spritesheets := make(map[string][]spritesheet)
	err := filepath.WalkDir(*srcPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.TrimPrefix(path, *srcPath+string(os.PathSeparator))
		name = filepath.Dir(name)
		name = strings.ReplaceAll(name, string(os.PathSeparator), "-")

		if _, ok := spritesheets[name]; !ok {
			spritesheets[name] = make([]spritesheet, 0)
		}
		spritesheets[name] = append(spritesheets[name], spritesheet{
			name:    path,
			content: content,
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*spritesPath, 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*componentsPath, 0755); err != nil {
		log.Fatal(err)
	}

	for name, sheets := range spritesheets {
		// Create SVG spritesheet
		content := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\">\n"
		ids := make([]string, 0)
		for _, sheet := range sheets {
			var svg struct {
				XMLName xml.Name `xml:"svg"`
				Width   string   `xml:"width,attr"`
				Height  string   `xml:"height,attr"`
				Content string   `xml:",innerxml"`
				ViewBox string   `xml:"viewBox,attr"`
			}
			xml.Unmarshal(sheet.content, &svg)

			id := strings.TrimSuffix(filepath.Base(sheet.name), filepath.Ext(sheet.name))
			ids = append(ids, id)
			content += fmt.Sprintf("<symbol id=\"%s\" viewBox=\"%s\" width=\"%s\" height=\"%s\">\n%s\n</symbol>\n", id, svg.ViewBox, svg.Width, svg.Height, svg.Content)
		}
		content += "</svg>"

		err := os.WriteFile(filepath.Join(*spritesPath, name+".svg"), []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Create React component
		publicPath := strings.TrimPrefix(*spritesPath, "public"+string(os.PathSeparator))
		content = "import * as React from \"react\";\n\nexport interface IconProps extends React.SVGProps<SVGSVGElement> {\n" +
			"\tname: \"" + strings.Join(ids, "\" | \"") + "\"\n" +
			"}\n\n" +
			"export function Icon(props: IconProps) {\n" +
			"\treturn <svg {...props}>\n" +
			"\t\t<use href={\"/" + publicPath + string(os.PathSeparator) + name + ".svg#\" + props.name} />\n" +
			"\t</svg>\n" +
			"}\n\n"

		err = os.WriteFile(filepath.Join(*componentsPath, name+".tsx"), []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
