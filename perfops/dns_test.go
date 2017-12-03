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

func RemainingCredits(t *testing.T) {
	ctx := context.Background()
	const exp = 5
	const body = `{"remaining-credits": 5}`
	tr := &respondingTransport{resp: dummyResp(201, "GET", body)}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	got, err := c.DNS.RemainingCredits(ctx)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
}
