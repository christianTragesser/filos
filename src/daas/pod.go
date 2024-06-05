package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func runFilosInstance(c Context) {
	var clientset *kubernetes.Clientset

	clientset, err := getKubeClient()
	if err != nil {
		var kubeconfig string
		// service account credentials not found, use local config
		kubeconfig, kubeconfigSet := os.LookupEnv("KUBECONFIG")
		if !kubeconfigSet {
			kubeconfig = filepath.Join(
				os.Getenv("HOME"), ".kube", "config",
			)
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Error("Failed to get cluster config: %v", err)
			return
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error("Failed to create clientset: %v", err)
			return
		}
	}

	// Create a pod definition
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filos-" + c.issueID,
			Namespace: c.namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "filos-" + c.issueID,
					Image: "docker.io/christiantragesser/filos",
					Env: []corev1.EnvVar{
						{
							Name:  "ALERT_ISSUE_ID",
							Value: c.issueID,
						},
						{
							Name:  "ALERT_NAMESPACE",
							Value: c.namespace,
						},
						{
							Name:  "ALERT_RESOURCE_TYPE",
							Value: c.resourceType,
						},
						{
							Name:  "ALERT_APP_NAME",
							Value: c.name,
						},
						{
							Name:  "ALERT_URL",
							Value: c.url,
						},
						{
							Name:  "ALERT_SUMMARY",
							Value: c.summary,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	// Create filos pod instance
	podsClient := clientset.CoreV1().Pods(c.namespace)
	podInstance, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Error("Failed to create pod: %v", err)
		return
	}

	// Wait for the pod to complete
	watch, err := podsClient.Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podInstance.Name),
	})
	if err != nil {
		log.Error("Failed to watch pod: %v", err)
		return
	}

	for event := range watch.ResultChan() {
		p, ok := event.Object.(*corev1.Pod)
		if !ok {
			log.Error("Failed to get pod: %v", err)
			return
		}
		switch event.Type {
		case "ADDED":
			log.Info("Filos pod for issue '" + c.issueID + "' has been created")
		case "MODIFIED":
			switch p.Status.Phase {
			case "Succeeded":
				log.Info("Filos pod for issue '" + c.issueID + "' has successfully completed")
				watch.Stop()
				return
			case "Failed":
				log.Error("Filos pod for issue '" + c.issueID + "' has encountered an error and failed")
				watch.Stop()
				return
			}
		}
	}
}
