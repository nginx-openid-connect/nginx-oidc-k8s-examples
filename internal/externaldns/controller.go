package externaldns

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	conf_v1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/configuration/v1"
	extdns_v1 "github.com/nginxinc/kubernetes-ingress/pkg/apis/externaldns/v1"
	k8s_nginx "github.com/nginxinc/kubernetes-ingress/pkg/client/clientset/versioned"
	listersV1 "github.com/nginxinc/kubernetes-ingress/pkg/client/listers/configuration/v1"
	extdnslisters "github.com/nginxinc/kubernetes-ingress/pkg/client/listers/externaldns/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	k8s_nginx_informers "github.com/nginxinc/kubernetes-ingress/pkg/client/informers/externalversions"
)

const (
	// ControllerName is the name of the externaldns controller.
	ControllerName = "externaldns"
)

// ExtDNSController represents ExternalDNS controller.
type ExtDNSController struct {
	sync          SyncFn
	ctx           context.Context
	mustSync      []cache.InformerSynced
	queue         workqueue.RateLimitingInterface
	recorder      record.EventRecorder
	client        k8s_nginx.Interface
	informerGroup map[string]*namespacedInformer
	resync        time.Duration
}

type namespacedInformer struct {
	vsLister              listersV1.VirtualServerLister
	sharedInformerFactory k8s_nginx_informers.SharedInformerFactory
	extdnslister          extdnslisters.DNSEndpointLister
}

// ExtDNSOpts represents config required for building the External DNS Controller.
type ExtDNSOpts struct {
	context       context.Context
	namespace     []string
	eventRecorder record.EventRecorder
	client        k8s_nginx.Interface
	resyncPeriod  time.Duration
}

// NewController takes external dns config and return a new External DNS Controller.
func NewController(opts *ExtDNSOpts) *ExtDNSController {
	ig := make(map[string]*namespacedInformer)
	c := &ExtDNSController{
		ctx:           opts.context,
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ControllerName),
		informerGroup: ig,
		recorder:      opts.eventRecorder,
		client:        opts.client,
		resync:        opts.resyncPeriod,
	}

	for _, ns := range opts.namespace {
		c.newNamespacedInformer(ns)
	}

	c.sync = SyncFnFor(c.recorder, c.client, c.informerGroup)
	return c
}

func (c *ExtDNSController) newNamespacedInformer(ns string) {
	nsi := namespacedInformer{sharedInformerFactory: k8s_nginx_informers.NewSharedInformerFactoryWithOptions(c.client, c.resync, k8s_nginx_informers.WithNamespace(ns))}
	nsi.vsLister = nsi.sharedInformerFactory.K8s().V1().VirtualServers().Lister()
	nsi.extdnslister = nsi.sharedInformerFactory.Externaldns().V1().DNSEndpoints().Lister()

	nsi.sharedInformerFactory.K8s().V1().VirtualServers().Informer().AddEventHandler(
		&QueuingEventHandler{
			Queue: c.queue,
		},
	)

	nsi.sharedInformerFactory.Externaldns().V1().DNSEndpoints().Informer().AddEventHandler(&BlockingEventHandler{
		WorkFunc: externalDNSHandler(c.queue),
	})

	c.mustSync = append(c.mustSync,
		nsi.sharedInformerFactory.K8s().V1().VirtualServers().Informer().HasSynced,
		nsi.sharedInformerFactory.Externaldns().V1().DNSEndpoints().Informer().HasSynced,
	)
	c.informerGroup[ns] = &nsi
}

// Run sets up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *ExtDNSController) Run(stopCh <-chan struct{}) {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	glog.Infof("Starting external-dns control loop")

	for _, ig := range c.informerGroup {
		go ig.sharedInformerFactory.Start(c.ctx.Done())
	}

	// wait for all informer caches to be synced
	glog.V(3).Infof("Waiting for %d caches to sync", len(c.mustSync))
	if !cache.WaitForNamedCacheSync(ControllerName, stopCh, c.mustSync...) {
		glog.Fatal("error syncing extDNS queue")
	}

	glog.V(3).Infof("Queue is %v", c.queue.Len())

	go c.runWorker(ctx)

	<-stopCh
	glog.V(3).Infof("shutting down queue as workqueue signaled shutdown")
	c.queue.ShutDown()
}

// runWorker is a long-running function that will continually call the processItem
// function in order to read and process a message on the workqueue.
func (c *ExtDNSController) runWorker(ctx context.Context) {
	glog.V(3).Infof("processing items on the workqueue")
	for {
		obj, shutdown := c.queue.Get()
		if shutdown {
			break
		}

		func() {
			defer c.queue.Done(obj)
			key, ok := obj.(string)
			if !ok {
				return
			}

			if err := c.processItem(ctx, key); err != nil {
				glog.V(3).Infof("Re-queuing item due to error processing: %v", err)
				c.queue.Add(obj)
				return
			}
			glog.V(3).Infof("finished processing work item")
		}()
	}
}

func (c *ExtDNSController) processItem(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return err
	}
	var vs *conf_v1.VirtualServer
	nsi := getNamespacedInformer(namespace, c.informerGroup)
	vs, err = nsi.vsLister.VirtualServers(namespace).Get(name)

	if err != nil {
		return err
	}
	glog.V(3).Infof("processing virtual server resource")
	return c.sync(ctx, vs)
}

func externalDNSHandler(queue workqueue.RateLimitingInterface) func(obj interface{}) {
	return func(obj interface{}) {
		ep, ok := obj.(*extdns_v1.DNSEndpoint)
		if !ok {
			runtime.HandleError(fmt.Errorf("not a DNSEndpoint object: %#v", obj))
			return
		}

		ref := metav1.GetControllerOf(ep)
		if ref == nil {
			// No controller should care about orphans being deleted or
			// updated.
			return
		}

		// We don't check the apiVersion
		// because there is no chance that another object called "VirtualServer" be
		// the controller of a DNSEndpoint.
		if ref.Kind != "VirtualServer" {
			return
		}

		queue.Add(ep.Namespace + "/" + ref.Name)
	}
}

// BuildOpts builds the externalDNS controller options
func BuildOpts(
	ctx context.Context,
	namespace []string,
	recorder record.EventRecorder,
	k8sNginxClient k8s_nginx.Interface,
	resync time.Duration,
) *ExtDNSOpts {
	return &ExtDNSOpts{
		context:       ctx,
		namespace:     namespace,
		eventRecorder: recorder,
		client:        k8sNginxClient,
		resyncPeriod:  resync,
	}
}

func getNamespacedInformer(ns string, ig map[string]*namespacedInformer) *namespacedInformer {
	var nsi *namespacedInformer
	var isGlobalNs bool
	var exists bool

	nsi, isGlobalNs = ig[""]

	if !isGlobalNs {
		// get the correct namespaced informers
		nsi, exists = ig[ns]
		if !exists {
			// we are not watching this namespace
			return nil
		}
	}
	return nsi
}
