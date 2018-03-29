PLUGINS = apache2 basic mysql nginx php mongod java8 postgresql10
MAIN_FOLDER = bin
PLUGINS_FOLDER = bin/plugins

all: main ${PLUGINS}

godeps:
	go get

vendorupdate:
	govendor update +external

main:
	go build -o ${MAIN_FOLDER}/configurator main.go

plugins: ${PLUGINS_FOLDER} ${PLUGINS}

${PLUGINS}:
	go build -buildmode=plugin -o ${PLUGINS_FOLDER}/$@.so plugins/$@/$@.go
