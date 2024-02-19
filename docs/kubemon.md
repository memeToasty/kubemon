# `KubeMon`'s
## What are `KubeMon`'s?
`KubeMon`'s are creatures that have specific characteristics, like strength, level or HP.
After spawning a `KubeMon`, using e.g. [this](../config/samples/kubemon_v1_kubemon1.yaml) manifest, it gets initalized by the game.
It could look something like this then:

```
$ kubectl get kubemon kubemon-sample1

NAME              SPECIES   LEVEL   HP
kubemon-sample1   test      1       10
```

or more detailed, by using e.g. `kubectl get kubemon kubemon-sample1 -oyaml`

```yaml
apiVersion: kubemon.memetoasty.github.com/v1
kind: KubeMon
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"kubemon.memetoasty.github.com/v1","kind":"KubeMon","metadata":{"annotations":{},"labels":{"app.kubernetes.io/created-by":"kubemon","app.kubernetes.io/instance":"kubemon-sample","app.kubernetes.io/managed-by":"kustomize","app.kubernetes.io/name":"kubemon","app.kubernetes.io/part-of":"kubemon"},"name":"kubemon-sample1","namespace":"default"},"spec":{"owner":"tobi","species":"test","strength":1}}
  creationTimestamp: "2024-02-19T16:10:21Z"
  generation: 1
  labels:
    app.kubernetes.io/created-by: kubemon
    app.kubernetes.io/instance: kubemon-sample
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: kubemon
    app.kubernetes.io/part-of: kubemon
  name: kubemon-sample1
  namespace: default
  resourceVersion: "33513"
  uid: 24355341-52a6-4f19-be41-5ea2804e32c1
spec:
  owner: tobi
  species: test
  strength: 1
status:
  hp: 10
  level: 1
```

## Healing
You can heal a kubemon, by adding the `KubeMon/action: "heal"` annotation to the `KubeMon` you wish to heal.

## Fighting
For Combat mechanics, please refer to [this](fights.md) document.