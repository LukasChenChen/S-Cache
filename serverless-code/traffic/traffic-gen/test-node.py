import http.client
import ssl
import json
import yaml 

ApiToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkduSXA4dG9BZExKMmExUXpNcWcwSG9QWHJFc0ZWZFZUTzFVdlVKWkh0OUUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQ0ZmE4ZGE2LThkOTktNDc0My1hOTk0LTAxODU3N2JlMGZkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.IXqZ2bfsqN41VWiQWvnRfi3SguSiBBBbYMjrpje4wcXMvnfdRzmJt9DepVFGfsgPEnGVuLGBAjeerStjj5GuREfXb8zlvnp_APLsDQLvZCR8ErsCZuvgR63DVW3p_Cl9K5clKu2ZzQiaTUA3J49b6xRdKKG8WHwuhxpt1hIvHmTRXciTtAdCnMr6DkFRu2WZ2aONzctTTn4LhSv2Ze7_6lAF7VKCUzOT2ZdBbkpB1p510s5vxRGyWkcDmLTFru65kw6prphPlR2DpkzzxGWcDFmyXQ6zEuJH9RPMRqVjJb0EsCzZ4wsP1BBWdM7OyVBv4eQYL7ISfjeXxlfKqKrrmw"
apiHost = "10.154.0.20:6443"
kourierPodEndPoint = "/api/v1/namespaces/kourier-system/pods"
kourierServiceEndPoint = "/api/v1/namespaces/kourier-system/services"

#kubectl get pods -n kourier-system -o yaml  | grep hostIP
def get_node_ip():
    header = {"authorization": "Bearer " + ApiToken, "Content-Type": "application/json"}

    #url = apiHost + ":" + node_endpoint
    conn = http.client.HTTPSConnection(apiHost, context = ssl._create_unverified_context())
    conn.request("GET", kourierPodEndPoint, headers = header)
    r1 = conn.getresponse()
    print(r1.status, r1.reason)
    if r1.status == 200:
        data_str = r1.read().decode()
        json_obj = json.loads(data_str)
        #print(data_str)   
        items = json_obj['items']
        kourier_pod = items[0]
        host_ip = items[0]['status']['hostIP']
        with open('config.yaml', 'w') as yaml_file:
                yaml.dump(json_obj, yaml_file)
    return host_ip
# kubectl get svc -n kourier-system kourier -o yaml | grep nodePort , pick first one for http, second for https
def get_node_port():
    header = {"authorization": "Bearer " + ApiToken, "Content-Type": "application/json"}

    
    conn = http.client.HTTPSConnection(apiHost, context = ssl._create_unverified_context())
    conn.request("GET", kourierServiceEndPoint, headers = header)
    r1 = conn.getresponse()
    kourier_node_port = ''
    
    print("ready to print")
    print(r1.status, r1.reason)
    if r1.status == 200:
        data_str = r1.read().decode()
        json_obj = json.loads(data_str) 
        #0 means http port, 1 means https port
    #    kourier_node_port = json_obj['spec']['ports'][0]['nodePort']
       # print(json_obj['items'][0]['spec']['ports'][0]['nodePort'])
        serviceList = json_obj['items']
        for item in serviceList:
            print(item['metadata']['name'])
            print(type(item['metadata']['name']))
            if item['metadata']['name'] == 'kourier':
                kourier_node_port = item['spec']['ports'][0]['nodePort']
                print('kourier node port is')
                print(kourier_node_port)
    #    print(json_obj)
    
    else:
       print(r1.status, r1.reason)

    return kourier_node_port

def main():
    ip = get_node_ip()
    port = get_node_port()
    print('ip is')
    print(ip)
    print(type(ip))

    print('port is')
    print(type(port))
if __name__ == "__main__":
                main()
