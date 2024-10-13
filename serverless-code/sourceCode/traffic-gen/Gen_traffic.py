import http.client
import ssl
import json

apiHost = "10.154.0.20:6443"

#ApiToken
ApiToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkduSXA4dG9BZExKMmExUXpNcWcwSG9QWHJFc0ZWZFZUTzFVdlVKWkh0OUUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQ0ZmE4ZGE2LThkOTktNDc0My1hOTk0LTAxODU3N2JlMGZkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.IXqZ2bfsqN41VWiQWvnRfi3SguSiBBBbYMjrpje4wcXMvnfdRzmJt9DepVFGfsgPEnGVuLGBAjeerStjj5GuREfXb8zlvnp_APLsDQLvZCR8ErsCZuvgR63DVW3p_Cl9K5clKu2ZzQiaTUA3J49b6xRdKKG8WHwuhxpt1hIvHmTRXciTtAdCnMr6DkFRu2WZ2aONzctTTn4LhSv2Ze7_6lAF7VKCUzOT2ZdBbkpB1p510s5vxRGyWkcDmLTFru65kw6prphPlR2DpkzzxGWcDFmyXQ6zEuJH9RPMRqVjJb0EsCzZ4wsP1BBWdM7OyVBv4eQYL7ISfjeXxlfKqKrrmw"
routesEndpoint = "/apis/serving.knative.dev/v1/namespaces/default/routes"
nodesEndpoint = "/api/v1/nodes"
podsEndpoint        = "/api/v1/pods"
watchPodsEndpoint   = "/api/v1/watch/pods"
defaultPodsEndPoint = "/api/v1/namespaces/default/pods"
deploymentEndpoint  = "/apis/apps/v1/namespaces/default/deployments"
knativeSvcEndpoint  = "/apis/serving.knative.dev/v1/namespaces/default/services"
kourierPodEndPoint = "/api/v1/namespaces/kourier-system/pods"
kourierServiceEndPoint = "/api/v1/namespaces/kourier-system/services"

#kubectl get routes
#url is servicename.default.example.com
def get_url(serviceName):

    url = serviceName + ".default.example.com"
    
    return url

#kubectl get pods -n kourier-system -o yaml  | grep hostIP
def get_node_ip():
   
    header = {"authorization": "Bearer " + ApiToken, "Content-Type": "application/json"}

    conn = http.client.HTTPSConnection(apiHost, context = ssl._create_unverified_context())
    conn.request("GET", kourierPodEndPoint, headers = header)
    r1 = conn.getresponse()
    kourier_node_ip = ''
    if r1.status == 200:
        data_str = r1.read().decode()
        json_obj = json.loads(data_str)  
        items = json_obj['items']
        #usually there is only one pod in this namespace
        kourier_node_ip = items[0]['status']['hostIP']
    
     # dump the json object to a yaml file
     #  with open('config.yaml', 'w') as yaml_file:
     #         yaml.dump(json_obj, yaml_file)
    else:
       print(r1.status, r1.reason)

    return kourier_node_ip

# kubectl get svc -n kourier-system kourier -o yaml | grep nodePort , pick first one for http, second for https
def get_node_port():
    header = {"authorization": "Bearer " + ApiToken, "Content-Type": "application/json"}

    
    conn = http.client.HTTPSConnection(apiHost, context = ssl._create_unverified_context())
    conn.request("GET", kourierServiceEndPoint, headers = header)
    r1 = conn.getresponse()
    kourier_node_port = ''
    
    if r1.status == 200:
        data_str = r1.read().decode()
        json_obj = json.loads(data_str) 
        #0 means http port, 1 means https port
        print(json_obj['items'][0]['spec']['ports'][0]['nodePort'])
        
        serviceList = json_obj['items']
        for item in serviceList:
            if item['metadata']['name'] == 'kourier':
                kourier_node_port = item['spec']['ports'][0]['nodePort']

        
    
    else:
       print(r1.status, r1.reason)

    return kourier_node_port

def send_get(serviceName):
    print("start send traffic")
    #kubectl get route, the url of the hello-world svc
    url = get_url(serviceName)
    url = url +":3333"
    
    #the node that host the kourier pod
    node_ip = get_node_ip()
    
    #the nodeport for kourier svc
    node_port = get_node_port()

    #url = "hello-world.default.example.com"
    #node_ip = "10.154.0.21"
    #node_port = "32429"
    header = {'Host': url, 'Content-type': 'application/json'}

    print("request node_ip",node_ip, "node_port", node_port,"header Host", url )
    
    conn = http.client.HTTPConnection(node_ip, node_port)
    # data = {'test':'send example data'}
    conn.request("GET", "/", headers= header)

    r1 = conn.getresponse()
    print(r1.status, r1.reason)
    print(r1.read().decode())
    conn.close()

#curl nodeip:nodeport -H 'Host: route-url'
#TODO: finish sending traffic then close connection
def send_traffic(serviceName, request):
    print("start send_traffic ", serviceName)
    #kubectl get route, the url of the hello-world svc
    url = get_url(serviceName)
    url = url +":3333"
    
    #the node that host the kourier pod
    node_ip = get_node_ip()
    
    #the nodeport for kourier svc
    node_port = get_node_port()

    print("request node_ip",node_ip, "node_port", node_port,"header Host", url )
    
    conn = http.client.HTTPConnection(node_ip, node_port)
    # data = {'test':'send example data'}
    json_body = json.dumps(request)

    #url = "hello-world.default.example.com"
    #node_ip = "10.154.0.21"
    #node_port = "32429"
    header = {'Host': url, 'Content-type': 'application/json', 'Content-length': len(json_body)}

    conn.request("POST", "/", headers = header, body = json_body)
    #conn.request("GET", "/", headers= header)

    r1 = conn.getresponse()
    print(r1.status, r1.reason)
    
    if r1.status == 200 :
        data = r1.read().decode()
    #    print("print data", data)
        # print(type(data))
        result = json.loads(data)
        processing_delay = result['processing_delay']
        conn.close()
        return processing_delay, True
    else:
        conn.close()
        return 0, False



if __name__ == "__main__":
   send_traffic("hello-world")
