# replpipe

# Usage

Create a clojure socket repl via either:

- https://clojure.org/guides/deps_and_cli#socket_repl
- https://clojure.org/reference/repl_and_main#_launching_a_socket_server

If the repl's running on localhost:5555, in one terminal, type:

    $ go run replpipe.go 127.0.0.1:5555

Then in another terminal:

    ~/replpipe> echo '(println "hello world")' >> .repl-pipe
    ~/replpipe> echo '(def a 1)' >> .repl-pipe
    ~/replpipe> echo 'a' >> .repl-pipe

The first terminal will print something like:

    2021/03/12 14:55:57 &{0xc00013e1e0}
    user=> (println "hello world")
    hello world
    nil
    user=> (def a 1)
    #'user/a
    user=> a
    1
    user=> ^Csignal: interrupt

Haven't finished up graceful termination for the program yet - defers aren't
run after ctrl-c.
