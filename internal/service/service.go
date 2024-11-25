package service

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// ServiceData represents the data for a Kubernetes service
type ServiceData struct {
	Name string
	IP   string
	Port int32
}

// K8sServiceClient is a wrapper around Kubernetes client for interacting with services
type K8sServiceClient struct {
	clientset kubernetes.Interface
	namespace string
}

// NewK8sServiceClient creates a new instance of K8sServiceClient
func NewK8sServiceClient(clientset kubernetes.Interface, namespace string) *K8sServiceClient {
	return &K8sServiceClient{
		clientset: clientset,
		namespace: namespace,
	}
}

// CreateService creates a service in the specified namespace
func (k *K8sServiceClient) CreateService(serviceData ServiceData) (*v1.Service, error) {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceData.Name,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": serviceData.Name},
			Ports: []v1.ServicePort{
				{
					Port: serviceData.Port,
				},
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}

	// Here, you could integrate logic to push information to MongoDB if needed
	// mongoClient, _ := database.CreateMongoClient()
	// err := mongoClient.InsertAPIData("testdb", "testcollection", service)
	// if err != nil {
	//     return nil, err
	// }

	// Retry logic in case of transient failures
	var createdService *v1.Service
	err := retry.OnError(retry.DefaultRetry, isRetryableError, func() error {
		var err error
		createdService, err = k.clientset.CoreV1().Services(k.namespace).Create(context.Background(), service, metav1.CreateOptions{})
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create service after retries: %v", err)
	}

	return createdService, nil
}

// GetService fetches a service by its name in the given namespace
func (k *K8sServiceClient) GetService(serviceName string) (*ServiceData, error) {
	service, err := k.clientset.CoreV1().Services(k.namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	return &ServiceData{
		Name: service.Name,
		IP:   service.Spec.ClusterIP,
		Port: service.Spec.Ports[0].Port,
	}, nil
}

// GetAllServices retrieves all services in the given namespace
func (k *K8sServiceClient) GetAllServices() ([]ServiceData, error) {
	serviceList, err := k.clientset.CoreV1().Services(k.namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	var services []ServiceData
	for _, svc := range serviceList.Items {
		services = append(services, ServiceData{
			Name: svc.Name,
			IP:   svc.Spec.ClusterIP,
			Port: svc.Spec.Ports[0].Port, // Extracting the Port number
		})
	}

	return services, nil
}

// isRetryableError checks if an error is retryable
func isRetryableError(err error) bool {
	// Customize this function based on the error types you want to retry
	// For example, retry on network errors or temporary service failures
	return err != nil
}
