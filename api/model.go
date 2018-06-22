/*
Copyright 2018 Intel Corporation.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

// CreateVnfRequest contains the VNF creation request parameters
type CreateVnfRequest struct {
	CsarID      string        `json:"csar_id"`
	CsarURL     string        `json:"csar_url"`
	OOFParams   OOFParameters `json:"oof_parameters"`
	ID          string        `json:"vnfdId"`
	Name        string        `json:"vnfInstanceName"`
	Description string        `json:"vnfInstanceDescription"`
}

// CreateVnfResponse contains the VNF creation response parameters
type CreateVnfResponse struct {
	DeploymentID string `json:"deployment_id"`
	Name         string `json:"name"`
}

// OOFParameters contains additional information required for the VNF instance
type OOFParameters struct {
	KeyValues map[string]string `json:"key_values"`
}

// UpdateVnfRequest contains the VNF creation parameters
type UpdateVnfRequest struct {
	CsarID      string        `json:"csar_id"`
	CsarURL     string        `json:"csar_url"`
	OOFParams   OOFParameters `json:"oof_parameters"`
	ID          string        `json:"vnfdId"`
	Name        string        `json:"vnfInstanceName"`
	Description string        `json:"vnfInstanceDescription"`
}

// UpdateVnfResponse contains the VNF update response parameters
type UpdateVnfResponse struct {
	DeploymentID string `json:"deployment_id"`
}

// GeneralResponse is a generic response
type GeneralResponse struct {
	Response string `json:"response"`
}