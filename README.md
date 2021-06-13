# Rewire (WIP)
Rewire makes it easy to create encrypted TCP connections across networks using a centralized transport server and client agents on both endpoints. It acts like a VPN software, but doesn't require elevated privileges to run since it doesn't add a tun/tap driver to the system.

## Overview
The Rewire architecture consists of three different components:

- **Control Client** - The host initiating the connection to a desired host on the remote agent network.
- **Remote Agent** - The tunnel connection endpoint connecting to the destination host.
- **Transport Server** - The server which is reachable from both the control client and remote agent that facilitates the communication between both client and agent.

## Remote Agent

This is the software that is deployed on a remote server where network access is required. It is run by specifying a server and port to connect, which can be either done via CLI arguments or a config file. All opened connections using this agent will appear to be originating from server which the agent is installed.

## Control Client

The control client is installed on the machine desiring to open a connection to a device on the remote agent or its network. This client also runs as an agent service which is not directly interacted with, but creates local tunnels to the transport server when it receives a command to do so from the transport server.

## Transport Server

The transport server is the connection point that both endpoint agents register with upon starting. Using the transport server's REST API, an end user can initiate a connection by specifying a registered control client and remote agent as well as the target host and port to connect. This command will then provide the local endpoint on the control client which the end user can use as the entry point to the tunnel. The transport server upon registration stores the agent keys locally to prevent spoofing attacks.

### TODO:
- [ ] Agent authentication via API tokens
- [ ] Transport server API authentication
- [ ] Rewrite the socket communication in Protobuf
- [ ] Finish rewire lib separation from individual modules
- [ ] Complete CLI interface for agents
- [ ] Wrap agents in an installable service with logging
- [ ] Key rotation for agents
- [ ] Documentation for API, Rewire lib and Agent CLIs