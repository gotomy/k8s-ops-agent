package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreClient "k8s.io/client-go/kubernetes"
	restClient "k8s.io/client-go/rest"
	cmdClient "k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	k8s, err := NewK8s()
	if err != nil {
		panic(err)
	}

	version, err := k8s.GetVersion()

	fmt.Printf("kubernetes version: %s\n", version)

	namespaces ,err := k8s.GetAllNamespaces()

	for _, item := range namespaces {
		fmt.Println(item)
	}
}

type K8s struct {
	Clientset coreClient.Interface
}

func NewK8s() (*K8s, error) {
	client := K8s{}
	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster == true {
		log.Infof("program running inside the cluster, picking the in-cluster configuration")

		config, err := restClient.InClusterConfig()
		if err != nil {
			return nil, err
		}
		client.Clientset, err = coreClient.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		return &client, nil
	}

	log.Info("Program running from outside of the cluster")
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := cmdClient.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	client.Clientset, err = coreClient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h;
	}
	return os.Getenv("USERPROFILE")
}

func (o *K8s) GetVersion() (string, error) {
	version, err := o.Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", version), nil
}

func (o *K8s) GetAllNamespaces() ([]string, error) {
	nss, err := o.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespalces []string
	for _, ns := range nss.Items {
		namespalces = append(namespalces, ns.Name)
	}

	return namespalces, nil
}
