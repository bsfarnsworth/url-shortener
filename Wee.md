# Wee

**A URL shortener application**

*Author’s note: I’m choosing the “we” narrative voice because I’m viewing this write up as a presentation to a team that includes the reader and myself. 
I chose the product name “Wee” simply because it’s the shortest word I could think of to mean shorter. I don’t mind that it also injects some whimsy.*

## Design

### Usage Model

**Operations**

* *shorten* - Create a wee URL from a given full URL
* *\<weeUrl\>* - Lookup a wee URL and follow its full URL
* *lengthen* - Lookup a wee URL to see its full URL 
* *retire* - Delete an existing wee URL

As spelled out in the assignment, the application's REST service will provide entry points to create new, shortened URLs and to follow those URLs to their destinations.  

We'll propose two additions: 

* *lengthen*, the capability to view the expansion of a shortened URL without redirecting to it; this will allow testing of the entire creation and lookup mechanism without incurring any of the actions involved in the redirect; and

* *retire*, for the URL owner to modify it.  We won't go so far as to allow editing or updates, but we'll make it possible to delete an existing URL; the user's update would then be to simply create a new one.  This allows for correction of mistakes and for removal of obsolete or dead links. We'll call this operation retirement.

With retirement there's a possibility of misuse: rivals or miscreants might want to delete URLs.  Protecting the modification capability could be complicated if it involved authentication.  We'll simplify it by issuing a token along with the wee URL.  The user will be directed to keep it in a safe place.  It won't expire, and when the user wishes to retire their URL they simply send it to the retirement endpoint.  Since the token itself requires no user convenience factor, we can make it arbitrarily large to reduce the likelihood of discovery by an attacker.

### Method of Redirect

The method of redirection could be a coded action performed by the browser, or it could make use of the HTTP redirect response.  We'll choose the latter because it complies with intended web design, and is uncomplicated to use:

https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections

The type of redirection that makes sense in our case is a *temporary* one, with method and body unchanged. Between codes 302 and 307 we'll choose 307 based on the advice that [the behavior with 307 is predictable](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/307).

### Construction of the Shortened URL

An appealing idea would be to construct it from the full URL in a reversible manner using only the text of the URL itself.  This would make it possible to shrink and expand the URLs without reference to a database.  A string compression scheme could be used.  The reality, however, is that even though a URL requires roughly only 6 bits of entropy per character, any encoding of the shrunken URL will occupy 8 bits per compressed character since the shortened URL will have to be represented in the same encoding (e.g. we can't use gzip binary output in the URL space); thus, there's no way to meaningfully leverage that unused space for compression.  Likewise, other text compression techniques without text-specific dictionaries rarely achieve better than 50% compression.  So, appealing or not, this approach wouldn't achieve significant shortening.

Instead we'll use the common approach of creating a tag that can be used to lookup the full URL in a hash table.  

#### Tagging

There are many possible approaches to creating unique tags: sequential characterization ('1' .. '9' .. '1A' .. 'zzzzzz' ..), license plate formatting, fixed-width labels, word groups ("yayPluto") to list a few.  A choice would depend on whether the application is a commercial product, an aid for sharing, or other use cases, and whether the desired result is simply to be as short as possible, customizable, legible, memorable, or other factors.  

For this application we'll (arbritrarily and maybe ridicuously) decide that our namespace should be large enough to allow one wee URL for every person on Earth, to encode them so that they can be unambiguously read to another person without the need to identify upper/lower case or puncuation, and to limit the character set to ASCII a-z.  That gives us a symbol set of 26 alpha + 10 digits.  We'll further decide that they will all have the same width.  Assuming a conservative maximum population of at least [10G](https://www.livescience.com/16493-people-planet-earth-support.html) (yikes), we need to use 7 characters in our tags (36\*\*6 is only 2G, 36\*\*7 gets us to 70G).

#### Hashing

Our design now has a potential capacity of 70G entries, indexed on our fixed-width, 7 character tag, and having variable-length, string type value fields to store the full URLs.

#### Macro Data

One other item to store with each entry is its retirement token.  Since we've decided these can be large and we want them to effectively act as be secure keys, we'll choose to implement them as UUIDs.

#### Storage

Required attributes:

* Persistence: yes
* Performance: capable of redirection with no discernable delay
* Access locking: we can reasonably assume that tag reads will be far more frequent than writes, so locking can be performed on the entire table and only when updates are needed
* Space per record (bytes): 
    * tag: 7
    * uuid: 16
    * longURL: N

The amount of storage per entry is the big variable because maximum N is [not well defined](https://www.rfc-editor.org/rfc/rfc9110.html#section-4.1).

Because the whole intent of a URL shortener is to make it simple to handle bulky URLs, we'll be generous and support the bulkiest.  The reference recommends that all "senders and recipients" (of which wee qualifies in both categories) "support, at a minimum, URIs with lengths of 8000 octets".  (The length requirements are more involved than this paraphrase. The section on "defensively" considering [length requirements](https://www.rfc-editor.org/rfc/rfc9110.html#section-2.3) contains more important discussion, but we'll consider that to be beyond the scope of this project.)

If we make the assumption for our product that very long URLs are far more likely than the concern of billions of users then we can choose reasonable deployment limits while keeping open the possibility of surprises.  

We'll choose limits for this deployment that are aspirational but also fit into the constraints of our deployment host (free tier!) and database choices (we're NOT calling for big data systems).  

These limits should be defined as configuration values.


### Web Service

#### Service URL

'wee' designates our service, and the domain name belongs to the deployment service.

  `wee.fly.dev/`

Normal usage -- for redirection -- looks like:

  `https://wee.fly.dev/<weeURL>`
  
#### APIs
All are appended to the Service URL

##### Summary

*  `/shorten/<fullURL>`
*  `/lengthen/<weeURL>`
*  `/retire/<token>`
*  `/<weeURL>`
  
