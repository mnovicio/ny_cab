package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	// service implemenation
	svc "mnovicio.com/nycab/server/service"

	// grpc server
	"mnovicio.com/nycab/server/grpc"

	// rest gateway
	"mnovicio.com/nycab/server/rest"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// GRPCPort is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "10001", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "10002", "HTTP port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "root", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "admin123", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "ny_cab_data", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	nyCabSvc := svc.GetServiceInstance(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	// run gRPC server
	return grpc.RunServer(ctx, nyCabSvc, cfg.GRPCPort)
}
