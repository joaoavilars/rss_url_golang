package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: gera_rss <caminho_para_arquivo.xml>")
		return
	}

	outputPath := os.Args[1]
	fmt.Println("Caminho do arquivo de saída:", outputPath)

	// URL da página HTML
	url := "https://www.sefaz.rs.gov.br/NFE/NFE-SVC.aspx"
	fmt.Println("Obtendo conteúdo HTML da página:", url)

	// Obtenha o conteúdo HTML da página
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Carregue o conteúdo HTML da resposta HTTP
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Crie um novo objeto RSS
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title: "Feed de Serviços de Contingência",
			Link:  url,
			Desc:  "Contingência SVC-RS",
		},
	}

	// Selecione a tabela com base no ID do div
	doc.Find("#painelConteudo table tbody tr").Each(func(i int, s *goquery.Selection) {
		// Extraia os dados de cada linha da tabela
		uf := strings.TrimSpace(s.Find("td").Eq(0).Text())
		situacao := strings.TrimSpace(s.Find("td").Eq(1).Text())

		// Verifique se os campos extraídos não estão vazios
		if uf != "" && situacao != "" {
			// Adicione um novo item ao feed RSS
			item := Item{
				Title: uf,
				Link:  url,
				Desc:  fmt.Sprintf("Situação: %s", situacao),
			}
			rss.Channel.Items = append(rss.Channel.Items, item)
		}
	})

	// Codifique o objeto RSS como XML
	xmlBytes, err := xml.MarshalIndent(rss, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	// Escreva o XML no arquivo
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	fmt.Println("Escrevendo XML no arquivo:", outputPath)

	_, err = outputFile.Write(xmlBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Arquivo XML gerado com sucesso: %s\n", outputPath)
}
