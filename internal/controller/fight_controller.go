/*
Copyright 2024.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubemonv1 "github.com/memeToasty/kubemon/api/v1"
)

// FightReconciler reconciles a Fight object
type FightReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	FightMessageMonNotFound = "Could not find KubeMon %s"
)

//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=fights,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=fights/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=fights/finalizers,verbs=update

func (r *FightReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var fight kubemonv1.Fight
	if err := r.Get(ctx, req.NamespacedName, &fight); err != nil {
		log.Error(err, "Unable to fetch Fight")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fightNamespace := fight.Namespace

	var mon1 kubemonv1.KubeMon
	var mon2 kubemonv1.KubeMon

	mon1Name := types.NamespacedName{
		Namespace: fightNamespace,
		Name:      fight.Spec.KubeMon1,
	}

	mon2Name := types.NamespacedName{
		Namespace: fightNamespace,
		Name:      fight.Spec.KubeMon2,
	}

	if err := r.Get(ctx, mon1Name, &mon1); err != nil {
		log.Error(err, "Unable to fetch KubeMon1")

		if client.IgnoreNotFound(err) == nil {
			fight.Status.LastMessage = fmt.Sprintf(FightMessageMonNotFound, mon1Name)

			if err := r.Status().Update(ctx, &fight); err != nil {
				log.Error(err, "Could not update Fight status")
			}
		}
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, mon2Name, &mon2); err != nil {
		log.Error(err, "Unable to fetch KubeMon1")

		if client.IgnoreNotFound(err) == nil {
			fight.Status.LastMessage = fmt.Sprintf(FightMessageMonNotFound, mon2Name)

			if err := r.Status().Update(ctx, &fight); err != nil {
				log.Error(err, "Could not update Fight status")
			}
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubemonv1.Fight{}).
		Complete(r)
}