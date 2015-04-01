# Build the go client

built_sources += golang/src/ola/Ola.pb.go \
								 golang/src/ola/OlaService.pb.go \
								 golang/src/ola/Rpc.pb.go

golang/src.ola/Ola.pb.go: golang/src/ola/Makefile.mk common/protocol/Ola.proto
	$(PROTOC) --go_out golang/src/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto

golang/src/ola/Rpc.pb.go: golang/src/ola/Makefile.mk common/rpc/Rpc.proto
	$(PROTOC) --go_out golang/src/ola --proto_path $(srcdir)/common/rpc $(srcdir)/common/rpc/Rpc.proto

golang/src/ola/OlaService.pb.go: golang/src/ola/Makefile.mk common/protocol/Ola.proto protoc/ola_protoc_golang_plugin$(EXEEXT)
	$(OLA_PROTOC) --goservice_out golang/src/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto
