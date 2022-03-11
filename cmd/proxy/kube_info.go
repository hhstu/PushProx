package main

import (
	"context"
	"github.com/go-kit/log/level"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"
	"net"
	"sync"
	"time"
)

var podCidrs sync.Map

func setKubeDiscover() error {

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		level.Error(logger).Log("build kubeconfig err: ", err)
		return err
	}

	kubeclient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		level.Error(logger).Log("get kube-client err：", err.Error())
		return err
	}
	nodes, err := kubeclient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		level.Error(logger).Log("get nodelist err：", err.Error())
		return err
	}

	for _, node := range nodes.Items {
		podCidrs.Store(node.Name, node.Spec.PodCIDR)
		level.Info(logger).Log("nodeName", node.Name, "nodeCidr", node.Spec.PodCIDR)
	}
	informersfc := informers.NewSharedInformerFactory(kubeclient, time.Second*30)
	informersfc.Core().V1().Nodes().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*corev1.Node)
			_, cidr, _ := net.ParseCIDR(node.Spec.PodCIDR)
			podCidrs.Store(node.Name, cidr)
		},
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*corev1.Node)
			oldDepl := old.(*corev1.Node)
			if newDepl.Spec.PodCIDR == oldDepl.Spec.PodCIDR {
				return
			}
			_, cidr, _ := net.ParseCIDR(newDepl.Spec.PodCIDR)
			podCidrs.Store(newDepl.Name, cidr)
		},
		DeleteFunc: func(obj interface{}) {
			node := obj.(*corev1.Node)
			podCidrs.Delete(node.Name)
		},
	})
	stopCh := signals.SetupSignalHandler()
	informersfc.Start(stopCh)
	if ok := cache.WaitForCacheSync(stopCh, func() bool { return true }); !ok {
		level.Error(logger).Log("msg", "failed to wait for caches to sync")
	}
	return nil
}
