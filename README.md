# go-stress-test
Desafio Técnico 2 do curso Pós Go Expert - Full Cycle

----

## stress test go

CLI para testes de carga em servicos web, desenvolvido em Go.

## parametros

| Flag            | Descricao                          | Obrigatorio |
|-----------------|------------------------------------|-------------|
| `--url`         | URL do servico a ser testado       | Sim         |
| `--requests`    | Numero total de requisicoes        | Sim         |
| `--concurrency` | Numero de chamadas simultaneas     | Sim         |

## executando com Docker

### build da imagem

```bash
docker build -t stress-test-go .
```

### executar teste

```bash
docker run stress-test-go --url=http://google.com --requests=1000 --concurrency=10
```

## executando localmente

```bash
go build -o stress-test .
./stress-test --url=http://google.com --requests=100 --concurrency=5
```

## testes

```bash
go test -v ./...
```

## exemplo de relatorio

```
========================================
       RELATORIO DE TESTE DE CARGA      
========================================
Tempo total:          1.234s
Total de requests:    1000
Requests com HTTP 200: 950
Erros de conexao:     10
----------------------------------------
Distribuicao de status HTTP:
  HTTP 200: 950
  HTTP 404: 20
  HTTP 500: 20
========================================
```
