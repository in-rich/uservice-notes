package main

import (
	"fmt"
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
	db, closeDB := deploy.OpenDB(config.App.Postgres.DSN)
	defer closeDB()

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
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

	listener, server := deploy.StartGRPCServer(fmt.Sprintf(":%d", config.App.Server.Port), "notes")
	defer deploy.CloseGRPCServer(listener, server)

	notes_pb.RegisterGetNoteServer(server, getNoteHandler)
	notes_pb.RegisterListNotesServer(server, listNotesHandler)
	notes_pb.RegisterListNotesByAuthorServer(server, listNotesByAuthorHandler)
	notes_pb.RegisterUpsertNoteServer(server, upsertNoteHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
