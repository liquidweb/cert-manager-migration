package project

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func kubeTest(conf Conf) {
	if conf.KubeConfigFile == "" {
		log.Fatal("No Kube Config file specified.")
	}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", conf.KubeConfigFile)
	if err != nil {
		log.Fatal("Failed to load client config: %v", err)
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal("Failed to create kubernetes client: %v", err)
	}

	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to retrieve pods: %v", err)
		return
	}

	for _, p := range pods.Items {
		log.Infof("Found pods: %s/%s", p.Namespace, p.Name)
	}

	secrets, err := client.CoreV1().Secrets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to retrieve secrets: %v", err)
		return
	}

	for _, secret := range secrets.Items {
		log.Infof("Found secrets: %s/%s", secret.Namespace, secret.Name)
	}
}
