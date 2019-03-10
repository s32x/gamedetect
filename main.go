package main /* import "s32x.com/gamedetect" */

import (
	"log"
	"math/rand"
	"os"
	"time"

	"s32x.com/gamedetect/service"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// Generate a new service
	s, err := service.NewService(
		getenv("MODEL_PATH", "graph/output_graph.pb"),    // The environment to run the server in
		getenv("LABELS_PATH", "graph/output_labels.txt"), // The host the server is running on)
	)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	s.Start(getenv("PORT", "8080")) // The port the server will run on
}

// getenv attempts to retrieve and return a variable from the environment. If it
// fails it will either crash or failover to a passed default value
func getenv(key string, def ...string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	if len(def) == 0 {
		log.Fatalf("%s not defined in environment", key)
	}
	return def[0]
}
