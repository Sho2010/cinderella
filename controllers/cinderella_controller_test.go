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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cinderellav1alpha1 "github.com/Sho2010/cinderella/api/v1alpha1"
)

var _ = Describe("Cinderella controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		CinderellaName         = "cinderella-test"
		ClusterRoleBindingName = "cinderella:test"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When Cinderella Resource", func() {
		It("Should create ", func() {
			By("By creating a new Cinderella")
			ctx := context.Background()

			c := &cinderellav1alpha1.Cinderella{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "cinderella.sho2010.dev/v1alpha1",
					Kind:       "Cinderella",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: CinderellaName,
				},
				Spec: cinderellav1alpha1.CinderellaSpec{
					Roles: []cinderellav1alpha1.Role{
						{
							Kind: "ClusterRole",
							Name: CinderellaName,
						},
					},
					Term: cinderellav1alpha1.Term{
						ExpiresAfter: new(int32),
					},
					Encryption: cinderellav1alpha1.Encryption{
						PublicKey: "test",
					},
				},
				// Status:     cinderellav1alpha1.CinderellaStatus{},
			}

			Expect(k8sClient.Create(ctx, c)).Should(Succeed())
		})
	})
})

// crb := &rbacv1.ClusterRoleBinding{
// 	TypeMeta:   metav1.TypeMeta{},
// 	ObjectMeta: metav1.ObjectMeta{},
// 	Subjects:   []rbacv1.Subject{},
// 	RoleRef:    rbacv1.RoleRef{},
// }
// cronJob := &cronjobv1.CronJob{
// 	TypeMeta: metav1.TypeMeta{
// 		APIVersion: "batch.tutorial.kubebuilder.io/v1",
// 		Kind:       "CronJob",
// 	},
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:      CronjobName,
// 		Namespace: CronjobNamespace,
// 	},
// 	Spec: cronjobv1.CronJobSpec{
// 		Schedule: "1 * * * *",
// 		JobTemplate: batchv1beta1.JobTemplateSpec{
// 			Spec: batchv1.JobSpec{
// 				// For simplicity, we only fill out the required fields.
// 				Template: v1.PodTemplateSpec{
// 					Spec: v1.PodSpec{
// 						// For simplicity, we only fill out the required fields.
// 						Containers: []v1.Container{
// 							{
// 								Name:  "test-container",
// 								Image: "test-image",
// 							},
// 						},
// 						RestartPolicy: v1.RestartPolicyOnFailure,
// 					},
// 				},
// 			},
// 		},
// 	},
// }
