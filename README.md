# muscat

muscat is a utility to use local machine's clipboard and browser from remote machine.
heavily inspired by [github.com/lemonade-command/lemonade](https://github.com/lemonade-command/lemonade)
muscat uses not tcp but unix domain socket to communicate between client and server.
you must use ssh remote forwarding of unix domain socket.

## Usage
### SSH config
```conf:~/.ssh/config
Host forward-workbench
  Hostname workbench.remote.local
  RemoteForward /home/<remote_user>/.muscat.socket /Users/<local_user>/.muscat.socket
  ExitOnForwardFailure yes
```

### Useful sshd config
```conf:/etc/ssh/sshd_config
StreamLocalBindUnlink yes
```

### start server
```sh
$ muscat server
```

### copy from stdin
```sh
$ echo foobar | muscat copy
```

### paste to stdout
```sh
$ muscat paste
```

### open URI with local browser
```sh
$ muscat open 'https://example.com'
```

## Tips
muscat can behave as pbcopy, pbpaste, xdg-open.
you can create symbolic link as their name and call with.
