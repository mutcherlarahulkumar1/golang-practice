package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("ConnectionString"))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}
	fmt.Println("Connected to DB")
	return db, nil
}
