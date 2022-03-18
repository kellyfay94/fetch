# Autify Backend Engineer Take-Home Test
https://www.notion.so/autifyhq/Backend-Engineer-Take-Home-Test-63032907b74341f8bd899018d685f03c

# Fetch Application
Application that, given a URL, fetches the HTML of the page at the URL.  Application stores metadata of each request, which can be returned by using the `-metadata` flag.

## Run Fetch
After built, run fetch from the terminal by just calling:
    `./fetch`
Or:
    `fetch` (if added to $PATH)

Arguments are as defined in the assignment:
* Fetch a Page: `fetch <urls>`
* Return the Metadata: `fetch -metadata <urls>`
* Verbose Output: `fetch -v <urls>` or `fetch -v -metadata <urls>`
Where `<urls>` is the space-delimited list of the URL with the protocol (ex: "https://", "http://")

## Build & Run Instructions
### Simple Docker
To run the `fetch` commands in the shell script "run_checks.sh" in a docker container, simply run:
    `sh run_checks_docker.sh`

This will:
1. Build the container
2. Run the container, executing "run_checks.sh" with the output going to "output.txt"
3. Cleanup the container and images
4. Clear the screen
5. Return `output.txt` to the shell, so it's easily readable

### Docker
The Dockerfile can be built into an image for most systems with a simple:
    `docker build -t fay_chris_fetch .`

The image can then be run with the below command.  This executes the commands in "run_checks.sh", and outputs them to console.
    `docker run --rm -it fay_chris_fetch`

To run the "./fetch" command inside the container, execute:
    `docker run --rm -it fay_chris_fetch sh`

And then use the prompt to execute the fetch application per the instructions in "Fetch Arguments"

### Local
Requirements: Go 1.17
1. Checkout this repo to the dir "`go env GOPATH`/github.com/kellyfay94/fetch"
2. Navigate a terminal to the directory mentioned in Step 1
3. Run the following command to checkout dependencies:
    `go mod vendor`
4. Run the following command to build the executable
    `go build`
5. The application can then be executed with the instructions in "Fetch Arguments"


## Note About Bug
Upon fetching the "www.autify.com" page, only 104 `<img>` tags are counted.  Via a search of the www.autify.com.html document, there appear to be 134 `<img>` tags.  I ran out of time while investigating this (both the true number of `<img>` tags and if there was an issue with the application).

I suspect that there are 134 images on the autify page, and the issue lies with the `html.Tokenizer`, with the potential that it is interpreting multiple tags together.  Given more time, I would like to use an html.Parser to sanitize the input first.

## Notes about Specifications
The specification that defined this assignment left a few elements undefined / vague.  In a professional setting, these requirements would have been clarified before development.  Below are some of the undefined cases for this assignment:
* URL -> Filename Parsing
    * Undefined Aspect: In the assignment, the only URLs provided point to a site's root / index.  No specification is provided for URLs including "/" or other filename-illegal characters.
    * Assumption for Assignment: We can transform the URL with some other legal character should transform the URLs to use a double underscore ('__').
    * Future Improvement: This assumption may cause problems in the future, so it would be preferred to not save pages according to their full URL path, but instead given an index and a lookup file if developed further for production
* Metadata Storage
    * Undefined Aspect: The format for storing metadata is not defined.  The assignment gives an example that indicates that upon running `ls` in a command prompt, only the "*.html" files are seen.  While this may imply that the metadata for the HTML Fetch should not in the same directory as the ".html", the "fetch" application is also not seen in this `ls` example.
    * Assumption for Assignment: The `ls` command is actually behaving like `ls *.html`, and storing the metadata in a file in the same folder 
    * Future Improvement: Strictly define the fetched file storage format, including how the metadata should be stored.  For example:
    ```
    fetch/
        fetch
        index.json
        <@foreach index>/
            metadata.json
            page/
                <name for page>.html
                <@foreach resources, use subdir as necessary>
    ```
* "Image" and "Link" Qualifications
    * Undefined Aspect: What is considered an "image" and a "link"?
    * Assumption for Assignment: Count `<a>` as links and `<img>` as images.
    * Future Improvement: Clarify the definition of "image" and "link" for this project.  Some possible inclusions would be:
        * Count images included via CSS classes
        * Count svgs as an "image"
        * Count images that may be added after a JavaScript is run (like the profile picture of the user on Google's Homepage)
        * Count "links" that may just be `<div>` with an "onClick()" or similar event
        * Count "links" of buttons / inputs  

## Notes for Further Improvement
This app is done for the 2-4 hour limit of the take home assignment.  If this were extended into a more production-ready app, then I would likely address the below points:
* Revision to the Docker Build / Run Pipeline
    * For the purpose of this assignment, I created a Dockerfile that would build and run the application.  Given more time, I would create a multi-stage Dockerfile (or separate build and run Dockerfiles) to build the application in a Dockerfile and run it exclusively in another
    * Additionally, I didn't use volumes in this assignment.  Generally, I would use them in a build-and-run-on-local pipeline to remove the need to rebuild the image each time.  However, volumes don't work the same on all platforms (and work differently with Dockerfile-compatible containerization systems like `podman` on Mac), so I just targeted for the system that should have the most compatibility.
* Address "FUTURE IMPROVEMENT" comments in the code
    * I left a few points to address in future iterations in the code, like:
        * Adding Timeouts and Proxies to the Fetch action
        * More graceful exit actions
        * Add Test Cases
            * I did not add any test cases to the application, as I ran out of time.  I structured the app to not have the logic write directly to console, but instead pass messages to later be written to console.  This is to make it more possible to write test cases around parsing the messages that /would/ be written to console.
            * See "fetcher/*_test.go" for planned test cases.
    * In a professional setting, I have left tags such as `// FUTURE TODO` in the code, but will follow it up with a ticket to address at a later time as well.
* 
* `Extra Credit`: Adding a system to fetch the additional resources
    * I ran out of time before I was able to address the `Extra Credit`, but I would address it by:
        * In the `(*Page).ExtractMetadata` function, a HTML Tokenizer iterates through the nodes on the page.  As noted in the 'FUTURE IMPROVEMENT', I would add a list of the resources to fetch to a new field in the `Page` struct.
            * Note: We would have to specify if this is just for resources of a relative path, or of any path.  If any path is an absolute resource (ex: an `<img src="http://subdomain.example.com/image.png"></img>` node if the URL is "http://www.example.com"),  we would have to update the `src` attribute as well.
        * After `(*Page).ExtractMetadata`, I would call a new function: `(*Fetcher).FetchPageResources(*Page)` which would iterate through each of the resource names and fetch them to the local drive