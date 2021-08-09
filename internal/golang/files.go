package golang

import (
	"fmt"

	"github.com/ezegrosfeld/cli/internal/util"
)

// createMain - creates the main.go file in the path
func createMain() error {
	path := "./" + Name + "/main.go"

	// Create the main.go file content
	main := fmt.Sprintf(`package main

import (
	"%s/cmd/server"
	"%s/pkg/logs"
)
	
func main() {
	logs.StartLogger()
		
	server := server.NewServer()

	server.StartServer()
}
	`, Repo, Repo)

	// Create the main.go file
	return util.CreateFile(path, main)
}

// createServiceFile - creates the service.go file
func createServerFile() error {
	cmdPath := "./" + Name + "/cmd"
	util.CreateFolder(cmdPath)

	serverPath := cmdPath + "/server"
	util.CreateFolder(serverPath)

	path := serverPath + "/server.go"

	// Create the service.go file content
	server := fmt.Sprintf(`package server
import (
	"log"
	"time"

	"github.com/ezegrosfeld/yoda"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

func (s *server) StartServer() {
	srv := yoda.NewServer(yoda.Config{
		Addr:         ":8080",
		Name:         "%s",
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	})

	// Routers
	r := srv.Group("/ping")
	r.Get("", func(c *yoda.Context) error {
		return c.JSON(200, "pong")
	})

	// Initialize server
	if err := srv.Start(); err != nil {
		log.Fatal(err.Error())
	}
}

`, Name)

	// Create the service.go file
	return util.CreateFile(path, server)
}

// createLogger - creates the logger.go file inside pkg/logs
func createLogger() error {
	pkgPath := "./" + Name + "/pkg"
	util.CreateFolder(pkgPath)

	pkgLogsPath := pkgPath + "/logs"
	util.CreateFolder(pkgLogsPath)

	path := pkgLogsPath + "/logger.go"

	// Create the logger.go file content
	logger := `package logs

import (
	"go.uber.org/zap"
)

var config zap.Config

func init() {
	config = zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		Development:       false,
		DisableStacktrace: true,
		Sampling:          zap.NewProductionConfig().Sampling,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
	}
}

func StartLogger() {
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	zap.ReplaceGlobals(logger)
}

`

	// Create the logger.go file
	return util.CreateFile(path, logger)
}
