package main /* import "s32x.com/gamedetect" */

import (
	"log"
	"os"

	"s32x.com/gamedetect/service"
)

func main() {
	// Generate a new service
	s, err := service.NewService(
		getenv("MODEL_PATH", "graph/output_graph.pb"),    // The trained output graph
		getenv("LABELS_PATH", "graph/output_labels.txt"), // The labels trained in the graph
	)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	s.Start(
		getenv("ENV", "dev"),              // The environment to run the server in
		getenv("DOMAIN", "gamedetect.io"), // The host the server is running on
		getenv("DEMO", "false"),           // Perform sanity tests and serve the web frontend
		getenv("PORT", "8080"),            // The port the server will run on
	)
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
