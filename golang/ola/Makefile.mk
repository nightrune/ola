# Build the go client

built_sources += golang/ola/Ola.pb.go \
								 golang/ola/OlaService.pb.go \
								 golang/ola/Rpc.pb.go

golang/ola/Ola.pb.go: golang/ola/Makefile.mk common/protocol/Ola.proto
	$(PROTOC) --go_out golang/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto

golang/ola/Rpc.pb.go: golang/ola/Makefile.mk common/rpc/Rpc.proto
	$(PROTOC) --go_out golang/ola --proto_path $(srcdir)/common/rpc $(srcdir)/common/rpc/Rpc.proto

golang/ola/OlaService.pb.go: golang/ola/Makefile.mk common/protocol/Ola.proto protoc/ola_protoc_golang_plugin$(EXEEXT)
	$(OLA_PROTOC) --goservice_out golang/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto
