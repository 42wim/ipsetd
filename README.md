# ipsetd
A daemon that exposes ipset in interactive mode and listens by default on :9999

Possible use case is to have several remote systems push live (big) ipset updates to ipsetd without running ```ipset``` for every call.


E.g create a test.set 
```
destroy abc
create abc hash:net
add abc 1.2.3.4
add abc 1.2.3.5
```

Now you can do ```cat test.set|nc yourhost 9999```

