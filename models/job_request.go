//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

package models

type JobRequest struct {
	Id         int                    `json:"id"`
	Delay      int                    `json:"delay"`
	Queue      string                 `json:"queue"`
	Class      string                 `json:"class"`
	Method     string                 `json:"method"`
	RetryCount int                    `json:"retry_count"`
	State      map[string]interface{} `json:"state,omitempty"`
}
