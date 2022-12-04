module github.com/patryk4815/usb-proxy

go 1.18

require (
	github.com/google/gousb v1.1.2
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/sirupsen/logrus v1.9.0
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace github.com/google/gousb => ./gohack/github.com/google/gousb
