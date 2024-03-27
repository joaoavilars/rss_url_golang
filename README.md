# Gerador de RSS para Serviços de Contingência

Este é um programa em Go que extrai informações de uma página HTML e gera um arquivo XML no formato RSS 2.0 contendo essas informações. O programa é útil para criar feeds de notícias ou atualizações a partir de dados estruturados da página https://www.sefaz.rs.gov.br/NFE/NFE-SVC.aspx.

o rss_nt.go é para obter uma lista das NTs que forem publicadas no site:
https://www.nfe.fazenda.gov.br/portal/listaConteudo.aspx?tipoConteudo=04BIflQt1aY=


## Requisitos

- Go 1.16 ou superior


## Instalação

1. Clone este repositório:

```bash
git clone https://github.com/joaoavilars/rss_url_golang.git
cd rss_url_golang
```
Instale a dependência:
```
go get github.com/PuerkitoBio/goquery
```

## Compilação:
```
GOARCH=amd64 GOOS=linux go build -o gera_rss -tags netgo -ldflags '-extldflags "-static"' main.go 
```

## Uso:
Para usar o programa, execute o binário compilado gera_rss passando o caminho para o arquivo XML de saída como argumento da linha de comando. Por exemplo:

```
./gera_rss /caminho/para/arquivo.xml
```

Substitua /caminho/para/arquivo.xml pelo caminho desejado para o arquivo XML que deseja criar.


