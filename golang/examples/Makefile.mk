# Build the go examples

built_sources += \
    golang/examples/ola_send_dmx

golang/examples/ola_send_dmx: golang/examples/ola_send_dmx.go
	go build -o golang/examples/ola_send_dmx golang/examples/ola_send_dmx.go
