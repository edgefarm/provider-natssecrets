package nkey

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func GetSeedFromSecret(namespace string, name string, key string) (string, error) {
	// get secret from kubernertes client
	restConfig, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", err
	}
	data, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	s, ok := data.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s/%s", key, namespace, name)
	}
	return string(s), nil
}
