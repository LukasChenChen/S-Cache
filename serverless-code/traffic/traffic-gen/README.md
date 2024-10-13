# curl nodeip:nodeport -H 'Host: route-url'
# nodeip: the node host kourier pod
kubectl get pods -n kourier-system -o yaml  | grep hostIP
# nodeport kubectl get svc -n kourier-system kourier -o yaml | grep -i nodeport
# route-url: kubectl get routes 
