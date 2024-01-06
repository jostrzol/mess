_cover name *pkgs:
  go test -v -coverpkg=$(echo {{pkgs}} | tr ' ' ',') -coverprofile=./out/{{name}} {{pkgs}}

cover-engine: (_cover "profile-engine.cov" "./pkg/mess" "./pkg/board" "./pkg/color" "./test")
cover-all: (_cover "profile-all.cov" "./...")
cover: cover-engine cover-all

covtotal:
  @for profile in ./out/*.cov; do \
    echo -n "$(basename "$profile"):{{"\t"}}"; \
    go tool cover -func "$profile" \
      | grep "total:" \
      | tr -s '\t' \
      | cut -f 3; \
  done

build-cli:
  go build ./cmd/mess

build-server:
  go build ./cmd/mess-server

build: build-cli build-server

test:
  go test ./...

fuzz-chess:
  go test ./test -fuzz=FuzzChess

fuzz-dobutsu:
  go test ./test -fuzz=FuzzDobutsuShogi -short

fuzz-halma:
  go test ./test -fuzz=FuzzHalma -short
