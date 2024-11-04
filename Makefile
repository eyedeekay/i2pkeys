
USER_GH=eyedeekay
VERSION=0.33.8
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

test-basic:
	go test -v -run Test_Basic

test-basic-lookup:
	go test -v -run Test_Basic_Lookup

test-newi2paddrfromstring:
	go test -v -run Test_NewI2PAddrFromString

test-i2paddr:
	go test -v -run Test_I2PAddr

test-desthashfromstring:
	go test -v -run Test_DestHashFromString

test-i2paddr-to-bytes:
	go test -v -run Test_I2PAddrToBytes

test-key-generation-and-handling:
	go test -v -run Test_KeyGenerationAndHandling

# Subtest targets
test-newi2paddrfromstring-valid:
	go test -v -run Test_NewI2PAddrFromString/Valid_base64_address

test-newi2paddrfromstring-invalid:
	go test -v -run Test_NewI2PAddrFromString/Invalid_address

test-newi2paddrfromstring-base32:
	go test -v -run Test_NewI2PAddrFromString/Base32_address

test-newi2paddrfromstring-empty:
	go test -v -run Test_NewI2PAddrFromString/Empty_address

test-newi2paddrfromstring-i2p-suffix:
	go test -v -run Test_NewI2PAddrFromString/Address_with_.i2p_suffix

test-i2paddr-base32-suffix:
	go test -v -run Test_I2PAddr/Base32_suffix

test-i2paddr-base32-length:
	go test -v -run Test_I2PAddr/Base32_length

test-desthashfromstring-valid:
	go test -v -run Test_DestHashFromString/Valid_hash

test-desthashfromstring-invalid:
	go test -v -run Test_DestHashFromString/Invalid_hash

test-desthashfromstring-empty:
	go test -v -run Test_DestHashFromString/Empty_hash

test-i2paddr-to-bytes-roundtrip:
	go test -v -run Test_I2PAddrToBytes/ToBytes_and_back

test-i2paddr-to-bytes-comparison:
	go test -v -run Test_I2PAddrToBytes/Direct_decoding_comparison

test-key-generation-and-handling-loadkeys:
	go test -v -run Test_KeyGenerationAndHandling/LoadKeysIncompat

test-key-generation-and-handling-storekeys-incompat:
	go test -v -run Test_KeyGenerationAndHandling/StoreKeysIncompat

test-key-generation-and-handling-storekeys:
	go test -v -run Test_KeyGenerationAndHandling/StoreKeys

test-key-storage:
	go test -v -run Test_KeyStorageAndLoading

# Individual key storage subtests
test-key-storage-file:
	go test -v -run Test_KeyStorageAndLoading/StoreAndLoadFile

test-key-storage-incompat:
	go test -v -run Test_KeyStorageAndLoading/StoreAndLoadIncompat

test-key-storage-nonexistent:
	go test -v -run Test_KeyStorageAndLoading/LoadNonexistentFile

# Aggregate targets
test-all:
	go test -v ./...

test-subtests: test-newi2paddrfromstring-valid test-newi2paddrfromstring-invalid test-newi2paddrfromstring-base32 test-newi2paddrfromstring-empty test-newi2paddrfromstring-i2p-suffix test-i2paddr-base32-suffix test-i2paddr-base32-length test-desthashfromstring-valid test-desthashfromstring-invalid test-desthashfromstring-empty test-i2paddr-to-bytes-roundtrip test-i2paddr-to-bytes-comparison test-key-generation-and-handling-loadkeys test-key-generation-and-handling-storekeys-incompat test-key-generation-and-handling-storekeys test-key-storage-file test-key-storage-incompat test-key-storage-nonexistent

test: test-basic test-basic-lookup test-newi2paddrfromstring test-i2paddr test-desthashfromstring test-i2paddr-to-bytes test-key-generation-and-handling test-key-storage test-subtests test-all