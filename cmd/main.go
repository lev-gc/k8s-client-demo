package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, _ := getClient(false)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		namespace := "sn-apps-autotest"
		ds := "lobby"
		_, err = clientset.AppsV1().Deployments(namespace).Get(context.TODO(), ds, metav1.GetOptions{})

		if errors.IsNotFound(err) {
			fmt.Printf("Deployment %s in namespace %s not found\n", ds, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting Deployment %s in namespace %s: %v\n", ds, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", ds, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}

// get client config in cluster or from current kubeconfig context
func getClient(inCluster bool) (*rest.Config, error) {
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
