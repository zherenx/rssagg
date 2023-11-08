# RSS Aggregator

## Some Notes

This project is part of the freeCodeCamp.org x boot.dev GoLang crash course by Lane Wagner. Personally, this project serves as a gateway for me to learn and practice the Go programming language, and I also take it as an opportunity to refresh my SQL knowledge.

This RSS Aggregator is a backend server that aggregates data from RSS feeds. The server allows users to add different RSS feeds to its database, and then, it will automatically collect all of the posts from those feeds (download and store them in the database) so that users can view them later. For example, a user can request to get the latest posts from feeds that the user followed.

While working on the project, I also consulted a diverse range of resources such as books, blog posts, and other materials (for links to some, see reference section) covering GoLang syntax, coding conventions, design patterns, and more. Whenever I encountered intriguing new ideas related to Go or faced confusion, these resources provided valuable insights, and some modifications and improvements were made to the codebase accordingly. As a result, the code may differ slightly from the tutorial at times. Also, please note that some non-documentation comments and notes were left in the code and those are just for my personal reference.

## Some of the Tools Involved

- Go (programming language)
- SQL
- PostgreSQL
- sqlc
- goose
- chi (for Go HTTP services)

## Potential Improvements/Extensions

- front-end
- apikey handling improvement (e.g. storing only the hash value in the database)
- logging, and implement a logger
- robustness improvement (e.g. different/more date formats should be supported)
- refactoring of some functions/duplicated code
- etc.

## References and Some Useful Links

- freeCodeCamp.org x boot.dev GoLang crash course by Lane Wagner, <https://www.youtube.com/watch?v=un6ZyFkqFKo>
- Go 101 book, by Tapir Liu, <https://go101.org/article/101.html>
- Go Generics 101 book, by Tapir Liu, <https://go101.org/generics/101.html>
- Useful posts, etc. on middleware pattern
  - <https://drstearns.github.io/tutorials/gomiddleware/>
  - <https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702#.e4k81jxd3>
  - <https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81>
- Sharing Values with Go Handlers, <https://drstearns.github.io/tutorials/gohandlerctx/>
- About http.handler interface, <https://lets-go.alexedwards.net/sample/02.09-the-http-handler-interface.html>
- JSON and Go, <https://go.dev/blog/json>
