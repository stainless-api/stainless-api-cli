package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/stainless-api/stainless-api-cli/internal/mockstainless"
)

func main() {
	port := flag.Int("port", 4010, "port to listen on")
	flag.Parse()

	mock := mockstainless.NewMock(
		mockstainless.WithDefaultOrg(),
		mockstainless.WithDefaultProject(),
		mockstainless.WithDefaultCompareBuild(),
		mockstainless.WithDeviceAuth(1),
		mockstainless.WithGitRepos(),
	)
	defer mock.Cleanup()
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Mock server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mock.Server()); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
