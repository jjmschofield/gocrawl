# GoCrawl
GoCrawl is a gopher powered web crawler for the internets! It exists largely as a learning tool for the author, it's only there second attempt at a Go app so be kind!

![Cute gopher holding a network cable](./docs/network-gopher.png)

If you find anything useful in here, feel free to use this code as you wish (within the bounds of applicable laws and the terms of service of any website you choose to run this code against). 

## Getting Started

## Objective 
So our crawler should:

* Crawl to find all pages from a starting absolute URL 
* Not follow external pages
* Print a simple site map, showing the links between pages

Sounds kind of simple? We'll make things a bit more interesting (and real worldish) with some additional non-functional requirements:

* As this operation is going to be waiting a lot for I/O operations - we will crawl pages in parallel to make better use of system resources and to get the job done faster
* Whilst DDOSing some poor webmasters site sounds like fun - we'll add a concurrency control to control how many parallel crawls we have running   
* We will leave the solution open for extension to store crawl results outside of memory - this could support crawling huge sites or many sites, but we are likely going to find writes tricky if we are running really fast parallel operations
* We will leave the solution open for extension to run as separate services - this would support scaling parts of the operation independently, but we are going to have to be quite careful with how we lay the code out

## Benchmarks
* **CPU:** i5-9600K (6 cores) @ ~4.3ghz
* **Mem:** 16GB DDR4 @ 2666ghz
* **Network:** ~80Mbps
* **Disk**: Samsung EVO SSD
* **OS**: Windows 10

## Solution  
### Project Structure
In an attempt ot arrive at an idiomatic we have followed the (non-official but fairly popular) standards [found here](https://github.com/golang-standards/project-layout).


                                                                               



## Thanks to
[Renee French](http://reneefrench.blogspot.com/) for the wonderful gopher icon from [this github repo](https://github.com/egonelbre/gophers).

