package main

import (
	"github.com/in-rich/lib-go/deploy"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/config"
	"github.com/in-rich/uservice-notes/migrations"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/handlers"
	"github.com/in-rich/uservice-notes/pkg/services"
	"log"
)

func main() {
	log.Println("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer closeDB()

	log.Println("Running migrations")
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	depCheck := func() map[string]bool {
		errDB := db.Ping()

		return map[string]bool{
			"GetNote":           errDB == nil,
			"ListNotes":         errDB == nil,
			"ListNotesByAuthor": errDB == nil,
			"UpsertNote":        errDB == nil,
			"":                  errDB == nil,
		}
	}

	createNoteDAO := dao.NewCreateNoteRepository(db)
	deleteNoteDAO := dao.NewDeleteNoteRepository(db)
	getNoteDAO := dao.NewGetNoteRepository(db)
	listNotesDAO := dao.NewListNotesRepository(db)
	listNotesByAuthorDAO := dao.NewListNotesByAuthorRepository(db)
	updateNoteDAO := dao.NewUpdateNoteRepository(db)

	getNoteService := services.NewGetNoteService(getNoteDAO)
	listNotesService := services.NewListNotesService(listNotesDAO)
	listNotesByAuthorService := services.NewListNotesByAuthorService(listNotesByAuthorDAO)
	upsertNoteService := services.NewUpsertNoteService(
		updateNoteDAO,
		createNoteDAO,
		deleteNoteDAO,
	)

	getNoteHandler := handlers.NewGetNoteHandler(getNoteService)
	listNotesHandler := handlers.NewListNotesHandler(listNotesService)
	listNotesByAuthorHandler := handlers.NewListNotesByAuthorHandler(listNotesByAuthorService)
	upsertNoteHandler := handlers.NewUpsertNoteHandler(upsertNoteService)

	log.Println("Starting to listen on port", config.App.Server.Port)
	listener, server, health := deploy.StartGRPCServer(config.App.Server.Port, depCheck)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	notes_pb.RegisterGetNoteServer(server, getNoteHandler)
	notes_pb.RegisterListNotesServer(server, listNotesHandler)
	notes_pb.RegisterListNotesByAuthorServer(server, listNotesByAuthorHandler)
	notes_pb.RegisterUpsertNoteServer(server, upsertNoteHandler)

	log.Println("Server started")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
