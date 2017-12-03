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
	"net/http"
)

type (
	// DNSService defines the interface for the DNS API
	DNSService service
)

// RemainingCredits retrieves the ramining credits from the server.
func (s *DNSService) RemainingCredits(ctx context.Context) (int, error) {
	u := s.client.BasePath + "/remaining-credits"
	req, _ := http.NewRequest("GET", u, nil)
	var v *struct {
		Val int `json:"remaining_credits"`
	}
	err := s.client.do(req, &v)
	if err != nil {
		return 0, err
	}
	return v.Val, nil
}
