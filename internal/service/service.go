package service

import (
	"context"
	"fmt"

	"example.com/m/internal/database"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceData struct to store the service details
type ServiceData struct {
	Name string
	IP   string //IP address
	Port int32  //Listening Port
}

// CreateService creates a service in the specified namespace
func CreateService(clientset kubernetes.Interface, namespace string, serviceData ServiceData) (*v1.Service, error) {
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

	//Push information into mongo DB
	mongoClient, _ := database.CreateMongoClient()
	err := mongoClient.InsertAPIData("testdb", "testcollection", service)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
}

// GetService fetches a service by its name in a given namespace
func GetService(clientset kubernetes.Interface, namespace string, serviceName string) (*ServiceData, error) {
	// Fetch the service by name from the Kubernetes cluster
	service, err := clientset.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	return &ServiceData{
		Name: service.Name,
		IP:   service.Spec.ClusterIP,
		Port: service.Spec.Ports[0].Port,
	}, nil
}

func GetAllServices(clientset kubernetes.Interface, namespace string) ([]ServiceData, error) {
	serviceList, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
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
