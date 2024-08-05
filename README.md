# Order Management System

Este projeto(Desafio) da FullCycle sobre um sistema de gerenciamento de pedidos com interfaces REST, gRPC e GraphQL.

## Requisitos

- Docker
- Docker Compose
- Evans CLI (para testes gRPC)

## Instalação

1. Clone o repositório
2. No diretório do projeto, execute: `docker-compose up -d`

## Uso

### REST API

Endpoints disponíveis:

- Listar pedidos: `GET http://localhost:8000/orders`
- Criar pedido: `POST http://localhost:8000/orders`

#### O arquivo create_order.http pode ser usado para testes na pasta api.

### GraphQL

Acesse `http://localhost:8080/` e use as seguintes queries/mutations:

```graphql
mutation createOrder {
    createOrder(input: {id:"6756", Price:34.00, Tax: 2.00}){
        id
        Price
        Tax
        FinalPrice
    }
}

query queryOrders {
    orders {
        id
        Price
        Tax
        FinalPrice
    }
}
```

### gRPC

Para testar com Evans CLI:

1. Instale Evans: `go install github.com/ktr0731/evans@latest`

2. Execute:

   `evans -r repl`

   `package pb`

   `service OrderService`

   Listar orders:

   `call ListOrders`

   Criar orders:

   `call CreateOrder` 