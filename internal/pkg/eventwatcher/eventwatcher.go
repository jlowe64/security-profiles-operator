package eventwatcher

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
)

// GenericEvent is a representation of any Kubernetes event
type Event interface {
	Type() string // Type of event (e.g., "Added", "Updated", "Deleted")
	Object() interface{}
}

// EventHandler is a type alias for an event handling function
type EventHandler func(event Event)

// EventController is our generic controller framework
type EventController struct {
	informerFactory informers.SharedInformerFactory
	eventHandlers   map[string][]EventHandler
}

// NewEventController creates a new EventController
func NewEventController(informerFactory informers.SharedInformerFactory) *EventController {
	return &EventController{
		informerFactory: informerFactory,
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
		return fmt.Errorf("failed to sync")
	}

	for {
		// ... Event dispatch here
	}
}

/*
// --- Example: Pod Event Handling ---

// PodEvent is a typed representation of a Pod-related event
type PodEvent struct {
    Type string
    Pod  *v1.Pod
}

// Type returns the event type
func (pe PodEvent) Type() string {
    return pe.Type
}

// Object returns the underlying object (Pod)
func (pe PodEvent) Object() interface{} {
    return pe.Pod
}

// Create Pod-specific event handler functions:
func podAddHandler(event GenericEvent) {
    podEvent := event.(PodEvent) // Typecast the generic event
    klog.Infof("POD CREATED: %s/%s", podEvent.Pod.Namespace, podEvent.Pod.Name)
}

// (Similar handler functions for podUpdateHandler, podDeleteHandler)

// ---- Usage (Within main function) -----

controller, err := NewEventController(factory)
if err != nil {
    klog.Fatal(err)
}

podInformer := factory.Core().V1().Pods()
podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: func(obj interface{}) {
        pod := obj.(*v1.Pod)
        controller.handleEvent(PodEvent{Type: "Added", Pod: pod})
    },
    // ... (Similar for UpdateFunc, DeleteFunc)
})

controller.RegisterHandler("Added", podAddHandler)
// ... Register other pod handlers

controller.Run(stop)
*/
