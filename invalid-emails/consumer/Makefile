APPNAME=invalid-emails-consumer
VERSION=$(shell cat version)
DEB_SOURCES=$(shell find dist/debian -type f)
ARCH=amd64
GOFLAGS=-ldflags "-X main.version $(VERSION) -X main.applicationName $(APPNAME)"

deb: $(APPNAME)_$(VERSION)_$(ARCH).deb

clean:
	rm -r debian $(APPNAME)_$(VERSION)_$(ARCH).deb $(APPNAME)

$(APPNAME)_$(VERSION)_$(ARCH).deb: $(APPNAME) $(DEB_SOURCES) Makefile
	# check invalid-emails-backend.git is x86-64
	sed s/Version\:.*$$/Version\:\ $(VERSION)/g dist\/debian\/DEBIAN\/control -i
	file -b $(APPNAME) | grep -e "^ELF 64-bit LSB executable, x86-64, version 1 (SYSV)"
	rm -fr debian
	cp -pR dist/debian debian
	mkdir -p debian/usr/bin
	cp $(APPNAME) debian/usr/bin/$(APPNAME)
	find debian -type f -name .DS_Store -exec rm '{}' ';'
	find debian -type d -exec chmod 0755 '{}' ';'
	fakeroot dpkg-deb --build debian $(APPNAME)_$(VERSION)_$(ARCH).deb

$(APPNAME): $(APPNAME).go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build  $(GOFLAGS) -o $(APPNAME)
