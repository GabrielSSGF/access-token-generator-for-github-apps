<h1 align="center">
  <br>
  <a href=""><img src="https://cdn-icons-png.flaticon.com/512/1068/1068671.png" alt="Keys" width="200"></a>
  <a href=""><img src="https://upload.wikimedia.org/wikipedia/commons/thumb/a/ae/Github-desktop-logo-symbol.svg/1024px-Github-desktop-logo-symbol.svg.png" alt="Github" width="200"></a>
  
  <br>
  Access Token Generator for Github Apps
  <br>
</h1>

<h4 align="center">API designed to generate a access token for your Github App.</h4>

<p align="center">
  <a href="#key-features">Key Features</a> •
  <a href="#how-to-use">How To Use</a>
  
</p>

## Key Features

* Generates a Access Token for your Github App
  - By using your Github App PEM, this application generates a jwt token and authenticates in the Github API, allowing for the access token generation.
* Multiple env storage options
  - Choose between using your Github App Pem stored on your dotenv file or in your AWS secret manager
* Dockerfile in case you decide to put it in production 

## How To Use

To clone and run this application, you'll need [Git](https://git-scm.com), [Go](https://go.dev/) and [Docker](https://docs.docker.com/engine/reference/commandline/cli/). From your command line:

```bash
# Clone this repository
$ git clone https://github.com/GabrielSSGF/access-token-generator-for-github-apps

# Go into the repository
$ cd access-token-generator-for-github-apps

# * Configure your .env file / Github vars with the necessary information

# Build and run the app
$ go build
$ ./access-token-generator

# If you decide to run within a container
$ docker run -P --env-file .env ubuntu bash

```
