from http.server import HTTPServer, BaseHTTPRequestHandler
from io import BytesIO
import json
import logging
from Gen_traffic import *
from queue import Queue
import time
from socketserver import ThreadingMixIn
import threading
import pandas as pd
import math
import datetime
# logging.basicConfig(filename='local.log', filemode='w', format='%(name)s - %(levelname)s - %(message)s')
logging.basicConfig(filename='server.log', filemode = 'w', level=logging.DEBUG)

SERVER_HOST = '0.0.0.0'
SERVER_PORT = 2222
FINAL_RESULT = {}
REQ_COUNT = 0
# queue = Queue()

# x = "5"
def getDist(coords_1, coords_2):
    """
    Calculate the great circle distance between two points 
    on the earth (specified in decimal degrees)
    """
    lat1 = coords_1[0]
    lon1 = coords_1[1]
    lat2 = coords_2[0]
    lon2 = coords_2[1]

    #if same coordinate, dist is 0.
    if lat1 == lat2 and lon1 == lon2:
        dist = float(0)
        return dist
    

    radlat1=math.pi*lat1/180
    radlat2=math.pi*lat2/180
	
    theta=lon1-lon2
    radtheta=math.pi*theta/180
    dist=math.sin(radlat1)*math.sin(radlat2)+math.cos(radlat1)*math.cos(radlat2)*math.cos(radtheta)

    if dist > 1:
        dist = 1
	
	
    dist=math.acos(dist)
    dist=dist*180/math.pi
    dist=dist*60*1.1515
	

    dist=dist*1.609344
 
    return dist
## check if the distance is right
def get_linkDelay(coords_1, coords_2, latencyPara):
    #read the topo file
    dist = getDist(coords_1, coords_2)
    # dist = dist/1000  ## in km size
    logging.warning("The distance is %f km" % dist)
    return dist*latencyPara*2



class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):

    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'Hello, world!')

    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length)
        my_json = body.decode('utf8')
        request = json.loads(my_json) ## this is just one request
        
        self.send_response(200)
        self.end_headers()
        # print("x is ", x)
        
        # for request in data:
        if request['Served'] == True:
    #        logging.warning(request)
            # logging.warning(type(request))
            # queue.put(request)
            send_request(request)
            message =  threading.currentThread().getName()
            self.wfile.write(str.encode(message))
           

def send_request(request):
    # threading.Timer(0.5, send_request).start()i
    threshold = 8
    global REQ_COUNT
    starttime = time.time()
    runtime = 0
    while True:
        service_name = request['Function']['Name']
        processing_delay, succ = send_traffic(service_name, request)
        if succ == False and time.time() - starttime < threshold:
            time.sleep(0.5)
        else:
            break
    REQ_COUNT += 1
    runtime = time.time() - starttime
    # logging.warning('request ', request['ID'],'runtime is ', runtime)
    coords_1 = (request['Ingress']['Lat'], request['Ingress']['Long'])
    coords_2 = (request['DeployNode']['Lat'], request['DeployNode']['Long'])
    link_delay = get_linkDelay(coords_1, coords_2, request['LatencyPara'])
    total_delay = runtime + link_delay
    FINAL_RESULT[request['ID']] = total_delay
    logging.warning( 'Request %s, link_delay %f, runtime is %f, total_delay %f' % (request['ID'],link_delay, runtime, total_delay))
    result = [{'time': request['ArriveTime'], 'ID': request['ID'],'type': request['Function']['Type'], 'linkDelay':link_delay, 'processing_delay': processing_delay, 'runtime':runtime, 'total_delay':total_delay, 'cold_start':request['IsColdStart']}]
    ##if total delay large means the resource not enough on node, the container is pending all the time.
    if total_delay < threshold:
        df = pd.DataFrame(result)  
        df.to_csv("./result.csv",header=False, mode='a', index=False)
  #  outfile = open("./result.csv", 'a')
  #  outfile.write("total count \n")
  #  outfile.write(str(REQ_COUNT) + "\n")
  #  outfile.close()

def readlines(filename):
    file = open(filename)
    # read the content of the file opened
    content = file.readlines()
    config = content[9:17]     
    return config
class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in a separate thread."""
       

def main():
    print('Listening on host %s port %s ...'% (SERVER_HOST , SERVER_PORT))
    with open("./result.csv", mode='a') as file:
        file.write(' recorded at %s.\n' %  datetime.datetime.now())
    file.close()
    config = readlines("./config/config-map.yaml")
    df_0 = pd.DataFrame({'count': config})
    df_0.to_csv ("./result.csv", index = None, header = False, mode = 'a')
    title = ['time', 'ID', 'type', 'linkDelay', 'processing_delay', 'runtime', 'total_delay', 'cold_start']
    with open("./result.csv", 'a') as fp:
        for item in title:
            # write each item on a new line
            fp.write("%s\n" % item)
   # httpd = HTTPServer((SERVER_HOST, SERVER_PORT), SimpleHTTPRequestHandler)
    httpd = ThreadedHTTPServer((SERVER_HOST, SERVER_PORT), SimpleHTTPRequestHandler)
    httpd.serve_forever()

if __name__ == "__main__":
   main()
