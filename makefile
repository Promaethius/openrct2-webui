requirements-plugin:
	cd plugin && \
	npm install && \
	npm run-script gettypes

build-plugin:
	cd plugin && \
	npm run-script build

build-openrct2-gui:
	cd openrct2-gui && \
	cp -r ../plugin/lib . && \
	docker build -t openrct2-gui .

openrct2-gui: requirements-plugin build-plugin build-openrct2-gui

openrct2-webui:
	docker build .

all: openrct2-gui openrct2-webui

clean:
	rm -r openrct2-gui/lib
	cd plugin && npm run-script clean