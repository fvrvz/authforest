// healthcheck is a minimal binary used by Docker's HEALTHCHECK instruction.
// It performs a single HTTP GET against the app's /health endpoint and
// exits 0 on HTTP 200, 1 on any other status or connection error.
// It is compiled as a fully static binary (CGO_ENABLED=0) so it can run
// inside the scratch runtime image with no libc or shell available.
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Fprintf(os.Stderr, "health check failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "health check returned %d\n", resp.StatusCode)
		os.Exit(1)
	}

	os.Exit(0)
}
