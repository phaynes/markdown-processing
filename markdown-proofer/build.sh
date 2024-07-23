export GOOS=darwin
export GOARCH=arm64
cd src
go build -ldflags="-s -w" -o ../mdp ./cmd/mdp
cd ..
cp mdp ../bin/
