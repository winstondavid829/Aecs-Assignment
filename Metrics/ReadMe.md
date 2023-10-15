////////////////////////////////////////////// Development - 35.228.251.106 /////////////////////////////////////////////////////////
gcloud endpoints services deploy openapi-infotech-dev.yaml
sudo docker network create --driver bridge esp_net_infotech

docker build . -t go-dock-aecsmetrics-dev:latest
docker tag go-dock-infotech-dev:latest gcr.io/inofetch-app/go-dock-infotech-dev:latest
docker push gcr.io/inofetch-app/go-dock-infotech-dev:latest
sudo docker pull gcr.io/inofetch-app/go-dock-infotech-dev:latest
sudo docker run -p 3000:8080 --detach --restart always --name=esp-infotech --net=esp_net_infotech gcr.io/inofetch-app/go-dock-infotech-dev:latest

// To run the code in local
docker run -p 3000:8080 --detach --restart always --name=esp-aecs-metrics go-dock-aecsmetrics-dev:latest


 #Creating Server Blocks in nginx
 For this user script 'create_gcpendpoint.sh' under /automate_nginx folder, the command line parameter is cloud end point url
sudo docker run --name=espssl-5001-infotech \
     --detach \
     --publish=5001:443 \
     --net=esp_net_infotech \
     -v /etc/letsencrypt/live/infotech-dev.endpoints.inofetch-app.cloud.goog/fullchain.pem:/etc/nginx/ssl/nginx.crt \
     -v /etc/letsencrypt/live/infotech-dev.endpoints.inofetch-app.cloud.goog/privkey.pem:/etc/nginx/ssl/nginx.key \
     --link=esp-infotech:esp-infotech \
     gcr.io/endpoints-release/endpoints-runtime:1 \
     --service=infotech-dev.endpoints.inofetch-app.cloud.goog \
     --rollout_strategy=managed \
     --backend=esp-infotech:8080 \
     --ssl_port=443
     
======================================================================================================================================================

/// ---------------------- Local Environment Setup ---------------------- ///
docker build -t aecsmetrics-dev .
docker tag aecsmetrics-dev:latest winstondavidj/aecsmetrics-dev:latest
docker push winstondavidj/aecsmetrics-dev:latest


docker run -p 3000:8080 --detach --restart always --name=esp-aecs-metrics go-dock-aecsmetrics-dev

// To run the code in local - with env
docker run -p 3000:8080 -e AWS_REGION=us-west-1 --detach --restart always --name=esp-aecs-metrics go-dock-aecsmetrics-dev:latest  

/// ---------------------- AWS Setup ---------------------- ///
Authenticating docker client to registry 
aws ecr get-login-password --region us-west-1 | docker login --username AWS --password-stdin 727433422324.dkr.ecr.us-west-1.amazonaws.com

// Build docker 
docker build -t aecs .

// Tag
docker tag aecs:latest 727433422324.dkr.ecr.us-west-1.amazonaws.com/aecs:latest

// Push image to repo
docker push 727433422324.dkr.ecr.us-west-1.amazonaws.com/aecs:latest

