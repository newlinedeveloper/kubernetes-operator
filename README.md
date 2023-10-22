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
mkdir cnl-operator
cd cnl-operator
```

Now, create a new Operator project using the Operator SDK:

```bash
operator-sdk init --plugins go/v3 --domain=cloudnloud.com --owner "Veera" --repo=github.com/newlinedeveloper/cnl-operator
```

This command initializes a new Operator project with a domain and repository information. 

## Step 2: Create a New Custom Resource Definition (CRD)

Next, let's create a Custom Resource Definition (CRD) that defines the custom resource your Operator will manage. For this example, let's create a simple "MyApp" CRD:

```bash
operator-sdk create api --group=api --version=v1alpha1 --kind=Cloudnloud --resource=true --controller=true
```

This command generates the necessary code for your `MyApp` CRD, including the API types and controller.

## Step 3: Define the Operator Logic

Edit the file `controllers/myapp_controller.go` to define the logic for your Operator. Here's a simplified example of an Operator that watches for changes to `MyApp` resources and logs a message when they are created:

update sample CRD spec file

```
spec:
  start: 5 
  end: 20 
  replicas: 5
  deployments:
    - name: nginx
      namespace: default

```

Define CRD Spec in go types file

```
// CloudnloudSpec defines the desired state of Cloudnloud
type CloudnloudSpec struct {
	Start       int              `json:"start"`
	End         int              `json:"end"`
	Replicas    int              `json:"replicas"`
	Deployments []NamespacedName `json:"deployments"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
```


```
package controllers

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/newlinedeveloper/cnl-operator/api/v1alpha1"
)

var logger = log.Log.WithName("cnl_controller")

// CloudnloudReconciler reconciles a Cloudnloud object
type CloudnloudReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.cloudnloud.com,resources=cloudnlouds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.cloudnloud.com,resources=cloudnlouds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.cloudnloud.com,resources=cloudnlouds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cloudnloud object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CloudnloudReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	log := logger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

	log.Info("CNL Reconcile called")
	scaler := &apiv1alpha1.Cloudnloud{}

	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Scaler resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed")
		return ctrl.Result{}, err
	}

	startTime := scaler.Spec.Start
	endTime := scaler.Spec.End

	// current time in UTC
	currentHour := time.Now().UTC().Hour()
	log.Info(fmt.Sprintf("current time in hour : %d\n", currentHour))

	if currentHour >= startTime && currentHour <= endTime {

		if err = cnlDeployment(scaler, r, ctx, int32(scaler.Spec.Replicas)); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

func cnlDeployment(cloudnloud *apiv1alpha1.Cloudnloud, r *CloudnloudReconciler, ctx context.Context, replicas int32) error {
	for _, deploy := range cloudnloud.Spec.Deployments {
		dep := &v1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: deploy.Namespace,
			Name:      deploy.Name,
		}, dep)
		if err != nil {
			return err
		}

		if dep.Spec.Replicas != &replicas {
			dep.Spec.Replicas = &replicas
			err := r.Update(ctx, dep)
			if err != nil {
				return err
			}
			err = r.Status().Update(ctx, cloudnloud)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CloudnloudReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Cloudnloud{}).
		Complete(r)
}

```

## Step 4: Create or update the manifest files

```
make manifests
```

## Step 5: Apply custom resource definition
```
kubectl apply -f config/crd/bases/api.cloudnloud.com_cloudnlouds.yaml

kubectl get crd
```

## Step 6: Deploy our application 
```
kubectl create deploy cnl --image=veerasolaiyappan/my-node-app:v2
```


## Step 7: Deploy our custom resource 
```
kubectl apply -f config/samples/api_v1alpha1_cloudnloud.yaml
```

## Step 8: Delete all the resources
```
kubectl delete deploy cnl

kubectl delete -f config/samples/api_v1alpha1_cloudnloud.yaml

kubectl delete -f config/crd/bases/api.cloudnloud.com_cloudnlouds.yaml
```


