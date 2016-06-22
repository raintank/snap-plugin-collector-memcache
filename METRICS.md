# snap collector plugin - memcache

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description
----------|-----------------------
/raintank/memcache/general/pid                  | Process id of this server process         
/raintank/memcache/general/uptime               | Number of secs since the server started   
/raintank/memcache/general/time                 | current UNIX time according to the server 
/raintank/memcache/general/version              | Version string of this server             
/raintank/memcache/general/pointer_size         | Default size of pointers on the host OS (generally 32 or 64)                      
/raintank/memcache/general/rusage_user          | Accumulated user time for this process (seconds:microseconds)                    
/raintank/memcache/general/rusage_system        | Accumulated system time for this process (seconds:microseconds)                    
/raintank/memcache/general/curr_items            | Current number of items stored            
/raintank/memcache/general/total_items           | Total number of items stored since the server started                        
/raintank/memcache/general/bytes                  | Current number of bytes used to store items                            
/raintank/memcache/general/curr_connections        | Number of open connections                
/raintank/memcache/general/total_connections       | Total number of connections opened since the server started running                
/raintank/memcache/general/rejected_connections     | Conns rejected in maxconns_fast mode      
/raintank/memcache/general/connection_structures     | Number of connection structures allocated by the server                             
/raintank/memcache/general/reserved_fds             | Number of misc fds used internally        
/raintank/memcache/general/cmd_get                  | Cumulative number of retrieval reqs       
/raintank/memcache/general/cmd_set                   | Cumulative number of storage reqs         
/raintank/memcache/general/cmd_flush                 | Cumulative number of flush reqs           
/raintank/memcache/general/cmd_touch                 | Cumulative number of touch reqs           
/raintank/memcache/general/get_hits              | Number of keys that have been requested and found present                         
/raintank/memcache/general/get_misses            | Number of items that have been requested and not found                             
/raintank/memcache/general/get_expired           | Number of items that have been requested but had already expired.                  
/raintank/memcache/general/delete_misses         | Number of deletions reqs for missing keys 
/raintank/memcache/general/delete_hits           | Number of deletion reqs resulting in an item being removed.                    
/raintank/memcache/general/incr_misses           | Number of incr reqs against missing keys. 
/raintank/memcache/general/incr_hits             | Number of successful incr reqs.           
/raintank/memcache/general/decr_misses           | Number of decr reqs against missing keys. 
/raintank/memcache/general/decr_hits             | Number of successful decr reqs.           
/raintank/memcache/general/cas_misses            | Number of CAS reqs against missing keys.  
/raintank/memcache/general/cas_hits              | Number of successful CAS reqs.            
/raintank/memcache/general/cas_badval            | Number of CAS reqs for which a key was found, but the CAS value did not match.   
/raintank/memcache/general/touch_hits            | Numer of keys that have been touched with a new expiration time                     
/raintank/memcache/general/touch_misses          | Numer of items that have been touched and not found                                 
/raintank/memcache/general/auth_cmds             | Number of authentication commands handled, success or failure.              
/raintank/memcache/general/auth_errors           | Number of failed authentications.         
/raintank/memcache/general/evictions             | Number of valid items removed from cache to free memory for new items              
/raintank/memcache/general/reclaimed             | Number of times an entry was stored using memory from an expired entry              
/raintank/memcache/general/bytes_read            | Total number of bytes read by this server from network                              
/raintank/memcache/general/bytes_written         | Total number of bytes sent by this server to network                                
/raintank/memcache/general/limit_maxbytes        | Number of bytes this server is allowed to use for storage.                          
/raintank/memcache/general/listen_disabled_num   | Number of times server has stopped accepting new connections (maxconns).     
/raintank/memcache/general/time_in_listen_disabled_us | Number of microseconds in maxconns.       
/raintank/memcache/general/threads               | Number of worker threads requested. (see doc/threads.txt)                     
/raintank/memcache/general/conn_yields           | Number of times any connection yielded to another due to hitting the -R limit.      
/raintank/memcache/general/hash_power_level      | Current size multiplier for hash table    
/raintank/memcache/general/hash_bytes            | Bytes currently used by hash tables                             
/raintank/memcache/general/expired_unfetched     | Items pulled from LRU that were never touched by get/incr/append/etc before expiring                                  
/raintank/memcache/general/evicted_unfetched     | Items evicted from LRU that were never touched by get/incr/append/etc.           
/raintank/memcache/general/slabs_moved           | Total slab pages moved                    
/raintank/memcache/general/crawler_reclaimed     | Total items freed by LRU Crawler          
/raintank/memcache/general/crawler_items-checked | Total items examined by LRU Crawler       
/raintank/memcache/general/lrutail_reflocked     | Times LRU tail was found with active ref. Items can be evicted to avoid OOM errors. 
/raintank/memcache/general/moves_to_cold         | Items moved from HOT/WARM to COLD LRU's   
/raintank/memcache/general/moves_to_warm         | Items moved from COLD to WARM LRU         
/raintank/memcache/general/moves_within_lru      | Items reshuffled within HOT or WARM LRU's 
/raintank/memcache/general/direct_reclaims       | Times worker threads had to directly reclaim or evict items.                   
/raintank/memcache/general/lru_crawler_starts    | Times an LRU crawler was started          
/raintank/memcache/general/lru_maintainer_juggles | Number of times the LRU bg thread woke up 
/raintank/memcache/general/slab_global_page_pool | Slab pages returned to global pool for reassignment to other slab classes.       
/raintank/memcache/general/slab_reassign_rescues | Items rescued from eviction in page move  
/raintank/memcache/general/slab_reassign_evictions_nomem  | Valid items evicted during a page move (due to no free memory in slab)           
/raintank/memcache/general/slab_reassign_inline_reclaim  | Internal stat counter for when the page mover clears memory from the chunk freelist when it wasn't expecting to.     
/raintank/memcache/general/slab_reassign_busy_items | Items busy during page move, requiring a retry before page can be moved.           
/raintank/memcache/general/log_worker_dropped    | Logs a worker never wrote due to full buf 
/raintank/memcache/general/log_worker_written    | Logs written by a worker, to be picked up 
/raintank/memcache/general/log_watcher_skipped   | Logs not sent to slow watchers.           
/raintank/memcache/general/log_watcher_sent      | Logs written to watchers.                 
/raintank/memcache/settings/maxbytes          | Maximum number of bytes allows in this cache 
/raintank/memcache/settings/maxconns          | Maximum number of clients allowed.           
/raintank/memcache/settings/tcpport           | TCP listen port.                             
/raintank/memcache/settings/udpport           | UDP listen port.                                                       
/raintank/memcache/settings/verbosity         | 0 = none, 1 = some, 2 = lots                 
/raintank/memcache/settings/oldest            | Age of the oldest honored object.                           
/raintank/memcache/settings/umask             | umask for the creation of the domain socket. 
/raintank/memcache/settings/growth_factor     | Chunk size growth factor.                    
/raintank/memcache/settings/chunk_size        | Minimum space allocated for key+value+flags. 
/raintank/memcache/settings/num_threads       | Number of threads (including dispatch).                
/raintank/memcache/settings/reqs_per_event    | Max num IO ops processed within an event.    
/raintank/memcache/settings/tcp_backlog       | TCP listen backlog.                                      
/raintank/memcache/settings/item_size_max     | maximum item size                                        
/raintank/memcache/settings/hashpower_init    | Starting size multiplier for hash table          
/raintank/memcache/settings/lru_crawler_sleep | Microseconds to sleep between LRU crawls     
/raintank/memcache/settings/lru_crawler_tocrawl | Max items to crawl per slab per run                  
/raintank/memcache/settings/hot_lru_pct       | Pct of slab memory reserved for HOT LRU       
/raintank/memcache/settings/warm_lru_pct      | Pct of slab memory reserved for WARM LRU      
/raintank/memcache/slabs/*/chunk_size      | The amount of space each chunk uses. One item will use one chunk of the appropriate size.                       
/raintank/memcache/slabs/*/chunks_per_page | How many chunks exist within one page. A page by default is less than or equal to one megabyte in size. Slabs are allocated by page, then broken into chunks.    
/raintank/memcache/slabs/*/total_pages     | Total number of pages allocated to the slab class.       
/raintank/memcache/slabs/*/total_chunks    | Total number of chunks allocated to the slab class.      
/raintank/memcache/slabs/*/get_hits        | Total number of get requests serviced by this class.     
/raintank/memcache/slabs/*/cmd_set         | Total number of set requests storing data in this class. 
/raintank/memcache/slabs/*/delete_hits     | Total number of successful deletes from this class.      
/raintank/memcache/slabs/*/incr_hits       | Total number of incrs modifying this class.              
/raintank/memcache/slabs/*/decr_hits       | Total number of decrs modifying this class.              
/raintank/memcache/slabs/*/cas_hits        | Total number of CAS commands modifying this class.       
/raintank/memcache/slabs/*/cas_badval      | Total number of CAS commands that failed to modify a value due to a bad CAS id.                               
/raintank/memcache/slabs/*/touch_hits      | Total number of touches serviced by this class.          
/raintank/memcache/slabs/*/used_chunks     | How many chunks have been allocated to items.            
/raintank/memcache/slabs/*/free_chunks     | Chunks not yet allocated to items, or freed via delete.  
/raintank/memcache/slabs/*/free_chunks_end | Number of free chunks at the end of the last allocated page.                                                    
/raintank/memcache/slabs/*/mem_requested   | Number of bytes requested to be stored in this slab[*].  
/raintank/memcache/slabs/total/active_slabs    | Total number of slab classes allocated.                  
/raintank/memcache/slabs/total/total_malloced  | Total amount of memory allocated to slab pages.        
/raintank/memcache/items/*/number             |    Number of items presently stored in this class. Expired items are not automatically excluded.
/raintank/memcache/items/*/number_hot         |    Number of items presently stored in the HOT LRU.
/raintank/memcache/items/*/number_warm        |    Number of items presently stored in the WARM LRU.
/raintank/memcache/items/*/number_cold        |    Number of items presently stored in the COLD LRU.
/raintank/memcache/items/*/number_noexp       |    Number of items presently stored in the NOEXP class.
/raintank/memcache/items/*/age                |    Age of the oldest item in the LRU.
/raintank/memcache/items/*/evicted            |    Number of times an item had to be evicted from the LRU before it expired.
/raintank/memcache/items/*/evicted_nonzero    |    Number of times an item which had an explicit expire time set had to be evicted from the LRU before it expired.
/raintank/memcache/items/*/evicted_time       |    Seconds since the last access for the most recent item evicted from this class. Use this to judge how recently active your evicted data is.
/raintank/memcache/items/*/outofmemory        |    Number of times the underlying slab class was unable to store a new item. This means you are running with -M or an eviction failed.
/raintank/memcache/items/*/tailrepairs        |    Number of times we self-healed a slab with a refcount leak. If this counter is increasing a lot, please report your situation to the developers.
/raintank/memcache/items/*/reclaimed          |    Number of times an entry was stored using memory from an expired entry.
/raintank/memcache/items/*/expired_unfetched  |    Number of expired items reclaimed from the LRU which were never touched after being set.
/raintank/memcache/items/*/evicted_unfetched  |    Number of valid items evicted from the LRU which were  never touched after being set.
/raintank/memcache/items/*/crawler_reclaimed  |    Number of items freed by the LRU Crawler.
/raintank/memcache/items/*/lrutail_reflocked  |    Number of items found to be refcount locked in the LRU tail.
/raintank/memcache/items/*/moves_to_cold      |    Number of items moved from HOT or WARM into COLD.
/raintank/memcache/items/*/moves_to_warm      |    Number of items moved from COLD to WARM.
/raintank/memcache/items/*/moves_within_lru   |    Number of times active items were bumped within HOT or WARM.
/raintank/memcache/items/*/direct_reclaims    |    Number of times worker threads had to directly pull LRU tails to find memory for a new item.


The list of available metrics might vary depending on the Memcache version or the system configuration.