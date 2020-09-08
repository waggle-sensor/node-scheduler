package main

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func LaunchPlugin(deployment *appsv1.Deployment) bool {
	var kubeconfig string
	kubeconfig = "/root/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments("sage-development")

	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return true
}

func TerminatePlugin(plugin string) bool {
	plugin = strings.ToLower(plugin)
	var kubeconfig string
	kubeconfig = "/root/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	deploymentsClient := clientset.AppsV1().Deployments("sage-development")
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), plugin, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("Deleted deployment %s.\n", plugin)
	return true
}

func CreateK3sPod(pluginConfig PluginConfig) *appsv1.Deployment {
	var pod appsv1.Deployment
	defaultNamespace := "sage-development"

	// Build containers
	var containers []apiv1.Container
	for _, plugin := range pluginConfig.Plugins {
		var container apiv1.Container
		container.Name = strings.ToLower(plugin.Name)
		container.Image = plugin.Image
		if len(plugin.Args) > 0 {
			container.Args = plugin.Args
		}
		if len(plugin.Env) > 0 {
			var envs []apiv1.EnvVar
			for k, v := range plugin.Env {
				var env apiv1.EnvVar
				env.Name = k
				env.Value = v
				envs = append(envs, env)
			}
			container.Env = envs
		}
		containers = append(containers, container)
	}

	// Set plugin name and namespace
	pod.ObjectMeta = metav1.ObjectMeta{
		Name:      strings.ToLower(pluginConfig.Name),
		Namespace: defaultNamespace,
	}

	pod.Spec = appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": strings.ToLower(pluginConfig.Name),
			},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": strings.ToLower(pluginConfig.Name),
				},
			},
			Spec: apiv1.PodSpec{
				Containers: containers,
			},
		},
	}
	fmt.Printf("%v", pod)
	return &pod
}

func int32Ptr(i int32) *int32 { return &i }
