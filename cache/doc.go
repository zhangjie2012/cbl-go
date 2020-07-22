/*
Package cache based on NoSQL Redis, encapsulated common scenario.

*Basic*

- Init/Close
- string/int/int64/float64/object Getter/Setter Delete
- TTL/PTTL
- compose redis key used appname/module prevent key repeat

*Distribute Lock*

support lock/unlock on distributed environment.
`ticket` for lock unique flag, avoid anther process unlock, make sure only one process lock, then unlock it.
`expire` lock timeout, avoid process dead forget unlock it.

Note: not consider redis server down caused deadlock.

*Message Queue*

based on redis data structure `list` map to a message queue. and `right push`, `left pop`.

*Counter*

a global counter.
*/
package cache
