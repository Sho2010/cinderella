/*


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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cinderellav1alpha1 "github.com/Sho2010/cinderella/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

var (
	ownerKey = "metadata.controller"
	apiGVStr = cinderellav1alpha1.GroupVersion.String()
)

// CinderellaReconciler reconciles a Cinderella object
type CinderellaReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=cinderella.sho2010.dev,resources=cinderellas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cinderella.sho2010.dev,resources=cinderellas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *CinderellaReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("cinderella", req.NamespacedName)
	log.Info("start reconciler")

	var c cinderellav1alpha1.Cinderella

	log.Info("fetching Cinderella Resource")
	if err := r.Get(ctx, req.NamespacedName, &c); err != nil {
		log.Error(err, "unable to fetch Cinderella")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// update cinderella.status
	if c.Status.ExpiredAt.IsZero() {
		// TODO Term を使ってちゃんと設定する
		c.Status.ExpiredAt = metav1.NewTime(time.Now().Add(1 * time.Minute))
		tmp := false
		c.Status.Expired = &tmp
	}

	now := metav1.Now()
	expired := c.Status.ExpiredAt.Before(&now)
	c.Status.Expired = &expired

	if err := r.Status().Update(ctx, &c); err != nil {
		log.Error(err, "cinderella status update failure")
		return ctrl.Result{}, err
	}

	crb, err := buildServiceAccount(c)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//TODO: resourceが変更されたときに古いリソースのcleanup処理を行う
	// if err := r.cleanupOwnedResources(ctx, log, &c); err != nil {
	// 	return ctrl.Result{}, err
	// }

	if !*c.Status.Expired {
		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, crb, func() error {
			if err := ctrl.SetControllerReference(&c, crb, r.Scheme); err != nil {
				log.Error(err, "unable to set ownerReference from Cinderella to Role")
				return err
			}
			return nil
		}); err != nil {
			log.Error(err, "unable to ensure RBAC is correct")
			return ctrl.Result{}, err
		}
	}

	if err := r.deleteExpiredResources(ctx, log, &c); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func buildServiceAccount(c cinderellav1alpha1.Cinderella) (*rbacv1.ClusterRoleBinding, error) {

	subject := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      "test",
		Namespace: "test",
	}

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cinderella:test",
		},
		Subjects: []rbacv1.Subject{subject},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "view",
		},
	}

	return crb, nil
}

func (r *CinderellaReconciler) deleteExpiredResources(ctx context.Context, log logr.Logger, c *cinderellav1alpha1.Cinderella) error {
	if !*c.Status.Expired {
		return nil
	}

	log.Info("cinderella expired. search delete targets", "expired-cinderella", c.Name)

	var roleList rbacv1.ClusterRoleBindingList
	if err := r.List(ctx, &roleList, client.MatchingFields(map[string]string{ownerKey: c.Name})); err != nil {
		return err
	}

	for _, role := range roleList.Items {
		if err := r.Delete(ctx, &role); err != nil {
			log.Error(err, "failed to delete resource")
			return err
		}

		log.Info("delete resource: " + role.Name)
		r.Recorder.Eventf(c, corev1.EventTypeNormal, "Deleted", "Deleted resource %q", role.Name)
	}
	return nil
}

func (r *CinderellaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Add OwnerKey index to ClusterRoleBinding object which cinderella resource owns
	if err := mgr.GetFieldIndexer().IndexField(&rbacv1.ClusterRoleBinding{}, ownerKey, func(rawObj runtime.Object) []string {
		crb := rawObj.(*rbacv1.ClusterRoleBinding)
		owner := metav1.GetControllerOf(crb)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "Cinderella" {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&cinderellav1alpha1.Cinderella{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Complete(r)
}
