package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	schedulerName = "my-scheduler"
)

// Scheduler represents our custom scheduler
type Scheduler struct {
	clientset  *kubernetes.Clientset
	podQueue   chan *v1.Pod
	nodeCache  map[string]*v1.Node
	stopCh     chan struct{}
}

// NewScheduler creates a new custom scheduler
func NewScheduler(clientset *kubernetes.Clientset) *Scheduler {
	return &Scheduler{
		clientset: clientset,
		podQueue:  make(chan *v1.Pod, 100),
		nodeCache: make(map[string]*v1.Node),
		stopCh:    make(chan struct{}),
	}
}

// Run starts the scheduler
func (s *Scheduler) Run() {
	log.Printf("Starting custom scheduler: %s", schedulerName)

	// Start watching for pods
	go s.watchPods()

	// Start watching for nodes
	go s.watchNodes()

	// Start scheduling loop
	go s.scheduleLoop()

	<-s.stopCh
}

// watchPods watches for unscheduled pods with our scheduler name
func (s *Scheduler) watchPods() {
	// Watch for pods assigned to our scheduler that are unscheduled
	watchlist := cache.NewListWatchFromClient(
		s.clientset.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				if pod.Spec.SchedulerName == schedulerName && pod.Spec.NodeName == "" {
					log.Printf("New pod to schedule: %s/%s", pod.Namespace, pod.Name)
					s.podQueue <- pod
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod := newObj.(*v1.Pod)
				if pod.Spec.SchedulerName == schedulerName && pod.Spec.NodeName == "" {
					// Only enqueue if pod was just assigned to our scheduler
					oldPod := oldObj.(*v1.Pod)
					if oldPod.Spec.SchedulerName != schedulerName {
						log.Printf("Pod updated to use our scheduler: %s/%s", pod.Namespace, pod.Name)
						s.podQueue <- pod
					}
				}
			},
		},
	)

	controller.Run(s.stopCh)
}

// watchNodes watches for node changes
func (s *Scheduler) watchNodes() {
	watchlist := cache.NewListWatchFromClient(
		s.clientset.CoreV1().RESTClient(),
		"nodes",
		v1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Node{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				node := obj.(*v1.Node)
				s.nodeCache[node.Name] = node
				log.Printf("Node added: %s", node.Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				node := newObj.(*v1.Node)
				s.nodeCache[node.Name] = node
			},
			DeleteFunc: func(obj interface{}) {
				node := obj.(*v1.Node)
				delete(s.nodeCache, node.Name)
				log.Printf("Node deleted: %s", node.Name)
			},
		},
	)

	controller.Run(s.stopCh)
}

// scheduleLoop continuously schedules pods from the queue
func (s *Scheduler) scheduleLoop() {
	for {
		select {
		case pod := <-s.podQueue:
			err := s.schedulePod(pod)
			if err != nil {
				log.Printf("Error scheduling pod %s/%s: %v", pod.Namespace, pod.Name, err)
			}
		case <-s.stopCh:
			return
		}
	}
}

// schedulePod schedules a single pod
func (s *Scheduler) schedulePod(pod *v1.Pod) error {
	log.Printf("Attempting to schedule pod: %s/%s", pod.Namespace, pod.Name)

	// Get the latest pod state
	pod, err := s.clientset.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod: %w", err)
	}

	// Check if already scheduled
	if pod.Spec.NodeName != "" {
		log.Printf("Pod %s/%s already scheduled to %s", pod.Namespace, pod.Name, pod.Spec.NodeName)
		return nil
	}

	// Filter nodes
	suitableNodes := s.filterNodes(pod)
	if len(suitableNodes) == 0 {
		return fmt.Errorf("no suitable nodes found for pod %s/%s", pod.Namespace, pod.Name)
	}

	// Score nodes and select the best one
	bestNode := s.selectBestNode(pod, suitableNodes)
	if bestNode == nil {
		return fmt.Errorf("failed to select best node for pod %s/%s", pod.Namespace, pod.Name)
	}

	// Bind pod to node
	err = s.bindPodToNode(pod, bestNode)
	if err != nil {
		return fmt.Errorf("failed to bind pod to node: %w", err)
	}

	log.Printf("✅ Successfully scheduled pod %s/%s to node %s", pod.Namespace, pod.Name, bestNode.Name)
	return nil
}

// filterNodes filters nodes that can run the pod
func (s *Scheduler) filterNodes(pod *v1.Pod) []*v1.Node {
	var suitableNodes []*v1.Node

	for _, node := range s.nodeCache {
		// Check if node is ready
		if !s.isNodeReady(node) {
			log.Printf("Node %s is not ready, skipping", node.Name)
			continue
		}

		// Check if node is schedulable
		if node.Spec.Unschedulable {
			log.Printf("Node %s is unschedulable, skipping", node.Name)
			continue
		}

		// Check resource requirements
		if !s.hasEnoughResources(node, pod) {
			log.Printf("Node %s doesn't have enough resources, skipping", node.Name)
			continue
		}

		// Check node selector
		if !s.matchesNodeSelector(node, pod) {
			log.Printf("Node %s doesn't match node selector, skipping", node.Name)
			continue
		}

		suitableNodes = append(suitableNodes, node)
	}

	log.Printf("Found %d suitable nodes for pod %s/%s", len(suitableNodes), pod.Namespace, pod.Name)
	return suitableNodes
}

// isNodeReady checks if a node is in Ready state
func (s *Scheduler) isNodeReady(node *v1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return condition.Status == v1.ConditionTrue
		}
	}
	return false
}

// hasEnoughResources checks if node has enough resources for the pod
func (s *Scheduler) hasEnoughResources(node *v1.Node, pod *v1.Pod) bool {
	// Get node allocatable resources
	allocatable := node.Status.Allocatable

	// Calculate total requested resources
	requestedCPU := int64(0)
	requestedMemory := int64(0)

	for _, container := range pod.Spec.Containers {
		if cpu, ok := container.Resources.Requests[v1.ResourceCPU]; ok {
			requestedCPU += cpu.MilliValue()
		}
		if memory, ok := container.Resources.Requests[v1.ResourceMemory]; ok {
			requestedMemory += memory.Value()
		}
	}

	// Check if node has enough resources
	availableCPU := allocatable.Cpu().MilliValue()
	availableMemory := allocatable.Memory().Value()

	return availableCPU >= requestedCPU && availableMemory >= requestedMemory
}

// matchesNodeSelector checks if node matches pod's node selector
func (s *Scheduler) matchesNodeSelector(node *v1.Node, pod *v1.Pod) bool {
	if pod.Spec.NodeSelector == nil || len(pod.Spec.NodeSelector) == 0 {
		return true
	}

	for key, value := range pod.Spec.NodeSelector {
		nodeValue, ok := node.Labels[key]
		if !ok || nodeValue != value {
			return false
		}
	}

	return true
}

// selectBestNode scores nodes and selects the best one
func (s *Scheduler) selectBestNode(pod *v1.Pod, nodes []*v1.Node) *v1.Node {
	if len(nodes) == 0 {
		return nil
	}

	// Simple scoring: prefer nodes with more available resources
	bestNode := nodes[0]
	bestScore := s.scoreNode(nodes[0], pod)

	for i := 1; i < len(nodes); i++ {
		score := s.scoreNode(nodes[i], pod)
		if score > bestScore {
			bestScore = score
			bestNode = nodes[i]
		}
	}

	log.Printf("Selected node %s with score %d", bestNode.Name, bestScore)
	return bestNode
}

// scoreNode scores a node based on available resources
func (s *Scheduler) scoreNode(node *v1.Node, pod *v1.Pod) int {
	// Score based on available resources (0-100)
	allocatable := node.Status.Allocatable

	availableCPU := allocatable.Cpu().MilliValue()
	availableMemory := allocatable.Memory().Value()

	// Simple scoring: higher available resources = higher score
	// Add some randomness for load balancing
	score := int(availableCPU/1000 + availableMemory/(1024*1024*1024))
	score += rand.Intn(10) // Random factor for load balancing

	return score
}

// bindPodToNode binds a pod to a node
func (s *Scheduler) bindPodToNode(pod *v1.Pod, node *v1.Node) error {
	binding := &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		Target: v1.ObjectReference{
			Kind: "Node",
			Name: node.Name,
		},
	}

	err := s.clientset.CoreV1().Pods(pod.Namespace).Bind(context.TODO(), binding, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Emit an event
	timestamp := time.Now()
	_, err = s.clientset.CoreV1().Events(pod.Namespace).Create(context.TODO(), &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pod.Name + "-",
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      "Pod",
			Name:      pod.Name,
			Namespace: pod.Namespace,
			UID:       pod.UID,
		},
		Reason:  "Scheduled",
		Message: fmt.Sprintf("Successfully assigned %s/%s to %s", pod.Namespace, pod.Name, node.Name),
		Source: v1.EventSource{
			Component: schedulerName,
		},
		FirstTimestamp:      metav1.NewTime(timestamp),
		LastTimestamp:       metav1.NewTime(timestamp),
		Type:                v1.EventTypeNormal,
		ReportingController: schedulerName,
	}, metav1.CreateOptions{})

	return err
}

// buildConfig builds the Kubernetes client config
func buildConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeconfig = home + "/.kube/config"
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Println("Building Kubernetes client configuration...")
	config, err := buildConfig()
	if err != nil {
		log.Fatalf("Failed to build config: %v", err)
	}

	log.Println("Creating Kubernetes clientset...")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	log.Println("Testing connection to Kubernetes API...")
	_, err = clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{Limit: 1})
	if err != nil {
		log.Fatalf("Failed to connect to Kubernetes API: %v", err)
	}
	log.Println("✅ Successfully connected to Kubernetes API")

	scheduler := NewScheduler(clientset)
	scheduler.Run()
}
