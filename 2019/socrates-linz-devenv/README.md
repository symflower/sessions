# Automated Development Environment (DevEnv)

Automated -> Installing/Updating is automated.

Question: Why do you want an automated DevEnv? My main goal -> automated reproducible DevEnv for every change of the repository.
Question: Who has an (partly) automated DevEnv? What was the experience?
Question: Who has a CI that automatically builds, runs tests, ...? -> If yes, you are almost there.

In my opinion there are phases you need to do:
- Automate the CI -> Then you can already build, run and test what you are developing!
- Do you have no/one service (one container is fine) or multiple services (I would use a VM with Vagrant then, but there is also "Docker in Docker")?
- Creating a Docker/Vagrant image is well documented. But how do you connect them to the host system?
	- Docker -> use forward ports/volumes (e.g. `docker run --name socrates --rm -it --publish 8080:8080 --volume $PWD:/go/src/github.com/symflower/sessions/2019/socrates-linz symflower/socrates-linz-2019:latest bash`)
	- Vagrant -> we use NFS. NFS was a few years ago the fastest solution, maybe there is now a faster one.
- DirEnv is one important tool for automatic DevEnv.
- Until now we only deployed in Containers/VMs what about the host system?
	- Packages/Configurations you need?
	- Editor?
		E.g. visual studio code
			You can install vscode via a script and then use --user-data-dir and --extensions-dir to make the configurations local for vs code. So you can install extensions per project.
	- Mail? Chat? ...
	- Maybe next step: Use "vagrant ssh --help"
	- https://code.visualstudio.com/docs/remote/remote-overview ?

One problem of such a setup: How do you do updates/rollbacks?
	- For us: rollbacks -> vagrant destroy && up
	- For us: updates: semi-automatic. Most things are automated but not **everything**. We are working on it when there is time.

One thing we did not automate yet: use images built by the CI e.g. Docker and Vagrant images instead of doing them on every developer machine.

@Symflower we have our DevEnv automated but not completely, e.g mail is still completely manually configured.
