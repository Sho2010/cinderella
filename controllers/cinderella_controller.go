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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cinderellav1alpha1 "github.com/Sho2010/cinderella/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	return ctrl.Result{}, nil
}

func (r *CinderellaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cinderellav1alpha1.Cinderella{}).
		Complete(r)
}
