package main

import (
	"context"
	"fmt"
	"k8s-client/common/client"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func main() {
	clientset, _ := client.GetClient()
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
