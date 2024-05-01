package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestGetReplicas(t *testing.T) {
	tests := []struct {
		scaleFlag int
		expected  int32
		wantErr   bool
	}{
		{scaleFlag: 1, expected: 1, wantErr: false},
		{scaleFlag: 0, expected: 0, wantErr: true},
		{scaleFlag: -1, expected: 0, wantErr: true},
	}

	for _, tt := range tests {
		replicas, err := getReplicas(&tt.scaleFlag)
		if (err != nil) != tt.wantErr {
			t.Errorf("getReplicas() error = %v, wantErr %v", err, tt.wantErr)
			continue
		}
		if replicas != tt.expected {
			t.Errorf("getReplicas() = %v, want %v", replicas, tt.expected)
		}
	}
}

func TestPrepareNginxVersion(t *testing.T) {
	tests := []struct {
		versionFlag string
		expected    string
		wantErr     bool
	}{
		{versionFlag: "1.13.12", expected: "nginx:1.13.12", wantErr: false},
		{versionFlag: "", expected: "", wantErr: true},
	}

	for _, tt := range tests {
		nginxImage, err := prepareNginxVersion(&tt.versionFlag)
		if (err != nil) != tt.wantErr {
			t.Errorf("prepareNginxVersion() error = %v, wantErr %v", err, tt.wantErr)
			continue
		}
		if nginxImage != tt.expected {
			t.Errorf("prepareNginxVersion() = %v, want %v", nginxImage, tt.expected)
		}
	}
}

func TestDeployNginx(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	deploymentsClient := clientset.AppsV1().Deployments("default")
	replicas := int32(1)
	image := "nginx:1.13.12"
	namespace := "default"

	err := deployNginx(deploymentsClient, replicas, image, namespace)
	assert.NoError(t, err, "deployNginx() error = %v, want nil")

	// Check if the deployment was created
	deployments, _ := deploymentsClient.List(context.Background(), metav1.ListOptions{})
	assert.Equal(t, 1, len(deployments.Items), "Expected 1 deployment, found %d", len(deployments.Items))
	assert.Equal(t, "nginx-deployment", deployments.Items[0].Name, "Expected deployment name %s, got %s", "nginx-deployment", deployments.Items[0].Name)
	assert.Equal(t, image, deployments.Items[0].Spec.Template.Spec.Containers[0].Image, "Expected container image %s, got %s", image, deployments.Items[0].Spec.Template.Spec.Containers[0].Image)
}
