SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=linux
DEL /Q bin\*
go build -o bin\rbq_anonymous_bot .
xz -z -e -9 -T 0 -v bin\rbq_anonymous_bot
