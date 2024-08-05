package main

import (
	"clean-architecture/configs"
	"clean-architecture/internal/infra/database"
	"clean-architecture/internal/infra/graph"
	"clean-architecture/internal/infra/grpc/pb"
	"clean-architecture/internal/infra/grpc/service"
	"clean-architecture/internal/infra/web"
	"clean-architecture/internal/infra/web/webserver"
	"clean-architecture/internal/usecase"
	"database/sql"
	"fmt"
	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	webServer := webserver.NewWebServer(":" + configs.WebServerPort)

	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)
	webOrderHandler := web.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase, orderRepository)

	webServer.AddHandler(http.MethodPost, "/orders", webOrderHandler.Create)
	webServer.AddHandler(http.MethodGet, "/orders", webOrderHandler.List)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go func() {
		err = webServer.Start()
		if err != nil {
			panic(err)
		}
	}()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}
