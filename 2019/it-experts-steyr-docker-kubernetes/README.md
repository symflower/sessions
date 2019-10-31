# Intro

This workshop gives you a very basic overview on how to use Docker and Kubernetes. After the workshop you should be able to move on to more advanced topics about Docker and Kubernetes and it should be a good starting point in your journey to use Docker and Kubernetes for your own applications.

The slides for this workshop can be found in [Symflower - IT-Experts - Docker and Kubernetes.pdf](./Symflower%20-%20IT-Experts%20-%20Docker%20and%20Kubernetes.pdf).

All the examples can be run on your local (Linux) machine using the described tools that need to be installed first (see chapter "Install"). The examples of course work with other Docker and Kubernetes instances.

# Install

- `mkdir bin`
- [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/) with `curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && mv minikube bin/`
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) with `curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl && chmod +x kubectl && mv kubectl bin/`
- [kubetail](https://github.com/johanhaleby/kubetail) with `curl -Lo kubetail https://raw.githubusercontent.com/johanhaleby/kubetail/master/kubetail && chmod +x kubetail && mv kubetail bin/`
- [Go](https://golang.org/doc/install) with `curl -Lo go.tar.gz https://dl.google.com/go/go1.13.linux-amd64.tar.gz && tar xvfz go.tar.gz && rm go.tar.gz`
- If you do not have "direnv" installed you have to run some settings from the file ".envrc".
- Install some Go packages:

```
go get github.com/codegangsta/negroni
go get github.com/gorilla/mux
go get github.com/lib/pq
```

- For kubectl auto-completion use `kubectl completion bash > .kubectl-completion && source .kubectl-completion`.
- Start minikube with `minikube start`. You should see something similar to:

	```
	😄  minikube v1.4.0 on "Opensuse-Tumbleweed"
	🔥  Creating virtualbox VM (CPUs=2, Memory=2000MB, Disk=20000MB) ...
	🐳  Preparing Kubernetes v1.16.0 on Docker 18.09.9 ...
	🚜  Pulling images ...
	🚀  Launching Kubernetes ...
	⌛  Waiting for: apiserver proxy etcd scheduler controller dns
	🏄  Done! kubectl is now configured to use "minikube"
	```

- Log into a new console and close the current once, because some environment variables have now changed.
- Check if you have a node available with `kubectl get nodes`.
- Explore `minikube help`.
- Check the IP of the VM with `minikube ip`.
- Explore `minikube dashboard`.

# Goals

- We have an old application (./src/itexperts-app/) which we currently deploy by hand on our server. The application uses an in-memory database.
- The application manages a list of timestamps with their unique ID.
- The application has two HTTP request handlers "GET /" to receive all timestamps and their IDs as well as "POST /add" to add a new entry.
- The "GET /" handler is very slow (it always takes at least 5 seconds!) and so we want to scale the application.
- The application should be deployed using Kubernetes so that we can scale it to more instances.
- If we can scale we also want to migrate to an external database.

# Some basics about Docker

- You can interact with the Docker instance of minikube.
- Beware that networking is inside of the minikube VM, so better login via `minikube ssh` to have more access.
- Explore `docker`
- Explore `docker info`
- Run one single container even if we do not have touched "images" yet. We run a terminal, the bash with a container.
	- `docker run ubuntu:18.04 bash` why does it exit immediatelly?
	- `docker ps` shows us running containers. Where is it?
	- `docker ps -a` shows us **all** containers.
	- `docker rm $CONTAINER_ID` remove all exited containers.
	- `docker run -it ubuntu:18.04 bash` run container in interactive (holds STDIN open) mode with a (pseudo) TTY. We should see a bash now.
	- `docker ps` can we see the running container?
	- `docker stop $CONTAINER_ID` or `docker kill $CONTAINER_ID` what is the difference?
	- `docker pause $CONTAINER_ID` what happens? Can we use the container? Why not?
	- `docker unpause $CONTAINER_ID` what happens?
	- Maybe we need some more information about the container. `docker inspect $CONTAINER_ID`.
	- `docker logs -f $CONTAINER_ID` what happens here if you type something in the container?
- Can we access files of another container?
	- Create a file in container A and try to access it in container B. Does that work? What is happening?
	- Let's create a directory and share the directory between containers. Connect to Minikube and run the following commands:
		- `mkdir shared`
		- Run `docker run -it -v $PWD/shared:/shared ubuntu:18.04 bash` on two different tabs.
		- Can you now create and access files from container A and container B in "/shared"? Why? How does that work?
- Can we ping a container?
	- From inside the container? Why?
	- From inside another container? Why? (this is one thing where Kubernetes is different, because it isolates more)
	- From inside the minikube VM? Why?
	- From inside the host? Why?
	- If you need the "ip" tool in a Ubuntu container, it is in the package "iproute2".
	- If you need the "ping" tool in a Ubuntu container, it is in the package "iputils-ping".
- All this installing is getting me tired, is there a way of creating a container with the tools already installed?
	- Yes, its called an "image".
	- Let's create Dockerfile to create a Docker image.
		- `mkdir nettools && cd nettools`
		- Create a file called "Dockerfile" with the following content

			```Docker
			FROM ubuntu:18.04

			RUN apt-get update
			RUN apt-get install -y iproute2
			RUN apt-get install -y iputils-ping
			```
		- Every line in this file holds a different command. What do you think the different commands mean?
		- Let's create the Docker image using `docker build --tag netutils .`.
		- What does the output mean? What are the steps? What are intermediate containres? What are this weird looking IDs?
		- What happens when we rerun the command? Why?
		- What happens when we switch commands e.g. the two install commands? Why?
		- During the build command we are using an argument "tag". What do you think that is good for?
	- Let's run our Docker image with `docker run -it netutils bash`.

# Steps for Kubernetes

## Deploy the old application (./src/itexperts-app/)

Let's have a look at our files:
- ./src/itexperts-app/main.go holds the application's source code
- ./src/itexperts-app/Dockerfile holds the Dockerfile to build a Docker image for the application
- ./src/itexperts-app/build.sh holds the commands for compiling the application and creating the Docker image
- ./src/itexperts-app/app.yml holds the YAML definition for our Kubernetes deployment of the application. As you can see we make the server available with the port 8080.

Let's deploy, use and debug our application:
- `./src/itexperts-app/build.sh` compile the application and create the Docker image.
- `watch kubectl get all` watch what happens to our k8s objects.
- `kubectl apply -f ./src/itexperts-app/app.yml` Apply the configuration for the deployment. Kubernetes will automatically download the Docker image for you, create the container and assign resources like CPU, memory and disk space.
- `kubectl get pods` List our pods in the current namespace.
- `kubectl get pods -l service=app` List all pods with the label "service" equal to "app" of the current namespace.
- `kubectl describe pod $POD_NAME` Detailly display all information of the pod with the name $POD_NAME. Notice the "IP". Why is it not possible to ping the this IP?
- `kubectl logs $POD_NAME` Output the log of the pod with the name $POD_NAME.
- `kubectl logs -f $POD_NAME` Tail the log of the pod with the name $POD_NAME.
- `kubectl exec -it $POD_NAME bash` We can even run a command inside the pod.

Let's make the deployment available from the outside:
- `kubectl expose deployment app --type=NodePort` Which creates a Kubernetes service for our deployment and opens a unique port for the container server port on every Kubernetes node.
- `kubectl get service` Display all services and their details. Somewhere in the list is the service "app" which has a port value such as "8080:31283/TCP". "31283" is the external port.
- `kubectl edit service $SERVICE_NAME` Allows us to directly edit the configuration of the service and have a look on how a service is defined.
- `minikube ip` we need the IP of minikube to access the services forwarded for external usage.
- Tail the log of the pod (e.g. with kubetail) in another console.
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

# Deploy an external database (./src/itexperts-db)

- ./src/itexperts-db/db.yml holds the YAML definition for our Kubernetes deployment of the database. As you can see we make the server available with the port 5432. We are also using a ConfigMap to configure the database and we can use this ConfigMap to later configure the application. There is also a readiness probe in our container definition which tells Kubernetes that the container is ready to be accessed. Since we do not want to access the database from outside of our cluster we do not expose the service to the outside but keep an internal service which is also defined in the same file.

Let's deploy, use and debug our database:
- `watch kubectl get all` watch what happens to our k8s objects.
- `kubectl apply -f ./src/itexperts-db/db.yml` Deploy our ConfigMap, our database deployment and our database service.
- Have a look if our pod is already up and tail the log of the database.

We are now ready to deploy a new version of our application which uses our external database.

# Deploy the new application (./src/itexperts-app-db/)

Let's have a look at our files:
- ./src/itexperts-app-db/main.go holds the application's source code
- ./src/itexperts-app-db/Dockerfile holds the Dockerfile to build a Docker image for the application
- ./src/itexperts-app-db/build.sh holds the commands for compiling the application and creating the Docker image
- ./src/itexperts-app-db/app-db.yml holds the YAML definition for our Kubernetes deployment of the application. As you can see we make the server available with the port 8080. And we use the ConfigMap that we defined earlier for our database!

Let's deploy, use and debug our application:
- `./src/itexperts-app-db/build.sh` compile the application and create the Docker image.
- `watch kubectl get all` watch what happens to our k8s objects.
- `kubectl apply -f ./src/itexperts-app-db/app-db.yml` Apply the configuration for the deployment.
- Have a look at our pods. Why do we suddenly have only ONE pod?!
- `curl --request POST http://$IP:$PORT/add` add some new entries.
- What happens now if we delete the pod?
- Scale the application and access the application three times at once. Do you see the same data for every response?

That's it! If you liked this small workshop and would like to know more or want to see more details contact us at hello@symflower.com and have a look at what else we have to offer on https://symflower.com/
