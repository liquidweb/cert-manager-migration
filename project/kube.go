package project

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

var SchemeGroupVersion = schema.GroupVersion{Group: "stable.liquidweb.com", Version: "v1"}

func kubeMigrationMain(conf Conf) {
	if conf.Kube.SourceConfigFile == "" || conf.Kube.DestConfigFile == "" {
		log.Fatal("No Kube Config file specified.")
	}

	if conf.Kube.SourceNamespace == "" || conf.Kube.DestNamespace == "" {
		log.Fatal("Source and/or Destination Namespaces missing.")
	}

	sourceClient := buildClients(conf.Kube.SourceConfigFile, conf.Kube.SourceNamespace)
	destClient := buildClients(conf.Kube.DestConfigFile, conf.Kube.DestNamespace)

	if sourceClient.KubeConfig.Host == destClient.KubeConfig.Host {
		log.Fatal("Source and Destination Kube hosts must be different.")
	}

	migrateSecrets(sourceClient, destClient)
	migrateCerts(sourceClient, destClient)
}

func buildClients(configFile string, namespace string) KubeClient {
	var client KubeClient
	var err error

	client.KubeConfig, err = clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		log.Fatalf("Failed to load client config: %v", err)
	}

	client.Client, err = kubernetes.NewForConfig(client.KubeConfig)
	if err != nil {
		log.Fatalf("Failed to create kubernetes client: %v", err)
	}

	client.CertClient, err = newCertClient(client.KubeConfig)
	if err != nil {
		log.Fatalf("Failed to create certificate client: %v", err)
	}

	client.Namespace = namespace

	return client
}

func migrateSecrets(sourceClient KubeClient, destClient KubeClient) {
	sourceSecrets := sourceClient.Client.CoreV1().Secrets(sourceClient.Namespace)

	secrets, err := sourceSecrets.List(metav1.ListOptions{LabelSelector: "creator=kube-cert-manager"})
	if err != nil {
		log.Fatal("Failed to retrieve secrets: %v", err)
		return
	}

	for _, secret := range secrets.Items {
		log.Infof("Copying Secret %s:%s/%s to %s:%s/%s", sourceClient.KubeConfig.Host, sourceClient.Namespace, secret.Name, destClient.KubeConfig.Host, destClient.Namespace, secret.Name)

		destSecrets := destClient.Client.CoreV1().Secrets(destClient.Namespace)
		secret.ResourceVersion = ""
		secret.Namespace = destClient.Namespace

		_, err := destSecrets.Create(&secret)
		if err != nil {
			log.Fatalf("Failed to create secret: %v", err)
		}
	}
}

func migrateCerts(sourceClient KubeClient, destClient KubeClient) {
	certs, err := getCertificates(sourceClient.CertClient, sourceClient.Namespace)
	if err != nil {
		log.Fatalf("Error while retrieving certificate: %v.", err)
	}

	for _, cert := range certs {
		log.Infof("Copying Certificate %s:%s/%s to %s:%s/%s", sourceClient.KubeConfig.Host, sourceClient.Namespace, cert.Metadata.Name, destClient.KubeConfig.Host, destClient.Namespace, cert.Metadata.Name)

		cert.Metadata.ResourceVersion = ""
		cert.Metadata.Namespace = destClient.Namespace

		_, err := createCertificate(destClient.CertClient, destClient.Namespace, &cert)
		if err != nil {
			log.Fatalf("Failed to create certificate: %v", err)
		}
	}
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Certificate{},
		&CertificateList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func newCertClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getCertificates(certClient *rest.RESTClient, namespace string) ([]Certificate, error) {
	rl := flowcontrol.NewTokenBucketRateLimiter(0.2, 3)
	for {
		rl.Accept()
		req := certClient.Get().Resource("certificates").Namespace(namespace)

		var certList CertificateList

		err := req.Do().Into(&certList)

		if err != nil {
			log.Printf("Error while retrieving certificate: %v. Retrying", err)
		} else {
			return certList.Items, nil
		}
	}
}

func createCertificate(certClient *rest.RESTClient, namespace string, obj *Certificate) (*Certificate, error) {
	result := &Certificate{}
	err := certClient.Post().
		Namespace(namespace).Resource("certificates").
		Body(obj).Do().Into(result)
	return result, err
}