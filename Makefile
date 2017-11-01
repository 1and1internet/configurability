PLUGINS = apache2 basic mysql nginx php
MAIN_FOLDER = bin
PLUGINS_FOLDER = bin/plugins

all: main ${PLUGINS}

godeps:
	go get

main: godeps
	go build -o ${MAIN_FOLDER}/configurator main.go

plugins: ${PLUGINS_FOLDER} ${PLUGINS}

${PLUGINS}: godeps
	go build -buildmode=plugin -o ${PLUGINS_FOLDER}/$@.so plugins/$@/$@.go
