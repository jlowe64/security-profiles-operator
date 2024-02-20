package eventwatcher

import {
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
}

// NodeLoggingController is a controller that watches for node events and logs them.
type NodeLoggingController struct {
	informerFactory informers.SharedInformerFactory
	nodeInformer    coreinformers.NodeInformer
}

func (c *NodeLoggingController) Run(stopCh chan struct{}) error {
	c.informerFactory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.nodeInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *NodeLoggingController) nodeDelete(obj interface{}) {
	node := obj.(*v1.Node)
	klog.Infof("Node Deleted: ", node.ObjectMeta.Name)
}

func NewNodeLoggingController(informerFactory informers.SharedInformerFactory) (*NodeLoggingController, error) {
	nodeInformer := informerFactory.Core().V1().Nodes()

	c := &NodeLoggingController{
		informerFactory: informerFactory,
		nodeInformer:     nodeInformer,
	}
	_, err := nodeInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: c.nodeDelete,
		},
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}
