package routes

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"
	"net/http"
	"net/rpc"
)

const (
	rpcPort  = "5001"
	gRPCPort = "50001"
)

type Config struct {
	Models data.Models
}

func (app *Config) Routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/log", app.WriteLog)

	return mux
}

func (app *Config) RPCListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

type LogServer struct {
	logs.UnimplementedLoggerServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.GetMessage(),
		Data: input.GetData(),
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		log.Println("Error inserting log entry: ", err)
		return &logs.LogResponse{
			Result: "failed",
		}, nil
	}

	return &logs.LogResponse{
		Result: "logged via GRPC",
	}, nil
}

func (app *Config) GRPCListen() error {
	log.Println("Starting GRPC server on port ", gRPCPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRPCPort))
	if err != nil {
		log.Println("Error starting GRPC server: ", err)
		return err
	}
	defer listen.Close()

	s := grpc.NewServer()

	logs.RegisterLoggerServer(s, &LogServer{Models: app.Models})

	log.Println("GRPC server started on port ", gRPCPort)

	if err = s.Serve(listen); err != nil {
		log.Println("Error serving GRPC server): ", err)
	}

	return nil
}
