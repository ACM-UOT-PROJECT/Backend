package server

import (
	"database/sql"
	"log/slog"

	d "backend/database"
	"backend/judge"

	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
)

type Server struct {
	server *fuego.Server
	logger *slog.Logger
	db     *d.DataService
	Judge  *judge.Judge
}

type NewServerArgs struct {
	Logger *slog.Logger
	Db     *sql.DB
	Judge  *judge.Judge
}

func (s *Server) Run() {
	s.logger.Info("Starting server...")
	s.logger.Info("Registering routes...")
	s.RegisterRoutes()
	s.logger.Info("Running server...")
	s.server.Addr = "0.0.0.0:9999"
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
	if args.Judge == nil {
		args.Logger.Error("No judge given")
		panic("must give judge")
	}

	server := fuego.NewServer(
		fuego.WithGlobalMiddlewares(cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}).Handler),
	)

	db := d.NewDataService(d.NewDataServiceArgs{
		Db:     args.Db,
		Logger: args.Logger,
	})

	return Server{
		logger: args.Logger,
		server: server,
		db:     db,
		Judge:  args.Judge,
	}
}
