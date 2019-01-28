# GoCrawl
GoCrawl is a gopher powered web crawler for the internets! It exists largely as a learning tool for the author, it's only there second attempt at a Go app so be kind!

<p align="center">
  <img alt="Cute gopher holding a network cable" src="docs/network-gopher.png>
</p>

If you find anything useful in here, feel free to use this code as you wish (within the bounds of applicable laws and the terms of service of any website you choose to run this code against). 

## Project Structure
In an attempt ot arrive at an idiomatic we have followed the (non-official but fairly popular) standards [found here](https://github.com/golang-standards/project-layout).

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

## Approach
### Language Selection
So first and foremost we are going to pick Go as our tool. 

Go is well suited to this operation due to it's praised parallel processing model - something which starting out the author hasn't worked with before. Hopefully through necessity comes much learnings!

We will try hard during the execution to achieve an idiomatic codebase, we can expect to make many missteps here (all code is bad code)  however hopefully the author will learn a thing or two!

### Extracting Links from a Page
First up we are only going to deal with HTML pages, and then we are only go deal with the response body from the server. This is to say we will not be executing any javascript. This means we will disadvantage any onClick js events or any sites that rely heavily on DOM manipulations (eg from a non SSR SPA). Google may have recently become capable of this, but for our simple crawler we will keep the process light and easy and avoid any virtual DOM or headless browser techniques - trusting that most sites are still using good SEO practices.  

Next we will need to get the links from the page efficiently, probably using an iterator to pluck out the anchor tags and strip the href attribute from them. If we encounter bad HTML, we should probably kill crawling of that page only - hanging on to any valid URLs we were able to grab. There may be more suitable business logic to handling failures - but we don't get much from the brief and can assume this is just something we need to see running to know what to do.

We are also going to need to determine between internal and external URLs - even when a content manager is referencing internal pages referencing any combination of http/https absolute/relative URLs. Probably normalizing all URLs to absolute urls (based on the originally provided URL) is the easiest way to handle this.  

Finally, we are going to get mailto's and tel's as well as internal and external links. We may as well hang on to these and provide them in the output - it probably doesn't affect complexity all that much (and we could answer why some people are getting huge amounts of spam).

### Concurrency
So we are going to want to make as many requests as possible to keep our box busy - but we are also going to want to control the number of requests being made to avoid hitting rate limits and generally to avoid being horrible to poor webmasters.

We will encounter a delay when requesting HTML from a site so we should plan for this to run in a gorouting.

In the future we will probably hit delays if we decide to write results to storage. We should probably plan for this (even if we are only storing things in memory) so that the solution doesn't grow to the point where making a change of this nature becomes needlessly complex. 

This will also be a fun exercise as we should also try to avoid locks and with many concurrent processes we can assume we are going to end up in race conditions. We can probably tolerate crawling a page twice (though we should avoid it for performance), but scenarios where we are loosing data to last write wins need to be avoided.

### Output
We haven't got much information on what an output should look like. We are asked for a sitemap but standard sitemaps don't involve a graph structure to relate the links between pages. 

We have many options here, dumping it to console, putting it into a JSON file, chucking everything into graph storage of some nature? We'll firm this up as we move forward.

### Deployments
We won't make an effort to put this into a suitably deployable format, but this might make quite a fun extension in the future if we want to index all the things! 


## Thanks to
[Renee French](http://reneefrench.blogspot.com/) for the wonderful gopher icon from [this github repo]( https://github.com/egonelbre/gophers ).