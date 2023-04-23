# SCP Challenge

This is called a challenge and not an assignment, because it will involve some research. The SSH Server doesn't have SCP (secure copy) support. How would we add that?

The OpenSSH client has a command "scp" that you can use to transfer files. This would transfer server.pub to the server:

```bash
scp -P 2022 -i mykey.pem server.pub 127.0.0.1:server.pub
```

When executing this command, exec is triggered and the payload shows the following output:

```bash
payload: scp -t server.pub
```

"-t" in scp is not documented in the man page, because it's for internal usage. You could execute 'scp -t server.pub' on your local machine and see what output it gives you. You could also read about how the SCP protocol works and implement it that way.
