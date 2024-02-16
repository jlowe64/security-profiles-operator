package nodewatcher

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type NodeDeletedCallback func(nodeName string)

type NodeWatcher struct {
	clientset *kubernetes.Clientset
	callback  NodeDeletedCallback
}

func NewNodeWatcher(clientset *kubernetes.Clientset, callback NodeDeletedCallback) *NodeWatcher {
	return &NodeWatcher{
		clientset: clientset,
		callback:  callback,
	}
}

func (w *NodeWatcher) Run(ctx context.Context) error {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(w.clientset, 0, informers.WithNamespace(v1.NamespaceAll))
	nodeInformer := informerFactory.Core().V1().Nodes()

	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			node, ok := obj.(*v1.Node)
			if !ok {
				return
			}
			if w.callback != nil {
				w.callback(node.GetName())
			}
		},
	})

	informerFactory.Start(ctx.Done())
	defer informerFactory.Stop()

	<-ctx.Done()
	return fmt.Errorf("context canceled")
}
