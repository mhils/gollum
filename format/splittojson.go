// Copyright 2015-2016 trivago GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package format

import (
	"bytes"
	"fmt"
	"github.com/trivago/gollum/core"
	"github.com/trivago/gollum/shared"
)

// SplitToJSON formatter plugin
// SplitToJSON is a formatter that splits a message by a given token and puts
// the result into a JSON object by using an array based mapping
// Configuration example
//
//  - "stream.Broadcast":
//    Formatter: "format.SplitToJSON"
//    SplitToJSONDataFormatter: "format.Forward"
//    SplitToJSONToken: "|"
//    SplitToJSONKeys:
//      - "timestamp"
//      - "server"
//      - "error"
//
// SplitToJSONDataFormatter defines the formatter to apply before executing
// this formatter. Set to "format.Forward" by default.
//
// SplitToJSONToken defines the separator character to use when processing a
// message. By default this is set to "|".
//
// SplitToJSONKeys defines an array of keys to apply to the tokens generated
// by splitting a message by SplitToJSONToken. The keys listed here are
// applied to the resulting token array by index.
// This list is empty by default.
type SplitToJSON struct {
	base  core.Formatter
	token []byte
	keys  []string
}

func init() {
	shared.TypeRegistry.Register(SplitToJSON{})
}

// Configure initializes this formatter with values from a plugin config.
func (format *SplitToJSON) Configure(conf core.PluginConfig) error {
	plugin, err := core.NewPluginWithType(conf.GetString("SplitToJSONDataFormatter", "format.Forward"), conf)
	if err != nil {
		return err
	}
	format.base = plugin.(core.Formatter)
	format.token = []byte(conf.GetString("SplitToJSONToken", "|"))
	format.keys = conf.GetStringArray("SplitToJSONKeys", []string{})
	return nil
}

// Format returns the splitted message payload as json
func (format *SplitToJSON) Format(msg core.Message) ([]byte, core.MessageStreamID) {
	data, streamID := format.base.Format(msg)

	components := bytes.Split(data, format.token)
	maxIdx := shared.MinI(len(format.keys), len(components))
	jsonData := ""

	switch {
	case maxIdx == 0:
	case maxIdx == 1:
		jsonData = fmt.Sprintf("{%s:\"%s\"}", format.keys[0], components[0])
	default:
		for i := 0; i < maxIdx; i++ {
			key := shared.EscapeJSON(format.keys[i])
			value := shared.EscapeJSON(string(components[i]))
			switch {
			case i == 0:
				jsonData = fmt.Sprintf("{\"%s\":\"%s\"", key, value)
			case i == maxIdx-1:
				jsonData = fmt.Sprintf("%s,\"%s\":\"%s\"}", jsonData, key, value)
			default:
				jsonData = fmt.Sprintf("%s,\"%s\":\"%s\"", jsonData, key, value)
			}
		}
	}

	return []byte(jsonData), streamID
}
