module github.com/gookit/gcli

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gookit/color v1.1.7
	github.com/gookit/filter v1.0.10
	github.com/gookit/goutil v0.1.3
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/sys v0.0.0-20190308023053-584f3b12f43e // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/sys v0.0.0-20190308023053-584f3b12f43e => github.com/golang/sys v0.0.0-20190308023053-584f3b12f43e
)
