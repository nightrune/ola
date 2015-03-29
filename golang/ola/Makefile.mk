# Build the go client

built_sources += golang/ola/Ola.pb.go golang/ola/Service.pb.go

golang/ola/Ola.pb.go: golang/ola/Makefile.mk common/protocol/Ola.proto
	$(PROTOC) --go_out golang/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto

golang/ola/Service.pb.go: golang/ola/Makefile.mk common/protocol/Ola.proto protoc/ola_protoc_golang_plugin$(EXEEXT)
	$(OLA_PROTOC) --goservice_out golang/ola --proto_path $(srcdir)/common/protocol $(srcdir)/common/protocol/Ola.proto
