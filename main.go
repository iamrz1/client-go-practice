package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	kutilcorev1 "github.com/appscode/kutil/core/v1"
	kutilappsv1 "github.com/appscode/kutil/apps/v1"
	. "k8s.io/apimachinery/pkg/apis/meta/v1"
	ct "github.com/iamrz1/client-go-practice-ho/pkg/client/clientset/versioned"
	//_ "k8s.io/code-generator"
	crontab "github.com/iamrz1/client-go-practice-ho/pkg/apis/examplecrd.com/v1"
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
	cs, err := ct.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	cron := &crontab.CronTab{
		ObjectMeta: ObjectMeta{
			Name:"my-cron-tab",
			Namespace:NamespaceDefault,
		},
		Spec: crontab.CronTabSpec{
			Replicas: 2,
			Template:crontab.CronTabPodTemplate{
				ObjectMeta: ObjectMeta{
					Name:"cron-pod",
					Namespace:NamespaceDefault,
				},
				Spec: corev1.PodSpec{


				},
			},
		},
	}
	fmt.Println("Creating cronTab")
	newct, err := cs.ExamplecrdV1().CronTabs(NamespaceDefault).Create(cron)
	fmt.Println("cronTab created")
	fmt.Println("cronTab = ",newct)

	fmt.Println("Deleting cronTab")
	err = cs.ExamplecrdV1().CronTabs(NamespaceDefault).Delete(newct.Name,NewDeleteOptions(0))
	if err != nil {
		panic(err)
	}
	fmt.Println("cronTab Deleted")


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
	fmt.Println("Verb = ",verb)

	// Patch deployment
	fmt.Println("Patching deployment...")
	res, verb, err = kutilappsv1.PatchDeployment(clientSet,deployment, func(in *appsv1.Deployment) *appsv1.Deployment {
		in.Spec.Template.Labels = map[string]string{
			"app": "book-server",
			"update": "v1",
		}
		return  in
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Deployment Patched")
	fmt.Println("*  ", res)
	fmt.Println("Verb = ",verb)
	//Delete Deployment
	err = kutilappsv1.DeleteDeployment(clientSet,deployment.ObjectMeta)
	if err != nil {
		fmt.Println("Deployment Delete failed")
		panic(err)
	}


	service := &corev1.Service{
		ObjectMeta: ObjectMeta{
			Name: "book-server",
			Namespace:NamespaceDefault,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": "book-server",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "book-server",
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}

	// Create Service
	fmt.Println("Creating Service...")
	serviceOutput, verb, err := kutilcorev1.CreateOrPatchService(clientSet,service.ObjectMeta, func(in *corev1.Service) *corev1.Service {
		in.Spec = service.Spec
		return in
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Service Created")
	fmt.Println("*  ", serviceOutput)
	fmt.Println("Verb = ",verb)

	// Patch Service
	fmt.Println("Patching service...")
	serviceOutput, verb, err = kutilcorev1.PatchService(clientSet,service, func(in *corev1.Service) *corev1.Service {
		in.Spec.Ports =  []corev1.ServicePort{
			{
				Name:       "book-server",
				Port:       81,
				TargetPort: intstr.FromInt(8080),
				NodePort: 30078,
			},
		}
		return in
	})
	if err != nil {
		fmt.Println("Verb = ",verb)
		panic(err)
	}
	fmt.Println("Service Patched")
	fmt.Println("*  ", res)

	//Delete Service
	//
	err = clientSet.CoreV1().Services(corev1.NamespaceDefault).Delete("book-server",NewDeleteOptions(0))

}

func int32Ptr(i int32) *int32 { return &i }