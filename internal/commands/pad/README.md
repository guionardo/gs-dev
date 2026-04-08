# PAD command

Pad is a form of transmit short information.

To send data:

```bash
> pad create "TEXT"
sending data to storage... DONE!

The access key is: ABCD-EFGH
```

To receive data:

```bash
> pad get ABCD-EFGH
TEXT
```

To delete data:

```bash
> pad delete ABCD-EFGH
OK
```

## The backend of this service is a simple GRPC API

