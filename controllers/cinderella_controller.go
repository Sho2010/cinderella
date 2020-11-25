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

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cinderellav1alpha1 "github.com/Sho2010/cinderella/api/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	ManagedLabel    = "cinderella.sho2010.dev/managed-by"
	CinderellaLabel = "cinderella.sho2010.dev/cinderella"
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

	var cinderella cinderellav1alpha1.Cinderella

	log.Info("fetching Cinderella Resource")
	if err := r.Get(ctx, req.NamespacedName, &cinderella); err != nil {
		log.Error(err, "unable to fetch Cinderella")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// update cinderella.status
	// 暫定対応、まだ設定されてないときのみ入れる
	if cinderella.Status.ExpiredAt.IsZero() {
		cinderella.Status.ExpiredAt = metav1.Now()
		if err := r.Status().Update(ctx, &cinderella); err != nil {
			log.Error(err, "cinderella status update failure")
			return ctrl.Result{}, err
		}
	}

	crb, err := create(cinderella)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Create or Update deployment object
	if _, err := ctrl.CreateOrUpdate(ctx, r.Client, crb, func() error {
		return nil
	}); err != nil {

		// error handling of ctrl.CreateOrUpdate
		log.Error(err, "unable to ensure RBAC is correct")
		return ctrl.Result{}, err
	}

	log.Info("fetched item", "cinderella", cinderella)

	return ctrl.Result{}, nil
}

func create(cinderella cinderellav1alpha1.Cinderella) (*rbacv1.ClusterRoleBinding, error) {
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

	labels := map[string]string{
		ManagedLabel:    c.Name,
		CinderellaLabel: "",
	}
	crb.SetLabels(labels)

	return crb, nil
}

func (r *CinderellaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//NOTE: Ownsを書いてるがガベージコレクションされないのを調べる
	return ctrl.NewControllerManagedBy(mgr).
		For(&cinderellav1alpha1.Cinderella{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Complete(r)
}
