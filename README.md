k8s-gitlab-runner
====


### Usage

docker image: `registry.cn-shenzhen.aliyuncs.com/juxuny-public/k8s-gitlab-runner:v1.1.1`

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: ${NAMESPACE}
  labels:
    name: ${NAMESPACE}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${RUNNER_NAME}-gitlab-runner
  namespace: ${NAMESPACE}
spec:
  replicas: 1
  selector:
    matchLabels:
      name: ${RUNNER_NAME}-gitlab-runner
  template:
    metadata:
      labels:
        name: ${RUNNER_NAME}-gitlab-runner
    spec:
      serviceAccountName: ${RUNNER_NAME}-gitlab-runner
      tolerations:
        - key: node-role.kubernetes.io/master
          operator: Exists
          effect: NoSchedule
      containers:
        - name: gitlab-runner
          image: registry.cn-shenzhen.aliyuncs.com/juxuny-public/k8s-gitlab-runner:v1.1.1
          imagePullPolicy: Always
          #securityContext:
          #  privileged: true
          env:
            - name: HOST
              value: ${GITLAB_HOST}
            - name: PRIVATE_TOKEN
              value: ${GITLAB_PRIVATE_TOKEN}
            - name: TOKEN
              value: ${REGISTRATION_TOKEN}
            - name: NAME
              value: ${RUNNER_NAME}
            - name: NAMESPACE
              value: ${NAMESPACE}
      restartPolicy: Always
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ${RUNNER_NAME}-gitlab-runner
  namespace: ${NAMESPACE}
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${RUNNER_NAME}-gitlab-runner
rules:
  - apiGroups: [""]
    resources: ["*"]
    verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ${RUNNER_NAME}-gitlab-runner
rules:
  - apiGroups: [""]
    resources: ["*"]
    verbs: ["*"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ${RUNNER_NAME}-gitlab-runner
subjects:
  - kind: ServiceAccount
    name: ${RUNNER_NAME}-gitlab-runner
    namespace: ${NAMESPACE}
roleRef:
  kind: ClusterRole
  name: ${RUNNER_NAME}-gitlab-runner
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ${RUNNER_NAME}-gitlab-runner-default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: default
    namespace: ${NAMESPACE}

---
```

`init-runner.sh`

```bash
#!/bin/sh
export GITLAB_HOST=https://${GITLAB_HOST}
export GITLAB_PRIVATE_TOKEN=${ACCESS_TOKEN}

# demo
NAMESPACE=runner \
RUNNER_NAME=demo \
REGISTRATION_TOKEN=${REGISTRATION_TOKEN} \
envsubst < deployment-template.yml | microk8s.kubectl apply -f -
```
