package main

import (
	"context"
	"log"

	"github.com/connor-davis/zingfibre-core/env"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func main() {
	context := context.Background()

	connection, err := pgx.Connect(context, string(env.POSTGRES_DSN))

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer connection.Close(context)

	log.Println("Connected to the database successfully")

	queries := postgres.New(connection)

	userId, err := uuid.Parse("f47ac10b-58cc-4372-a567-0e02b2c3d479")

	if err != nil {
		log.Fatalf("Failed to parse UUID: %v\n", err)
	}

	user, err := queries.GetUser(context, userId)

	if err != nil {
		log.Fatalf("Failed to get user: %v\n", err)
	}

	log.Printf("User retrieved: %+v\n", user)
}
