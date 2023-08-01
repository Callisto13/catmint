This is a spike repo demonstrating how to use the Pure1 Manage API to get
effective capacity usage.

First iteration of this is [here](https://gist.github.com/Callisto13/0464ba4ce05a8e77399612a489d928f4).

This one uses [my fork](https://github.com/Callisto13/pugo) of an old Go client.

Pre-req:

You need admin access to Pure1.
These are linux instructions, bluescreens you are on your own.

1. Create a private key `openssl genrsa -out private.pem`
2. Create a public key `openssl rsa -in private.pem -outform PEM -pubout -out public.pem`
3. Log into PureOne, go to Administration, API Registration and fill in the details: paste in your public key and select Viewer role
4. Copy the generated Application ID

Usage:

```
go build
```

Use the subscription ID to get effective use data about all licenses in that sub:
```
./catmint --app-id <YOUR APP ID> \
  --private-key <PATH TO PRIVATE SSH> \
  --sub <SUBSCRIPTION ID>
```

Use the license name to get effective use data about just one license:
```
./catmint --app-id <YOUR APP ID> \
  --private-key <PATH TO PRIVATE SSH> \
  --license <LICENSE NAME>
```

Example output:
```
Data for License 'Beetroot' in Subscription '12345678'
Effective used: 5.58TB out of 10.00TB
```

There is a lot of other junk available but I have formatted the output to demo
the main data we want to see.
