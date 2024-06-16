# Lock-free data structures

Collection of various lock-free data structures.

Currently contains three data-structures:

* [Array N](array_n.go)
* [Treiber Stack](treiber_stack.go)
* [Michael-Scott Queue](michael_scott_queue.go)

## Array N

The main idea: lock-free access/update operations over array(slide). Atomic counter provides access to the cell of the
array.

Caution! Only for cases when number of threads is immutable and known before start of processing. Don't forget to call
before processing Reserve() method with number of threads.

## Treiber Stack

The main idea: apply CAS-loops on push/pop operations till item will pushed/popped.

## Michael-Scott Queue

The main idea: use linked list as storage; apply CAS-loops at enqueue to the tail and deque from the head. Use assistance
of parallel threads to solve ABA-problem.
