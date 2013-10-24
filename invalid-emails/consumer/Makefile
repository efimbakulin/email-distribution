VERSION=0.1
DEB_SOURCES=$(shell find dist/debian -type f)
ARCH=amd64

deb: invalid-emails-consumer_$(VERSION)_$(ARCH).deb

clean:
	rm -r debian invalid-emails-consumer_$(VERSION)_$(ARCH).deb consumer

invalid-emails-consumer_$(VERSION)_$(ARCH).deb: invalid-emails-consumer $(DEB_SOURCES) Makefile
	# check invalid-emails-consumer.git is x86-64
	sed s/Version\:.*$$/Version\:\ $(VERSION)/g dist\/debian\/DEBIAN\/control -i
	file -b consumer | grep -e "^ELF 64-bit LSB executable, x86-64, version 1 (SYSV)"
	rm -fr debian
	cp -pR dist/debian debian
	mkdir -p debian/usr/bin
	cp consumer debian/usr/bin/invalid-emails-consumer
	find debian -type f -name .DS_Store -exec rm '{}' ';'
	find debian -type d -exec chmod 0755 '{}' ';'
	fakeroot dpkg-deb --build debian invalid-emails-consumer_$(VERSION)_$(ARCH).deb

invalid-emails-consumer: invalid-emails-consumer.go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build
