// Copyright (c) 2014 - The Event Horizon authors.
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

package eventhorizon

import (
	"fmt"
	"reflect"
)

// CompareConfig is a config for the ComparEvents function.
type CompareConfig struct {
	ignoreTimestamp bool
	ignoreVersion   bool
}

// CompareOption is an option setter used to configure comparing of events.
type CompareOption func(*CompareConfig)

// IgnoreTimestamp ignores the timestamps of events when comparing.
func IgnoreTimestamp() CompareOption {
	return func(o *CompareConfig) {
		o.ignoreTimestamp = true
	}
}

// IgnoreVersion ignores the versions of events when comparing.
func IgnoreVersion() CompareOption {
	return func(o *CompareConfig) {
		o.ignoreVersion = true
	}
}

// CompareEvents compares two events, with options for ignoring timestamp,
// version etc.
func CompareEvents(e1, e2 Event, options ...CompareOption) error {
	var opts CompareConfig
	for _, o := range options {
		if o == nil {
			continue
		}
		o(&opts)
	}

	if e1.EventType() != e2.EventType() {
		return fmt.Errorf("incorrect event type: %s (should be %s)", e1.EventType(), e2.EventType())
	}
	if !reflect.DeepEqual(e1.Data(), e2.Data()) {
		return fmt.Errorf("incorrect event data: %s (should be %s)", e1.Data(), e2.Data())
	}
	if !opts.ignoreTimestamp {
		if !e1.Timestamp().Equal(e2.Timestamp()) {
			return fmt.Errorf("incorrect timestamp: %s (should be %s)", e1.Timestamp(), e2.Timestamp())
		}
	}
	if e1.AggregateID() != e2.AggregateID() {
		return fmt.Errorf("incorrect aggregate ID: %s (should be %s)", e1.AggregateID(), e2.AggregateID())
	}
	if e1.AggregateType() != e2.AggregateType() {
		return fmt.Errorf("incorrect aggregate type: %s (should be %s)", e1.AggregateType(), e2.AggregateType())
	}
	if !opts.ignoreVersion {
		if e1.Version() != e2.Version() {
			return fmt.Errorf("incorrect aggregate version: %d (should be %d)", e1.Version(), e2.Version())
		}
	}
	if !reflect.DeepEqual(e1.Metadata(), e2.Metadata()) {
		return fmt.Errorf("incorrect event metadata: %s (should be %s)", e1.Metadata(), e2.Metadata())
	}
	return nil
}

// CompareEventSlices compares two slices of events, using options.
func CompareEventSlices(evts1, evts2 []Event, opts ...CompareOption) bool {
	if len(evts1) != len(evts2) {
		return false
	}
	for i, e1 := range evts1 {
		if err := CompareEvents(e1, evts2[i], opts...); err != nil {
			return false
		}
	}
	return true
}
