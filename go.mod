module github.com/gookit/gcli/v3

go 1.13

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gookit/color v1.3.7
	github.com/gookit/goutil v0.3.7
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
)

// for develop
replace github.com/gookit/goutil => ../goutil
