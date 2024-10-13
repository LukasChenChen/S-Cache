curl 10.154.0.6:30434 --header "Host: container-4-18.default.svc.cluster.local:3333" \
	-X POST \
	-d '{ "title":"foo","body":"bar", "id": 1}' \
	-H 'Content-Type: application/json' \
