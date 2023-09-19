# kubernetes-operator
Kubernetes Operator Basics and Advance

# How to Write Your First Kubernetes Operator

Kubernetes is a powerful container orchestration platform known for its flexibility and extensibility. One of the key features that make Kubernetes so versatile is its ability to be extended through Custom Resource Definitions (CRDs) and Operators. Operators are applications that are built on top of Kubernetes to automate complex, application-specific tasks. In this guide, we will walk you through the process of writing your first Kubernetes Operator.

## Prerequisites

Before we begin, ensure you have the following prerequisites:

1. **Kubernetes Cluster**: You need access to a Kubernetes cluster to deploy and test your Operator.

2. **kubectl**: Make sure you have `kubectl` installed and configured to connect to your Kubernetes cluster.

3. **Operator SDK**: You'll use the Operator SDK to generate boilerplate code for your Operator. Install it using:

   ```bash
   curl -LO https://github.com/operator-framework/operator-sdk/releases/latest/download/operator-sdk_linux_amd64
   chmod +x operator-sdk_linux_amd64
   sudo mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk
   ```

4. **Go Programming Language**: Operators are typically written in Go, so ensure you have Go installed on your development machine.

## Step 1: Create a New Operator Project

First, create a new directory for your Operator project and navigate into it:

```bash
mkdir my-first-operator
cd my-first-operator
```

Now, create a new Operator project using the Operator SDK:

```bash
operator-sdk init --domain=example.com --repo=github.com/example/my-first-operator
```

This command initializes a new Operator project with a domain and repository information. Replace `example.com` and `github.com/example/my-first-operator` with your own domain and repository details.

## Step 2: Create a New Custom Resource Definition (CRD)

Next, let's create a Custom Resource Definition (CRD) that defines the custom resource your Operator will manage. For this example, let's create a simple "MyApp" CRD:

```bash
operator-sdk create api --group=app --version=v1alpha1 --kind=MyApp --resource=true --controller=true
```

This command generates the necessary code for your `MyApp` CRD, including the API types and controller.

## Step 3: Define the Operator Logic

Edit the file `controllers/myapp_controller.go` to define the logic for your Operator. Here's a simplified example of an Operator that watches for changes to `MyApp` resources and logs a message when they are created:

```go
package controllers

import (
	"context"
	"log"

	appv1alpha1 "github.com/example/my-first-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new MyApp Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMyApp{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// ReconcileMyApp reconciles a MyApp object
type ReconcileMyApp struct {
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MyApp object and makes changes based on the state read
// and what is in the MyApp.Spec
func (r *ReconcileMyApp) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Fetch the MyApp instance
	instance := &appv1alpha1.MyApp{}
	err := r.client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		log.Printf("Error getting MyApp: %v", err)
		return reconcile.Result{}, err
	}

	// Log a message
	log.Printf("MyApp created: %s/%s", instance.Namespace, instance.Name)

	return reconcile.Result{}, nil
}

// ...

// SetupWithManager sets up the controller with the Manager.
func (r *ReconcileMyApp) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.MyApp{}).
		Complete(r)
}
```

## Step 4: Build and Deploy the Operator

Build and package your Operator using Docker:

```bash
operator-sdk build my-first-operator
```

Create a new namespace for your Operator:

```bash
kubectl create namespace my-first-operator
```

Deploy your Operator to the cluster:

```bash
operator-sdk run bundle my-first-operator:latest --namespace my-first-operator
```

## Step 5: Create a Custom Resource

Create a YAML file for your custom resource (`myapp_cr.yaml`):

```yaml
apiVersion: app.example.com/v1alpha1
kind: MyApp
metadata:
  name: example-myapp
spec:
  # Define your custom resource fields here
```

Apply the custom resource to your cluster:

```bash
kubectl apply -f myapp_cr.yaml
```

## Step 6: Verify Your Operator

Check the logs of your Operator to ensure it's working:

```bash
kubectl logs -n my-first-operator -l name=my-first-operator
```

You should see log messages indicating that the Operator detected the creation of your `MyApp` resource.

Congratulations! You've written and deployed your first Kubernetes Operator. This example is a basic introduction; in practice, Operators can be much more complex and automate complex tasks related to your application. As you become more familiar with Kubernetes and Operators, you can explore advanced features and functionality to manage your applications more effectively.
