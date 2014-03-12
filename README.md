# GoHLL

[![build status](https://travis-ci.org/mynameisfiber/gohll.png?branch=master)](https://travis-ci.org/mynameisfiber/gohll)
[![GoDoc](https://godoc.org/github.com/mynameisfiber/gohll?status.png)](https://godoc.org/github.com/mynameisfiber/gohll)

> [HLL++][1] for gophers

## What is this?

Have you ever had a large set of data (or maybe even a never ending stream of
data) and wanted to know how many unique items there were?  Or maybe you had
two sets of data, and you wanted to know how many unique items there were in
the two sets combined?  Or maybe how many items appeared in both datasets?
Well, `gohll` is for you!

HLL is a probabilistic counting algorithm that can tell you how many unique
items you have added to it.  In addition, you can perform union and
intersection operations between multiple HLL objects.  It's easy!  Let me show you:

```
// First we make an HLL with an error rate of ~0.1%
h := NewHLLByError(0.001)

// Now it's time to start adding things to it!
for i := 0; i < 100000; i += 1 {
    h.Add(fmt.Sprintf("%d", math.Uint32())
}

uniqueitems := h.Cardinality()
```

Wow! That was so easy!  But wait a second, you may be saying... What about
those set operations you were talking about?  Well, that can be done quite
easily as well!

```
// let's make 2 hll's... they must have the same error rate!
h1 := NewHLLByError(0.001)
h2 := NewHLLByError(0.001)

// now let's add different things to each one
for i := 0; i < 100000; i += 1 {
    h1.Add(fmt.Sprintf("%d", math.Uint32())
    h2.Add(fmt.Sprintf("%d", math.Uint32())
}

uniqueItemsH1 := h1.Cardinality() // |h1|
uniqueItemsH2 := h2.Cardinality() // |h2|

uniqueItemsEither := h1.CardinalityUnion(h2) // |h1 u h2|
uniqueItemsBoth   := h1.CardinalityIntersection(h2) // |h1 n h2|

h3 := h1.Union(h2)
```

In this example, all the `Cardinality*` queries return a `float64` with the
size of the set under that operation.  That is to say, the result of
`h1.CardinalityUnion(h2)` is the number of unique items in either h1 and h2.
So, if h1 and h2 both only contain the item "foo", then the cardinality of the
union is 1 -- there is only one unique item between them.  The intersection
finds items that exist in both sets.  Finally, the `h1.Union(h2)` call creates
a new HLL that represents both sets h1 and h2 unioned together.

NOTE: Intersections are not natively supported in HLL's so we simply use the
inclusion–exclusion principle which has completely different error bounds than
any other operation on the HLL (generally much worse)

## Can't I just use a `map[string]bool` object to do that?

Sure, you could.  But I could also try to keep track of the unique items in
your set in my head but that also wouldn't work out too well.  The problem with
using a `map[string]bool` object is the amount of space it takes... if you had
1e10 unique items, each represented by a 10 digit string, you are looking at
480GB of memory being used.  On the other hand, an HLL++ with an error rate of
0.01% would have a memory footprint of 1.06GB and would stay at that size even
if you keep adding more items!

This does come with some conditions though, you don't know what the original
items are -- you can simply see how many unique items there are and do
intersection and union operations with other sets.  Also, there is an error
bound on the results (the lower error you set, the more memory the hll uses),
however the trade-offs are quite reasonable.

So basically the question you should be asking yourself is: what questions do I
need to answer with my data?  Also, is it worthwhile to need potentially vast
amounts of memory in order to get exact solutions?  Typically the answer is no.

## HLL vs HLL++

I've been throwing around the words HLL and HLL++ as if they were the same
thing.  Let's talk a bit about how they are different.

HLL++ is an extension to HLL (first talked about in [this][1] paper) that gives
it better biasing properties and _much_ better error rates for small set sizes
without increasing memory usage.  The biasing issue is addressed by some
experiments that were run that gave quantitative numbers as to how the HLL's
were being biased for different values.  With this knowledge, we are able to
adjust for the biasing effects (this is done in the `EstimateBias` function).

On the other hand, for small set sizes HLL++ uses a smart way of encoding
integers to create a miniature HLL with much higher precision.  HLL's have a
nice property of doing better when you give it more data, however this
miniature HLL (the `SparseList` in our implementation) is designed such that it
gives very low errors in this regime (giving errors in the range of 0.018%).
In addition, this list _could_ be compressed easily to allow us to use this
encoding much longer.  Once enough items have been placed into the HLL, the
integer encoding is reversed and we insert the old data into a classic HLL
structure.

## Speed

This library is fast!  With an error rate of __0.1%__ (ie: `p=20`), while in
the sparse regime (ie: small number of items compared to the error rate
chosen), we can accomplish about __269,000 insertions per second__ on a 2011
Macbook Air.  For the normal regime (ie: large number of items) we accomplish
__1,880,000 insertions per second__ on the same hardware!  Furthermore,
cardinality queries (ie: asking "how many unique elements are in this HLL?")
are quite quick.  While they depend on many subtleties regarding the state of
the HLL, they have an upper bound of __7ms__ per query (average of __2ms__) for
the same setup.

It may seem that strange that it is slower to add items to the HLL while it is
less full, but this makes sense since we go to extra lengths and have a
different insertion model when there aren't many items in order to fulfill the
error guarantees.

If you care more about insertion speed than you do having good error bounds
when the HLL is still relatively empty, simply call `h.ToNormal()` on the HLL
once you have instantiated it in order to skip the sparse phase.  This would be
desirable if you are pre-loading the HLL with a lot of data and know that once
the loading is done, it will be out of the sparse phase anyways (so you may as
well get 4x faster insertion speeds to that your loading procedure finishes
faster!).

Benchmarks can be run with `go test --bench=.`

## Hashing functions

"Isn't the entropy of your hashing function very important for the error
estimates?" you may be asking.  Why, yes it is!  We choose to use murmurhash3's
64bit hashing function by default, but this can be changed.  Simply create a
new hashing function that takes in a string and outputs a `uint64` and set your
HLL object's `Hasher` property to it.  This should only be done before you have
inserted any items into the object!

This is quite useful if you know a priori some properties of the data you will
insert and can thus pick a more appropriate hashing function.

## Resources

* [Original Paper][1]
* [Great blog post on HLL++][2]

[1]: http://static.googleusercontent.com/external_content/untrusted_dlcp/research.google.com/en/us/pubs/archive/40671.pdf
[2]: http://blog.aggregateknowledge.com/2013/01/24/hyperloglog-googles-take-on-engineering-hll/
