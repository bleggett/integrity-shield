#!/usr/bin/env bash

HELM_VER="3.6.0"

helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts                                                                                                                 130 ↵ bleggett@bleggett

helm install gatekeeper/gatekeeper --name-template=gatekeeper --namespace gatekeeper-system --create-namespace --version $HELM_VER

kubectl create -f https://raw.githubusercontent.com/stolostron/integrity-shield/master/integrity-shield-operator/deploy/integrity-shield-operator-latest.yaml

kubectl create -f https://raw.githubusercontent.com/stolostron/integrity-shield/master/integrity-shield-operator/config/samples/apis_v1_integrityshield.yaml -n integrity-shield-operator-system
