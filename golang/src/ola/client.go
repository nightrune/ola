/**
 * @file ola.go
 *
 *
 *
 */
package ola

type DmxData struct {
}

type SendDMXArgs struct {
}

// RDM Data
type UID struct {
}

type SendRDMArgs struct {
}

type RDMMetadata struct {
}

type RDMRequest struct {
}

type RDMResponse struct {
}

type PluginState struct {
	_name                string
	_enabled             bool
	_active              bool
	_preferences_source  string
	_conflicting_plugins []Plugin
}

type Plugin struct {
	_id     uint
	_name   string
	_active bool
}

type PriorityCapability struct {
}

type PriorityMode struct {
}

type Port struct {
	_id                  uint
	_universe            uint
	_active              bool
	_description         string
	_priority_capability PriorityCapability
	_priority_mode       PriorityMode
	_priority            uint8
	_supports_rdm        bool
}

type InputPort struct {
	Port
}

type OutputPort struct {
	Port
}

type PortDirection struct {
}

type Device struct {
	_id           string
	_alias        int
	_name         string
	_plugin_id    int
	_input_ports  []InputPort
	_output_ports []OutputPort
}

type MergeMode struct {
}

type Universe struct {
	_id                uint
	_merge_mode        MergeMode
	_name              string
	_input_port_count  uint
	_output_port_count uint
	_rdm_device_count  uint
}

type PatchAction struct {
}

type DiscoveryType struct {
}

type Client struct {
	_channel RpcChannel
}

// Plugin methods
func (self *Client) FetchPlugins() ([]Plugin, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) ReloadPlugins() (bool, error) {
	return false, NewNotImplemented("Not implemented in client")
}

func (self *Client) FetchPluginDescription(plugin_id int) (string, error) {
	return "", NewNotImplemented("Not implemented in client")
}

func (self *Client) FetchPluginState(plugin_id int) (*PluginState, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

// Device Methods
func (self *Client) FetchDeviceInfo(plugin_id int) ([]Device, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) ConfigureDevice(device_alias uint, msg []byte) ([]byte,
	error) {
	return nil, NewNotImplemented("Not implemented in client")
}

// Port Methods
func (self *Client) SetPortPriorityInherit(device_alias uint, port uint,
	direction PortDirection) error {
	return NewNotImplemented("Not implemented in client")
}

func (self *Client) SetPortPriorityOverride(device_alias uint, port uint,
	direction PortDirection, value uint8) error {
	return NewNotImplemented("Not implemented in client")
}

// Universe Methods
func (self *Client) FetchUniverseList() ([]Universe, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) FetchUniverseInfo(universe uint) (*Universe, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) SetUniverseName(universe uint, name string) error {
	return NewNotImplemented("Not implemented in client")
}

func (self *Client) SetUniverseMergeMode(universe uint, mode MergeMode) error {
	return NewNotImplemented("Not implemented in client")
}

func (self *Client) Patch(device_alias uint,
	port uint,
	direction PortDirection,
	action PatchAction,
	universe uint) error {
	return NewNotImplemented("Not implemented in client")
}

// Dmx Functions
/**
 * @notes We'll stray from the normal api here since a channel is more
 * idiomatic for golang
 */
func (self *Client) RegisterUniverse(universe uint) (chan *DmxData, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) DeregisterUniverse(universe uint) error {
	return NewNotImplemented("Not implemented in client")
}

func (self *Client) SendDMX(universe uint, data *DmxData, args SendDMXArgs) {
}

func (self *Client) FetchDMX(universe uint) (*DmxData, error) {
	return nil, NewNotImplemented("Not implemented in client")
}

// RDM Functions
// No sets in go?!
func (self *Client) RunDiscovery(universe uint, disc_type DiscoveryType) ([]UID,
	error) {
	return nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) SetSourceUID(uid UID) error {
	return NewNotImplemented("Not implemented in client")
}

/**
 * @notes We stray from the C++ api since we don't need callbacks
 */
func (self *Client) RDMGet(universe uint,
	uid UID,
	sub_device uint16,
	pid uint16,
	data []byte,
	args SendRDMArgs) (*RDMMetadata, *RDMResponse, error) {
	return nil, nil, NewNotImplemented("Not implemented in client")
}

func (self *Client) RDMSet(universe uint,
	uid UID,
	sub_device uint16,
	pid uint16,
	data []byte,
	args SendRDMArgs) (*RDMMetadata, *RDMResponse, error) {
	return nil, nil, NewNotImplemented("Not implemented in client")
}
