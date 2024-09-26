module github.com/ConnorsApps/poweredge-fanctl/fanctl

go 1.23.0

replace github.com/ConnorsApps/poweredge-fanctl => ../

require (
	github.com/ConnorsApps/poweredge-fanctl v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.9.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
