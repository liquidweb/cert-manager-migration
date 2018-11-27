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

	certClient, err := newCertClient(kubeConfig)
	if err != nil {
		log.Fatalf("Failed to create certificate client: %v", err)
	}

	migrateCerts(certClient, conf)
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

func migrateCerts(certClient *rest.RESTClient, conf Conf) {
	certs, err := getCertificates(certClient, conf.Kube.SourceNamespace)
	if err != nil {
		log.Fatalf("Error while retrieving certificate: %v.", err)
	}

	for _, cert := range certs {
		log.Infof("Copying Certificate %s from %s to %s", cert.Metadata.Name, conf.Kube.SourceNamespace, conf.Kube.DestNamespace)

		cert.Metadata.Namespace = conf.Kube.DestNamespace
		cert.Metadata.ResourceVersion = ""

		_, err := createCertificate(certClient, conf.Kube.DestNamespace, &cert)
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