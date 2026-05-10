package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"time"

	_DeliveryHTTP "face-recognition-fyp/facerpca/delivery/http"
	_RepoJira "face-recognition-fyp/facerpca/repository/jira"
	_RepoMinio "face-recognition-fyp/facerpca/repository/minio"
	_RepoSQL "face-recognition-fyp/facerpca/repository/sql"
	_Usecase "face-recognition-fyp/facerpca/usecase"

	"github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initDB() *sql.DB {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.database"),
		viper.GetString("database.ssl_mode"),
	)

	dbConn, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("FATAL: Failed to init DB driver: %v", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("FATAL: Failed to CONNECT to DB: %v", err)
	}

	log.Info("Database Connection Success.")
	return dbConn
}

func initMinio() *minio.Client {
	endpoint := viper.GetString("minio.endpoint")
	accessKeyID := viper.GetString("minio.access_key")
	secretAccessKey := viper.GetString("minio.secret_key")
	useSSL := viper.GetBool("minio.use_ssl")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("FATAL: Failed to init Minio client: %v", err)
	}

	// Ping Minio to ensure it's reachable at startup
	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("FATAL: Failed to CONNECT to Minio (%s): %v", endpoint, err)
	}

	log.Infof("Minio Connection (%s) Success.", endpoint)
	return minioClient
}

func main() {
	configFile := flag.String("c", "config.yaml", "Config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Could not read config file, relying on env vars: %v", err)
	}

	log.SetLevel(log.InfoLevel)

	// ── 1. Infrastructure ──
	db := initDB()
	minioClient := initMinio()
	bucketName := viper.GetString("minio.bucket_name")

	// ── 2. Repositories ──
	repoSQL := _RepoSQL.NewFaceRPCASQLRepository(db)
	repoMinio := _RepoMinio.NewStorageRepo(minioClient, bucketName)
	repoJira := _RepoJira.NewService(
		viper.GetString("jira.url"),
		viper.GetString("jira.email"),
		viper.GetString("jira.token"),
		viper.GetString("jira.project_key"),
	)

	// ── 3. Usecase ──
	mlServiceURL := viper.GetString("ml_service.url")
	facerpcaUsecase := _Usecase.NewFaceRPCAUsecase(repoMinio, repoSQL, repoJira, mlServiceURL)

	// ── 4. Scheduler ──
	log.Info("Initializing scheduler...")
	scheduler := gocron.NewScheduler(time.Local)

	// 05:00 — Initialize daily status (Menunggu Kehadiran for all users)
	_, _ = scheduler.Every(1).Day().At("05:00").Do(func() {
		log.Info("Running Daily Init (Status 0: Menunggu Kehadiran)...")
		if err := facerpcaUsecase.HandleDailyInit(context.Background()); err != nil {
			log.Errorf("Daily init failed: %v", err)
		}
	})

	// 07:00 — Overnight Jira catch-up (poll yesterday's tasks for overtime workers)
	_, _ = scheduler.Every(1).Day().At("07:00").Do(func() {
		log.Info("Running Overnight Jira Catch-up (yesterday)...")
		yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
		if err := facerpcaUsecase.PollJiraTasksForDate(context.Background(), yesterday); err != nil {
			log.Errorf("Overnight Jira poll failed: %v", err)
		}
	})

	// 09:00 — Mark absent (status 0 → status 2)
	_, _ = scheduler.Every(1).Day().At("09:00").Do(func() {
		log.Info("Running Morning Cutoff (Mark Absent)...")
		if err := facerpcaUsecase.HandleMorningCutoff(context.Background()); err != nil {
			log.Errorf("Morning cutoff failed: %v", err)
		}
	})

	// Every 30 min — Jira polling for active users
	_, _ = scheduler.Every(30).Minutes().Do(func() {
		log.Info("Running Jira Task Polling...")
		if err := facerpcaUsecase.PollJiraTasks(context.Background()); err != nil {
			log.Errorf("Jira polling failed: %v", err)
		}
	})

	// 20:00 — Force clock out (status 1,3 → status 4) + KPI accumulation
	_, _ = scheduler.Every(1).Day().At("20:00").Do(func() {
		log.Info("Running Evening Cutoff (Force Clock Out)...")
		if err := facerpcaUsecase.HandleEveningCutoff(context.Background()); err != nil {
			log.Errorf("Evening cutoff failed: %v", err)
		}
	})

	scheduler.StartAsync()
	log.Info("Scheduler started: DailyInit (05:00), JiraCatchup (07:00), Absent (09:00), Jira (30m), ClockOut (20:00)")

	app := _DeliveryHTTP.SetupRouter(facerpcaUsecase)

	port := viper.GetString("server.port")
	if port == "" {
		port = "5000"
	}

	log.WithFields(log.Fields{
		"component": "http",
		"framework": "fiber",
		"port":      port,
	}).Info("HTTP VPS server is running")

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}