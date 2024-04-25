package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Desc    string   `xml:"description"`
	Guid    string   `xml:"guid"`
}

type Channel struct {
	XMLName       xml.Name `xml:"channel"`
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Desc          string   `xml:"description"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Items         []Item   `xml:"item"`
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
			Title: "Contingência SVC-RS",
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
			// gerar guid
			guid := generateGUID()
			// Adicione um novo item ao feed RSS
			item := Item{
				Title: uf,
				Link:  url,
				Desc:  fmt.Sprintf("Situação: %s", situacao),
				Guid:  guid,
			}
			rss.Channel.Items = append(rss.Channel.Items, item)
		}
	})

	// Atualize a data de LastBuildDate para o momento atual
	rss.Channel.LastBuildDate = time.Now().Format(time.RFC1123)

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

func generateGUID() string {
	// Crie um slice de 16 bytes para armazenar o UUID
	uuid := make([]byte, 16)

	// Preencha o slice com bytes aleatórios
	_, err := rand.Read(uuid)
	if err != nil {
		// Em caso de erro, retorne uma string vazia
		return ""
	}

	// Defina os bits específicos de versão e variante do UUID
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Versão 4 (random)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variante RFC 4122

	// Codifique o UUID como uma string hexadecimal com hífens
	return hex.EncodeToString(uuid)
}
