/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stabledwkv1 "stable.dwk/api/v1"
)

// DummySiteReconciler reconciles a DummySite object
type DummySiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DummySite object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *DummySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	ds := &stabledwkv1.DummySite{}
	if err := r.Get(ctx, req.NamespacedName, ds); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ds.Name + "-dep",
			Namespace: ds.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": ds.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": ds.Name},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{{
						Name:         "fetch",
						Image:        "curlimages/curl",
						Command:      []string{"sh", "-c", "curl -L " + ds.Spec.WebsiteUrl + " > /data/index.html"},
						VolumeMounts: []corev1.VolumeMount{{Name: "webdata", MountPath: "/data"}},
					}},
					Containers: []corev1.Container{{
						Name:         "nginx",
						Image:        "nginx:alpine",
						VolumeMounts: []corev1.VolumeMount{{Name: "webdata", MountPath: "/usr/share/nginx/html"}},
					}},
					Volumes: []corev1.Volume{{
						Name:         "webdata",
						VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
					}},
				},
			},
		}

		return controllerutil.SetControllerReference(ds, deployment, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ds.Name + "-svc",
			Namespace: ds.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
		svc.Spec = corev1.ServiceSpec{
			Selector: map[string]string{"app": ds.Name},
			Ports: []corev1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt(80),
			}},
		}
		return controllerutil.SetControllerReference(ds, svc, r.Scheme)
	})

	if err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Successfully reconciled DummySite", "name", ds.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummySiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stabledwkv1.DummySite{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named("dummysite").
		Complete(r)
}
