package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	workflow "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func connectToK8s(kubeconfig string) *workflow.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := workflow.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func createArgoWorkflow(clientset *workflow.Clientset, name, imageName string) {
	workflows := clientset.ArgoprojV1alpha1().Workflows("argo")
	wf := &wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "s8-k8s-golang-client-",
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: "helloworld",
			Templates: []wfv1.Template{
				{
					Name: "helloworld",
					Container: &corev1.Container{
						Image: "busybox",
						Command: []string{
							"echo",
						},
						Args: []string{
							"hello world",
						},
					},
				},
			},
		},
	}

	_, err := workflows.Create(context.TODO(), wf, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create argo workflow.")
	}

	log.Printf("Created argo workflow successfully %v", wf.Name)
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	clientset := connectToK8s(*kubeconfig)

	createArgoWorkflow(clientset, "s8-k8s-golang-client", "helloworld")
}
