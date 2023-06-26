SET NAME=rbq_anonymous_bot_v1.1.0_
DEL /Q bin\*
SET CGO_ENABLED=0
SET GOARCH=amd64
go generate
SET GOOS=windows
go build -o bin\%NAME%Windows64.exe .
DEL /Q *.syso
SET GOOS=linux
go build -o bin\%NAME%Linux64 .
SET GOOS=darwin
go build -o bin\%NAME%macOS64 .
SET GOARCH=386
go generate
SET GOOS=windows
go build -o bin\%NAME%Windows32.exe .
DEL /Q *.syso
SET GOOS=linux
go build -o bin\%NAME%Linux32 .
CD bin
MAKECAB /d compressiontype=lzx /d compressionmemory=21 %NAME%Windows32.exe
RENAME %NAME%Windows32.ex_ %NAME%Windows32.exe.cab
MAKECAB /d compressiontype=lzx /d compressionmemory=21 %NAME%Windows64.exe
RENAME %NAME%Windows64.ex_ %NAME%Windows64.exe.cab
DEL /Q *.exe
7z a -txz -mx9 %NAME%Linux64.xz %NAME%Linux64
DEL /Q %NAME%Linux64
7z a -txz -mx9 %NAME%Linux32.xz %NAME%Linux32
DEL /Q %NAME%Linux32
7z a -tzip -mx9 %NAME%macOS64.zip %NAME%macOS64
DEL /Q %NAME%macOS64
CD ..
SET NAME=
SET CGO_ENABLED=
SET GOARCH=
SET GOOS=
