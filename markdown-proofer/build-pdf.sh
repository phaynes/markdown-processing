export GOOS=darwin
export GOARCH=arm64
cd src
# go build -ldflags="-s -w" -o ../pdf2md ./cmd/pdf2md
go build -o ../pdf2md ./cmd/pdf2md
cd .
# cp pdf2md ../bin/
