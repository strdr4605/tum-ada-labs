docker container prune -f
docker run --name openmpi  -it -v "$(pwd)":/home/student/lab1/openmpi openmpi