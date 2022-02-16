
# Atlas

- Database
    - Stop/restart
- Update data
- Queue system
- Plugins
    - Prysm

Remove all containers:

```
$ docker rm -f $(docker ps -a -q)
```

// Put all the user input data that can be changed (size of volume, cpu..) under the same path
// differentiate between the parts that are backend supplied and the ones that not.
// Note that there are two things that can be updated:
    // one for the instance (cpu)
    // another one for the node (cache flag)

// Bring AWS to build the system, do we need the queue?
    // with the queue create the system to stop, update data and restart.
// Create Snapshot system.

// Plugin system to execute commands
    // differentiate between stopped and running commands.
    // exec in the machine?

// Provision
    // Node state => Stopped, Terminated, (Latency), Provisioning.

# Development

## Using ngrok
