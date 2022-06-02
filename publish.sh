VERSION=1.1.2
docker build -t registry.cn-shenzhen.aliyuncs.com/juxuny-public/k8s-gitlab-runner:${VERSION} -f Dockerfile . && \
docker push registry.cn-shenzhen.aliyuncs.com/juxuny-public/k8s-gitlab-runner:${VERSION}
