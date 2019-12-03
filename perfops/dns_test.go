// Copyright 2017 Prospect One https://prospectone.io/. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package perfops

import (
	"context"
	"testing"
)

func TestRemainingCredits(t *testing.T) {

	testCases := map[string]struct {
		expected interface{}
		body     string
	}{
		"numeric": {5, `{"remaining_credits": 5}`},
		"string":  {"unlimited", `{"remaining_credits": "unlimited"}`},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tr := &respondingTransport{resp: dummyResp(201, "GET", tc.body)}
			c, err := newTestClient(tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.DNS.RemainingCredits(ctx)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if got != tc.expected {
				t.Fatalf("expected %v; got %v", tc.expected, got)
			}
		})
	}
}
