package main

import (
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//kutilCoreV1 "github.com/appscode/kutil/core/v1"
	kutilappsv1 "github.com/appscode/kutil/apps/v1"
	. "k8s.io/apimachinery/pkg/apis/meta/v1"
)
func main(){
	kubeFlag := flag.String("kubeconfig",filepath.Join(homedir.HomeDir(),".kube","config"),"Path to kubeconfig")
	flag.Parse()
	config , err := clientcmd.BuildConfigFromFlags("",*kubeFlag)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
		panic(err)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: ObjectMeta{
			Name: "book-server",
			Namespace:NamespaceDefault,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &LabelSelector{
				MatchLabels: map[string]string{
					"app": "book-server",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: ObjectMeta{
					Labels: map[string]string{
						"app": "book-server",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "book-server",
							Image: "rezoan/api_server:1.0.1",
							Args:  []string{"-v", "-b"},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
	// Create Deployment
	fmt.Println("Creating deployment...")
	//result, err := deploymentsClient.Create(deployment)
	res, verb, err := kutilappsv1.CreateOrPatchDeployment(clientSet,deployment.ObjectMeta, func(in *appsv1.Deployment) *appsv1.Deployment {
		in.Spec = deployment.Spec
		return in
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Deployment Created")
	fmt.Println("*  ", res)

	// Patch deployment
	res, verb, err = kutilappsv1.PatchDeployment(clientSet,deployment, func(in *appsv1.Deployment) *appsv1.Deployment {
		in.Spec.Template.Labels = map[string]string{
			"app": "book-server",
			"update": "v1",
		}
		return  in
	})
	if err != nil {
		fmt.Println("Verb = ",verb)
		panic(err)
	}
	fmt.Println("Deployment Updated")
	fmt.Println("*  ", res)

}
func int32Ptr(i int32) *int32 { return &i }