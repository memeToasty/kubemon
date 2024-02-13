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
	"math"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
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
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons,verbs=get;list;watch;create;update;patch;delete

func (r *FightReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Reconciling fight!")

	var fight kubemonv1.Fight
	if err := r.Get(ctx, req.NamespacedName, &fight); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "Unable to fetch Fight")
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Got Fight object")

	if fight.DeletionTimestamp != nil {
		log.V(1).Info("Fight is marked for deletion, stop reconciling")
		return ctrl.Result{Requeue: false}, nil
	}
	log.Info("Fight not marked for deletion")

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

	// Both KubeMons exist

	if *mon1.Status.HP == 0 {
		mon2.Status.Level = ptr.To(int32(*mon2.Status.Level + 1))
		if err := r.Status().Update(ctx, &mon2); err != nil {
			log.Error(err, "Could not update status of KubeMon")

			return ctrl.Result{}, err
		}

		if err := r.Delete(ctx, &fight); err != nil {
			log.Error(err, "Could not delete Fight")

			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if *mon2.Status.HP == 0 {
		mon1.Status.Level = ptr.To(int32(*mon1.Status.Level + 1))
		if err := r.Status().Update(ctx, &mon1); err != nil {
			log.Error(err, "Could not update status of KubeMon")

			return ctrl.Result{}, err
		}

		if err := r.Delete(ctx, &fight); err != nil {
			log.Error(err, "Could not delete Fight")

			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if fight.Status.NextMon == 1 {
		// Mon1 attacks
		mon2.Status.HP = ptr.To(int32(math.Max(0, float64(*mon2.Status.HP-mon1.Spec.Strength))))
		if err := r.Status().Update(ctx, &mon2); err != nil {
			log.Error(err, "Could not update status of KubeMon")

			return ctrl.Result{}, err
		}

	} else {
		// Mon2 attacks
		mon1.Status.HP = ptr.To(int32(math.Max(0, float64(*mon1.Status.HP-mon2.Spec.Strength))))
		if err := r.Status().Update(ctx, &mon1); err != nil {
			log.Error(err, "Could not update status of KubeMon")

			return ctrl.Result{}, err
		}
	}

	fight.Status.TurnNumber += 1
	if fight.Status.NextMon == 1 {
		fight.Status.NextMon = 2
	} else {
		fight.Status.NextMon = 1
	}
	if err := r.Status().Update(ctx, &fight); err != nil {
		log.Error(err, "Could not update status of Fight")

		return ctrl.Result{}, err
	}

	log.Info("Got through reconcile! requeuing")
	return ctrl.Result{RequeueAfter: time.Second * 20}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubemonv1.Fight{}).
		Complete(r)
}
