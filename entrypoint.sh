#!/bin/bash
if [ -z ${NAME} ]; then echo "runner name cannot be empty"; exit -1; fi
if [ -z ${TOKEN} ]; then echo "TOKEN cannot be empty"; exit -1; fi
if [ -z ${HOST} ]; then echo "HOST cannot be empty"; exit -1; fi
if [ -z ${NAMESPACE} ]; then echo "NAMESPACE cannot be empty"; exit -1; fi
gitlab-runner register --non-interactive --url ${HOST} --registration-token ${TOKEN} --executor docker --docker-image juxuny/kubectl-envsubst:v1.18.10  --name ${NAME}
cat /etc/gitlab-runner/config.toml
export REGISTRATION_TOKEN=$(cat /etc/gitlab-runner/config.toml| yj -tj | jq -r ".runners[0].token")
gitlab-runner stop
cat config-template.toml | envsubst > /etc/gitlab-runner/config.toml
cat /etc/gitlab-runner/config.toml
/entrypoint run
