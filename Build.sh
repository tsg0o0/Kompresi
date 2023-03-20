rm -rf Build/
mkdir Build/

mkdir -p Build/mac/resources/
mkdir -p Build/windows/resources/
mkdir -p Build/linux/resources/

# macOS
cp Resources/config.json Build/mac/config.json
cp -r Resources/resources/mac/ Build/mac/resources/mac/
cd Main/
GOOS=darwin GOARCH=amd64 go build kompresi.go
cp kompresi ../Build/mac/kompresi
rm kompresi
cd ../Configure/
GOOS=darwin GOARCH=amd64 go build KompresiConfigure.go
cp KompresiConfigure ../Build/mac/KompresiConfigure
rm -r KompresiConfigure
cd ../

# windows
cp Resources/config.json Build/windows/config.json
cp -r Resources/resources/win/ Build/windows/resources/win/
cd Main/
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 go build kompresi.go
cp kompresi.exe ../Build/windows/kompresi.exe
rm kompresi.exe
cd ../Configure/
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 go build -ldflags "-H windowsgui" KompresiConfigure.go
cp KompresiConfigure.exe ../Build/windows/KompresiConfigure.exe
rm KompresiConfigure.exe
cd ../

# linux
cp Resources/config.json Build/linux/config.json
cp -r Resources/resources/linux/ Build/linux/resources/linux/
cd Main/
GOOS=linux GOARCH=amd64 go build kompresi.go
cp kompresi ../Build/linux/kompresi
rm kompresi
cd ../
