export GOOS=darwin
export GOARCH=arm64
go build -ldflags="-s -w" -o mdp ./cmd/mdp
cp mdp ../bin/
