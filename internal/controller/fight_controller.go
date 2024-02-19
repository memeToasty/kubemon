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

	kubemon "github.com/memeToasty/kubemon/internal/kubemon"
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
		if client.IgnoreNotFound(err) == nil {
			log.Info("Could not find Fight")
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

	mon1Name := types.NamespacedName{
		Namespace: fightNamespace,
		Name:      fight.Spec.KubeMon1,
	}

	mon2Name := types.NamespacedName{
		Namespace: fightNamespace,
		Name:      fight.Spec.KubeMon2,
	}

	mon1, err := r.getKubeMon(ctx, mon1Name)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			if err := r.updateStatusMessage(ctx, &fight, fmt.Sprintf(FightMessageMonNotFound, mon1Name)); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	mon2, err := r.getKubeMon(ctx, mon2Name)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			if err := r.updateStatusMessage(ctx, &fight, fmt.Sprintf(FightMessageMonNotFound, mon2Name)); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Both KubeMons exist

	// Death logic
	if mon1.IsDead() {
		if err := mon2.LevelUp(); err != nil {
			log.Error(err, "Could not level up KubeMon", "KubeMon", mon2.Name())
			return ctrl.Result{}, err
		}

		if err := r.Delete(ctx, &fight); err != nil {
			log.Error(err, "Could not delete Fight")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if mon2.IsDead() {
		if err := mon1.LevelUp(); err != nil {
			log.Error(err, "Could not level up KubeMon", "KubeMon", mon1.Name())
			return ctrl.Result{}, err
		}

		if err := r.Delete(ctx, &fight); err != nil {
			log.Error(err, "Could not delete Fight")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if fight.Status.NextMon == 1 {
		if err := mon1.Attack(mon2); err != nil {
			log.Error(err, "Could not execute attack", "Attacker", mon1.Name(), "Defender", mon2.Name())
			return ctrl.Result{}, err
		}
	} else {
		if err := mon2.Attack(mon1); err != nil {
			log.Error(err, "Could not execute attack", "Attacker", mon2.Name(), "Defender", mon1.Name())
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
	return ctrl.Result{Requeue: true}, nil
}

func (r *FightReconciler) getKubeMon(ctx context.Context, name types.NamespacedName) (*kubemon.KubeMon, error) {
	apiMon := &kubemonv1.KubeMon{}
	if err := r.Get(ctx, name, apiMon); err != nil {
		return nil, err
	}
	mon, err := kubemon.New(ctx, r.Client, r.Status(), apiMon)
	if err != nil {
		return nil, err
	}
	return mon, nil
}

func (r *FightReconciler) updateStatusMessage(ctx context.Context, fight *kubemonv1.Fight, message string) error {
	fight.Status.LastMessage = message

	if err := r.Status().Update(ctx, fight); err != nil {
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubemonv1.Fight{}).
		Complete(r)
}
