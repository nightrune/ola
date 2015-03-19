/*
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA
 *
 * GoGenerator.cpp
 * Copyright (C) 2013 Simon Newton, Sean Sill
 */

#include <google/protobuf/descriptor.h>
#include <google/protobuf/io/printer.h>
#include <google/protobuf/io/zero_copy_stream.h>

#include <memory>
#include <string>

#include "protoc/GolangGenerator.h"
#include "protoc/GolangFileGenerator.h"
#include "protoc/GeneratorHelpers.h"
#include "protoc/GolangServiceGenerator.h"
#include "protoc/StrUtil.h"

namespace ola {

using google::protobuf::FileDescriptor;
using google::protobuf::ServiceDescriptor;
using google::protobuf::compiler::OutputDirectory;
using google::protobuf::io::Printer;
using std::auto_ptr;
using std::string;

bool GolangGenerator::Generate(const FileDescriptor *file,
                            const string&,
                            OutputDirectory *generator_context,
                            string*) const {
  string basename = StripProto(file->name()) + "Service";

  string code_name = basename + ".pb.go";
  GolangFileGenerator file_generator(file, basename);

  auto_ptr<google::protobuf::io::ZeroCopyOutputStream> code_output(
    generator_context->Open(code_name));
  Printer code_printer(code_output.get(), '$');
  file_generator.GenerateImplementation(&code_printer);

  return true;
}
}  // namespace ola
