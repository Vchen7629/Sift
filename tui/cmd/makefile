test_all: integration unit

PKGS := . ../internal/api/... ../internal/service/... ../internal/types/... ../internal/ui/...

format: 
	go fmt ${PKGS} .

lint:
	golangci-lint run ${PKGS}