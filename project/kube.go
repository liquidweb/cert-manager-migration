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

	defaultSecrets := client.CoreV1().Secrets("default")
	russtestSecrets := client.CoreV1().Secrets("russtest")

	secrets, err := defaultSecrets.List(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to retrieve secrets: %v", err)
		return
	}

	for _, secret := range secrets.Items {
		if secret.Type == "kubernetes.io/tls" {
			log.Infof("Copying secret %s from default to russtest", secret.Name)
			secret.Namespace = "russtest"
			secret.ResourceVersion = ""
			_, err := russtestSecrets.Create(&secret)
			if err != nil {
				log.Fatalf("Failed to create secret: %v", err)
			}
		}
	}
}
