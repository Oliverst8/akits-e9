- It appears to work perfectly without Lamport timestamps. I think this is due to the hieracy of ports, being used as the comparison. In that case, the lamport timestamps are only neccesary to fulfill the Liveliness Criterion
- The MutexNode passed around (known as "s") is not passed by pointer (ie. by value). Changing this breaks the functionality.
* * One of the problems caused by passing by pointer, was that the list of requests to respond to was never cleared. This can be easily fixed, but more problems are remain.
- When the MutexNode is passed by value, the changes to its Lamport time are (unsurprisingly) unsaved. Changes are just made to the local value.

