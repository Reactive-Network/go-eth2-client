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

	"github.com/attestantio/go-eth2-client/spec/electra"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blindedBeaconBlockBodyJSON is the spec representation of the struct.
type blindedBeaconBlockBodyJSON struct {
	RANDAOReveal           string                                `json:"randao_reveal"`
	ETH1Data               *phase0.ETH1Data                      `json:"eth1_data"`
	Graffiti               string                                `json:"graffiti"`
	ProposerSlashings      []*phase0.ProposerSlashing            `json:"proposer_slashings"`
	AttesterSlashings      []*electra.AttesterSlashing           `json:"attester_slashings"`
	Attestations           []*electra.Attestation                `json:"attestations"`
	Deposits               []*phase0.Deposit                     `json:"deposits"`
	VoluntaryExits         []*phase0.SignedVoluntaryExit         `json:"voluntary_exits"`
	SyncAggregate          *altair.SyncAggregate                 `json:"sync_aggregate"`
	ExecutionPayloadHeader *deneb.ExecutionPayloadHeader         `json:"execution_payload_header"`
	BLSToExecutionChanges  []*capella.SignedBLSToExecutionChange `json:"bls_to_execution_changes"`
	BlobKZGCommitments     []string                              `json:"blob_kzg_commitments"`
	ExecutionRequests      *electra.ExecutionRequests            `json:"execution_requests"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlindedBeaconBlockBody) MarshalJSON() ([]byte, error) {
	blobKZGCommitments := make([]string, len(b.BlobKZGCommitments))
	for i := range b.BlobKZGCommitments {
		blobKZGCommitments[i] = b.BlobKZGCommitments[i].String()
	}

	return json.Marshal(&blindedBeaconBlockBodyJSON{
		RANDAOReveal:           fmt.Sprintf("%#x", b.RANDAOReveal),
		ETH1Data:               b.ETH1Data,
		Graffiti:               fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings:      b.ProposerSlashings,
		AttesterSlashings:      b.AttesterSlashings,
		Attestations:           b.Attestations,
		Deposits:               b.Deposits,
		VoluntaryExits:         b.VoluntaryExits,
		SyncAggregate:          b.SyncAggregate,
		ExecutionPayloadHeader: b.ExecutionPayloadHeader,
		BLSToExecutionChanges:  b.BLSToExecutionChanges,
		BlobKZGCommitments:     blobKZGCommitments,
		ExecutionRequests:      b.ExecutionRequests,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlindedBeaconBlockBody) UnmarshalJSON(input []byte) error {
	var data blindedBeaconBlockBodyJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

//nolint:gocyclo
func (b *BlindedBeaconBlockBody) unpack(data *blindedBeaconBlockBodyJSON) error {
	if data.RANDAOReveal == "" {
		return errors.New("RANDAO reveal missing")
	}
	randaoReveal, err := hex.DecodeString(strings.TrimPrefix(data.RANDAOReveal, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for RANDAO reveal")
	}
	if len(randaoReveal) != phase0.SignatureLength {
		return errors.New("incorrect length for RANDAO reveal")
	}
	copy(b.RANDAOReveal[:], randaoReveal)
	if data.ETH1Data == nil {
		return errors.New("ETH1 data missing")
	}
	b.ETH1Data = data.ETH1Data
	if data.Graffiti == "" {
		return errors.New("graffiti missing")
	}
	graffiti, err := hex.DecodeString(strings.TrimPrefix(data.Graffiti, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for graffiti")
	}
	if len(graffiti) != phase0.GraffitiLength {
		return errors.New("incorrect length for graffiti")
	}
	copy(b.Graffiti[:], graffiti)
	if data.ProposerSlashings == nil {
		return errors.New("proposer slashings missing")
	}
	for i := range data.ProposerSlashings {
		if data.ProposerSlashings[i] == nil {
			return fmt.Errorf("proposer slashings entry %d missing", i)
		}
	}
	b.ProposerSlashings = data.ProposerSlashings
	if data.AttesterSlashings == nil {
		return errors.New("attester slashings missing")
	}
	for i := range data.AttesterSlashings {
		if data.AttesterSlashings[i] == nil {
			return fmt.Errorf("attester slashings entry %d missing", i)
		}
	}
	b.AttesterSlashings = data.AttesterSlashings
	if data.Attestations == nil {
		return errors.New("attestations missing")
	}
	for i := range data.Attestations {
		if data.Attestations[i] == nil {
			return fmt.Errorf("attestations entry %d missing", i)
		}
	}
	b.Attestations = data.Attestations
	if data.Deposits == nil {
		return errors.New("deposits missing")
	}
	for i := range data.Deposits {
		if data.Deposits[i] == nil {
			return fmt.Errorf("deposits entry %d missing", i)
		}
	}
	b.Deposits = data.Deposits
	if data.VoluntaryExits == nil {
		return errors.New("voluntary exits missing")
	}
	for i := range data.VoluntaryExits {
		if data.VoluntaryExits[i] == nil {
			return fmt.Errorf("voluntary exits entry %d missing", i)
		}
	}
	b.VoluntaryExits = data.VoluntaryExits
	if data.SyncAggregate == nil {
		return errors.New("sync aggregate missing")
	}
	b.SyncAggregate = data.SyncAggregate
	if data.ExecutionPayloadHeader == nil {
		return errors.New("execution payload header missing")
	}
	b.ExecutionPayloadHeader = data.ExecutionPayloadHeader
	if data.BLSToExecutionChanges == nil {
		b.BLSToExecutionChanges = make([]*capella.SignedBLSToExecutionChange, 0)
	} else {
		for i := range data.BLSToExecutionChanges {
			if data.BLSToExecutionChanges[i] == nil {
				return fmt.Errorf("bls to execution changes entry %d missing", i)
			}
		}
		b.BLSToExecutionChanges = data.BLSToExecutionChanges
	}
	if data.BlobKZGCommitments == nil {
		return errors.New("blob KZG commitments missing")
	}
	for i := range data.BlobKZGCommitments {
		if data.BlobKZGCommitments[i] == "" {
			return fmt.Errorf("blob KZG commitments entry %d missing", i)
		}
	}
	b.BlobKZGCommitments = make([]deneb.KZGCommitment, len(data.BlobKZGCommitments))
	for i := range data.BlobKZGCommitments {
		data, err := hex.DecodeString(strings.TrimPrefix(data.BlobKZGCommitments[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "failed to parse blob KZG commitment")
		}
		if len(data) != deneb.KZGCommitmentLength {
			return errors.New("incorrect length for blob KZG commitment")
		}
		copy(b.BlobKZGCommitments[i][:], data)
	}
	if data.ExecutionRequests == nil {
		return errors.New("execution requests missing")
	}
	b.ExecutionRequests = data.ExecutionRequests

	return nil
}
