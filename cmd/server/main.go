package main

import (
	"database/sql"
	"fmt"
	notes_pb "github.com/in-rich/proto/proto-go/notes"
	"github.com/in-rich/uservice-notes/config"
	"github.com/in-rich/uservice-notes/migrations"
	"github.com/in-rich/uservice-notes/pkg/dao"
	"github.com/in-rich/uservice-notes/pkg/handlers"
	"github.com/in-rich/uservice-notes/pkg/services"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.App.Postgres.DSN)))
	db := bun.NewDB(sqldb, pgdialect.New())

	defer func() {
		_ = db.Close()
		_ = sqldb.Close()
	}()

	err := db.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	defer func() {
		server.GracefulStop()
		_ = listener.Close()
	}()

	notes_pb.RegisterGetNoteServer(server, getNoteHandler)
	notes_pb.RegisterListNotesServer(server, listNotesHandler)
	notes_pb.RegisterListNotesByAuthorServer(server, listNotesByAuthorHandler)
	notes_pb.RegisterUpsertNoteServer(server, upsertNoteHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
