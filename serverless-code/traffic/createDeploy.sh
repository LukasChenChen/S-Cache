curl https://10.154.0.20:6443/apis/apps/v1/namespaces/default/deployments \
  --header "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Im9fQ1pLR0NFTGhmN0pMNGNDeGJYc3JkTWhIVHV3dEVPcnVFVGFPemdRbGMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjQ4MmEzNTcyLTM4Y2EtNGJhOS04Y2Q1LTg4Yzc4MjU3Mzk1YSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.Qvr6g9_4TU9ilFF-xQXRlGVGK3AtLBjnidvA5yYUC77768O8gCSSQUcDh_Iu2QWJu08su2NclevvwcNpXpa3ZYxtfmMu41tnJF8VFAoubxzn6zTIsLAC4rJqCJsp-GSMI3LypOp5zh4Mywa1CBFup1AHumb8Y27zGJDdIx_mxkAX4UZRcyCP3wWw7KTXOETTec5avwQ8QSKb7KPJcbD3ktH2IWZHoO3TtyQ7JTtesDwzLbSy4XhzIFX-JTIUhYkRRJJRHV0B2z6xBUBpNXkN455sUzWps82kombUiKtnNmO3ehNPzsGbyW0f9ooo1qfutYKmWVe4oeUE8j-892NAjQ" \
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
