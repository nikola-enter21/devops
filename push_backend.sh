#!/bin/bash

docker build -t europe-west4-docker.pkg.dev/devops-fmi-course-476112/devops-fmi-course-repo/backend-app:latest ./backend
docker push europe-west4-docker.pkg.dev/devops-fmi-course-476112/devops-fmi-course-repo/backend-app:latest
