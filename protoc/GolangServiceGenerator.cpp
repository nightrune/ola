// Protocol Buffers - Google's data interchange format
// Copyright 2008 Google Inc.  All rights reserved.
// http://code.google.com/p/protobuf/
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Author: kenton@google.com (Kenton Varda)
//  Based on original Protocol Buffers design by
//  Sanjay Ghemawat, Jeff Dean, and others.
//
// Based on Simon Newton's updated CppServiceGenerator

#include <google/protobuf/io/printer.h>

#include <map>
#include <string>

#include "protoc/GeneratorHelpers.h"
#include "protoc/GolangServiceGenerator.h"
#include "protoc/StrUtil.h"

namespace ola {

using google::protobuf::ServiceDescriptor;
using google::protobuf::MethodDescriptor;
using google::protobuf::io::Printer;
using std::map;
using std::string;

GolangServiceGenerator::GolangServiceGenerator(
    const ServiceDescriptor* descriptor, const Options& options)
  : descriptor_(descriptor) {
  vars_["classname"] = descriptor_->name();
  vars_["full_name"] = descriptor_->full_name();
  (void)options;
}

GolangServiceGenerator::~GolangServiceGenerator() {}

void GolangServiceGenerator::GenerateInterface(Printer* printer) {
  printer->Print(vars_, "type $classname$ struct {\n");
  printer->Print("}\n\n");
}

void GolangServiceGenerator::GenerateStubDefinition(Printer* printer) {
  printer->Print(vars_,
    "class $dllexport$$classname$_Stub : public $classname$ {\n"
    " public:\n");

  printer->Indent();

  printer->Print(vars_,
    "$classname$_Stub(ola::rpc::RpcChannel* channel);\n"
    "$classname$_Stub(ola::rpc::RpcChannel* channel,\n"
    "                 ::google::protobuf::Service::ChannelOwnership ownership"
    ");\n"
    "~$classname$_Stub();\n"
    "\n"
    "inline ola::rpc::RpcChannel* channel() { return channel_; }\n"
    "\n"
    "// implements $classname$ ------------------------------------------\n"
    "\n");

  GenerateMethodSignatures(printer);

  printer->Outdent();
  printer->Print(vars_,
    " private:\n"
    "  ola::rpc::RpcChannel* channel_;\n"
    "  bool owns_channel_;\n"
    "  GOOGLE_DISALLOW_EVIL_CONSTRUCTORS($classname$_Stub);\n"
    "};\n"
    "\n");
}

void GolangServiceGenerator::GenerateMethodSignatures(Printer* printer) {
  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["name"] = method->name();
    sub_vars["input_type"] = method->input_type()->full_name();
    sub_vars["output_type"] = method->output_type()->full_name();

    printer->Print(sub_vars,
      "$name$(request *$input_type$) (\n"
      "    response *$output_type$, err error)\n");
  }
}

// ===================================================================

void GolangServiceGenerator::GenerateDescriptorInitializer(
    Printer* printer, int index) {
  map<string, string> vars;
  vars["classname"] = descriptor_->name();
  vars["index"] = SimpleItoa(index);

  printer->Print(vars,
    "$classname$_descriptor_ = file->service($index$);\n");
}

// ===================================================================

void GolangServiceGenerator::GenerateImplementation(Printer* printer) {
  // Generate methods of the interface.
  GenerateNotImplementedMethods(printer);
  GenerateCallMethod(printer);
  GenerateGetPrototype(REQUEST, printer);
  GenerateGetPrototype(RESPONSE, printer);

  // Generate stub implementation.
  printer->Print(vars_,
    "type $classname$Stub struct {\n"
    "  ola.rpc.Channel _channel\n"
    "}\n"
    "\n"
    "func (m *$classname$Stub) SetChannel(channel *ola.rpc.Channel) {\n"
    "  m._channel = channel\n"
    "}\n\n");

  GenerateStubMethods(printer);
}

void GolangServiceGenerator::GenerateNotImplementedMethods(Printer* printer) {
  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["classname"] = descriptor_->name();
    sub_vars["name"] = method->name();
    sub_vars["index"] = SimpleItoa(i);
    sub_vars["input_type"] = method->input_type()->full_name();
    sub_vars["output_type"] = method->output_type()->full_name();

    printer->Print(sub_vars,
      "func (m *$classname$) $name$(\n"
      "    request *$input_type$) (\n"
      "    response *$output_type$, err error) {\n"
      "  return (nil, new(ola.rpc.NotImplemented))\n"
      "}\n\n");
  }
}

void GolangServiceGenerator::GenerateCallMethod(Printer* printer) {
  printer->Print(vars_,
    "func (m *$classname$) CallMethod(method *MethodDescriptor,\n"
    "    request *proto.Message) (\n"
    "    response *proto.Message, err error)){\n"
    "  switch method.index() {\n");

  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["name"] = method->name();
    sub_vars["index"] = SimpleItoa(i);
    sub_vars["input_type"] = method->input_type()->full_name();
    sub_vars["output_type"] = method->output_type()->full_name();

    // Note:  down_cast does not work here because it only works on pointers,
    //   not references.
    printer->Print(sub_vars,
      "    case $index$:\n"
      "      return $name$(request)\n");
  }

  printer->Print(vars_,
    "    default:\n"
    "      //TODO(Sean) Add some logging here\n"
    "      break;\n"
    "  }\n"
    "  return (nil, new(InvalidMethod))\n"
    "}\n"
    "\n");
}

void GolangServiceGenerator::GenerateGetPrototype(RequestOrResponse which,
                                            Printer* printer) {
  if (which == REQUEST) {
    printer->Print(vars_,
      "func (m *$classname$) GetRequestPrototype(\n");
  } else {
    printer->Print(vars_,
      "const (m *$classname$) GetResponsePrototype(\n");
  }

  printer->Print(vars_,
    "    method *MethodDescriptor) *proto.Message {\n"
    "  switch method.index() {\n");

  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    const google::protobuf::Descriptor* type =
      (which == REQUEST) ? method->input_type() : method->output_type();

    map<string, string> sub_vars;
    sub_vars["index"] = SimpleItoa(i);
    sub_vars["type"] = type->full_name();

    printer->Print(sub_vars,
      "    case $index$:\n"
      "      return new($type$)\n");
  }

  printer->Print(vars_,
    "    default:\n"
    "      //TODO(Sean) Add some logging\n"
    "      return nil\n"
    "  }\n"
    "}\n"
    "\n");
}

void GolangServiceGenerator::GenerateStubMethods(Printer* printer) {
  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["classname"] = descriptor_->name();
    sub_vars["name"] = method->name();
    sub_vars["index"] = SimpleItoa(i);
    sub_vars["input_type"] = method->input_type()->full_name();
    sub_vars["output_type"] = method->output_type()->full_name();

    printer->Print(sub_vars,
      "func (m *$classname$Stub) $name$(\n"
      "    request *$input_type$) (\n"
      "    response *$output_type$, err error) {\n"
      "  c := channel->CallMethod(GetMethodDescriptor($index$),\n"
      "      request);\n"
      "  respData <- c\n"
      "  return respData\n"
      "}\n\n");
  }
}
}  // namespace ola
