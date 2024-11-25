package service

import (
	"testing"
)

func TestService(t *testing.T) {
	// // Create a fake Kubernetes clientset
	// clientset := fake.NewSimpleClientset()
	// // Create an instance of K8sServiceClient
	// k8sClient := NewK8sServiceClient(clientset, "default")

	// t.Run("TestCreateService", func(t *testing.T) {
	// 	serviceData := ServiceData{
	// 		Name: "test-service",
	// 		IP:   "10.0.0.1",
	// 		Port: 2332,
	// 	}

	// 	service, err := k8sClient.CreateService(serviceData)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, service)
	// 	assert.Equal(t, "test-service", service.Name)
	// })

	// t.Run("TestGetService", func(t *testing.T) {
	// 	fetchedService, err := k8sClient.GetService("test-service")

	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, fetchedService)
	// 	assert.Equal(t, "test-service", fetchedService.Name)
	// 	assert.Equal(t, int32(2332), fetchedService.Port)
	// })

	// t.Run("TestGetAllServices", func(t *testing.T) {
	// 	service1 := &v1.Service{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name: "service1",
	// 		},
	// 		Spec: v1.ServiceSpec{
	// 			ClusterIP: "10.0.0.2",
	// 			Ports: []v1.ServicePort{
	// 				{
	// 					Port:     8080,
	// 					Protocol: v1.ProtocolTCP,
	// 				},
	// 			},
	// 		},
	// 	}
	// 	service2 := &v1.Service{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name: "service2",
	// 		},
	// 		Spec: v1.ServiceSpec{
	// 			ClusterIP: "10.0.0.3",
	// 			Ports: []v1.ServicePort{
	// 				{
	// 					Port:     9090,
	// 					Protocol: v1.ProtocolTCP,
	// 				},
	// 			},
	// 		},
	// 	}

	// 	// Add the services to the fake clientset
	// 	clientset.CoreV1().Services("default").Create(context.Background(), service1, metav1.CreateOptions{})
	// 	clientset.CoreV1().Services("default").Create(context.Background(), service2, metav1.CreateOptions{})

	// 	services, err := k8sClient.GetAllServices()

	// 	assert.NoError(t, err)
	// 	assert.Len(t, services, 3)
	// 	assert.Equal(t, "service1", services[0].Name)
	// 	assert.Equal(t, "10.0.0.2", services[0].IP)
	// 	assert.Equal(t, int32(8080), services[0].Port)
	// 	assert.Equal(t, "service2", services[1].Name)
	// 	assert.Equal(t, "10.0.0.3", services[1].IP)
	// 	assert.Equal(t, int32(9090), services[1].Port)
	// })

}
