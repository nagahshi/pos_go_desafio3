# api-orders

__API de pedidos tipos de chamadas suportadas REST/GRAPHQL/GRPC__

### Recursos

* docker
* [migrate](https://github.com/golang-migrate/migrate)
* Rest Client no vscode


### Configuracoes

* Crie um arquivo .env com as configuracoes de ambiente usando 
modelo .env-example
* use make caso queria rodar a aplicação sem o build
* aplicação está escutando nas seguintes portas:
    - 8080 REST
    - 8081 GRPC
    - 8082 GRAPHQL

### Uso
* após  a instalação dos pacotes necessários
abra o terminal digite `make run`.
    - endpoints : [GET] /orders e [POST] /orders
    - gRPC: [OrderService] listOrder e [OrderService] CreateOrder
    - graphql: [query] listOrder e [mutation] createOrder

### Fontes

* [boilerplate](https://github.com/devfullcycle/goexpert/tree/main/20-CleanArch)