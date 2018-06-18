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

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	pkgerrors "github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	"github.com/shank7485/k8-plugin-multicloud/krd"
	"github.com/shank7485/k8-plugin-multicloud/utils"
)

// VNFInstanceService communicates the actions to Kubernetes deployment
type VNFInstanceService struct {
	Client VNFInstanceClientInterface
}

// VNFInstanceClientInterface has methods to work with VNF Instance resources.
type VNFInstanceClientInterface interface {
	Create(deployment *appsV1.Deployment) (string, error)
	List(limit int64) (*appsV1.DeploymentList, error)
	Delete(name string, options *metaV1.DeleteOptions) error
}

// NewVNFInstanceService creates a client that comunicates with a Kuberentes Cluster
func NewVNFInstanceService(kubeConfigPath string) (*VNFInstanceService, error) {
	client, err := GetVNFClient(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return &VNFInstanceService{
		Client: client,
	}, nil
}

// GetVNFClient retrieve the client used to communicate with a Kubernetes Cluster
var GetVNFClient = func(kubeConfigPath string) (VNFInstanceClientInterface, error) {
	var client VNFInstanceClientInterface

	client, err := krd.NewClient(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return client, err
}

// Create is the POST method creates a new VNF instance resource.
func (s *VNFInstanceService) Create(w http.ResponseWriter, r *http.Request) {
	var resource CreateVnfRequest

	if r.Body == nil {
		http.Error(w, "Body empty", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	uuid := uuid.NewUUID()
	// Persist in AAI database.
	log.Println(resource.CsarID + "_" + string(uuid))

	deployment, err := utils.GetDeploymentInfo(resource.CsarURL)
	if err != nil {
		werr := pkgerrors.Wrap(err, "Get Deployment information error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
		return
	}

	name, err := s.Client.Create(deployment)
	if err != nil {
		werr := pkgerrors.Wrap(err, "Create VNF error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
		return
	}

	resp := GeneralResponse{
		Response: "Created Deployment:" + name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		werr := pkgerrors.Wrap(err, "Parsing output of new VNF error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}
}

// List the existing VNF instances created in a given Kubernetes cluster
func (s *VNFInstanceService) List(w http.ResponseWriter, r *http.Request) {
	_, err := s.Client.List(int64(10)) // TODO (electrocucaracha): export this as configuration value
	if err != nil {
		werr := pkgerrors.Wrap(err, "Get VNF list error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
		return
	}

	resp := GeneralResponse{
		Response: "Listing:",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		werr := pkgerrors.Wrap(err, "Parsing output VNF list error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}

}

// Delete existing VNF instances created in a given Kubernetes Cluster
func (s *VNFInstanceService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	deletePolicy := metaV1.DeletePropagationForeground
	deleteOptions := &metaV1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	err := s.Client.Delete(name, deleteOptions)
	if err != nil {
		werr := pkgerrors.Wrap(err, "Delete VNF error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
		return
	}

	resp := GeneralResponse{
		Response: "Deletion complete:" + name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		werr := pkgerrors.Wrap(err, "Parsing output of delete VNF error")
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}
}
