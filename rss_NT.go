package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Desc    string   `xml:"description"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Desc    string   `xml:"description"`
	Items   []Item   `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

func getHref(n *html.Node) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return "https://www.nfe.fazenda.gov.br/portal/" + strings.TrimSpace(attr.Val), true
		}
	}
	return "", false
}

func extractItems(n *html.Node, rss *RSS) {
	if n.Type == html.ElementNode && n.Data == "div" && n.Attr != nil {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == "conteudoDinamico" {
				extractLinks(n, rss)
				return // Apenas processa esta div e retorna
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractItems(c, rss)
	}
}

// Extrai links especificamente dentro da div conteudoDinamico
func extractLinks(n *html.Node, rss *RSS) {
	// Verifica se o nó atual é um link dentro de um parágrafo (<p>)
	if n.Type == html.ElementNode && n.Data == "a" {
		link, _ := getHref(n) // Pega o link correto

		var title, desc string
		// Percorre os filhos do link (<a>) para encontrar o <span> com a classe "tituloConteudo"
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "span" {
				for _, a := range c.Attr {
					if a.Key == "class" && a.Val == "tituloConteudo" {
						title = strings.TrimSpace(c.FirstChild.Data)
						break
					}
				}
			}
		}

		// Percorre os nós irmãos do <a> para encontrar a descrição após o primeiro <br>
		for s := n.NextSibling; s != nil; s = s.NextSibling {
			if s.Type == html.ElementNode && s.Data == "br" {
				if s.NextSibling != nil && s.NextSibling.Type == html.TextNode {
					desc = strings.TrimSpace(s.NextSibling.Data)
					break
				}
			}
		}

		// Se conseguiu extrair o título e o link
		if title != "" && link != "" {
			item := Item{
				Title: title,
				Link:  link,
				Desc:  desc,
			}
			rss.Channel.Items = append(rss.Channel.Items, item)
		}
	}

	// Continua a percorrer a árvore DOM recursivamente
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, rss)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: ./rss_nt <arquivo_de_saida.xml>")
		return
	}
	outputPath := os.Args[1]

	url := "https://www.nfe.fazenda.gov.br/portal/listaConteudo.aspx?tipoConteudo=04BIflQt1aY="
	cmd := exec.Command("wget", "-qO-", url)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(strings.NewReader(string(output)))
	if err != nil {
		log.Fatal(err)
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title: "Últimas NTs",
			Link:  url,
			Desc:  "Últimas Notas Técnicas",
		},
	}

	extractItems(doc, &rss)

	xmlBytes, err := xml.MarshalIndent(rss, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(outputPath, xmlBytes, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Arquivo XML gerado com sucesso: %s\n", outputPath)
}
