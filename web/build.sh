#!/bin/bash 

build() {
  rm app
  GOOS=linux go build -o app
  docker rmi webcrawler
  docker build . -t webcrawler
}

push() {
  send_command "docker system prune -f"
  aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 732407346024.dkr.ecr.us-east-1.amazonaws.com
  docker tag webcrawler:latest 732407346024.dkr.ecr.us-east-1.amazonaws.com/webcrawler:latest
  docker push 732407346024.dkr.ecr.us-east-1.amazonaws.com/webcrawler:latest
}

send_command() {
  aws ssm send-command \
    --document-name "AWS-RunShellScript" \
    --document-version "1" \
    --targets "Key=instanceids,Values=i-09006bb3612ab4df6" \
    --parameters "commands='$1'" \
    --timeout-seconds 600 \
    --max-concurrency "50" \
    --max-errors "0" \
    --region us-east-1
}

pull() {
  echo "Pulling docker image from ECR to EC2 Instance"
  send_command "aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 732407346024.dkr.ecr.us-east-1.amazonaws.com"
  send_command "docker pull 732407346024.dkr.ecr.us-east-1.amazonaws.com/webcrawler:latest"
  send_command "docker image tag 732407346024.dkr.ecr.us-east-1.amazonaws.com/webcrawler:latest webcrawler:latest"
  echo "Successfully pulled docker image from ECR to EC2 Instance"
} 

run() {
  send_command "docker network ls | grep discord || docker network create discord"
  send_command "docker run -d -p 9090:9090 --network=discord --name webcrawler webcrawler"
}

stop() {
  send_command "docker stop webcrawler"
  send_command "docker container prune -f"
}

init_ec2(){
  send_command "sudo chmod 666 /var/run/docker.sock"
}

echo "Starting build"
case $1 in
  build_and_deploy)
    echo "Building go executable"
    build
    push
    pull
    stop
    run
    ;;
  init)
    echo "Init EC2"
    init_ec2
    ;;
  run)
    echo "Running docker container"
    run
    ;;
  stop)
    echo "Stoping EC2 instance"
    stop
    ;;
  *)
    echo "No flags passed, doing nothing"
    ;;
esac
echo "Finished build"
