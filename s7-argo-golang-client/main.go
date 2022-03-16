package main

import (
	"context"
	"log"

	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	ctx, apiClient := client.NewAPIClient(context.TODO())
	serviceClient := apiClient.NewWorkflowServiceClient()
	wf := wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "s7-golang-client-",
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
	options := &metav1.CreateOptions{}
	createdWorkflow, err := serviceClient.CreateWorkflow(ctx, &workflowpkg.WorkflowCreateRequest{
		Namespace:     "argo",
		Workflow:      &wf,
		CreateOptions: options,
	})

	if err != nil {
		log.Fatalf("Failed to submit workflow: %v", err)
	}

	log.Printf("created workflow: %v", createdWorkflow.ObjectMeta.Name)
}
