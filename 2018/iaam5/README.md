# Intro

You need:
- A working Kubernetes cluster to deploy our application.
- Your kubectl command should be able to interact with that cluster.
- This repository so you can interact with all files of the repository.

Some additions:
- If you want to have auto completion in your console for the kubectl command have a look at https://kubernetes.io/docs/tasks/tools/install-kubectl/#enabling-shell-autocompletion
- If you like to do tail the log of more than one pod at the same time have a look at https://github.com/johanhaleby/kubetail or just install it using `sudo curl -o /usr/local/bin/kubetail https://raw.githubusercontent.com/johanhaleby/kubetail/master/kubetail && sudo chmod +x /usr/local/bin/kubetail`

# Sidenote

There is an extensive documentation for the API we will use online https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/

# Goals

- We have an old application (./iaam5-app/) which we currently deploy by hand on our server. The application uses an in-memory database.
- The application manages a list of timestamps with their unique ID.
- The application has two HTTP request handlers "GET /" to receive all timestamps and their IDs as well as "POST /add" to add a new entry.
- The "GET /" handler is very slow (it always takes at least 5 seconds!) and so we want to scale the application.
- The application should be deployed using Kubernetes so that we can scale it to more instances.
- If we can scale we also want to migrate to an external database.

# Steps

## Create and use your own namespace

- `kubectl create namespace mz1987` This is "m" because of my firstname and "z" because of my lastname and "1987" because of my birthday.
- `kubectl config set-context $(kubectl config current-context) --namespace=mz1987`

If there are no errors you are now using your own namespace where you can deploy your own applications.

## Deploy the old application (./iaam5-app/)

Let's have a look at our files:
- ./iaam5-app/main.go holds the application's source code
- ./iaam5-app/Dockerfile holds the Dockerfile to build a Docker image for the application
- ./iaam5-app/build.sh holds the commands for compiling the application and creating the Docker image
- We already uploaded a Docker image for you to the official Docker registry with the name "symflower/iaam5:app"
- ./iaam5-app/app.yml holds the YAML definition for our Kubernetes deployment of the application. As you can see we make the server available with the port 8080.

Let's deploy, use and debug our application:
- `kubectl apply -f app.yml` Apply the configuration for the deployment. Kubernetes will automatically download the Docker image for you, create the container and assign resources like CPU, memory and disk space.
- `kubectl get pods` List our pods in the current namespace.
- `kubectl get pods -l service=app` List all pods with the label "service" equal to "app" of the current namespace.
- `kubectl describe pod $POD_NAME` Detailly display all information of the pod with the name $POD_NAME. Notice the "IP". Why is it not possible to ping the this IP?
- `kubectl logs $POD_NAME` Output the log of the pod with the name $POD_NAME.
- `kubectl logs -f $POD_NAME` Tail the log of the pod with the name $POD_NAME.

Let's make the deployment available from the outside:
- `kubectl expose deployment app --type=NodePort` Which creates a Kubernetes service for our deployment and opens a unique port for the container server port on every Kubernetes node.
- `kubectl get service` Display all services and their details. Somewhere in the list is the service "app" which has a port value such as "8080:31283/TCP". "31283" is the external port.
- `kubectl describe nodes | grep -i extern` we need the external IP addresses of our Kubernetes nodes too.
- `curl http://$IP:$PORT/` Do an HTTP request to our current pod. Tail the log to see what is going on. Each request should take about 5 minutes.
- `curl --request POST http://$IP:$PORT/add` Add an entry to the internal database.
- `(time curl http://$IP:$PORT/ &); (time curl http://$IP:$PORT/ &); (time curl http://$IP:$PORT/)` Query the service three times at once.

Let's scale our service:
- Edit our deployment YAML file. Uncomment the "replicas" line. The value of it tells Kubernetes how many instances we want to have of our pod.
- Apply the configuration again. Kubernetes will only transfer and execute the changes.
- `watch kubectl get pods` Have a look how these pods are started by Kubernetes.
- `kubectl delete pod $POD_NAME` We can even delete one of the pods and Kubernetes just starts new pods until we reach our requested instance number.
- Since we now have more than one service our request time of the second and third request should be much faster if we do three requests at once.
- `(time curl http://$IP:$PORT/ &); (time curl http://$IP:$PORT/ &); (time curl http://$IP:$PORT/)`
- If you have installed "kubetail" you can look at the logs of all pods at once with: `kubetail -l service=app`

We have successfully deployed our old application in Kubernetes and scaled it but there are two problem: We loose the database if a pod fails and we do not share data among the instances. We need an external database!

# Deploy an external database (./iaam5-db)

- ./iaam5-db/db.yml holds the YAML definition for our Kubernetes deployment of the database. As you can see we make the server available with the port 5432. We are also using a ConfigMap to configure the database and we can use this ConfigMap to later configure the application. There is also a readiness probe in our container definition which tells Kubernetes that the container is ready to be accessed. Since we do not want to access the database from outside of our cluster we do not expose the service to the outside but keep an internal service which is also defined in the same file.

Let's deploy, use and debug our database:
- `kubectl apply -f db.yml` Deploy our ConfigMap, our database deployment and our database service.
- Have a look if our pod is already up and tail the log of the database.

We are now ready to deploy a new version of our application which uses our external database.

# Deploy the new application (./iaam5-app-db/)

Let's have a look at our files:
- ./iaam5-app-db/main.go holds the application's source code
- ./iaam5-app-db/Dockerfile holds the Dockerfile to build a Docker image for the application
- ./iaam5-app-db/build.sh holds the commands for compiling the application and creating the Docker image
- We already uploaded a Docker image for you to the official Docker registry with the name "symflower/iaam5:app-db"
- ./iaam5-app-db/app-db.yml holds the YAML definition for our Kubernetes deployment of the application. As you can see we make the server available with the port 8080. And we use the ConfigMap that we defined earlier for our database!

Let's deploy, use and debug our application:
- `kubectl apply -f app-db.yml` Apply the configuration for the deployment.
- Have a look at our pods. Why do we suddenly have only ONE pod?!
- `curl --request POST http://$IP:$PORT/add` add some new entries.
- What happens now if we delete the pod?
- Scale the application and access the application three times at once. Do you see the same data for every response?

That's it! If you liked this small workshop and would like to know more or want to see more details contact us at hello@symflower.com and have a look at what else we have to offer on https://symflower.com/
