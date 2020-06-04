# treee
Fast indexing engine for hashed data

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/cyrildever/treee)
![GitHub last commit](https://img.shields.io/github/last-commit/cyrildever/treee)
![GitHub issues](https://img.shields.io/github/issues/cyrildever/treee)
![GitHub](https://img.shields.io/github/license/cyrildever/treee)

This is the Go implementation of the Treee&trade; indexing engine.
It's both a [library](#usage) to use as a module in your Go projects and an [executable](#executable) to run as a micro-service.


### Motivation

The challenge was to set up a powerful and safe search engine to use when the data is some linked list of items that could be themselves connected to each other in subchains, indexed through their identifiers that are only made of hashed values (like SHA-256 string representations), and all stored in an immutable file.

Its best application is for a Blockchain file where an item is a transaction embedding a smart contract, and each subchain of items the subsequent uses and/or modifications of this smart contract.
As such, the Treee&trade; index is currently used in the [Rooot&trade;](https://rooot.io) blockchain.


### Formal description

We define here an algorithm for indexing the system based on identifiers that are hashed values which is at the same time very powerful to the writing and the reading.

#### 1) How the index works

The Treee&trade; index is constructed as an acyclic graph (a tree). Each node contains either the node address (its sons) or a set of *Leafs*, a *Leaf* corresponding to the information helping to retrieve one or more linked items.

The number of sons of a node is deterministic and depends on the depth of the tree. We denote by <img src="https://latex.codecogs.com/gif.latex?p_{1}" /> the number of sons of the nodes at depth 1, <img src="https://latex.codecogs.com/gif.latex?p_{2}" /> the number of sons of the nodes at depth 2, ..., <img src="https://latex.codecogs.com/gif.latex?p_{k}" /> the number of sons at depth <img src="https://latex.codecogs.com/gif.latex?k" />.

The goal is to create a balanced tree whose width is adaptive to decrease depth and optimize performance. We are looking to index numbers, in this case the numerical value of the items' unique identifiers.

Let's explain the course of the index.

For the binary tree, we write the number in binary form (for example, <img src="https://latex.codecogs.com/gif.latex?4\text{&space;=&space;}" /> `100`) which indicates its position in the tree.

At the step <img src="https://latex.codecogs.com/gif.latex?i" />, we pass to the child <img src="https://latex.codecogs.com/gif.latex?0" /> if the <img src="https://latex.codecogs.com/gif.latex?i^{th}" /> bit is `0`, otherwise we pass to the child <img src="https://latex.codecogs.com/gif.latex?1" /> if the <img src="https://latex.codecogs.com/gif.latex?i^{th}" /> bit is `1`. We stop when the node is a *Leaf*.

For the tree, we build a representation of this number by a sequence of numbers and we traverse the tree in the same way. At the step <img src="https://latex.codecogs.com/gif.latex?i" /> of depth <img src="https://latex.codecogs.com/gif.latex?k" />, we pass to the child <img src="https://latex.codecogs.com/gif.latex?0" /> if the <img src="https://latex.codecogs.com/gif.latex?i^{th}" /> representative is `0`, we pass to the child <img src="https://latex.codecogs.com/gif.latex?1" /> if the <img src="https://latex.codecogs.com/gif.latex?i^{th}" /> representative is `1`, ..., we pass to the child <img src="https://latex.codecogs.com/gif.latex?p_{k_i}&space;-&space;1" /> if the <img src="https://latex.codecogs.com/gif.latex?i^{th}" /> representative is <img src="https://latex.codecogs.com/gif.latex?p_{k_i}&space;-&space;1" />. We stop when the node is a *Leaf*.

To construct the representation of a number, we will successively take the modulo of prime numbers. According to the theorem of the Chinese remains, each number has a unique representative that could be written as the continuation of these modulos. Indeed, a number <img src="https://latex.codecogs.com/gif.latex?n" /> can be written in the following form: <img src="https://latex.codecogs.com/gif.latex?n&space;\mapsto&space;(n\&space;mod\&space;p_{i},\text{&space;where&space;}p_{i}\text{&space;is&space;the&space;}i^{th}\text{&space;prime&space;number})" />

The value of the identifier of the item (its number) is denoted <img src="https://latex.codecogs.com/gif.latex?n" /> and the modulating number <img src="https://latex.codecogs.com/gif.latex?M" />. Modulos are calculated in <img src="https://latex.codecogs.com/gif.latex?O(1)" /> for fixed-sized integers. Since the multiplication is faster than the division (necessary for the calculation of the modulo), one can use multiplications by means of floating: <img src="https://latex.codecogs.com/gif.latex?M&space;\times&space;\left(n-\text{&space;int}\left(n\times&space;1&space;/&space;M\right)\right)" />.

Given the random nature of the numbers (or pseudo-random, since the identifiers of the items are generated by cryptographic hashing technologies), the tree is balanced. To unbalance the tree in a malicious way, it would be necessary to be able to generate hashes whose modulo follows a particular trajectory. However, the difficulty of such an operation increases exponentially (of the order of <img src="https://latex.codecogs.com/gif.latex?exp(plog(p))" /> where <img src="https://latex.codecogs.com/gif.latex?p" /> is the depth).
As a reminder, the product of the first 16 prime numbers equals <img src="https://latex.codecogs.com/gif.latex?32,589,158,477,190,044,730&space;\simeq&space;3&space;\times&space;10^{19}" />.
Therefore, as soon as the index contains a reasonably large amount of data, unbalancing the tree in a malicious way would become more and more impossible.

A *Leaf* contains the following list of information about an item:
 * Identifier of the current item;
 * Position: start address of the current item in the file;
 * Size: the size (in bytes) of the saved item in the file;
 * Origin: unique identifier of the item that is at the origin of the item's subchain;
 * Previous: unique identifier of the previous item chained to it;
 * Next: optionally, unique identifier of the next item chained.

 A *Leaf* whose next item field is empty is the last item in the subchain.
 
 A *Leaf* whose origin item field is equal to the identifier of the current item is necessarily the origin of the subchain. As such, it has a particular operating since, if there were to be one or more items thereafter, the last item of the subchain will be identified here as the previous item.

 The last three fields of the *Leaf* therefore correspond to a circular linked list.

#### 2) Using the index

To add an element in the index:
* The new *Leaf* is written in the index;
* We update the 'Next' field of the *Leaf* that previously corresponded to the last item of the subchain;
* We modify the 'Previous' field of the *Leaf* of the original item by writing the identifier of the current item.

To read/search an item in the index:
* We find in the tree the *Leaf* corresponding to the identifier of the searched item:
  * If the 'Next' field of the *Leaf* is empty, this is the last item of the subchain;
  * Otherwise, we go to the next step;
* We find the *Leaf* corresponding to the identifier of the field 'Origin';
* We use the 'Previous' field of this *Leaf* to find the last item of the subchain.

When using the index, we can seen that we would perform at most 3 reads or 3 writes and index runs of <img src="https://latex.codecogs.com/gif.latex?O(log(n))" /> order, where <img src="https://latex.codecogs.com/gif.latex?n" /> is the number of items in the index.


### Usage

```console
$ go get github.com/cyrildever/treee
```

```golang
import "github.com/cyrildever/treee"


```


### Executable

You can simply download the executable and start an instance of the Treee&trade; indexing engine.

```console
$ ./treee
```

##### API

// TODO


### License

This module is distributed under an MIT license.
See the [LICENSE](LICENSE) file.


<hr />
&copy; 2019-2020 Cyril Dever. All rights reserved.