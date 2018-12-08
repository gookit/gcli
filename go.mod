module github.com/gookit/gcli

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gookit/color v1.1.4
	github.com/gookit/filter v1.0.8
	github.com/gookit/goutil v0.0.0-20181022075743-4a30479feb86
	github.com/stretchr/testify v1.2.2
	golang.org/x/crypto v0.0.0-20180802221118-56440b844dfe
	golang.org/x/sys v0.0.0-20180806143827-98c5dad5d1a0 // indirect
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20180802221118-56440b844dfe
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180806143827-98c5dad5d1a0
)
