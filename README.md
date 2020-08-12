
Build docker image for leader:
docker build . --build-arg SERVICE=leader -t leader

Launch docker container for leader:
docker run -p external_port_here:3000 leader

Build docker image for follower:
docker build . --build-arg SERVICE=follower -t follower

Launch docker container for follower:
docker run -p external_port_here:3000 follower