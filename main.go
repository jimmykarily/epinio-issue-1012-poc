package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jimmykarily/crossplane-marketplace/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	config, err := KubeConfig()
	handleError(err)

	ctx := context.Background()

	// clientset, err := kubernetes.NewForConfig(config)
	// handleError(err)
	// podList, err := clientset.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
	// handleError(err)

	client, err := XDRClient(config)
	handleError(err)

	list, err := client.List(ctx, metav1.ListOptions{})
	handleError(err)

	for _, xrd := range list.Items {
		kind, found, err := unstructured.NestedString(xrd.Object, "spec", "claimNames", "kind")
		handleError(err)
		if found {
			fmt.Println(kind)
		}
	}
}

func XDRClient(config *rest.Config) (dynamic.NamespaceableResourceInterface, error) {
	cs, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.crossplane.io",
		Version:  "v1",
		Resource: "compositeresourcedefinitions",
	}

	return cs.Resource(gvr), nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func KubeConfig() (*rest.Config, error) {
	restConfig, err := config.NewGetter().Get(viper.GetString("kubeconfig"))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't fetch kubeconfig; ensure kubeconfig is present to continue")
	}
	if err := config.NewChecker().Check(restConfig); err != nil {
		return nil, errors.Wrap(err, "couldn't check kubeconfig; ensure kubeconfig is correct to continue")
	}
	return restConfig, nil
}
