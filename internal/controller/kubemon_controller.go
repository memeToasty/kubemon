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
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubemonv1 "github.com/memeToasty/kubemon/api/v1"
	kubemon "github.com/memeToasty/kubemon/internal/kubemon"
)

// KubeMonReconciler reconciles a KubeMon object
type KubeMonReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	ErrKubeMonGone = errors.New("kubeMon is marked for deletion")
)

const (
	ActionKubeMonHeal = "heal"
)

//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubemon.memetoasty.github.com,resources=kubemons/finalizers,verbs=update

func (r *KubeMonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	kubemon, err := r.getKubeMon(ctx, req.NamespacedName)
	if err != nil {
		if err == ErrKubeMonGone {
			log.Info("KubeMon is marked for deletion, stop reconciling")
			return ctrl.Result{Requeue: false}, nil
		}

		if client.IgnoreNotFound(err) == nil {
			log.Info("Could not find KubeMon")
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Got KubeMon object")

	if kubemon.GetAction() == ActionKubeMonHeal {
		if err := kubemon.AddHealth(10); err != nil {
			return ctrl.Result{}, err
		}

		if err := kubemon.ResetAction(); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *KubeMonReconciler) getKubeMon(ctx context.Context, name types.NamespacedName) (*kubemon.KubeMon, error) {
	apiMon := &kubemonv1.KubeMon{}
	if err := r.Get(ctx, name, apiMon); err != nil {
		return nil, err
	}

	if apiMon.DeletionTimestamp != nil {
		return nil, ErrKubeMonGone
	}
	mon, err := kubemon.New(ctx, r.Client, r.Status(), apiMon)
	if err != nil {
		return nil, err
	}
	return mon, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeMonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubemonv1.KubeMon{}).
		Complete(r)
}
