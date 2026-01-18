## Aplicação de exemplo para SQS com GO

Essa é uma aplicação extremamente simples usada apenas para exemplificar o consumo de mensagens usando SQS e a serialização e deserialização de dados em GO


:warning: Não leve essa arquitetura como a melhor para produção sem entender o contexto de onde trabalha! Isso foi feito por alguém que tem pouco contato com Golang e nunca trabalhou em um time com essa linguagem para fins que não fossem aplicações de gerenciamento de infraestrutura e templates de projetos.

**Pontos relevantes:**

- No pacote `messaging` são definidos templates para criação de consumers e producers do SQS além da configuração;
- No pacote `consumers` são definidos os consumers em si (Talvez isso tenha ficado meio confuso?);
- No pacote `router` a classe que dá inicio ao GIN Server (Servidor http para receber requisições) é definida;
- O formato para consumo de mensagens descritor aqui pode não ser o ideal para produção, afinal estamos usando apenas uma thread para processar todas as mensagens.


### Configurando

Adicione as variáveis abaixo:

- SQS_QUEUE_URL: Endereço da fila do SQS para onde as mensagens serão enviadas e recebidas
- AWS_ACCESS_KEY_ID: ID da chave de acesso da AWS
- AWS_SECRET_ACCESS_KEY: Chave de acesso da AWS
- AWS_REGION: Região da AWS

### Usando

1. Inicie a aplicação

```sh
go run main.go
```

2. Envie uma requisição http

```sh
curl --location 'http://localhost:8080/deposito' \
--header 'Content-Type: application/json' \
--data '{
    "contaOrigem": "100000",
     "contaDestino": "400001",
     "valor": 1000
}'
```