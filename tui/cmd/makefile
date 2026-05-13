test_all: integration unit

PKGS := . ../internal/api/... ../internal/service/... ../internal/types/... ../internal/ui/...

format: 
	go fmt ${PKGS} .

lint:
	golangci-lint run ${PKGS}

unit:
	go test -tags unit ${PKGS}

generate_coverage:
	go test -tags=unit,integration ${PKGS} -v -p 1 -coverprofile=coverage.out -covermode=atomic

check_coverage:
	go tool cover -func=coverage.out