package service

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// CreateK8sClientset creates a Kubernetes clientset from the kubeconfig file.
// TODO: mock testing for the CreateK8sClientset
// "k8s.io/client-go/kubernetes/fake"
func CreateK8sClientset() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	if value, exists := os.LookupEnv("KUBECONFIG"); exists {
		kubeconfig = value
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
