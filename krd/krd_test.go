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

package krd

import (
	"testing"

	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mockClient struct {
	create func() (*appsV1.Deployment, error)
	list   func() (*appsV1.DeploymentList, error)
	delete func() error
}

func (c *mockClient) Create(deployment *appsV1.Deployment) (*appsV1.Deployment, error) {
	if c.create != nil {
		return c.create()
	}
	return nil, nil
}

func (c *mockClient) List(opts metaV1.ListOptions) (*appsV1.DeploymentList, error) {
	if c.create != nil {
		return c.list()
	}
	return nil, nil
}

func (c *mockClient) Delete(name string, options *metaV1.DeleteOptions) error {
	if c.delete() != nil {
		return c.delete()
	}
	return nil
}

func TestClientCreateMethod(t *testing.T) {
	t.Run("Succesful deployment creation", func(t *testing.T) {
		expected := "sise-deploy"
		input := &appsV1.Deployment{
			ObjectMeta: metaV1.ObjectMeta{
				Name: expected,
			},
		}
		GetKubeClient = func(configPath string) (ClientDeploymentInterface, error) {
			return &mockClient{
				create: func() (*appsV1.Deployment, error) {
					return input, nil
				},
			}, nil
		}
		client, err := NewClient("")
		if err != nil {
			t.Fatalf("TestDeploymentCreation returned an error (%s)", err)
		}
		result, err := client.Create(input)
		if result != expected {
			t.Fatalf("TestDeploymentCreation returned:\n result=%v\n expected=%v", result, expected)
		}
	})
}

func TestClientDeleteMethod(t *testing.T) {
	t.Run("Succesful deployment deletion", func(t *testing.T) {
		expected := "test-deployment"
		input := &appsV1.Deployment{
			ObjectMeta: metaV1.ObjectMeta{
				Name: expected,
			},
		}
		GetKubeClient = func(configPath string) (ClientDeploymentInterface, error) {
			return &mockClient{
				create: func() (*appsV1.Deployment, error) {
					return input, nil
				},
				delete: func() error {
					return nil
				},
			}, nil
		}
		client, err := NewClient("")
		if err != nil {
			t.Fatalf("TestDeploymentDeletion returned an error (%s)", err)
		}
		result, err := client.Create(input)
		if result != expected {
			t.Fatalf("TestDeploymentDeletion returned:\n result=%v\n expected=%v", result, expected)
		}

		deleteOpts := &metaV1.DeleteOptions{}
		err = client.Delete("", deleteOpts)
		if err != nil {
			t.Fatalf("TestDeploymentDeletion returned an error (%s)", err)
		}
	})
}
