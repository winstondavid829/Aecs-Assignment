////////////////////////////////////////////// Development - 35.228.251.106 /////////////////////////////////////////////////////////
gcloud endpoints services deploy openapi-infotech-dev.yaml
sudo docker network create --driver bridge esp_net_infotech
docker build . -t go-dock-infotech-dev:latest
docker tag go-dock-infotech-dev:latest gcr.io/inofetch-app/go-dock-infotech-dev:latest
docker push gcr.io/inofetch-app/go-dock-infotech-dev:latest
sudo docker pull gcr.io/inofetch-app/go-dock-infotech-dev:latest
sudo docker run --detach --restart always --name=esp-infotech --net=esp_net_infotech gcr.io/inofetch-app/go-dock-infotech-dev:latest


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
docker build . -t go-dock-aecspulls-dev:latest

docker build -t aecspulls-dev .
docker tag aecspulls-dev:latest winstondavidj/aecspulls-dev:latest
docker push winstondavidj/aecspulls-dev:latest

docker pull winstondavidj/aecspulls-dev:latest
docker run -p 3001:8080 --detach --restart always --name=esp-aecs-pulls winstondavidj/aecspulls-dev:latest

docker run -p 3001:8080 --detach --restart always --name=esp-aecs-pulls go-dock-aecspulls-dev:latest

// --------------- aws deployment ------------ //
docker build -t aecspulls .

docker tag aecspulls:latest 727433422324.dkr.ecr.us-west-1.amazonaws.com/aecspulls:latest

docker push 727433422324.dkr.ecr.us-west-1.amazonaws.com/aecspulls:latest