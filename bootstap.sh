#!/bin/bash

printf "\nEnter the site FQDN: "
read FQDN

printf "\nEnter your OpenAI API key: "
read -s OPENAI_API_KEY

helm repo add coroot https://coroot.github.io/helm-charts
helm repo update coroot

# Ingress controller
printf "\nChecking for Nginx ingress controller:\n"
helm ls -n ingress-nginx | grep ingress-nginx > /dev/null || \
helm upgrade --install ingress-nginx ingress-nginx --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace

# coroot
printf "\nChecking for coroot services:\n"
helm ls -n coroot | grep coroot > /dev/null || \
helm install --namespace coroot --create-namespace --set node-agent.tracesEndpoint="" \
    --set corootCE.ingress.enabled=true --set corootCE.ingress.hostname=coroot.${FQDN} \
    --set corootCE.ingress.className=nginx coroot coroot/coroot

# daas
printf "\nChecking for daas services:\n"
kubectl create namespace daas || true
kubectl apply -f ./manifests/daas

kubectl create secret generic gptscript-key -n daas \
    --from-literal=key="${OPENAI_API_KEY}" || true

# myapp-mem
printf "\nChecking for myapp-mem:\n"
helm ls | grep myapp-mem > /dev/null || \
helm install myapp-mem ./charts/myapp -f ./charts/myapp/mem-values.yaml --set fqdn=${FQDN}

# myapp-cpu
printf "\nChecking for myapp-cpu:\n"
helm ls | grep myapp-cpu > /dev/null || \
helm install myapp-cpu ./charts/myapp -f ./charts/myapp/cpu-values.yaml --set fqdn=${FQDN}

# for i in $(seq 1 200); do curl -I http://myapp-mem.${FQDN}; curl -I http://myapp-cpu.${FQDN}; sleep 2; done