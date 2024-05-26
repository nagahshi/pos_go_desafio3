package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"

	graphqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nagahshi/pos_go_desafio3/configs"
	"github.com/nagahshi/pos_go_desafio3/internal/event/handler"
	"github.com/nagahshi/pos_go_desafio3/internal/infra/graph"
	"github.com/nagahshi/pos_go_desafio3/internal/infra/grpc/pb"
	"github.com/nagahshi/pos_go_desafio3/internal/infra/grpc/service"
	"github.com/nagahshi/pos_go_desafio3/internal/infra/web/webserver"
	"github.com/nagahshi/pos_go_desafio3/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := getMysqlConn(configs.DBDriver, configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName, configs.PathOfMigrations)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs.RMQUser, configs.RMQPassword, configs.RMQHost, configs.RMQServerPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	webserver.AddHandler("POST", "/orders", webOrderHandler.Create)
	webserver.AddHandler("GET", "/orders", webOrderHandler.List)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphqlHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

// realiza conexao com o rmq
func getRabbitMQChannel(RMQUser, RMQPassword, RMQHost, RMQServerPort string) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", RMQUser, RMQPassword, RMQHost, RMQServerPort))
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}

// realiza a conexão com mysql
func getMysqlConn(DBDriver, DBUser, DBPassword, DBHost, DBPort, DBName, PathOfMigrations string) (*sql.DB, error) {
	db, err := sql.Open(DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/sys", DBUser, DBPassword, DBHost, DBPort))
	if err != nil {
		return nil, err
	}

	// cria banco de dados caso não exista
	row, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + DBName)
	if err != nil {
		return nil, err
	}

	_, err = row.RowsAffected()
	if err != nil {
		return nil, err
	}

	db.Close()

	// abre uma conexão com o banco de dados da aplicacao
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, DBHost, DBPort, DBName)
	db, err = sql.Open(DBDriver, DSN)
	if err != nil {
		return nil, err
	}

	// checa conexao
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// busca pasta das migracoes
	pathMigrationFiles, err := filepath.Abs(PathOfMigrations)
	if err != nil {
		return nil, err
	}

	// executa migracoes
	// garanta que golang-migrate esteja instalado corretamente
	// estou rodando dessa forma pois em meu ambiente não estava
	// funcionando o esquema de diretorios usando a lib golang-migrate
	_, _, err = Shellout("migrate -database \"mysql://" + DSN + "\" -path " + pathMigrationFiles + " up")
	if err != nil {
		return nil, err
	}

	// retorno uma conexão ativa
	return db, nil
}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
