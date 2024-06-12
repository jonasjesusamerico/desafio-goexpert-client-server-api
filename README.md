# Projeto de Cotação do Dólar

Este projeto consiste em dois sistemas escritos em Go: um servidor (`server.go`) e um cliente (`client.go`). O servidor consulta a cotação do dólar em uma API externa, armazena a cotação em um banco de dados SQLite e retorna a cotação para o cliente. O cliente solicita a cotação do servidor e salva o resultado em um arquivo `cotacao.txt`.

## Pré-requisitos

- Go 1.16 ou superior
- Acesso à internet para consultar a API de cotação
- Permissões de escrita no diretório onde o projeto será executado

## Configuração do Banco de Dados

O servidor usa um banco de dados SQLite para armazenar as cotações. O arquivo do banco de dados será criado automaticamente se não existir.

## Executando o Projeto

### Passo 1: Clone o Repositório

Clone o repositório para o seu diretório local:

```sh
git clone https://github.com/jonasjesusamerico/desafio-goexpert-client-server-api.git
cd desafio-goexpert-client-server-api
```

### Passo 2: Inicie o Servidor

Execute o arquivo `server.go` para iniciar o servidor HTTP na porta 8080:

```sh
go run server.go
```

Você deve ver uma saída indicando que o servidor está rodando:

```sh
Listening on :8080
```

### Passo 3: Execute o Cliente

Em outra janela de terminal, execute o arquivo `client.go` para solicitar a cotação do servidor e salvar o resultado em um arquivo `cotacao.txt`:

```sh
go run client.go
```

### Passo 4: Verifique o Resultado

Após a execução do cliente, verifique o arquivo `cotacao.txt` no diretório do projeto. O conteúdo do arquivo deve ser algo como:

```sh
Dólar: 5.0000
```

## Estrutura do Projeto

- `server.go`: Implementa o servidor HTTP que consulta a cotação do dólar, armazena no banco de dados e retorna a cotação.
- `client.go`: Implementa o cliente que solicita a cotação do servidor e salva em um arquivo.
- `cambio.db`: Arquivo de banco de dados SQLite criado pelo servidor para armazenar as cotações.

## Tratamento de Erros

- O servidor usa contextos para definir timeouts na consulta à API e na inserção no banco de dados.
- O cliente usa um contexto para definir um timeout na solicitação ao servidor.
- Em caso de erro (timeout ou outros problemas), mensagens de erro detalhadas são registradas nos logs.

