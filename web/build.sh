#!/bin/bash 

build_docker() {
  build_linux
  docker rmi webcrawler
  docker build -t webcrawler .
}

build_windows() {
  rm app
  GOOS=windows go build 
}

build_linux() {
  rm app
  GOOS=linux go build -o app
}

upload_to_s3() {
  file_to_upload=$1
  bucket_name=$2
  aws s3 cp $1 $2 --recursive
}

push() {
  send_command "docker system prune -f"
  aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com
  docker tag webcrawler:latest ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/webcrawler:latest 
  docker push ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/webcrawler:latest

}

send_command() {
  aws ssm send-command \
    --document-name "AWS-RunShellScript" \
    --document-version "1" \
    --targets "Key=instanceids,Values=${INSTANCE_ID}" \
    --parameters "commands='$1'" \
    --timeout-seconds 600 \
    --max-concurrency "50" \
    --max-errors "0" \
    --region ${REGION}
}

pull() {
  echo "Pulling docker image from ECR to EC2 Instance"
  send_command "aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"
  send_command "docker pull ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/webcrawler:latest"
  send_command "docker image tag ${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/webcrawler:latest webcrawler:latest"
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

echo "Getting flags"
S3_BUCKET="S"
while getopts :a:r:i:s:b: opt; do
  case $opt in
    a)
      ACCOUNT_ID=${OPTARG}
      ;;
    r)
      REGION=${OPTARG}
      ;;
    i)
      INSTANCE_ID=${OPTARG}
      ;;
    s)
      S3_BUCKET=${OPTARG}
      ;;
    b)
      BUILD_STEP=${OPTARG}
      ;;
  *)
    echo "$opt Flag not supported"
    ;;
 esac
done
echo "Successfully got all flags"

echo ${S3_BUCKET}
echo "Starting build step: $BUILD_STEP"
case $BUILD_STEP in
  build_docker)
    build_docker
    ;;
  build_windows)
    build_windows
    ;;
  build_linux)
    build_linux
    ;;
  build_and_deploy)
    echo "Building go executable"
    build_docker
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
