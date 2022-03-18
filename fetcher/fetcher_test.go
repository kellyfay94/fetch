package fetcher

// FUTURE IMPROVEMENT
// Test Cases:
//	- Start *http.Server, hosting a simple page, then fetch the page; ensure that the page was fetched and stored properly
//	- Don't start server, fetch page known not to exist; print error
// < I wouldn't include any tests that require dependence on an outside resources, as then an external force can have a negative effect on the performance of the test >
