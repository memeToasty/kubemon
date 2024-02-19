package kubemon

import (
	"context"

	kubemonv1 "github.com/memeToasty/kubemon/api/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubeMon struct {
	client       client.Client
	statusClient client.SubResourceWriter
	ctx          context.Context

	apiKubeMon *kubemonv1.KubeMon
}

const (
	KubeMonActionAnnotation = "KubeMon/action"
	KubeMonActionHeal       = "heal"
)

func New(ctx context.Context, c client.Client, sc client.SubResourceWriter, apiKubeMon *kubemonv1.KubeMon) (*KubeMon, error) {
	k := KubeMon{}

	k.client = c
	k.statusClient = sc
	k.ctx = ctx
	k.apiKubeMon = apiKubeMon

	if err := k.init(); err != nil {
		return nil, err
	}

	return &k, nil
}

func (k *KubeMon) init() error {
	if k.apiKubeMon.Status.HP == nil {
		if err := k.SetHealth(10); err != nil {
			return err
		}
	}
	if k.apiKubeMon.Status.Level == nil {

		if err := k.SetLevel(1); err != nil {
			return err
		}
	}
	return nil
}

func (k *KubeMon) Name() string {
	return k.apiKubeMon.Name
}

func (k *KubeMon) GetAction() string {
	return k.apiKubeMon.Annotations[KubeMonActionAnnotation]
}

func (k *KubeMon) ResetAction() error {
	delete(k.apiKubeMon.Annotations, KubeMonActionAnnotation)
	if err := k.update(); err != nil {
		return err
	}
	return nil
}

func (k *KubeMon) IsDead() bool {
	return *k.apiKubeMon.Status.HP == 0
}

func (k *KubeMon) SetHealth(health int32) error {
	k.apiKubeMon.Status.HP = ptr.To(health)
	if err := k.updateStatus(); err != nil {
		return err
	}

	return nil
}

func (k *KubeMon) AddHealth(health int32) error {
	k.apiKubeMon.Status.HP = ptr.To(*k.apiKubeMon.Status.HP + health)
	if err := k.updateStatus(); err != nil {
		return err
	}

	return nil
}

func (k *KubeMon) GetDamage(damage int32) error {
	newHP := *k.apiKubeMon.Status.HP - damage
	if newHP < 0 {
		newHP = 0
	}

	k.apiKubeMon.Status.HP = ptr.To(newHP)
	if err := k.updateStatus(); err != nil {
		return err
	}
	return nil
}

func (k *KubeMon) Attack(k2 *KubeMon) error {
	if err := k2.GetDamage(k.apiKubeMon.Spec.Strength); err != nil {
		return err
	}
	return nil
}

func (k *KubeMon) SetLevel(level int32) error {
	k.apiKubeMon.Status.Level = ptr.To(level)
	if err := k.updateStatus(); err != nil {
		return err
	}

	return nil
}

func (k *KubeMon) LevelUp() error {
	k.apiKubeMon.Status.Level = ptr.To(int32(*k.apiKubeMon.Status.Level + 1))
	if err := k.updateStatus(); err != nil {
		return err
	}
	return nil
}

func (k *KubeMon) updateStatus() error {
	if err := k.statusClient.Update(k.ctx, k.apiKubeMon); err != nil {
		return err
	}
	return nil
}

func (k *KubeMon) update() error {
	if err := k.client.Update(k.ctx, k.apiKubeMon); err != nil {
		return err
	}
	return nil
}
