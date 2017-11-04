# `rawhttp` format

This is a simple test data format that we use to store raw http requests and responses.

You can write raw http in the files with escape sequences in it.
We use six equal signs (`======`) as begin/end notifiers.
Each case should have a name. The name defined as a comment after begin notifier and starts with a hash sign (`======#`)

