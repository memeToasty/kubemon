# `Fight`s
## What are `Fight`s
Fights are a way to measure `KubeMon`'s with each other.
## Creating a `Fight`
A fight can be created by specifying two `KubeMon`'s which should fight each other.
It could look something like [this](../config/samples/kubemon_v1_fight.yaml):

```yaml
apiVersion: kubemon.memetoasty.github.com/v1
kind: Fight
metadata:
  name: fight-sample
spec:
  kubemon1: kubemon-sample1
  kubemon2: kubemon-sample2

```

> [!NOTE]  
> Currently, both `KubeMon`'s have to reside in the same namespace to be able to fight each other

Each round the `KubeMon` which's turn it is, attacks the opponent. It deals the damage that is specified in its `.spec.strength` field, until one `KubeMon`'s health reaches `0`.

The winning `KubeMon` receives a level.