package project

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func kubeMigrationMain(conf Conf) {
	if conf.Kube.ConfigFile == "" {
		log.Fatal("No Kube Config file specified.")
	}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", conf.Kube.ConfigFile)
	if err != nil {
		log.Fatal("Failed to load client config: %v", err)
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal("Failed to create kubernetes client: %v", err)
	}

	migrateSecrets(client, conf)
}

func migrateSecrets(client *kubernetes.Clientset, conf Conf) {
	sourceSecrets := client.CoreV1().Secrets(conf.Kube.SourceNamespace)
	destSecrets := client.CoreV1().Secrets(conf.Kube.DestNamespace)

	secrets, err := sourceSecrets.List(metav1.ListOptions{LabelSelector: "creator=kube-cert-manager"})
	if err != nil {
		log.Fatal("Failed to retrieve secrets: %v", err)
		return
	}

	for _, secret := range secrets.Items {
		log.Infof("Copying secret %s from %s to %s", secret.Name, conf.Kube.SourceNamespace, conf.Kube.DestNamespace)

		secret.Namespace = conf.Kube.DestNamespace
		secret.ResourceVersion = ""

		_, err := destSecrets.Create(&secret)
		if err != nil {
			log.Fatalf("Failed to create secret: %v", err)
		}
	}
}