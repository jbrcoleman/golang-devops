package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func deploy(ctx context.Context, client *kubernetes.Clientset, appFile []byte) (map[string]string, int32, error) {
	var deployment *v1.Deployment

	appFile, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		return nil, 0, fmt.Errorf("readfile error: %s", err)
	}

	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)

	switch obj.(type) {
	case *v1.Deployment:
		deployment = obj.(*v1.Deployment)
	default:
		return nil, 0, fmt.Errorf("Unrecognized type: %s\n", groupVersionKind)
	}

	_, err = client.AppsV1().Deployments("default").Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		deploymentResponse, err := client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return nil, 0, fmt.Errorf("deployment error: %s", err)
		}
		return deploymentResponse.Spec.Template.Labels, 0, nil
	} else if err != nil && !errors.IsNotFound(err) {
		return nil, 0, fmt.Errorf("deployment get error: %s", err)
	}
	deploymentResponse, err := client.AppsV1().Deployments("default").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, 0, fmt.Errorf("deployment error: %s", err)
	}

	return deploymentResponse.Spec.Template.Labels, *deploymentResponse.Spec.Replicas, nil
}

func getClient(inCluster bool) (*kubernetes.Clientset, error) {
	var (
		err    error
		config *rest.Config
	)

	if inCluster {

	} else {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
		if err != nil {
			return nil, err
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getGitHubClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func main() {
	var (
		err error
	)
	ctx := context.Background()
	s := server{
		webhookSecretKey: os.Getenv("WEBHOOK_SECRET"),
	}
	if s.client, err = getClient(false); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	if s.githubClient = getGitHubClient(ctx, os.Getenv("GITHUB_TOKEN")); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	http.HandleFunc("/webhook", s.webhook)
	http.ListenAndServe(":8080", nil)

}
