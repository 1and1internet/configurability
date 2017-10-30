PLUGINS = apache2 basic mysql nginx php
MAIN_FOLDER = bin
PLUGINS_FOLDER = bin/plugins

all: godeps main ${PLUGINS}

godeps:
	go get

main:
	go build -o ${MAIN_FOLDER}/configurator main.go

plugins: ${PLUGINS_FOLDER} ${PLUGINS}

${PLUGINS}:
	go build -buildmode=plugin -o ${PLUGINS_FOLDER}/$@.so plugins/$@/$@.go
