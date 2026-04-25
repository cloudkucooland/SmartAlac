module github.com/cloudkucooland/SmartAlac

go 1.23

toolchain go1.24.1

require (
	github.com/Sorrow446/go-mp4tag v0.0.0-20240130220823-68ce31d53e37
	github.com/ebitengine/purego v0.10.0
	github.com/jo-hoe/chromaprint v0.0.0-20260413105333-99a143cad505
	github.com/kr/pretty v0.3.1
	github.com/michiwend/gomusicbrainz v0.0.0-20181012083520-6c07e13dd396
	github.com/urfave/cli/v3 v3.8.0
	go.uber.org/ratelimit v0.3.1
)

require (
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/michiwend/golang-pretty v0.0.0-20141116172505-8ac61812ea3f // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
)

replace github.com/Sorrow446/go-mp4tag => github.com/cloudkucooland/go-mp4tag v0.0.0-20260424214855-e79d6c04b1c0
