#compile code
go buiid

#push to the remote dockerhub
docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7  -t cocc/scheduler-demo:v3.0 .

#main.go main function

#kubernetes.go   get, read resources from api
#scaler.go       scale the pods
#bestprice.go    pod assignment
#types.go        customised data structure           
