PLUGINS = apache2 basic mysql nginx php mongod java8 postgresql10 php_opcache
MAIN_FOLDER = bin
PLUGINS_FOLDER = bin/plugins

all: main ${PLUGINS}

godeps:
	go get

vendorupdate:
	dep ensure

main:
	go build -o ${MAIN_FOLDER}/configurator main.go

plugins: ${PLUGINS_FOLDER} ${PLUGINS}

${PLUGINS}:
	go build -buildmode=plugin -o ${PLUGINS_FOLDER}/$@.so plugins/$@/$@.go
