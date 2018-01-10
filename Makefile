.PHONY: dist

vendor: server/glide.lock
	cd server && go get github.com/Masterminds/glide
	cd server && $(shell go env GOPATH)/bin/glide install

dist: vendor plugin.json
	@echo Building plugin

	# Clean old dist
	rm -rf dist
	# rm -rf webapp/dist
	rm -f server/plugin.exe

	# Build and copy files from webapp
	# cd webapp && npm run build
	# mkdir -p dist/sample/webapp
	mkdir -p dist/sample
	# cp webapp/dist/* dist/sample/webapp/

	# Build files from server
	cd server && go get github.com/mitchellh/gox
	$(shell go env GOPATH)/bin/gox -osarch='darwin/amd64 linux/amd64 windows/amd64' -output 'dist/intermediate/plugin_{{.OS}}_{{.Arch}}' ./server

	# Copy plugin files
	cp plugin.json dist/sample/

	# Copy server executables & compress plugin
	mkdir -p dist/sample/server
	mv dist/intermediate/plugin_darwin_amd64 dist/sample/server/plugin.exe
	cd dist && tar -zcvf mattermost-sample-plugin-darwin-amd64.tar.gz sample/*
	mv dist/intermediate/plugin_linux_amd64 dist/sample/server/plugin.exe
	cd dist && tar -zcvf mattermost-sample-plugin-linux-amd64.tar.gz sample/*
	mv dist/intermediate/plugin_windows_amd64.exe dist/sample/server/plugin.exe
	cd dist && tar -zcvf mattermost-sample-plugin-windows-amd64.tar.gz sample/*

	# Clean up temp files
	rm -rf dist/sample
	rm -rf dist/intermediate

	@echo MacOS X plugin built at: dist/mattermost-sample-plugin-darwin-amd64.tar.gz
	@echo Linux plugin built at: dist/mattermost-sample-plugin-linux-amd64.tar.gz
	@echo Windows plugin built at: dist/mattermost-sample-plugin-windows-amd64.tar.gz

clean:
	@echo Cleaning plugin
	rm -fr dist
	rm -f server/plugin.exe