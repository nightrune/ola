# Build the go client

if BUILD_GO_LIBS

built_sources += golang/src/ola/ola_proto/Ola.pb.go \
								 golang/src/ola/OlaService.pb.go \
								 golang/src/ola/ola_rpc/Rpc.pb.go

golang/src/ola/ola_proto/Ola.pb.go: golang/src/ola/Makefile.mk common/protocol/Ola.proto
	$(PROTOC) --go_out golang/src/ola/ola_proto --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto

golang/src/ola/ola_rpc/Rpc.pb.go: golang/src/ola/Makefile.mk common/rpc/Rpc.proto
	$(PROTOC) --go_out golang/src/ola/ola_rpc --proto_path $(srcdir)/common/rpc $(srcdir)/common/rpc/Rpc.proto

golang/src/ola/OlaService.pb.go: golang/src/ola/Makefile.mk common/protocol/Ola.proto protoc/ola_protoc_golang_plugin$(EXEEXT)
	$(OLA_PROTOC) --goservice_out golang/src/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto

else

endif
