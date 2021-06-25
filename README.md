$ curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzEK"}}'
$ curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzIK"}}'
$ curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzMK"}}'

$ curl -X GET localhost:8080 -d '{"offset": 0}'
$ curl -X GET localhost:8080 -d '{"offset": 1}'
$ curl -X GET localhost:8080 -d '{"offset": 2}'


$ sudo kubectl create -f proglog.yml
sudo kubectl delete -f proglog.yml

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

