package server

import (
	"database/sql"
	"log/slog"

	d "github.com/CollCaz/UniSite/database"
	"github.com/go-fuego/fuego"
)

type Server struct {
	server *fuego.Server
	logger *slog.Logger
	db     *d.DataService
}

type NewServerArgs struct {
	Logger *slog.Logger
	Db     *sql.DB
}

func (s *Server) Run() {
	s.logger.Info("Starting server...")
	s.logger.Info("Registering routes...")
	s.RegisterRoutes()
	s.logger.Info("Running server...")
	err := s.server.Run()
	s.logger.Error(err.Error())
}

func InitServer(args NewServerArgs) Server {
	if args.Logger == nil {
		args.Logger = slog.Default()
	}
	if args.Db == nil {
		args.Logger.Error("No db connection given")
		panic("must provide db connection")
	}

	server := fuego.NewServer(
		fuego.WithLogHandler(args.Logger.Handler()),
	)

	db := d.NewDataService(d.NewDataServiceArgs{
		Db:     args.Db,
		Logger: args.Logger,
	})

	return Server{
		logger: args.Logger,
		server: server,
		db:     db,
	}
}
