browser := "chrome"

out               := "./out"
cover_file_engine := out / "profile-engine.cov"
cover_file_all    := out / "profile-all.cov"

cover: cover-engine cover-all
cover-engine: (_cover cover_file_engine "./pkg/mess" "./pkg/board" "./pkg/color" "./test")
cover-all: (_cover cover_file_all "./...")
_cover file *pkgs:
  go test -v -coverpkg=$(echo {{pkgs}} | tr ' ' ',') -coverprofile={{file}} {{pkgs}}

covtotal:
  @for profile in ./out/*.cov; do \
    echo -n "$(basename "$profile"):{{"\t"}}"; \
    go tool cover -func "$profile" \
      | grep "total:" \
      | tr -s '\t' \
      | cut -f 3; \
  done

covhtml-engine: (_covhtml cover_file_engine)
covhtml-all: (_covhtml cover_file_all)
_covhtml src:
  go tool cover -html {{src}} -o {{without_extension(src) + ".html"}}
  {{browser}} {{without_extension(src) + ".html"}}

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

lint:
  golangci-lint run -E revive,godot
