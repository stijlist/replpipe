# replpipe

# Problem statement

If you want a program (e.g. your text editor) to send data to evaluate to a Clojure (or other) repl, your options are to copy and paste the text yourself or use a language or editor specific plugin.

e.g.

- tslime
- slimv
- vim-fireplace

slimv and vim-fireplace use specialized protocols (SWANK and nREPL, respectively) to start and communicate with the REPL.

tslime is blessedly simple but execs `tmux-send-keys` for every write.

The Clojure team gave us an interesting option by adding some nice default tooling (expose a stream-based REPL over a socket server).
This is almost perfect - your program can start, connect to the socket server, write to a new socket, then exit.

Evaluating code in the REPL is now almost as easy as writing to a file.

However, there's a small hitch - each time you connect, you get a new REPL session.
This makes some sense because the socket server serves you a new socket every time you connect.

But from the user's perspective, they want to be able to incrementally do stateful things in the REPL - e.g. define a function in one evaluation and then use it in the next.

You could manage these REPL sessions from your editor. That would require keeping the socket somewhere as state, and if you exited your editor your REPL session would go away.

Ideally, though, given that it's possible to just type text into the REPL, we should be able to preserve that simplicity and allow your program or editor to just write to a file.

This program is an adapter between a FIFO (letting any program write data to evaluate and exit) and a socket (a stateful stream), letting you decouple the lifetime of your REPL from the editor.

The factoring could be even better - it'd be great to have a simple way to split a socket into its two flows - `split-socket localhost:5555 -read=.repl-in -write=.repl-out` which creates two fifos with those names.

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

## Vim integration

...can be a one-liner.

    vnoremap <leader>e :!tee .repl-pipe<CR>

Assumes you've run replpipe already in Vim's working directory.

So you can open an empty buffer, type `(println "hello world")`, visually select it with `V`, and then hit leader-e.

The replpipe process will print in response:

    (println "hello world")
    hello world
    nil
    user=>
