package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var (
		versionFlag = flag.String("version", "", "Version of Nginx to deploy")
		scaleFlag   = flag.Int("scale", 1, "Number of replicas to scale to")
		kubeConfig  = flag.String("kubeconfig", "", "Path to kubeconfig file")
		namespace   = flag.String("namespace", apiv1.NamespaceDefault, "Kubernetes namespace to deploy into")
	)

	flag.Parse()

	// Validate kubeconfig file path
	if *kubeConfig == "" {
		log.Fatal("Error: kubeconfig cannot be empty")
	}

	// Validate and get replicas count
	replicas, err := getReplicas(scaleFlag)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Prepare nginx version
	nginxImage, err := prepareNginxVersion(versionFlag)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Load kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	deploymentsClient := clientSet.AppsV1().Deployments(*namespace)

	// Call deployNginx with correct arguments
	if err := deployNginx(deploymentsClient, replicas, nginxImage, *namespace); err != nil {
		log.Fatalf("Error deploying Nginx: %v", err)
	}
	fmt.Printf("Nginx deployed successfully in namespace %s.\n", *namespace)
}

// getReplicas validates scale flag and returns the expected pod count for the deployment set
func getReplicas(scaleFlag *int) (int32, error) {
	replicas := int32(*scaleFlag)
	if replicas <= 0 {
		return 0, fmt.Errorf("scale must be greater than zero")
	}
	return replicas, nil
}

// prepareNginxVersion prepares nginx version based on version flag and returns error if any
func prepareNginxVersion(versionFlag *string) (string, error) {
	if *versionFlag == "" {
		return "", fmt.Errorf("version is required")
	}
	return fmt.Sprintf("nginx:%s", *versionFlag), nil
}

// deployNginx creates the deployment for nginx
func deployNginx(deploymentsClient v1.DeploymentInterface, replicas int32, image, namespace string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}
	fmt.Printf("Created deployment %q in namespace %q.\n", result.GetObjectMeta().GetName(), namespace)
	return nil
}
