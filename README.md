
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

TODO
- Custom specs.
- Queue system for agent.
- On/off params field in spec.
- Figure out cli and interactions.

# Development

## Using ngrok

Tasks:

- Aws provider, split schema and state, so that we remove a lot of load from node.
