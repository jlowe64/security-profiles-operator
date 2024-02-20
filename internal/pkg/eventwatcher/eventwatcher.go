package eventwatcher

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type Event interface {
	Type() string        // Type of event (e.g., "Added", "Updated", "Deleted")
	Object() interface{} // The underlying Kubernetes object
}

// EventHandler is a type alias for an event handling function
type EventHandler func(event Event)

// EventController is our generic controller framework
type EventController struct {
	informerFactory informers.SharedInformerFactory
	nodeInformer    coreinformers.NodeInformer
	eventHandlers   map[string][]EventHandler
}

// NewEventController creates a new EventController
func NewEventController(informerFactory informers.SharedInformerFactory) *EventController {
	return &EventController{
		informerFactory: informerFactory,
		nodeInformer:    informerFactory.Core().V1().Nodes(),
		eventHandlers:   make(map[string][]EventHandler),
	}
}

// RegisterHandler registers an event handler for a given event type
func (c *EventController) RegisterHandler(eventType string, handler EventHandler) {
	c.eventHandlers[eventType] = append(c.eventHandlers[eventType], handler)
}

// Run starts the controller's informers and listens for events
func (c *EventController) Run(stopCh chan struct{}) error {
	c.informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, c.informerFactory.WaitForCacheSync()) {
		klog.V(4).Info("Failed to sync")
		return fmt.Errorf("Failed to sync")
	}

	for {
		select {
		case <-stopCh:
			return nil // Exit if the stop signal is received
		default:
			// Iterate over informers (you'll likely have multiple of these)
			for informerType, informer := range c.informers {
				// Check if the informer's cache has synced
				if !informer.HasSynced() {
					klog.V(4).Info("Informer not synced, waiting...", informerType)
					continue
				}

				// Use informer.GetStore().List() or a similar method to get the events
				events := informer.GetStore().List()

				// Dispatch each event
				for _, obj := range events {
					node, ok := obj.(*v1.Node)
					if !ok {
						continue // Skip if not a Node event
					}

					if node.DeletionTimestamp != nil {
						// This is a deleted Node event; proceed with handling
						controller.handleEvent(Event{
							Type:   "Deleted",
							Object: node,
						})
						klog.Infof("Node Deleted: %s", node.Name)
					}
				}
			}
		}
	}
}

// handleEvent is a helper to dispatch events to registered handlers
func (c *EventController) handleEvent(event Event) {
	handlers, ok := c.eventHandlers[event.Type()]
	if !ok {
		return // No handlers registered for this event type
	}

	for _, handler := range handlers {
		handler(event)
	}
}
