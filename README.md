# file_serv



## Building the Image and container

Run the following command in the root directory of the project (where the Dockerfile is located):

```bash

docker build -t file-serv:1.0.0 .
```
### Running the Container with Required Environment Variables
```bash

docker run \
  -p 3000:3000 \
  -e SERV_PORT=3000 \
  -e BUCKET_NAME="" \
  -e S3_ACCESS_KEY="" \
  -e S3_SECRET_ACCESS_KEY="" \
  --name file-serv-cnt \
  file-serv:1.0.0
```

*  ```-p``` 3000:3000 maps the container's port 3000 to the host's port 3000.
* The ```-e``` flags specify the environment variables for the container.
* The ```-d``` flag runs the container in detached mode (in the background).
* ```--name``` file-serv-cnt assigns a custom name to the container for easier management.

### Restarting a Stopped or Crashed Container
```bash

docker start file-serv-cnt
```