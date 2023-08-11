$ curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzEK"}}'<br/>
$ curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzIK"}}'<br/>
$ curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzMK"}}'<br/>

$ curl -X GET localhost:8080 -d '{"offset": 0}'<br/>
$ curl -X GET localhost:8080 -d '{"offset": 1}'<br/>
$ curl -X GET localhost:8080 -d '{"offset": 2}'<br/>

. Record—the data stored in our log.<br/>
• Store—the file we store records in.<br/>
• Index—the file we store index entries in.<br/>
• Segment—the abstraction that ties a store and an index together.<br/>
• Log—the abstraction that ties all the segments together.<br/>

• C—country<br/>
• L—locality or municipality (such as city)<br/>
• ST—state or province<br/>
• O—organization<br/>
• OU—organizational unit (such as the department responsible for owning the key)<br/>

$ sudo kubectl create -f proglog.yml<br/>
sudo kubectl delete -f proglog.yml<br/>

# to see rendered template
$ helm template proglog

#remove extra files
$ rm proglog/templates/**/*.yaml proglog/templates/NOTES.tx

# to see rendered helm
$ helm template proglog deploy/proglog

# Now, install the Chart by running this command:
$ helm install proglog deploy/proglog

# We can tell Kubernetes to forward a pod or a Service’s port to a port on your
# computer so you can request a service running inside Kubernetes without a
# load balancer:
$ kubectl port-forward pod/proglog-0 8400 8400

# run the command to request our service to get and print the list of
# servers:
$ go run cmd/getservers/main.go
