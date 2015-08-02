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

void GolangServiceGenerator::GenerateInit(Printer* printer) {
  printer->Print(vars_,
      "var method_descriptors map[uint32]*MethodDescriptor\n\n");
  printer->Print("func init() {\n");
  printer->Indent();
  printer->Print(vars_,
      "method_descriptors = make(map[uint32]*MethodDescriptor)\n");
  GenerateMethodDescriptors(printer);
  printer->Outdent();
  printer->Print("}\n\n");
}

void GolangServiceGenerator::GenerateType(Printer* printer) {
  printer->Print(vars_, "type $classname$ struct {\n"
  "}\n\n"
  "func (m *$classname$) GetMethodDescriptor(method_index uint32) *MethodDescriptor {\n"
  "  // We handle this because otherwise the golang map would return the 0 value\n"
  "  if method_index >= uint32(len(method_descriptors)) {\n"
  "    return NewMethodDescriptor(uint32(len(method_descriptors)),\n"
  "        \"Invalid Method\", \"\", \"\")\n"
  "  }\n"
  "  return method_descriptors[method_index]\n"
  "}\n\n");
}

void GolangServiceGenerator::GenerateMethodDescriptors(Printer *printer) {
  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["name"] = method->name();
    sub_vars["index"] = SimpleItoa(i);
    sub_vars["output_type"] = method->output_type()->name();
    sub_vars["input_type"] = method->input_type()->name();
    printer->Print(sub_vars,
        "method_descriptors[$index$] = NewMethodDescriptor($index$, \"$name$\","
        "\"$output_type$\", \"$input_type$\")\n");
  }
}

void GolangServiceGenerator::GenerateMethodSignatures(Printer* printer) {
  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["name"] = method->name();
    sub_vars["input_type"] = GoTypeName(method->input_type());
    sub_vars["output_type"] = GoTypeName(method->output_type());

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
  GenerateType(printer);
  GenerateNotImplementedMethods(printer);
  GenerateCallMethod(printer);

  // Generate stub implementation.
  printer->Print(vars_,
    "type $classname$Stub struct {\n"
    "  $classname$\n"
    "  _channel *RpcChannel\n"
    "}\n"
    "\n"
    "func (m *$classname$Stub) SetChannel(channel *RpcChannel) {\n"
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
    sub_vars["input_type"] = GoTypeName(method->input_type());
    sub_vars["output_type"] = GoTypeName(method->output_type());

    printer->Print(sub_vars,
      "func (m *$classname$) $name$(\n"
      "    request *$input_type$) (\n"
      "    response *$output_type$, err error) {\n"
      "  return nil, NewNotImplemented(\"$name$ not implemented!\")\n"
      "}\n\n");
  }
}

void GolangServiceGenerator::GenerateCallMethod(Printer* printer) {
  printer->Print(vars_,
    "func (m *$classname$) CallMethod(method *MethodDescriptor,\n"
    "    requestData []byte) (\n"
    "    responseData []byte, err error) {\n"
    "  switch method.Index() {\n");

  for (int i = 0; i < descriptor_->method_count(); i++) {
    const MethodDescriptor* method = descriptor_->method(i);
    map<string, string> sub_vars;
    sub_vars["name"] = method->name();
    sub_vars["index"] = SimpleItoa(i);
    sub_vars["input_type"] = GoTypeName(method->input_type());
    sub_vars["output_type"] = GoTypeName(method->output_type());

    // Note:  down_cast does not work here because it only works on pointers,
    //   not references.
    printer->Print(sub_vars,
      "    case $index$:\n"
      "      request := new($input_type$)\n"
      "      err = proto.Unmarshal(requestData, request)\n"
      "      if err != nil {\n"
      "        return nil, err\n"
      "      }\n"
      "      response, err := m.$name$(request)\n"
      "      if err != nil {\n"
      "        return nil, err\n"
      "      }\n"
      "      responseData, err = proto.Marshal(response)\n"
      "      return responseData, err\n");
  }

  printer->Print(vars_,
    "    default:\n"
    "      //TODO(Sean) Add some logging here\n"
    "      break;\n"
    "  }\n"
    "  return nil, new(InvalidMethod)\n"
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
    sub_vars["input_type"] = GoTypeName(method->input_type());
    sub_vars["output_type"] = GoTypeName(method->output_type());

    printer->Print(sub_vars,
      "func (m *$classname$Stub) $name$(\n"
      "    request *$input_type$) (\n"
      "    response *$output_type$, err error) {\n"
      "  var reqData []byte\n"
      "  reqData, err = proto.Marshal(request)\n"
      "  if err != nil {\n"
      "    return nil, err\n"
      "  }\n"
      "  c := make(chan *ResponseData, 1)\n"
      "  go m._channel.CallMethod(m.GetMethodDescriptor($index$),\n"
      "      reqData, c);\n"
      "  var respData *ResponseData\n"
      "  respData = <- c\n"
      "  if respData.err != nil {\n"
      "    return nil, respData.err\n"
      "  }\n"
      "  response = new($output_type$)\n"
      "  err = proto.Unmarshal(respData.data, response)\n"
      "  if err != nil {\n"
      "    return nil, err\n"
      "  }\n"
      "  return response, nil\n"
      "}\n\n");
  }
}

}  // namespace ola