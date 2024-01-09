
USER_GH=eyedeekay
VERSION=0.33.7
packagename=i2pkeys

echo:
	@echo "type make version to do release $(VERSION)"

version:
	github-release release -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "version $(VERSION)"

del:
	github-release delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION)

tar:
	tar --exclude .git \
		--exclude .go \
		--exclude bin \
		-cJvf ../$(packagename)_$(VERSION).orig.tar.xz .

copier:
	echo '#! /usr/bin/env sh' > deb/copy.sh
	echo 'for f in $$(ls); do scp $$f/*.deb user@192.168.99.106:~/DEBIAN_PKGS/$$f/main/; done' >> deb/copy.sh

fmt:
	find . -path ./.go -prune -o -name "*.go" -exec gofmt -w {} \;
	find . -path ./.go -prune -o -name "*.i2pkeys" -exec rm {} \;

upload-linux:
	github-release upload -R -u $(USER_GH) -r "$(packagename)" -t $(VERSION) -l `sha256sum ` -n "$(packagename)" -f "$(packagename)"

