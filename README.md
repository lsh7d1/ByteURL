## Record1
The expiration strategy based on TimeWheel **seriously affects** the qps of `Cache.set`
```bash
redis-benchmark -p 6380 -t set,get -n 10000000 -q -P 512 -c 32

# keep TimeWheel
SET: 846650.25 requests per second
GET: 9944811.00 requests per second

# remove TimeWheel (simply remove the logic in Cache.Set and Cache.Del methods)
SET: 4336714.00 requests per second
GET: 10335736.00 requests per second
```

It is worth mentioning that **during stress testing**, simply using `Cache.Take` 
instead of `Cache.Set` can significantly improve the qps of the set.
Testing environment: Ubuntu 22.04 AMD 5800H 8C32G
```bash
redis-benchmark -p 6380 -t set,get -n 100000000 -q -P 512 -c 32

SET: 9036758.00 requests per second
GET: 10232666.00 requests per second
```

And this is the data for the redis test
```bash
redis-benchmark -p 6379 -t set,get -n 10000000 -q -P 512 -c 32

SET: 2436910.75 requests per second
GET: 2987981.50 requests per second
```