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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// signedBeaconBlockJSON is the spec representation of the struct.
type signedBeaconBlockJSON struct {
	Message   *BeaconBlock `json:"message"`
	Signature string       `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBeaconBlockJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBeaconBlock) UnmarshalJSON(input []byte) error {
	var data signedBeaconBlockJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&data)
}

func (s *SignedBeaconBlock) unpack(data *signedBeaconBlockJSON) error {
	if data.Message == nil {
		return errors.New("message missing")
	}
	s.Message = data.Message
	if data.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(data.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return fmt.Errorf("incorrect length %d for signature", len(signature))
	}
	copy(s.Signature[:], signature)

	return nil
}
