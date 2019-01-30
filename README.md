# GoCrawl
GoCrawl is a gopher powered web crawler for the internets! It exists largely as a learning tool for the author, it's only there second attempt at a Go app so be kind!


![Cute gopher holding a network cable](./docs/network-gopher.png)

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

## Benchmarks
* **CPU:** i5-9600K (6 cores) @ ~4.3ghz
* **Mem:** 16GB DDR4 @ 2666ghz
* **Network:** ~80Mbps
* **Disk**: Samsung EVO SSD
* **OS**: Windows 10

### First Attempt
First we are going to do a spread of worker counts to asses if our concurrency model is working. We aren't going to take averages yet - so the numbers may spike.

| Site  | Pages | 1 Worker | 10 Workers | 25 Workers  | 50 Workers | 100 Workers | 1000 Workers | 10000 Workers |
|---|---|---|---|---|---|---|---|---|
| https://www.magicleap.com | 109 | 2167ms  | 712ms | 830ms |  576ms | 469ms  | 525ms  | 761ms |
| https://www.akqa.com | 417 | 27348ms  | 2722ms | 2476ms  | 2603ms  | 2405ms   | 1674ms  | 2521ms |
| https://monzo.com | 1351 | 747346ms  | 90327ms | 36993ms |  21567ms | 17062ms  | 16335ms  | 18041ms |

So we definitely have parallel operations taking place and getting some benefits from them!

Entertainingly our crawls seem to be limited more by the number of urls we have discovered then anything else. For example when we hit empty pages like http://monzo.com/blog/authors/kate-hollowood/11 we seem to only discover one link each time (the next link) - which means that we only end up using a couple of workers in parallel 

To get some real performance insights we need to really stress the solution. When we crawl a big complex site like https://www.bbc.co.uk we get some interesting results:

* CPU sits around 30% utilization
* Memory eventually gets exceeded and the system starts paging
* Bandwidth sits constantly at 75Mbps - 80Mbps

Now this is more like the data that we need! 

### Optimizations
So some observations:

* Bandwidth is our key limiting factor when we encounter a big site that responds quickly - excellent news!
* For small sites - we hit a point where the number of parallel workers doesn't matter so much, or worse degrades performance due to resource allocation times 
* We are blowing up memory
    * This is largely because we are storing everything in memory and not streaming out our results
    * Paging doesn't seem to hurt us (CPU doesn't go through the roof) - but then we are using a physically attached SSD
    * This will cause problems on a VM where memory is more expensive and page file may not be available
    * The use of caches and streaming results to storage (rather then hanging on to them) should clear this up
    * We could probably be more efficient with memory we are using by keeping only what we need

#### Do less work
##### Only treat pages on the same protocol as internal
We are currently treating both http and https for the host as being internal links, [sitemaps.org](https://www.sitemaps.org/protocol.html) tells us this is not right.

We can see a ton of what looks like duplicated pages for several of the sites under test so we can safely assume we should see an improvement if we do this.  

What happens when we did this?
 
| Site  | Pages | 1000 Workers Run 1 | 1000 Workers Run 2 | 1000 Workers Run 3 | 1000 Workers Avg | Page Count | Performance Change | 
|---|---|---|---|---|---|---|---|
| https://www.magicleap.com | 109 | 526ms | 635ms | 509ms | 556ms | 0 (0%) | +31ms (+0.06%) |   
| https://www.akqa.com | 177 | 928ms | 1078ms | 938ms | 981ms | -294 (-70%) | -693ms (-41%) |
| https://monzo.com | 688 | 14334ms  | 14341ms | 12684ms | 13786ms | -663 (-49)%  | -2549ms (-16%) | 
 
How exciting! Interestingly two of the sites had a large number of mixed links, not so surprisingly crawling less pages dramatically reduced time.

We may also orphan some pages if the site is badly linked together sadly - but them's the breaks. Really no-one should be linking to http, and we can probably quickly use the results from the crawler to discover and quickly fix the naughty pages.

##### Avoid chasing files
We are currently following urls which are not html pages at all (eg pdf's, png's and other common file types), these will probably never give us a `text/http` response so why try? 

Anyways, we can see an error being printed to screen whenever we encounter one of these non-html entities and we are incorrectly treating them as broken pages. It's probably a bug so lets just fix it and see what happens.

What happened when we did this?

Of the sites under test, only Monzo had any noticeable change from the number of pages crawled with -13, https://www.akqa.com had 0 and https://www.magicleap.com had 1. Crawl times in all cases appear to have not really noticed the impact.

In order to understand the benefits of this we should take a look at a site which has a large number of download links. https://www.cas.org.uk have many reports offered up as pdf - so it would seem like an excellent choice. 
 
So first a benchmark without the file filter:
 
| Site  | Pages | 1000 Workers Run 1 | 1000 Workers Run 2 | 1000 Workers Run 3 | 1000 Workers Avg |  
|---|---|---|---|---|---|
| https://www.cas.org.uk | 1,329 | 40,512ms | 42,160ms | 37,021ms | 39,897ms |

And when we apply the filter:

| Site  | Pages | 1000 Workers Run 1 | 1000 Workers Run 2 | 1000 Workers Run 3 | 1000 Workers Avg |  
|---|---|---|---|---|---|
| https://www.cas.org.uk | 520 | 1398ms | 1413ms | 1389ms | 1430ms |

Wow! The difference for download heavy sites is staggering 61% less "pages" to crawl resulting in a 96% less crawl time.

TODO 
Noticed cas sites cache goes cold quick- maybe the benchmark is not quite as good as I thought (how slow is that site without the cache!)

![Dancing gopher](./docs/gopher-dance-long-3x.gif)

What's the fun of working in Go if you can't break out an animated gopher when things go well?

#### Store less things
##### Can we reduce memory usage by structuring things a bit better?
TODO
* Stopped hanging onto big link structures for any longer then nesc by adding only the raw url to outlinks on pages
* Added id for outpages
* Carried out page construction in worker so main thread doesn't get concerned with it
* Refactored a race condition (accessing a non locked var outside of a go routine) in eneque system

Akqa site seems to be super quick now 475ms - real or just the time of night? Doubt it would make a diff for a small site. 

Other sites not impacted.

Need to work out a way to bench this against a huge site using a profiling tool.

#### Other Optimizations
##### Is our approach to queuing work holding us back?
We aren't currently buffering channels, and work around the our circular structure can cause by spinning up a goroutine to push to the channel as soon as it can.

What kind of impact does this have? It sounds bad but does it have an impact?  

It seems to me (maybe incorrectly) that we need to hang onto our pushing goroutines so that we don't deadlock or find a channel we can't push to because all the consumers are busy. Maybe if we used buffered channels we would be a bit more resource efficient though?

What happened when we did this?

None of the sites under test seemed to crawl more quickly when we added a buffer. We left it in because it seems nicer then keeping tons of goroutines waiting but the impact seems virtually non-existent. Maybe some deeper profiling would yield more conclusive results.

#### Others
* Do less work
  * Pay attention to robot.txt
* Store less things  
  * We could probably get away with reducing the amount of things we are tracking and being a bit and still hit our brief
* Stream our results somewhere
  * We are hanging onto all of our results for a big bang tada moment at the end
  * Easy to implement but with really big sites (or if we followed external links) we would eventually exhaust every resource available to us
  * Dumping the results somewhere would probably free up memory really quickly        
* Preload urls
  * If we could get a head start on the pages of a site we could improve how much concurrent work we have on small sites
  * It's almost like this is what sitemaps were invented for...
  * Given the brief this is probably cheating (why not just return the sitemap and goto the pub?)
* Are our log statements blocking?
    * At high concurrency our logs to stdout are actually impacting performance quite noticeably
    * Maybe we could stream the log statements in a non-blocking way or make use of a logging library to reduce the impact of this?

## Thanks to
[Renee French](http://reneefrench.blogspot.com/) for the wonderful gopher icon from [this github repo]( https://github.com/egonelbre/gophers ).