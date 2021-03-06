# acker

Acker is a cli application to consume/produce messages from/to AMQP servers, e.g. RabbitMQ. It can be used to ack unacked messages on the queue, for the messages not to build up on the queue.

## Compile from Source

You can compile Acker from source with the following commands:

```bash
$ go get github.com/mostafa/acker
$ go $GOPATH/github.com/mostafa/acker
$ go get
$ go build
$ go install
```

Alternatively you can cross-compile it using these commands:

```bash
$ go get github.com/mostafa/acker
$ cd $GOPATH/github.com/mostafa/acker
$ ./build.sh
```

## Running Acker Consumer

For recovering a queue from unacked messages piling up and slowing the processing, which is usually due to a non-responding consumer, run the following command. It also acts as a consumer that runs forever, until killed by `CTRL+C`.

```bash
$ acker consume --server=<AMQP-URL> --queue=<QUEUE> --autoack=true --recover=true --current-consumer=true
```

The `--server` and `--queue` flags set server and queue configuration to connect to. The `--autoack` flag causes each message consumed by Acker consumer to be automatically acknowledged. Setting `--autoack` to `false` causes unacked consumed messages to pile up the queue, which is only used for testing purposes. The `--recover` flag will recover unacknowledged messages on the queue and mark them as ready, so that they can be processed again. Setting the `--recover` flags also means that the messages will be processed by the Acker consumer.

## Running Acker Producer

For producing messages on the queue, run the following command:

```bash
$ acker produce --server=<AMQP-URL> --queue=<QUEUE> --body="<BODY>" --count=10
```

Just like the Acker consumer, the `--server` and `--queue` flags set server and queue configuration to connect to. The `--body` specifies the body of the message to be published on the queue. The `--count` specifies the number of messages to be produced.
