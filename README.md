# Consistent Hash

This is an implementation of the consistent hashing algorithm.

It is made in order to better understand the method of managing nodes of a distributed NOSQL DB. 

## Functional Requirements
The controller would support 
- adding nodes
- removing nodes
- assigning storage to node(s) 
- data redundancy
- adjusting Consistency and Availability via R/W Quorum configuration
- each node will receive 100 virtual nodes
- use the first 64 bits of the SHA-256 hash algorithm
- retrieve the number of nodes currently distributed

## Limitations
The controller will not be responsible for identifying failed or unresponsive nodes, spinning up replacements nodes, or dealing with such reliability concerns because this simulation does not work with actual nodes and thus won't have a heartbeat or health check implemented.

No data will be stored, for this simulation, storage data will not be saved; the only goal will be to identify which nodes the key/value pair would be stored at. 

Nodes will be stored in memory and will not persist on server restart.

## Non-functional Requirements

The controller should be highly performant, forwarding the requests in <5ms. 

Methods can be reached via HTTP request. 

We can think of the controller being utilized by a team managing a distributed DB. The controller should be responsible for all nodes and the status. 

## API 
POST /nodes -> node ID string
body: {
    url: string // internal url 
}

DELETE /nodes/{node-id}

POST /data -> list of node-ids
body: {
    key: string,
    value: string 
}

POST /config -> 200 OK
body: {
    readQuorum: number,
    writeQuorum: number,
    redundancy: number // < total nodes
}

