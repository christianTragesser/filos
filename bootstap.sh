#!/bin/bash

printf "\nEnter the site FQDN: "
read FQDN

if kubectl get secret gptscript-key -n daas > /dev/null; then
  printf "\n Found GPTScript APIkey secret.\n"
else
  printf "\nEnter your OpenAI API key: "
  read -s OPENAI_API_KEY
  kubectl create namespace daas > /dev/null || true
  kubectl create secret generic gptscript-key -n daas \
      --from-literal=key="${OPENAI_API_KEY}"
fi

printf "\n"
helm repo add coroot https://coroot.github.io/helm-charts
helm repo update coroot

# Ingress controller
printf "\nChecking for Nginx ingress controller:\n"
helm ls -n ingress-nginx | grep ingress-nginx > /dev/null || \
helm upgrade --install ingress-nginx ingress-nginx --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
kubectl wait --timeout=120s --for=condition=Available deploy/ingress-nginx-controller -n ingress-nginx > /dev/null

# coroot
printf "\nChecking for coroot services:\n"
helm ls -n coroot | grep coroot > /dev/null || \
helm install --namespace coroot --create-namespace --set node-agent.tracesEndpoint="" \
    --set corootCE.ingress.enabled=true --set corootCE.ingress.hostname=coroot.${FQDN} \
    --set corootCE.ingress.className=nginx coroot coroot/coroot

# daas
printf "\nChecking for daas services:\n"
kubectl apply -f ./manifests/daas


## myapp-mem
printf "\nChecking for myapp-mem:\n"
helm ls | grep myapp-mem > /dev/null || \
helm install myapp-mem ./charts/myapp -f ./charts/myapp/mem-values.yaml --set fqdn=${FQDN}

# myapp-cpu
printf "\nChecking for myapp-cpu:\n"
helm ls | grep myapp-cpu > /dev/null || \
helm install myapp-cpu ./charts/myapp -f ./charts/myapp/cpu-values.yaml --set fqdn=${FQDN}

printf "\nWhen ready, use the following bash command to simulate application request traffic:\n\n'for i in \$(seq 1 200); do curl -I http://myapp-mem.${FQDN}; curl -I http://myapp-cpu.${FQDN}; sleep 1; done'\n"