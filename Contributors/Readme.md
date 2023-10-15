docker build -t aecscontributors .

docker tag aecscontributors:latest 727433422324.dkr.ecr.us-west-1.amazonaws.com/aecscontributors:latest

docker push 727433422324.dkr.ecr.us-west-1.amazonaws.com/aecscontributors:latest


// =================== Local ==================//
docker build -t aecscontributors .
docker tag aecscontributors:latest winstondavidj/aecscontributors:latest
docker push winstondavidj/aecscontributors:latest
