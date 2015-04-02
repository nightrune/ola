# Build the go examples
# We force autotools to build go binaries with the built_sources hack
# but you'll need to specify how to build
# Should get updated in configure.ac later

# These are actually binaries to be built
if BUILD_GO_LIBS

built_sources += \
    golang/examples/ola_send_dmx

golang/examples/ola_send_dmx: golang/examples/ola_send_dmx.go
	$(OLA_GO) build -o golang/examples/ola_send_dmx golang/examples/ola_send_dmx.go

endif
