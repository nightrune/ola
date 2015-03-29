# Build the go examples
# We force autotools to build go binaries with the built_sources hack
# but you'll need to specify how to build
# Should get updated in configure.ac later

# These are actually binaries to be built
built_sources += \
    golang/examples/ola_send_dmx

golang/examples/ola_send_dmx: golang/examples/ola_send_dmx.go
	./go_helper.sh -o golang/examples/ola_send_dmx golang/examples/ola_send_dmx.go
