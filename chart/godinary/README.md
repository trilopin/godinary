# Godinary Helm Chart

## Prerequisites Details

* `autoscaling/v2beta1` enabled

## Charts Details

This chart will do the following:

* Create an Ingress, Service, Horizontal Pod Autoscaller and Deploy
* Expose Godinary on port 80

## Installing the Chart

To install the chart with release name `my-release`:

```
$ helm install --name my-release .
```

## Uninstall the Chart

To uninstall the chart with release name `my-release`:

```
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the Datadog chart and their default values.

Parameter | Description | Default
--- | --- | ---
`namespace` | Define a kubernetes namespace | `default`
`deploy.replicaCount` | Define deploy's replicas | `1`
`deploy.restartPolicy` | Define restart policy | `1`
`deploy.rollout.strategy` | Define deploy's strategy | `RollingUpdate`
`deploy.rollout.maxUnavailable` | Define max unavailable replicas | `1`
`deploy.image` | The image repository to pull from | `trilopin/godinary`
`deploy.imageVersion` | The image tag to pull | Chart version
`deploy.resources.requests.memory` | Define requestes memory | `1Gi`
`deploy.resources.requests.cpu` | Define requestes cpu | `1`
`deploy.resources.limits.memory` | Define memory limit | `3.5Gi`
`deploy.resources.limits.cpu` | Define cpu limit | `2`
`deploy.probe.path` | Define the path for healtcheck | `/up`
`deploy.probe.scheme` | Define the schema for healtcheck | `HTTP`
`deploy.probe.timeoutSeconds` | Define the response timeout for healtcheck | `1`
`deploy.probe.initialDelaySeconds` | Define the initial delay healtcheck | `5`
`deploy.env` | Add Godinary environment variables | `GOMAXPROCS: 2`, `GODINARY_ALLOW_HOSTS: example.com,`, `GODINARY_FS_BASE: /tmp/`
`deploy.volumes` | Define a persistence cache with kubernetes volumes | `nil`
`service.type` | Define the type of Kubernetes service | `NodePort`
`service.sessionAffinity` | Define the session affinity of Kubernetes service | `None`
`service.internalPort` | Define the internal port of Godinary | `3002`
`service.frontalPort` | Define the exposed port of kubernetes service | `80`
`service.protocol` | Define the protocol of the exposed port of kubernetes service | `TCP`
`ingress.annotations` | Define annotations of the kubernetes ingress | `kubernetes.io/ingress.allow-http: "true"`
`hpa.minReplicas` | Define the min replicas for horizontal pod autoscaler | `1`
`hpa.maxReplicas` | Define the max replicas for horizontal pod autoscaler | `8`
`hpa.resourceName` | Define the kubernetes metric for horizontal pod autoscaler | `memory`
`hpa.resourceAverage` | Define the kubernetes average metric for horizontal pod autoscaler | `memory`

Specify each parameter using the --set key=value[,key=value] argument to helm install. For example,

$ helm install --name my-release \
    --set deploy.imageVersion=GODINARY-TAG
    .


## Image tags

You could find `GODINARY-TAG` on https://hub.docker.com/r/trilopin/godinary/tags/

## Persistence

If you want persistence you need to define  `deploy.volume` section configuration like:

```
volumes:
  godinary-data:
    type: hostPath
    mountPath: /app/cache
    readOnly: false
    values:
      path: /app/cache
      type: Directory
```
