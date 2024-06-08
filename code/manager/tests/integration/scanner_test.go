package integration

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/CSPF-Founder/shieldsup/common/manager/db/models"
// 	"github.com/CSPF-Founder/shieldsup/common/manager/enums/jobstatus"
// 	"github.com/CSPF-Founder/shieldsup/common/manager/internal/scanner"
// 	"github.com/CSPF-Founder/shieldsup/common/manager/logger"
// )

// func TestScanner_GetJobsToScan(t *testing.T) {
// 	// Create a new Docker pool
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		t.Fatalf("Could not construct pool: %s", err)
// 	}

// 	// Uses pool to try to connect to Docker
// 	err = pool.Client.Ping()
// 	if err != nil {
// 		t.Fatalf("Could not connect to Docker: %s", err)
// 	}

// 	// Start a MongoDB container
// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "mongo",
// 		Tag:        "latest",
// 	})
// 	if err != nil {
// 		t.Fatalf("Could not start resource: %s", err)
// 	}
// 	defer func() {
// 		if err := pool.Purge(resource); err != nil {
// 			t.Fatalf("Could not purge resource: %s", err)
// 		}
// 	}()

// 	// MongoDB connection string
// 	dbURI := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
// 	if err != nil {
// 		t.Fatalf("Could not connect to MongoDB: %s", err)
// 	}
// 	defer client.Disconnect(context.Background())

// 	testDB := client.Database("testDB")

// 	if testDB == nil {
// 		t.Fatalf("Could not connect to database")
// 	}

// 	lgr, err := logger.NewLogger(&logger.Config{
// 		Level:    "debug",
// 		Filename: "/tmp/manager.log",
// 	})
// 	if err != nil {
// 		t.Fatalf("Could not create logger: %s", err)
// 	}

// 	dbService := models.New(testDB)
// 	jobService := &dbService.Job
// 	scanner := scanner.NewScanner("scanner", jobService, lgr)

// 	ctx := context.Background()
// 	jobs := scanner.GetJobsToScan(ctx)

// 	if len(jobs) != 0 {
// 		t.Fatalf("Expected 0 jobs, got %d", len(jobs))
// 	}

// 	createdAt := time.Now()

// 	// Create a new job
// 	job := models.Job{
// 		ID:        1,
// 		Status:    int(jobstatus.Default),
// 		APIURL:    "http://localhost:8080",
// 		CreatedAt: createdAt,
// 		UserID:    1,
// 	}

// 	// Insert the job
// 	_, err = jobService.Create(ctx, job)

// 	if err != nil {
// 		t.Fatalf("Could not insert job: %s", err)
// 	}

// 	jobs = scanner.GetJobsToScan(ctx)
// 	if len(jobs) != 1 {
// 		t.Fatalf("Expected 1 job, got %d", len(jobs))
// 	}
// }
