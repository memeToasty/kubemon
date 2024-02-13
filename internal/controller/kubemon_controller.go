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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubemonv1 "github.com/memeToasty/kubemon/api/v1"
)

// KubeMonReconciler reconciles a KubeMon object
type KubeMonReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons/finalizers,verbs=update

func (r *KubeMonReconciler) UpdateKubeMonHealth(ctx context.Context, kubemon *kubemonv1.KubeMon, hp int32) error {
	log := log.FromContext(ctx)

	kubemon.Status.HP = ptr.To(hp)
	if err := r.Status().Update(ctx, kubemon); err != nil {
		log.Error(err, "Could not update status of KubeMon")

		return err
	}
	return nil
}

func (r *KubeMonReconciler) UpdateKubeMonLevel(ctx context.Context, kubemon *kubemonv1.KubeMon, level int32) error {
	log := log.FromContext(ctx)

	kubemon.Status.Level = ptr.To(level)
	if err := r.Status().Update(ctx, kubemon); err != nil {
		log.Error(err, "Could not update status of KubeMon")

		return err
	}
	return nil
}

func (r *KubeMonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var kubemon kubemonv1.KubeMon
	if err := r.Get(ctx, req.NamespacedName, &kubemon); err != nil {
		log.Error(err, "Unable to fetch KubeMon")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if kubemon.DeletionTimestamp != nil {
		log.V(1).Info("KubeMon is marked for deletion, stop reconciling")

		return ctrl.Result{Requeue: false}, nil
	}

	if kubemon.Status.HP == nil {
		log.V(1).Info("KubeMon does not have any health")

		if err := r.UpdateKubeMonHealth(ctx, &kubemon, 10); err != nil {
			return ctrl.Result{}, err
		}
	}

	if kubemon.Status.Level == nil {
		log.V(1).Info("KubeMon does not have any Level")

		if err := r.UpdateKubeMonLevel(ctx, &kubemon, 1); err != nil {
			return ctrl.Result{}, err
		}
	}

	if kubemon.Spec.Action == "heal" {
		if err := r.UpdateKubeMonHealth(ctx, &kubemon, *kubemon.Status.HP+10); err != nil {
			return ctrl.Result{}, err
		}

		kubemon.Spec.Action = ""
		if err := r.Update(ctx, &kubemon); err != nil {
			log.Error(err, "Could not update KubeMon object")

			return ctrl.Result{}, err
		}
	}

	if err := r.Status().Update(ctx, &kubemon); err != nil {
		log.Error(err, "Could not update status of KubeMon")

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeMonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubemonv1.KubeMon{}).
		Complete(r)
}
