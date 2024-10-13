import http.client
import ssl
import json

apiHost = "10.154.0.20:6443"
#ApiToken
ApiToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkduSXA4dG9BZExKMmExUXpNcWcwSG9QWHJFc0ZWZFZUTzFVdlVKWkh0OUUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQ0ZmE4ZGE2LThkOTktNDc0My1hOTk0LTAxODU3N2JlMGZkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.IXqZ2bfsqN41VWiQWvnRfi3SguSiBBBbYMjrpje4wcXMvnfdRzmJt9DepVFGfsgPEnGVuLGBAjeerStjj5GuREfXb8zlvnp_APLsDQLvZCR8ErsCZuvgR63DVW3p_Cl9K5clKu2ZzQiaTUA3J49b6xRdKKG8WHwuhxpt1hIvHmTRXciTtAdCnMr6DkFRu2WZ2aONzctTTn4LhSv2Ze7_6lAF7VKCUzOT2ZdBbkpB1p510s5vxRGyWkcDmLTFru65kw6prphPlR2DpkzzxGWcDFmyXQ6zEuJH9RPMRqVjJb0EsCzZ4wsP1BBWdM7OyVBv4eQYL7ISfjeXxlfKqKrrmw"
route_endpoint = "/apis/serving.knative.dev/v1/namespaces/default/routes"
nodes_endpoint = "/api/v1/nodes"
podsEndpoint        = "/api/v1/pods"
watchPodsEndpoint   = "/api/v1/watch/pods"
defaultPodsEndPoint = "/api/v1/namespaces/default/pods"
deploymentEndpoint  = "/apis/apps/v1/namespaces/default/deployments"
knativeSvcEndpoint  = "/apis/serving.knative.dev/v1/namespaces/default/services"
kourierPodEndPoint = "/api/v1/namespaces/kourier-system/pods"
kourierServiceEndPoint = "/api/v1/namespaces/kourier-system/services"

#curl nodeip:nodeport -H 'Host: route-url'
def send_traffic(serviceName):

    # url = get_url(serviceName)

    # node_ip = get_node_port()

    # node_port = get_node_port()

    #url = "hello-world.default.example.com"
    url = "container-12-1.default.example.com:3333"
    node_ip = "10.154.0.21"
    #node_ip = "34.142.71.212"
    node_port = "32429"
    header = {'Host': url}
    #header = {'content-type': 'application/json'}


    #header = {"Content-Type": "application/json"}
    #header = {"authorization": "Bearer " + ApiToken, "Content-Type": "application/json"}
    #conn = http.client.HTTPConnection(node_ip, node_port)
    conn = http.client.HTTPConnection(node_ip + ":"+node_port)
    #conn.request("POST", url, headers = header, body = data)
    #conn.request("GET", "/")
    conn.request("GET","/", headers=header)
    #conn.request("GET", url)
    r1 = conn.getresponse()
    
    print(r1.status, r1.reason)
    data = r1.read()
    print("data:")
    print(data.decode())
   # responseObject = json.loads(data)
    #print(responseObject)
#    print("Status: {} and reason: {}".format(r1.status, r1.reason))
    if r1.status == 200:
        print(r1.status, r1.reason)
def main():
    send_traffic("hello-world")

if __name__ == "__main__":
    main()
