package client

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

// get client config in cluster or from current kubeconfig context
func getKubeConfig(inCluster bool) (*rest.Config, error) {
	var config *rest.Config
	var err error
	if inCluster {
		// creates the in-cluster config
		if config, err = rest.InClusterConfig(); err != nil {
			panic(err.Error())
		}
	} else {
		// creates the out-of-cluster config
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}
	return config, nil
}

// get k8s client
func GetClient() (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(false)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}
