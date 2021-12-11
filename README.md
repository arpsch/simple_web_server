1. Steps to run the server
   from this folder, run
   $ go run main.go 
       OR
   $ go run main.go -port=<port of your choice>
    [if not port sent, 8080 is the default port]
2. Unit tests are under api/http folder
   go test ./... -short
3. Integration tests are under api/http
   go test ./...
   [This command runs all the tests - unit and integration]
4. Client: A go based client program tests/client.go
   $ go run client.go -port=<the same port as the server is listening>
    [defaults to :8080]

/***************************************************************/
I have tried to finish all the bonus requirements along with the mandatory ones
a. Basic Auth with hardcoded username-idt and password-idt123
b. GetUsers endpoint
c. Reading port from command line
d. Simple logger middleware to log request start and end

//
Known issues:
1. GET /users/:id  call where :id is nill is resolving to GET /users.
This is issue is leading to failure of one unit test.
[I guess that the router package is doing this.] As I had to jump to part 2, I skipped this one.

2. In unit tests, I've not verified the error string. I went too far to do it for all test cases.
   But the user is getting the right errors, verfied by the client program.
3. I was meaning to provide a simple html interface, but ran out of time
   You could see the Error if you submit it as the form is not adding credentials
