package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const dsn = "postgres://postgres:postgres@localhost:5432/rutadeorigen_poc?sslmode=disable"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("uso: go run ./cmd <seed N | list>")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("no se pudo conectar: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := ensureTable(ctx, db); err != nil {
		log.Fatalf("no se pudo crear la tabla: %v", err)
	}

	switch os.Args[1] {
	case "seed":
		n := 1
		if len(os.Args) > 2 {
			n, _ = strconv.Atoi(os.Args[2])
		}
		seed(ctx, db, n)
	case "list":
		list(ctx, db)
	default:
		fmt.Println("comando desconocido, usa 'seed' o 'list'")
	}
}

func ensureTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS registros_trazabilidad (
			id SERIAL PRIMARY KEY,
			descripcion TEXT NOT NULL,
			creado_en TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}

func seed(ctx context.Context, db *sql.DB, n int) {
	for i := 0; i < n; i++ {
		var id int
		var creadoEn time.Time
		err := db.QueryRowContext(ctx,
			`INSERT INTO registros_trazabilidad (descripcion) VALUES ($1) RETURNING id, creado_en`,
			fmt.Sprintf("registro de prueba #%d insertado en %s", i+1, time.Now().Format(time.RFC3339)),
		).Scan(&id, &creadoEn)
		if err != nil {
			log.Fatalf("no se pudo insertar: %v", err)
		}
		fmt.Printf("insertado id=%d creado_en=%s\n", id, creadoEn.Format(time.RFC3339))
	}
}

func list(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, `SELECT id, descripcion, creado_en FROM registros_trazabilidad ORDER BY id`)
	if err != nil {
		log.Fatalf("no se pudo consultar: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var descripcion string
		var creadoEn time.Time
		if err := rows.Scan(&id, &descripcion, &creadoEn); err != nil {
			log.Fatalf("error leyendo fila: %v", err)
		}
		fmt.Printf("id=%d creado_en=%s descripcion=%q\n", id, creadoEn.Format(time.RFC3339), descripcion)
		count++
	}
	fmt.Printf("\ntotal de registros: %d\n", count)
}
