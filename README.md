This repository is divided into two main components, each with their own subcomponents.
The first main component pertains to Battlesnake-specific logic, such as encoding and decoding 
board states, handling requests from the Battlesnake server, and expanding or evaluating 
a board state.
The second main component pertains to Minimax-specific logic, such as coordinating minimax tree
expansion and evaluation, as well as gRPC connections with follower machines.

The following descriptions show how to deploy a leader machine and how to deploy follower machines.
Everything is configured in docker so, if you have docker installed, deployment should be easy. Simply 
add your desired external port numbers to the commands below.

Build docker image for leader:
docker build . --build-arg SERVICE=leader -t leader

Launch docker container for leader:
docker run -p external_port_1:3000 -p external_port_2:3001 leader

Build docker image for follower:
docker build . --build-arg SERVICE=follower -t follower

Launch docker container for follower:
docker run -p external_port_here:3000 follower

Note that in order to have the leader generate gRPC connections with the follower machines, the 
addresses of the follower machines must be sent in a REST POST request to the leader's /followers endpoint
once the leader has launched, as shown below:
{
    "addresses": [
        "HOST1:PORT1",
        "HOST2:PORT2",
        ...
    ]
}

After the gRPC connections have been established, the Battlesnake server's endpoints are ready to service requests!
