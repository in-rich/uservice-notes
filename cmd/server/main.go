package main

import (
	"fmt"
	"github.com/in-rich/lib-go/deploy"
	"github.com/in-rich/lib-go/monitor"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/config"
	"github.com/in-rich/uservice-notes/migrations"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/handlers"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/rs/zerolog"
	"os"
)

func getLogger() monitor.GRPCLogger {
	if deploy.IsReleaseEnv() {
		return monitor.NewGCPGRPCLogger(zerolog.New(os.Stdout), "uservice-notes")
	}

	return monitor.NewConsoleGRPCLogger()
}

func main() {
	logger := getLogger()

	logger.Info("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Fatal(err, "failed to connect to database")
	}
	defer closeDB()

	logger.Info("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		logger.Fatal(err, "failed to migrate")
	}

	depCheck := deploy.DepsCheck{
		Dependencies: func() map[string]error {
			return map[string]error{
				"Postgres": db.Ping(),
			}
		},
		Services: deploy.DepCheckServices{
			"GetNote":           {"Postgres"},
			"ListNotes":         {"Postgres"},
			"ListNotesByAuthor": {"Postgres"},
			"UpsertNote":        {"Postgres"},
		},
	}

	createNoteDAO := dao.NewCreateNoteRepository(db)
	deleteNoteDAO := dao.NewDeleteNoteRepository(db)
	getNoteDAO := dao.NewGetNoteRepository(db)
	listNotesDAO := dao.NewListNotesRepository(db)
	listNotesByAuthorDAO := dao.NewListNotesByAuthorRepository(db)
	updateNoteDAO := dao.NewUpdateNoteRepository(db)
	getNoteByIDDAO := dao.NewGetNoteByIDRepository(db)

	getNoteService := services.NewGetNoteService(getNoteDAO)
	listNotesService := services.NewListNotesService(listNotesDAO)
	listNotesByAuthorService := services.NewListNotesByAuthorService(listNotesByAuthorDAO)
	upsertNoteService := services.NewUpsertNoteService(
		updateNoteDAO,
		createNoteDAO,
		deleteNoteDAO,
	)
	getNoteByIDService := services.NewGetNoteByIDService(getNoteByIDDAO)

	getNoteHandler := handlers.NewGetNoteHandler(getNoteService, logger)
	listNotesHandler := handlers.NewListNotesHandler(listNotesService, logger)
	listNotesByAuthorHandler := handlers.NewListNotesByAuthorHandler(listNotesByAuthorService, logger)
	upsertNoteHandler := handlers.NewUpsertNoteHandler(upsertNoteService, logger)
	getNoteByIDHandler := handlers.NewGetNoteByIDHandler(getNoteByIDService)

	logger.Info(fmt.Sprintf("Starting to listen on port %v", config.App.Server.Port))
	listener, server, health := deploy.StartGRPCServer(logger, config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	notes_pb.RegisterGetNoteServer(server, getNoteHandler)
	notes_pb.RegisterListNotesServer(server, listNotesHandler)
	notes_pb.RegisterListNotesByAuthorServer(server, listNotesByAuthorHandler)
	notes_pb.RegisterUpsertNoteServer(server, upsertNoteHandler)
	notes_pb.RegisterGetNoteByIDServer(server, getNoteByIDHandler)

	logger.Info("Server started")
	if err := server.Serve(listener); err != nil {
		logger.Fatal(err, "failed to serve")
	}
}
