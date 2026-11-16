# Debugging guide for debugging on the VM
See the following [link](https://github.com/go-delve/delve/blob/master/Documentation/cli) for an exhaustive set of commands, configuration options and more. In this guide `<text>` is a placeholder for actual input, so do not add the duckbills in your actual input. `[something]` is to be understood as an optional argument.
### How to debug with `delve`
1. The `delve` debugger is shipped with the Go extension for VSCode, but we installed it on the VM as well with `sudo dnf install delve`.
2. Clone the `main` branch of the repository, or a different commit on it, if that's what you need, and `cd` into it.
3. Run `DB_PATH="<path>" COOKIE_SECURE="false" dlv debug ./cmd/main`. If you need any extra flags, make sure you add `--` between the first part and the flags (with spaces after the first and before the second). Example: `dlv debug ./cmd/main -- -arg1 <value>...`.
4. Set one or more breakpoints by writing `break X` where `X` is either a memory address (a hexadecimal integer), a `<filename>:<line>` pair, a specific line in the *current* file (just an integer) or one of the other options that may be found [here](https://github.com/go-delve/delve/blob/master/Documentation/cli/locspec.md).
    - You may add conditional breakpoints by first running the previous command, which gives you a breakpoint number (integer) which starts from 1. To make a breakpoint conditional, do `cond <breakpoint-number> <expression>`, where `<expression>` is a most often a predicate like `i == 6`, to only trigger the breakpoint in a loop or method with a variable `i` which may attain the value 6. See `expr.md` in the first link of this document for other valid expression types.
5. Type `continue` then start interacting with the code through Postman, curl or your browser and hit a breakpoint.
6. Inspect program values:
    - You may print the value of a specific variable with `print <variable-name-or-expression>` (or `p` instead of `print`). 
    - You may print the value of a method argument with `args` (no arguments). 
    - You may print the values of local variables with `locals` (no arguments).
    - You may also print all variables in the package with `vars [-v] [<regex>]`. If you add a regex, only package variables which satisfy the regex will be shown. -v shows you additional information about the package.
7. When on a specific breakpoint, either skip to the line `[count]` places ahead with `next [count]`, run until the breakpoint is hit again with `continue` or perform a single step with `step`. 

### How to debug a container
Delve must be installed inside the container, and the container must be running before you attach to its application process. Delve’s remote interface is unauthenticated and can allow remote code execution, so keep it inaccessible from the public network.

- Start the container and application. Publish the debug port on the VM’s loopback interface when you create the container:
```sh
docker run -p 127.0.0.1:<debug-port>:<debug-port> ...
```
If the container is already running without this port mapping, you’ll need to recreate it with the mapping.

- Find the application’s PID inside the running container:
```sh
docker exec <container> ps
```

- Attach Delve to the application process inside the container:
```sh
docker exec <container> dlv attach <PID> --headless --listen=:<debug-port> --api-version=2
```
Delve must be installed in the container and have permission to trace the application.

- In another local terminal, connect to Delve:
```sh
dlv connect 127.0.0.1:<debug-port>
```
Use the same port number throughout. In VS Code, configure a remote Go debug session to connect through the forwarded port. 