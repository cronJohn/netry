# About
A CLI tool that attempts to build a network topology visualization

# Q&A
## Why not use SNMP?
This tool is built around a 'black box' approach to network discovery — treating the network as an unknown system where internal configurations, such as SNMP, are either unavailable or cannot be modified.
Since SNMP requires prior configuration and access credentials on each device, it isn’t compatible with this approach.
Instead, the tool uses passive and active discovery techniques — such as those provided by nmap — to infer network topology and gather host information without relying on pre-configured management protocols.
