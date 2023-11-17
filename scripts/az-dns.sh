#!/bin/sh

ip=$(kubectl get service mess-proxy -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
id=$(az network public-ip list --query="[?ipAddress=='$ip'].id | [0]" -o tsv)
az network public-ip update --ids "$id" --dns-name mess
