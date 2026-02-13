# Como executar o projeto

## Pre-requisitos

- [Go](https://go.dev/dl/) versao 1.24 ou superior instalado

## Estrutura do projeto

```
├── server/server.go   # Servidor HTTP (porta 8080)
├── client/client.go   # Cliente HTTP
├── go.mod             # Modulo Go
└── go.sum             # Checksums das dependencias
```

## Instalando dependencias

Na raiz do projeto, execute:

```bash
go mod tidy
```

## Executando

### 1. Iniciar o servidor

Abra um terminal e execute:

```bash
go run ./server/
```

O servidor ira iniciar na porta 8080 e ficara aguardando requisicoes.
Voce vera a mensagem: `Servidor iniciado na porta 8080`

### 2. Executar o cliente

Em **outro terminal**, execute:

```bash
go run ./client/
```

O cliente ira:
1. Fazer uma requisicao para `http://localhost:8080/cotacao`
2. Receber o valor do dolar (campo `bid`)
3. Salvar o resultado no arquivo `cotacao.txt` no formato `Dolar: {valor}`
4. Exibir o valor no terminal

## Testando

### Teste basico

1. Inicie o servidor em um terminal
2. Execute o cliente em outro terminal
3. Verifique que:
   - O terminal do cliente mostra `Cotacao salva com sucesso! Dolar: X.XXXX`
   - O arquivo `cotacao.txt` foi criado na raiz do projeto com o conteudo `Dolar: X.XXXX`
   - O arquivo `cotacoes.db` (SQLite) foi criado na raiz do projeto

### Teste de multiplas execucoes

Execute o cliente varias vezes. Cada execucao:
- Atualiza o arquivo `cotacao.txt` com a cotacao mais recente
- Insere um novo registro no banco SQLite `cotacoes.db`

### Teste de timeout (comportamento esperado)

Os seguintes timeouts estao configurados:
- **Server -> API externa**: 200ms (se a API demorar mais, o log mostrara erro de timeout)
- **Server -> Banco SQLite**: 10ms (se a gravacao demorar mais, o log mostrara erro de timeout)
- **Client -> Server**: 300ms (se o servidor demorar mais, o log mostrara erro de timeout)

Se algum timeout ocorrer, uma mensagem de erro sera exibida nos logs do respectivo servico.

## Arquivos gerados em tempo de execucao

| Arquivo | Descricao |
|---|---|
| `cotacao.txt` | Cotacao atual do dolar, gerada pelo cliente |
| `cotacoes.db` | Banco SQLite com historico de cotacoes, gerado pelo servidor |
