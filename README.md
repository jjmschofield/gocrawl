# gocrawl ![Cute gopher holding a network cable](./docs/network-gopher.png)
gocrawl is a gopher powered web crawler for the internets!

gocrawl will happily make it's way through a website and gather all of the data into a JSONL (line delimited JSON) file for you to do what you wish with.

It will capture a few things: 

* Every internal page it can find
* Every link on those pages, classified as internal/external/file/tel/mailto
* Errors getting those pages
 
It's pretty fast (check out the benchmarks) and capable  and is capable of crawling websites beyond 2.5M pages in size with a 16GB machine over a couple of hours.

gocrawl was created largely as a learning tool for the author, it's only their second attempt at a Go app so be kind!

If your interested in the journey, take a peek at the following:

* [Objective & Approach](./APPROACH.md)
* [A Tale of Web Crawler Optimization](./OPTIMIZATION.md)
* [Solution Discussion]() 

## Getting Started
Now before we run it - a note on politeness. It turned out that this tool is pretty fast, and very good at consuming all of your bandwidth. 

Please be respectful to site owner,  it's possible for you to encounter DDOS defences on sensitive sites. It's unlikely you'll hurt a site with a single instance of this tool from a domestic network link, but there is a remote possibility if the website is very fragile.  

By default the tool will open up 50 concurrent connections which should be plenty fast and shouldn't give most sites any trouble.

First of go get this tool:

```
$ go get https://github.com/jjmschofield/gocrawl
```

Then run it:

```
$ gocrawl -url=https://<your website>
``` 

## Options
Everyone likes options, here you go:
```
$ gocrawl -h
Usage of gocrawl:
  -dir string
        A file path to send results to (default "data")
  -url string
        an absolute url, including protocol and hostname (default "https://monzo.com")
  -workers int
        Number of crawl workers to run (default 50)

``` 

## Benchmarks
* **CPU:** i5-9600K (6 cores) @ ~4.3ghz
* **Mem:** 16GB DDR4 @ 2666ghz
* **Network:** ~80Mbps
* **Disk**: Samsung EVO SSD
* **OS**: Windows 10

In our benchmarks make use of 100 workers against a selection of different sites picked to view performance in slightly different scenarios.

All sites used in the benchmarks had their caches pre-warmed with an initial run.

### [Magic Leap](https://www.magicleap.com) (~297 pages/second)
The [Magic Leap](https://www.magicleap.com) website makes an excellent choice for testing small sites - owing to its structure and performance.
```
 Ran 3 samples:
  runtime:
    Fastest Time: 0.271s
    Slowest Time: 0.386s
    Average Time: 0.318s ± 0.049s
  pages discovered:
    Smallest: 115.000
     Largest: 115.000
     Average: 115.000 ± 0.000
  pages crawled:
    Smallest: 115.000
     Largest: 115.000
     Average: 115.000 ± 0.000
  queue peak:
    Smallest: 2.000
     Largest: 2.000
     Average: 2.000 ± 0.000
  parallel scrape peak:
    Smallest: 67.000
     Largest: 71.000
     Average: 68.333 ± 1.886
------------------------------
```                                   

### [AKQA](https://www.akqa.com) (~240 pages/second)
The [AKQA](https://www.akqa.com) website is an interesting site, as on one hand it is small and highly performant - yet on the other hand it features a lot of https -> http internal links which helped us during [optimization](./OPTIMIZATION.md) make some gains by being lazy **cough** more correct in what we crawl.
```
 Ran 3 samples:
  runtime:
    Fastest Time: 0.425s
    Slowest Time: 1.062s
    Average Time: 0.735s ± 0.260s
  pages discovered:
    Smallest: 177.000
     Largest: 177.000
     Average: 177.000 ± 0.000
  pages crawled:
    Smallest: 177.000
     Largest: 177.000
     Average: 177.000 ± 0.000
  queue peak:
    Smallest: 2.000
     Largest: 2.000
     Average: 2.000 ± 0.000
  parallel scrape peak:
    Smallest: 15.000
     Largest: 25.000
     Average: 20.333 ± 4.110
------------------------------
```

### [Monzo](https://www.monzo.com) (~33 pages / second)
Our middling size benchmark [Monzo](https://www.monzo.com) performs rather badly in our testing - despite being of a size to actually cap out our concurrent connections.

Entertainingly this isn't really due that much to the performance of the site, or the performance of the crawler. It really is down to [Kate Hollowood's blog](http://monzo.com/blog/authors/kate-hollowood/11)! 

Kate's blog posts are linked together by single page pagination only, and we always seem to hit it late in the crawl. The net result is that towards the end of the crawl we loose all parallelism as `gocrawl` slowly makes its way through one page at a time.

Of all the sites on test, this site would probably benefit most from making use of a `sitemap.xml` to preload the crawler with a bunch of links helping us maintain consistent parallelism.  

```
 Ran 3 samples:
  runtime:
    Fastest Time: 19.587s
    Slowest Time: 20.167s
    Average Time: 19.930s ± 0.249s
  pages discovered:
    Smallest: 665.000
     Largest: 665.000
     Average: 665.000 ± 0.000
  pages crawled:
    Smallest: 665.000
     Largest: 665.000
     Average: 665.000 ± 0.000
  queue peak:
    Smallest: 56.000
     Largest: 66.000
     Average: 62.000 ± 4.320
  parallel scrape peak:
    Smallest: 100.000
     Largest: 100.000
     Average: 100.000 ± 0.000
------------------------------
```             
                                            

### [Citizens Advice Scotland](https://www.cas.org.uk) (351 pages / second)
The first of our larger sites, the [CAS website](https://www.cas.org.uk) is interesting for a few reasons.

Firstly, it is really slow but has a very fast cache when warmed up - not much we can do about this but it's entertaining to watch the site perform so much better when hot!

Secondly, the site features a large number of downloadable documents which, we discovered filtering out made a huge difference (93-96%) during [optimization](./OPTIMIZATION.md).

Finally, the site has a very variable number of pages discovered - more investigation is required on this. Is it a very dynamic site? Are we finding broken pages? Do we have some race condition? Only questions - no answers right now...    
 
```
 Ran 3 samples:
  runtime:
    Fastest Time: 6.580s
    Slowest Time: 8.310s
    Average Time: 7.362s ± 0.716s
  pages discovered:
    Smallest: 2553.000
     Largest: 2567.000
     Average: 2558.000 ± 6.377
  pages crawled:
    Smallest: 2553.000
     Largest: 2567.000
     Average: 2558.000 ± 6.377
  queue peak:
    Smallest: 800.000
     Largest: 916.000
     Average: 870.667 ± 50.632
  parallel scrape peak:
    Smallest: 100.000
     Largest: 100.000
     Average: 100.000 ± 0.000
------------------------------
```

# [Golang](https://golang.org) (533 pages/second)
Finally the site that really makes `gocrawl` shine is non other then the official [Go](https://golang.org) website. 

The largest site in the test at 13.5k pages crawls blisteringly quickly to the point were I'm actually a bit afraid.

```
 Ran 3 samples:
  runtime:
    Fastest Time: 25.182s
    Slowest Time: 25.529s
    Average Time: 25.334s ± 0.145s
  pages discovered:
    Smallest: 13514.000
     Largest: 13516.000
     Average: 13515.000 ± 0.816
  pages crawled:
    Smallest: 13514.000
     Largest: 13516.000
     Average: 13515.000 ± 0.816
  queue peak:
    Smallest: 6658.000
     Largest: 7028.000
     Average: 6853.000 ± 151.712
  parallel scrape peak:
    Smallest: 100.000
     Largest: 100.000
     Average: 100.000 ± 0.000
------------------------------
```

## Contributing
## Tests
You can run the unit tests with coverage with:
```
$ go test ./internal/... -cover
ok      github.com/jjmschofield/gocrawl/internal/caches (cached)        coverage: 100.0% of statements
ok      github.com/jjmschofield/gocrawl/internal/counters       (cached)        coverage: 100.0% of statements
ok      github.com/jjmschofield/gocrawl/internal/crawl  (cached)        coverage: 96.8% of statements
ok      github.com/jjmschofield/gocrawl/internal/links  (cached)        coverage: 100.0% of statements
ok      github.com/jjmschofield/gocrawl/internal/md5    (cached)        coverage: 100.0% of statements
ok      github.com/jjmschofield/gocrawl/internal/pages  (cached)        coverage: 100.0% of statements
ok      github.com/jjmschofield/gocrawl/internal/scrape (cached)        coverage: 0.0% of statements [no tests to run]
ok      github.com/jjmschofield/gocrawl/internal/writers        (cached)        coverage: 0.0% of statements
```

There are some coverage gaps there - so feel free to help raise it if you wish :)

If you want to run the benchmark tests, do:
```
$ go test
```
This will execute 3 runs against the default website to benchmark the tool to give you a flavor for speed
                    
                    




## Thanks to
[Renee French](http://reneefrench.blogspot.com/) for the wonderful gopher icon from [this github repo](https://github.com/egonelbre/gophers).

