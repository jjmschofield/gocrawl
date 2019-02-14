# gocrawl ![Cute gopher holding a network cable](docs/network-gopher.png)
`gocrawl` is a gopher powered web crawler for the internets!

`gocrawl` will happily make it's way through a website and gather all of the data into a JSONL (line delimited JSON) file for you to do with as you wish.

It will capture a few things: 

* Every internal page it can find
* Every link on those pages, classified as internal/external/file/tel/mailto
* Errors getting those pages
 
It's pretty fast (check out the benchmarks) and capable of crawling websites beyond 2.5M pages in size with a 16GB machine over a couple of hours.

`gocrawl` was created largely as a learning tool for the author, it's really my first attempt at production quality Go application - so no doubt it could do with a review or two!

If your interested in the journey, take a peek at the following:

* [Objective & Approach](./APPROACH.md)
* [A Tale of Web Crawler Optimization](./OPTIMIZATION.md)

## Getting Started
Now before we run it - a note on politeness. It turned out that this tool is pretty fast, and very good at consuming all of your bandwidth. 

**Please be respectful to site owner**. It's possible for you to encounter DDOS defences on sensitive sites but it's unlikely you'll hurt a site with a single instance of this tool from a domestic network link - unless the website is very fragile.  

By default the tool will open up 50 concurrent connections which should be plenty fast and shouldn't give most sites any trouble.

First of go get this tool:

```
$ go get https://github.com/jjmschofield/gocrawl
```

Then run it:

```
$ gocrawl -url=https://<your website>
``` 

### Options
Everyone likes to have options, here are yours:
```
$ gocrawl -h
-dir string
        A relative file path to send results to (default "data")
  -redis string
        An optional redis address to make use of redis rather then in memory queues and caches eg: localhost:6379
  -url string
        an absolute url, including protocol and hostname (default "https://monzo.com")
  -workers int
        Number of crawl workers to run (default 50)
```

## Results
You will get results in a JSONL file called `pages.jsonl`. If you didn't specify `-dir` this will be in `./data/` relative to wherever you ran `gocrawl`

A single page looks a bit like:
```
{
  "id": "81146cd6051bbfc62fe29dac5b078e11",
  "url": "https://monzo.com",
  "outPages": {
    "internal": {
      "1242e724ade9b927626e0160b06292a2": "https://monzo.com/press",
      "28ea4e4813c2a033e4e0971a50cee585": "https://monzo.com/legal/terms-and-conditions",
      "3456cc237ec3a4e5bc13f42cdb2670a2": "https://monzo.com/download",
      "563051fd7215aa58591a05afd3bd6c3e": "https://monzo.com/careers",
      "5a54d0d0e550a83a534c3cee45021b41": "https://monzo.com/-play-store-redirect",
      "5d82719fab0026f709c46acf892c2f32": "https://monzo.com/community",
      "70d7b336b5c9b1a879913b322d0eb0ea": "https://monzo.com/features/overdrafts",
      "7be3d6b130f2911636007962c08ba57d": "https://monzo.com/blog",
      "81146cd6051bbfc62fe29dac5b078e11": "https://monzo.com",
      "910e1ef1fb9a8636501269edd8deffc5": "https://monzo.com/legal/privacy-policy",
      "9730974b7ed77209da6ab00fbd3332e9": "https://monzo.com/legal/cookie-policy",
      "a2b33f749fedc804add1f7e2e834a3d0": "https://monzo.com/cdn-cgi/l/email-protection",
      "ae4e2c3e2419f7784d1c66d49f87d4cc": "https://monzo.com/features/travel",
      "bcb6ea59d7d8f64730640addfe8bb8c3": "https://monzo.com/about",
      "be93233b884615d2c31fba361565459b": "https://monzo.com/features/switch",
      "c5384afe3366bfdbe3111a09fac556b6": "https://monzo.com/legal/fscs-information",
      "d337313b51b94e92a7e1bd4e9ba2e2f3": "https://monzo.com/tone-of-voice",
      "d39543e6b964b23ff552d9d6d478643c": "https://monzo.com/transparency",
      "d3aada3f1e292a1c951ee77a9bb8a254": "https://monzo.com/community/making-monzo",
      "d65f5efe28855c0ab6aa546c6ede407b": "https://monzo.com/blog/how-money-works",
      "decce6152e5ae74e46254fe14b24081d": "https://monzo.com/faq",
      "ecca2d96ec11476c91aa2d2b50404c37": "https://monzo.com/features/google-pay",
      "fbb1d7a29f1dd568760180db067ea520": "https://monzo.com/features/apple-pay"
    }
  },
  "outLinks": {
    "05ef754dd277a008bc6dfdba10a65341": {
      "id": "05ef754dd277a008bc6dfdba10a65341",
      "toUrl": "https://twitter.com/monzo",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "079661da7b74461c5e19945d415bb57e": {
      "id": "079661da7b74461c5e19945d415bb57e",
      "toUrl": "https://www.youtube.com/monzobank",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "1242e724ade9b927626e0160b06292a2": {
      "id": "1242e724ade9b927626e0160b06292a2",
      "toUrl": "https://monzo.com/press",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "2545906c28821b13c74d484b824ac20c": {
      "id": "2545906c28821b13c74d484b824ac20c",
      "toUrl": "https://www.linkedin.com/company/monzo-bank",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "28ea4e4813c2a033e4e0971a50cee585": {
      "id": "28ea4e4813c2a033e4e0971a50cee585",
      "toUrl": "https://monzo.com/legal/terms-and-conditions",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "3456cc237ec3a4e5bc13f42cdb2670a2": {
      "id": "3456cc237ec3a4e5bc13f42cdb2670a2",
      "toUrl": "https://monzo.com/download",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "563051fd7215aa58591a05afd3bd6c3e": {
      "id": "563051fd7215aa58591a05afd3bd6c3e",
      "toUrl": "https://monzo.com/careers",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "5a54d0d0e550a83a534c3cee45021b41": {
      "id": "5a54d0d0e550a83a534c3cee45021b41",
      "toUrl": "https://monzo.com/-play-store-redirect",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "5d82719fab0026f709c46acf892c2f32": {
      "id": "5d82719fab0026f709c46acf892c2f32",
      "toUrl": "https://monzo.com/community",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "635db89b9ce7485b2e19458a7af7b9ab": {
      "id": "635db89b9ce7485b2e19458a7af7b9ab",
      "toUrl": "https://web.monzo.com",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "6c051e3b66551d656f56794784870fa7": {
      "id": "6c051e3b66551d656f56794784870fa7",
      "toUrl": "https://www.thetimes.co.uk/article/tom-blomfield-the-man-who-made-monzo-g8z59dr8n",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "6deb00a2eae3bcf79e5989d4b55d99b1": {
      "id": "6deb00a2eae3bcf79e5989d4b55d99b1",
      "toUrl": "https://monzo.com/cdn-cgi/l/email-protection#7018151c00301d1f1e0a1f5e131f1d",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "70adfde50f4666cfda9eb4ac800dea64": {
      "id": "70adfde50f4666cfda9eb4ac800dea64",
      "toUrl": "https://www.instagram.com/monzo",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "70d7b336b5c9b1a879913b322d0eb0ea": {
      "id": "70d7b336b5c9b1a879913b322d0eb0ea",
      "toUrl": "https://monzo.com/features/overdrafts",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "7be3d6b130f2911636007962c08ba57d": {
      "id": "7be3d6b130f2911636007962c08ba57d",
      "toUrl": "https://monzo.com/blog",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "823c8ed14c5e75c6464f590c2fd0edbd": {
      "id": "823c8ed14c5e75c6464f590c2fd0edbd",
      "toUrl": "https://www.telegraph.co.uk/personal-banking/current-accounts/monzo-atom-revolut-starling-everything-need-know-digital-banks/",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "8c6672a993504cbebbe653e547cd5d09": {
      "id": "8c6672a993504cbebbe653e547cd5d09",
      "toUrl": "https://www.theguardian.com/technology/2017/dec/17/monzo-facebook-of-banking",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "910e1ef1fb9a8636501269edd8deffc5": {
      "id": "910e1ef1fb9a8636501269edd8deffc5",
      "toUrl": "https://monzo.com/legal/privacy-policy",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "9730974b7ed77209da6ab00fbd3332e9": {
      "id": "9730974b7ed77209da6ab00fbd3332e9",
      "toUrl": "https://monzo.com/legal/cookie-policy",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "a9d6b806b24f7d62ea76504a1eca9436": {
      "id": "a9d6b806b24f7d62ea76504a1eca9436",
      "toUrl": "https://monzo.com/cdn-cgi/l/email-protection#0f676a637f4f6260617560216c6062",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "ae4e2c3e2419f7784d1c66d49f87d4cc": {
      "id": "ae4e2c3e2419f7784d1c66d49f87d4cc",
      "toUrl": "https://monzo.com/features/travel",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "b2af5f82028b38bb897547708d7b7a67": {
      "id": "b2af5f82028b38bb897547708d7b7a67",
      "toUrl": "https://www.facebook.com/monzobank",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "b7d9de4bac141413c1901fd690de35f6": {
      "id": "b7d9de4bac141413c1901fd690de35f6",
      "toUrl": "https://monzo.com/",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "bcb6ea59d7d8f64730640addfe8bb8c3": {
      "id": "bcb6ea59d7d8f64730640addfe8bb8c3",
      "toUrl": "https://monzo.com/about",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "be93233b884615d2c31fba361565459b": {
      "id": "be93233b884615d2c31fba361565459b",
      "toUrl": "https://monzo.com/features/switch",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "c5384afe3366bfdbe3111a09fac556b6": {
      "id": "c5384afe3366bfdbe3111a09fac556b6",
      "toUrl": "https://monzo.com/legal/fscs-information",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "c7da1ee3bee5489ac02cd81fcbcf1d4a": {
      "id": "c7da1ee3bee5489ac02cd81fcbcf1d4a",
      "toUrl": "https://monzo.com/community",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "d337313b51b94e92a7e1bd4e9ba2e2f3": {
      "id": "d337313b51b94e92a7e1bd4e9ba2e2f3",
      "toUrl": "https://monzo.com/tone-of-voice",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "d39543e6b964b23ff552d9d6d478643c": {
      "id": "d39543e6b964b23ff552d9d6d478643c",
      "toUrl": "https://monzo.com/transparency",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "d3aada3f1e292a1c951ee77a9bb8a254": {
      "id": "d3aada3f1e292a1c951ee77a9bb8a254",
      "toUrl": "https://monzo.com/community/making-monzo",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "d65f5efe28855c0ab6aa546c6ede407b": {
      "id": "d65f5efe28855c0ab6aa546c6ede407b",
      "toUrl": "https://monzo.com/blog/how-money-works",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "decce6152e5ae74e46254fe14b24081d": {
      "id": "decce6152e5ae74e46254fe14b24081d",
      "toUrl": "https://monzo.com/faq",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "dfc645d45e30aaf47f9a012daed000c4": {
      "id": "dfc645d45e30aaf47f9a012daed000c4",
      "toUrl": "https://itunes.apple.com/gb/app/mondo/id1052238659",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "ecca2d96ec11476c91aa2d2b50404c37": {
      "id": "ecca2d96ec11476c91aa2d2b50404c37",
      "toUrl": "https://monzo.com/features/google-pay",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    },
    "ecf1bfbf83d9e9dd15c4f1cd7710eec1": {
      "id": "ecf1bfbf83d9e9dd15c4f1cd7710eec1",
      "toUrl": "https://www.standard.co.uk/tech/monzo-prepaid-card-current-accounts-challenger-bank-a3805761.html",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "ed415fae124bd226ed8ab99270e31b00": {
      "id": "ed415fae124bd226ed8ab99270e31b00",
      "toUrl": "https://www.fscs.org.uk/",
      "fromUrl": "https://monzo.com",
      "type": "external"
    },
    "fbb1d7a29f1dd568760180db067ea520": {
      "id": "fbb1d7a29f1dd568760180db067ea520",
      "toUrl": "https://monzo.com/features/apple-pay",
      "fromUrl": "https://monzo.com",
      "type": "internal"
    }
  },
  "error": null
}
```

There are a lot of things you can do with this data, for example here is a directed graph for the whole of the [Monzo](https://monzo.com) website.

![Confusing circular graph](docs/monzomap.png)

This graph shows every relation between pages, it's pretty hard to read but it demonstrates some of the power of the data.

If you want to get started with your own one of these there is a [demo](demo)(for small sites) you can spin up by just opening the HTML page. Just run `gocrawl` to populate with data eg:

```
$ gocrawl -dir=./demo/data -url=https://www.akqa.com
``` 

## Things gocrawl doesn't do but should...
* Try and grab a `sitemap.xml` and enqueue everything in it to get started
  * This would help us with small or badly linked sites (from our point of view) by helping us achieve better parallelization
* Pay attention to `robots.txt` 
  * This would help us avoid doing more work then we need to and would also be a bit more respectful to site owners   

## I heard you wanted to go bigger...
The main limiting factor for the size of a site is memory - currently we keep a track of where we have been and where we are going inside in memory caches.

You'll find that the [PageCrawler](internal/crawl/crawler.go) lets you inject in these cashes - satisfy the [ThreadSafeCache interface](internal/caches/thread_safe.go) and swap them out for something backed by disk or magic cloud memory and you'll go further.

By default we have provided an example of a redis cache (unauthenticated) which is available from the command line. You'll need a redis server for this - the included `docker-compose.yml` can get you started.

## I heard you wanted to query lots of data...
`gocrawl` itself intentionally doesn't waste resources on how you may want to query the data later.

Our JSONL format is pretty handy - even for big files we should hopefully be able to go line by line and stream it into something else. 

You can take this approach by creating a parser, but if you want to go really big your can create your own [Writer](internal/writers/writer.go) which simply ranges over the out channel of [PageCrawler](internal/crawl/crawler.go). The crawler fans in at this stage - so it's up to you whether you want the writer to block the crawl. 

You **should** only get unique pages once - but it's not guaranteed. ID's are deterministic (md5 hashes of normalized page urls) so with a huge data set you may well get a collision. If you do, or you have a better idea for a fast hashing algo that gives us a relatively small and consistently sized has let me know (I'm sure there are many)!
  
## Benchmarks
* **CPU:** i5-9600K (6 cores) @ ~4.3ghz
* **Mem:** 16GB DDR4 @ 2666ghz
* **Network:** ~80Mbps
* **Disk**: Samsung EVO SSD
* **OS**: Windows 10

In our benchmarks we make use of 100 workers against a selection of different sites. The sites are all picked to analyze our performance in slightly different scenarios and were used during the [optimization](./OPTIMIZATION.md) of `gocrawl`.

All sites used in the benchmarks had their caches pre-warmed with an initial run before running 3 samples in quick succession.

### [Magic Leap](https://www.magicleap.com) (~297 pages/second)
The [Magic Leap](https://www.magicleap.com) website makes an excellent choice for testing our performance against small sites - owing to its well connected structure and excellent performance.
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

Kate's blog posts are linked together by single page pagination only, which we always seem to hit late in the crawl. The net result is that towards the end of the crawl we loose all parallelism as `gocrawl` slowly makes its way through one page at a time.

Of all the sites on test, this is pointing to the most important furture [optimization](./OPTIMIZATION.md), making use of a `sitemap.xml` to preload the crawler with a bunch of links which will help us achieve more consistent parallelism with elongated site structures.  

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
                                            

### [Citizens Advice Scotland](https://www.cas.org.uk) (~351 pages / second)
The first of our larger sites, the [CAS website](https://www.cas.org.uk) is interesting for a few reasons.

Firstly, it is really slow but has a very fast cache when warmed up - not much we can do about this but it's entertaining to watch the site perform so much better when hot!

Secondly, the site features a large number of downloadable documents which, we discovered filtering out made a huge difference (93-96%) during [optimization](./OPTIMIZATION.md).

Finally, the site has a very variable number of pages discovered - more investigation is required on this. Is it a very dynamic site? Are we finding broken pages? Do we have some race condition? Only questions right now...    
 
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

### [Golang](https://golang.org) (~533 pages/second)
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
### Tests
We use the [ginkgo](https://github.com/onsi/ginkgo) test framework to achieve BDD tests of a similar vein to Jest / Jasmine / Mocha. 

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

We've also used [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) to provide fake generation. I'm still in two minds as to whether or not this is helping or hurting.

If you get failing tests after changing an interface regenerate the fakes with the following:

```
$ go generate ./...
```

If you want to run the benchmark tests, do:
```
$ go test
```
This will execute 3 runs against the default website to benchmark the tool outputting the metrics above to give you a flavor for speed. If you want profiling do:
```
$ go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
```
                    
                    
## Thanks to
[Renee French](http://reneefrench.blogspot.com/) for the wonderful gopher icon from [this github repo](https://github.com/egonelbre/gophers).

