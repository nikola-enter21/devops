#!/bin/bash

docker build -t europe-west4-docker.pkg.dev/devops-fmi-course-2025/devops-fmi-course-repo/backend-app ./backend
docker push europe-west4-docker.pkg.dev/devops-fmi-course-2025/devops-fmi-course-repo/backend-app
