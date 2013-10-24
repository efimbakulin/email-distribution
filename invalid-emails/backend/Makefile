VERSION=0.1
DEB_SOURCES=$(shell find dist/debian -type f)
ARCH=amd64

deb: invalid-emails-backend_$(VERSION)_$(ARCH).deb

clean:
	rm -r debian invalid-emails-backend_$(VERSION)_$(ARCH).deb backend

invalid-emails-backend_$(VERSION)_$(ARCH).deb: invalid-emails-backend $(DEB_SOURCES) Makefile
	# check invalid-emails-backend.git is x86-64
	sed s/Version\:.*$$/Version\:\ $(VERSION)/g dist\/debian\/DEBIAN\/control -i
	file -b backend | grep -e "^ELF 64-bit LSB executable, x86-64, version 1 (SYSV)"
	rm -fr debian
	cp -pR dist/debian debian
	mkdir -p debian/usr/bin
	cp backend debian/usr/bin/invalid-emails-backend
	find debian -type f -name .DS_Store -exec rm '{}' ';'
	find debian -type d -exec chmod 0755 '{}' ';'
	fakeroot dpkg-deb --build debian invalid-emails-backend_$(VERSION)_$(ARCH).deb

invalid-emails-backend: invalid-emails-backend.go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build
