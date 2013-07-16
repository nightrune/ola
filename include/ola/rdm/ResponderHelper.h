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
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 *
 * ResponderHelper.h
 * Copyright (C) 2013 Simon Newton
 */

/**
 * @addtogroup rdm_resp
 * @{
 * @file ResponderHelper.h
 * @brief Helper methods for building RDM responders.
 * @}
 */
#ifndef INCLUDE_OLA_RDM_RESPONDERHELPER_H_
#define INCLUDE_OLA_RDM_RESPONDERHELPER_H_

#include <string>
#include <vector>
#include "ola/rdm/RDMCommand.h"
#include "ola/rdm/ResponderPersonality.h"

namespace ola {
namespace rdm {

/**
 * Helper methods for building RDM responders. These don't check that the
 * request is for the specified pid, so be sure to get it right!
 */
class ResponderHelper {
  public:
    // Request Parsing methods
    static bool ExtractUInt8(const RDMRequest *request, uint8_t *output);
    static bool ExtractUInt16(const RDMRequest *request, uint16_t *output);
    static bool ExtractUInt32(const RDMRequest *request, uint32_t *output);

    // Response Generation methods
    static const RDMResponse *GetDeviceInfo(
        const RDMRequest *request,
        uint16_t device_model,
        rdm_product_category product_category,
        uint32_t software_version,
        uint16_t dmx_footprint,
        uint8_t current_personality,
        uint8_t personality_count,
        uint16_t dmx_start_address,
        uint16_t sub_device_count,
        uint8_t sensor_count,
        uint8_t queued_message_count = 0);

    static const RDMResponse *GetDeviceInfo(
        const RDMRequest *request,
        uint16_t device_model,
        rdm_product_category product_category,
        uint32_t software_version,
        const PersonalityManager *personality_manager,
        uint16_t dmx_start_address,
        uint16_t sub_device_count,
        uint8_t sensor_count,
        uint8_t queued_message_count = 0);

    static const RDMResponse *GetProductDetailList(
        const RDMRequest *request,
        const std::vector<rdm_product_detail> &product_details,
        uint8_t queued_message_count = 0);

    static const RDMResponse *GetPersonality(
        const RDMRequest *request,
        const PersonalityManager *personality_manager,
        uint8_t queued_message_count = 0);
    static const RDMResponse *SetPersonality(
        const RDMRequest *request,
        PersonalityManager *personality_manager,
        uint16_t start_address,
        uint8_t queued_message_count = 0);
    static const RDMResponse *GetPersonalityDescription(
        const RDMRequest *request,
        const PersonalityManager *personality_manager,
        uint8_t queued_message_count = 0);
    static const RDMResponse *GetDmxAddress(
        const RDMRequest *request,
        const PersonalityManager *personality_manager,
        uint16_t start_address,
        uint8_t queued_message_count = 0);
    static const RDMResponse *SetDmxAddress(
        const RDMRequest *request,
        const PersonalityManager *personality_manager,
        uint16_t *dmx_address,
        uint8_t queued_message_count = 0);

    static const RDMResponse *GetRealTimeClock(
        const RDMRequest *request,
        uint8_t queued_message_count = 0);

    static const RDMResponse *GetString(const RDMRequest *request,
                                        const std::string &value,
                                        uint8_t queued_message_count = 0);
    static const RDMResponse *GetBoolValue(
        const RDMRequest *request, bool value,
        uint8_t queued_message_count = 0);
    static const RDMResponse *SetBoolValue(
        const RDMRequest *request, bool *value,
        uint8_t queued_message_count = 0);

    static const RDMResponse *GetUInt8Value(
        const RDMRequest *request, uint8_t value,
        uint8_t queued_message_count = 0);
    static const RDMResponse *GetUInt16Value(
        const RDMRequest *request, uint16_t value,
        uint8_t queued_message_count = 0);
    static const RDMResponse *GetUInt32Value(
        const RDMRequest *request, uint32_t value,
        uint8_t queued_message_count = 0);
    static const RDMResponse *SetUInt8Value(
        const RDMRequest *request, uint8_t *value,
        uint8_t queued_message_count = 0);
    static const RDMResponse *SetUInt16Value(
        const RDMRequest *request, uint16_t *value,
        uint8_t queued_message_count = 0);
    static const RDMResponse *SetUInt32Value(
        const RDMRequest *request, uint32_t *value,
        uint8_t queued_message_count = 0);
};
}  // namespace rdm
}  // namespace ola
#endif  // INCLUDE_OLA_RDM_RESPONDERHELPER_H_
