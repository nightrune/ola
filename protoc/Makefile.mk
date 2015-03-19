# Programs
#########################
if BUILD_OLA_PROTOC_PLUGIN
noinst_PROGRAMS += protoc/ola_protoc_plugin protoc/ola_protoc_golang_plugin

protoc_ola_protoc_plugin_SOURCES = \
    protoc/CppFileGenerator.cpp \
    protoc/CppFileGenerator.h \
    protoc/CppGenerator.cpp \
    protoc/CppGenerator.h \
    protoc/GeneratorHelpers.cpp \
    protoc/GeneratorHelpers.h \
    protoc/CppServiceGenerator.cpp \
    protoc/CppServiceGenerator.h \
    protoc/StrUtil.cpp \
    protoc/StrUtil.h \
  protoc/ola-protoc-generator-plugin.cpp
protoc_ola_protoc_plugin_LDADD = $(libprotobuf_LIBS) -lprotoc

protoc_ola_protoc_golang_plugin_SOURCES = \
    protoc/GolangFileGenerator.cpp \
    protoc/GolangFileGenerator.h \
    protoc/GolangGenerator.cpp \
    protoc/GolangGenerator.h \
    protoc/GeneratorHelpers.cpp \
    protoc/GeneratorHelpers.h \
    protoc/GolangServiceGenerator.cpp \
    protoc/GolangServiceGenerator.h \
    protoc/StrUtil.cpp \
    protoc/StrUtil.h \
  protoc/ola-protoc-golang-generator-plugin.cpp

protoc_ola_protoc_golang_plugin_LDADD = $(libprotobuf_LIBS) -lprotoc

else

# If we're using a different ola_protoc_plugin, we need to provide a rule to
# create this file since the generated service configs depend on it.
protoc/ola_protoc_plugin$(EXEEXT):
  touch protoc/ola_protoc_plugin$(EXEEXT)

endif
