#image resize

from http.server import HTTPServer, BaseHTTPRequestHandler
from io import BytesIO
import json
import logging
import math
from checkout import checkout
import time
from socketserver import ThreadingMixIn
import threading

# logging.basicConfig(filename='local.log', filemode='w', format='%(name)s - %(levelname)s - %(message)s')


SERVER_HOST = '0.0.0.0'
SERVER_PORT = 3333

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
        logging.warning('Come in a get...')
        checkout(12335)
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'This is checkout app...')

    def do_POST(self):
        start = time.time()
        logging.warning('Come in a  post...')
        checkout(12335)
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length)
        my_json = body.decode('utf8')
        request = json.loads(my_json) ## this is just one request
        logging.warning('print request %f', content_length)
        # logging.warning('print request %s', request)
        
        self.send_response(200)
        self.end_headers()
        response = BytesIO()
        time.sleep(0.02)
        # results = []
        total_delay = 0
        # for request in data:
        if request['Served'] == True:
            logging.warning('Request served')
            # coords_1 = (request['Ingress']['Lat'], request['Ingress']['Long'])
            # coords_2 = (request['DeployNode']['Lat'], request['DeployNode']['Long'])
            # logging.warning('Ingress ID:', request['Ingress']['ID'], "Deploynode:", request['DeployNode']['ID'])
            # link_delay = get_linkDelay(coords_1, coords_2, request['LatencyPara'])

            # if request['IsColdStart'] == True:
            #     total_delay = link_delay + request['Function']['ProcessTime'] + request['Function']['ColdStartTime']
            # else:
            #     total_delay = link_delay + request['Function']['ProcessTime']
            processing_delay = time.time() - start
            logging.warning('Processing time %f', time.time() - start)
            result = {
                 "ID": request['ID'],
                 "processing_delay": processing_delay,
            }
            #     "cold_start_delay": request['Function']['ColdStartTime'],
            #     "processing_delay": request['Function']['ProcessTime'],
            #     "link_delay": link_delay,
            #     "total_delay": total_delay,
            # }
        
            logging.warning('send latency result %s', result)
            response.write(str.encode(json.dumps(result)))
            self.wfile.write(response.getvalue())
        
        logging.warning('Finish response to checkout request...')
        
class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
    """Handle requests in a separate thread."""

def main():
    logging.warning('Listening on host %s port %s ...'% (SERVER_HOST , SERVER_PORT))
    #httpd = HTTPServer((SERVER_HOST, SERVER_PORT), SimpleHTTPRequestHandler)
    httpd = ThreadedHTTPServer((SERVER_HOST, SERVER_PORT), SimpleHTTPRequestHandler)
    httpd.serve_forever()



if __name__ == "__main__":
   main()
