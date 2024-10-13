curl https://10.154.0.20:6443/apis/apps/v1/namespaces/default/deployments \
  --header "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IkduSXA4dG9BZExKMmExUXpNcWcwSG9QWHJFc0ZWZFZUTzFVdlVKWkh0OUUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQ0ZmE4ZGE2LThkOTktNDc0My1hOTk0LTAxODU3N2JlMGZkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.IXqZ2bfsqN41VWiQWvnRfi3SguSiBBBbYMjrpje4wcXMvnfdRzmJt9DepVFGfsgPEnGVuLGBAjeerStjj5GuREfXb8zlvnp_APLsDQLvZCR8ErsCZuvgR63DVW3p_Cl9K5clKu2ZzQiaTUA3J49b6xRdKKG8WHwuhxpt1hIvHmTRXciTtAdCnMr6DkFRu2WZ2aONzctTTn4LhSv2Ze7_6lAF7VKCUzOT2ZdBbkpB1p510s5vxRGyWkcDmLTFru65kw6prphPlR2DpkzzxGWcDFmyXQ6zEuJH9RPMRqVjJb0EsCzZ4wsP1BBWdM7OyVBv4eQYL7ISfjeXxlfKqKrrmw"  \
  --insecure \
  -X POST \
  -H 'Content-Type: application/yaml' \
  -d '---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sleep
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sleep
  template:
    metadata:
      labels:
        app: sleep
    spec:
      containers:
      - name: sleep
        image: curlimages/curl
        command: ["/bin/sleep", "365d"]
'
