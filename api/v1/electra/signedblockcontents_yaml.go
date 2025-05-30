// Copyright © 2024 Attestant Limited.
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

package electra

import (
	"bytes"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/electra"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// signedBlockContentsYAML is the spec representation of the struct.
type signedBlockContentsYAML struct {
	SignedBlock *electra.SignedBeaconBlock `yaml:"signed_block"`
	KZGProofs   []deneb.KZGProof           `yaml:"kzg_proofs"`
	Blobs       []deneb.Blob               `yaml:"blobs"`
}

// MarshalYAML implements yaml.Marshaler.
func (s *SignedBlockContents) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedBlockContentsYAML{
		SignedBlock: s.SignedBlock,
		KZGProofs:   s.KZGProofs,
		Blobs:       s.Blobs,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedBlockContents) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var unmarshaled signedBlockContentsJSON
	if err := yaml.Unmarshal(input, &unmarshaled); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}

	marshaled, err := json.Marshal(unmarshaled)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return s.UnmarshalJSON(marshaled)
}
