/*
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Library General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 *
 * E133StatusPDU.h
 * The E133StatusPDU
 * Copyright (C) 2013 Simon Newton
 */

#ifndef PLUGINS_E131_E131_E133STATUSPDU_H_
#define PLUGINS_E131_E131_E133STATUSPDU_H_

#include <ola/io/IOStack.h>
#include <string>
#include "ola/e133/E133Enums.h"
#include "plugins/e131/e131/PDU.h"

namespace ola {
namespace plugin {
namespace e131 {

class E133StatusPDU : private PDU {
 public:
    static void PrependPDU(ola::io::IOStack *stack,
                           ola::e133::E133StatusCode status_code,
                           const std::string &status);
};
}  // namespace e131
}  // namespace plugin
}  // namespace ola
#endif  // PLUGINS_E131_E131_E133STATUSPDU_H_
