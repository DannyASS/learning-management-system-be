package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/DannyAss/users/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// --- Load config ---
	cfg := config.InitConfig()

	// --- Ambil argumen CLI ---
	create := flag.Bool("create", false, "Create migration")
	dbType := flag.String("dbtype", "postgres", "Database type is required")
	dbName := flag.String("dbname", "Master", "Database name as registered in config (field name in DBConfig)")
	folderName := flag.String("foldername", "user", "folder name migrate")
	migrationName := flag.String("migrationname", "", "Migration name is required")
	direction := flag.String("dir", "up", "Migration direction: up or down")
	status := flag.Bool("status", false, "Check migration status")
	steps := flag.Int("steps", 0, "Number of migration steps to run")
	force := flag.Int("force", -1, "Force set migration version")
	flag.Parse()

	// --- Tentukan path migrations ---
	migrationsPath := fmt.Sprintf("internal/migrate/%s", strings.ToLower(*folderName))

	// --- Buat migration file ---
	if *create {
		if *migrationName == "" {
			log.Fatal("You must provide -migrationname argument")
		}

		createMigrationFile(migrationsPath, dbType, migrationName)
		return
	}

	// --- Build map DB otomatis dari DBConfig ---
	dbMap := buildDBMap(cfg.DBConnnect)

	// --- Ambil DB sesuai nama flag ---
	dbConn, ok := dbMap[*dbName]
	if !ok {
		log.Fatalf("Unknown db name: %s", *dbName)
	}

	// --- Buat connection string dari config ---
	var dbURL string
	//pakai ini kalau password ada special character
	encodedPass := url.QueryEscape(dbConn.DBPassword)
	switch *dbType {
	case "mysql":
		dbURL = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConn.DBUsername,
			encodedPass,
			dbConn.DBHostWrite,
			dbConn.DBPort,
			dbConn.DBName,
		)
	default:
		log.Fatalf("Unsupported db type: %s", *dbType)
	}

	// --- Init migrate ---
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to init migrate: %v", err)
	}

	// --- Run migration ---
	switch {
	case *status:
		checkMigrationStatus(m, migrationsPath)
	case *steps != 0:
		runMigrationSteps(m, *steps)
	case *force >= 0:
		forceMigrationVersion(m, *force)
	case *direction == "up":
		runMigrationUp(m)
	case *direction == "down":
		runMigrationDown(m)
	default:
		log.Fatalf("Unknown command")
	}
}

// --- Fungsi bantu: build map semua DB dari DBConfig ---
func buildDBMap(cfg config.DBManager) map[string]config.DBConnection {
	dbMap := make(map[string]config.DBConnection)
	v := reflect.ValueOf(cfg)
	t := reflect.TypeOf(cfg)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name // misal "Master", "B2B"
		fieldValue := v.Field(i).Interface().(config.DBConnection)
		dbMap[fieldName] = fieldValue
	}
	return dbMap
}

func createMigrationFile(folderPath string, dbType, migrationName *string) {
	if isMigrationNameExists(folderPath, *migrationName) {
		log.Fatalf("Migration name '%s' already exists in %s", *migrationName, folderPath)
	}

	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS

	err_ := os.MkdirAll(folderPath, 0755)
	if err_ != nil {
		panic(err_)
	}

	upFile := fmt.Sprintf("%s/%s_%s.up.sql", folderPath, timestamp, *migrationName)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", folderPath, timestamp, *migrationName)

	err := os.WriteFile(upFile, []byte("-- Migration up\n"), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(downFile, []byte("-- Migration down\n"), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created: %s, %s\n", upFile, downFile)
}

func isMigrationNameExists(folderPath, migrationName string) bool {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("Warning: Error reading directory %s: %v", folderPath, err)
		return false
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), migrationName) {
			return true
		}
	}

	return false
}

func checkMigrationStatus(m *migrate.Migrate, migrationsPath string) {
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get version: %v", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations applied yet")
	} else {
		log.Printf("Current version: %d, Dirty: %t", version, dirty)
	}

	// List available migrations
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Printf("Cannot read migrations directory: %v", err)
		return
	}

	log.Printf("Available migration files (%d):", len(files))
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".up.sql") || strings.HasSuffix(file.Name(), ".down.sql")) {
			log.Printf(" - %s", file.Name())
		}
	}
}

func runMigrationSteps(m *migrate.Migrate, steps int) {
	if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration steps failed: %v", err)
	}
	log.Printf("Migration steps completed: %d steps", steps)
}

func forceMigrationVersion(m *migrate.Migrate, version int) {
	if err := m.Force(version); err != nil {
		log.Fatalf("Force version failed: %v", err)
	}
	log.Printf("Forced migration version to: %d", version)
}

func runMigrationUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration up failed: %v", err)
	}
	log.Println("Migration up completed")
}

func runMigrationDown(m *migrate.Migrate) {
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration down failed: %v", err)
	}
	log.Println("Migration down completed")
}
