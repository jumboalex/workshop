# Custom Kubernetes Scheduler - `my-scheduler`

A custom Kubernetes scheduler implementation in Go that schedules pods to nodes based on resource availability and node selection criteria.

## Features

✅ **Pod Watching** - Monitors unscheduled pods assigned to `my-scheduler`
✅ **Node Filtering** - Filters nodes based on:
  - Node readiness status
  - Schedulability (not cordoned)
  - Resource availability (CPU/Memory)
  - Node selectors

✅ **Node Scoring** - Ranks nodes based on available resources
✅ **Pod Binding** - Binds pods to selected nodes
✅ **Event Generation** - Creates Kubernetes events for scheduling actions
✅ **In-cluster & Out-of-cluster** - Works both inside and outside the cluster

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  Custom Scheduler                       │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Pod Watcher  │  │ Node Watcher │  │  Scheduler   │  │
│  │              │  │              │  │    Loop      │  │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  │
│         │                 │                 │          │
│         └─────────────────┴─────────────────┘          │
│                           │                            │
└───────────────────────────┼────────────────────────────┘
                            │
                            ▼
                  ┌──────────────────┐
                  │  Kubernetes API  │
                  └──────────────────┘
```

## Scheduling Algorithm

### 1. Filter Phase
```
For each node:
  ✓ Is node ready?
  ✓ Is node schedulable?
  ✓ Does node have enough CPU/Memory?
  ✓ Does node match pod's node selector?
```

### 2. Score Phase
```
Score = (Available CPU in cores) + (Available Memory in GB) + Random(0-10)
```

### 3. Select Phase
```
Select node with highest score
```

### 4. Bind Phase
```
Bind pod to selected node via Kubernetes API
```

## Prerequisites

- Go 1.21+
- Access to a Kubernetes cluster
- Valid kubeconfig file (default: `~/.kube/config`)

## Installation

### 1. Download Dependencies

```bash
cd /home/jumbo/workspace/workshop/practice/simulate/k8s-scheduler
go mod download
```

### 2. Build the Scheduler

```bash
go build -o my-scheduler main.go
```

## Usage

### Option 1: Run Locally (Out-of-Cluster)

```bash
# Make sure your kubeconfig is properly configured
export KUBECONFIG=~/.kube/config

# Run the scheduler
./my-scheduler
```

### Option 2: Deploy to Kubernetes (In-Cluster)

Create a deployment for the scheduler:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-scheduler
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-scheduler
  template:
    metadata:
      labels:
        app: my-scheduler
    spec:
      serviceAccountName: my-scheduler
      containers:
      - name: scheduler
        image: your-registry/my-scheduler:latest
        imagePullPolicy: Always
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-scheduler
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: my-scheduler-binding
subjects:
- kind: ServiceAccount
  name: my-scheduler
  namespace: kube-system
roleRef:
  kind: ClusterRole
  name: system:kube-scheduler
  apiGroup: rbac.authorization.k8s.io
```

## Testing the Scheduler

### 1. Start the Scheduler

```bash
./my-scheduler
```

You should see output like:
```
Building Kubernetes client configuration...
Creating Kubernetes clientset...
Testing connection to Kubernetes API...
✅ Successfully connected to Kubernetes API
Starting custom scheduler: my-scheduler
```

### 2. Create a Test Pod

```bash
kubectl apply -f example-pod.yaml
```

### 3. Watch the Scheduler Logs

You'll see logs like:
```
New pod to schedule: default/test-pod-custom-scheduler
Attempting to schedule pod: default/test-pod-custom-scheduler
Found 3 suitable nodes for pod default/test-pod-custom-scheduler
Selected node worker-node-1 with score 45
✅ Successfully scheduled pod default/test-pod-custom-scheduler to node worker-node-1
```

### 4. Verify Pod Scheduling

```bash
# Check pod status
kubectl get pod test-pod-custom-scheduler

# Check events
kubectl describe pod test-pod-custom-scheduler

# You should see an event like:
# Normal  Scheduled  1s    my-scheduler  Successfully assigned default/test-pod-custom-scheduler to worker-node-1
```

## Example Pods

### Basic Pod
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  schedulerName: my-scheduler
  containers:
  - name: app
    image: nginx
```

### Pod with Node Selector
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app-specific-node
spec:
  schedulerName: my-scheduler
  nodeSelector:
    disktype: ssd
  containers:
  - name: app
    image: nginx
```

### Pod with Resource Requests
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app-with-resources
spec:
  schedulerName: my-scheduler
  containers:
  - name: app
    image: nginx
    resources:
      requests:
        memory: "128Mi"
        cpu: "500m"
```

## Configuration

The scheduler can be configured by modifying constants in `main.go`:

```go
const (
    schedulerName = "my-scheduler"  // Change this to use a different name
)
```

## How It Works

1. **Pod Watching**: The scheduler watches for pods with `spec.schedulerName: my-scheduler` and `spec.nodeName: ""`
2. **Node Watching**: Maintains a cache of all available nodes
3. **Filtering**: Filters nodes based on readiness, schedulability, resources, and selectors
4. **Scoring**: Scores remaining nodes based on available resources
5. **Selection**: Selects the highest-scoring node
6. **Binding**: Binds the pod to the selected node via the Kubernetes API
7. **Event**: Creates a Kubernetes event for the scheduling action

## Scheduling Decisions

The scheduler makes decisions based on:

| Factor | Weight | Description |
|--------|--------|-------------|
| Node Ready | Binary | Node must be in Ready state |
| Schedulable | Binary | Node must not be cordoned |
| Resources | Binary | Node must have enough CPU/Memory |
| Node Selector | Binary | Node labels must match pod selector |
| Available Resources | Score | More available = higher score |
| Randomness | Score | Random 0-10 for load balancing |

## Troubleshooting

### Pods Not Being Scheduled

1. **Check scheduler is running**: `ps aux | grep my-scheduler`
2. **Check pod has correct scheduler name**: `kubectl get pod <pod> -o yaml | grep schedulerName`
3. **Check scheduler logs**: Look for error messages
4. **Check node availability**: `kubectl get nodes`

### Permission Errors

Make sure the scheduler has proper RBAC permissions. If running in-cluster, the service account needs:
- Get, List, Watch permissions on Pods and Nodes
- Create permission on Bindings
- Create permission on Events

### No Suitable Nodes

Check:
- Are nodes ready? `kubectl get nodes`
- Do nodes have enough resources? `kubectl describe node <node>`
- Do nodes match the pod's node selector?

## Advanced Features to Add

- **Affinity/Anti-affinity** - Support pod and node affinity rules
- **Taints and Tolerations** - Respect node taints
- **Priority and Preemption** - Support pod priorities
- **Multiple Scheduling Policies** - Different algorithms for different workloads
- **Metrics Integration** - Use Prometheus metrics for smarter decisions
- **Predicates and Priorities** - Plugin architecture for extensibility

## License

MIT

## Contributing

Feel free to submit issues and enhancement requests!
