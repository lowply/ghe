# ghe

Manage GitHub Enterprise Server instances on AWS.

## Usage

First you need the following:

- The config file (*~/.config/ghe.yaml*)
- Domain you use for the instance
- SSL certificate for the domain
- dnsmasq shuold be running on your computer
- AWS credentials
- SSH key in your AWS console
- Security group and subnet id configured on the VPC

### ami

Prints the ami id of the target version in your region.

```
$ ghe ami 2.20.4
Getting https://github-enterprise.s3.amazonaws.com/release/release-2.20.4.yml
ami-0678075aefcdce317
```

### launch

Launches a new instance.

```
$ ghe launch 2.18.4
Getting https://github-enterprise.s3.amazonaws.com/release/release-2.18.4.yml
Waiting for i-0f1e5ab2dde71cf13 ...
[...]
Launched: i-0f1e5ab2dde71cf13
Updating dns...
Restarting dnsmasq...
+-------------------------------+---------------------+-----------------+-------------------------------------------------------+--------------+---------+---------+
|             DATE              |     INSTANCE ID     |       NAME      |                      PUBLIC DNS                       |  PUBLIC IP   |  STATE  | VERSION |
+-------------------------------+---------------------+-----------------+-------------------------------------------------------+--------------+---------+---------+
| 2019-10-16 05:32:09 +0000 UTC | i-0ba7a56b7312aa9fb | ghe.example.com | ec2-54-248-53-72.ap-northeast-1.compute.amazonaws.com | 54.48.53.72  | running | 2.18.4  |
| 2019-10-17 03:22:34 +0000 UTC | i-0f1e5ab2dde71cf13 | ghe.example.com | ec2-3-112-56-237.ap-northeast-1.compute.amazonaws.com | 3.112.56.237 | running | 2.18.4  |
+-------------------------------+---------------------+-----------------+-------------------------------------------------------+--------------+---------+---------+
```

If you provide the `replica` as a second argument, it launches an instance but doesn't update your dns configuration.

```
$ ghe launch 2.20.4 replica
```

### start

Starts the instance. Accepts multiple arguments.

```
$ ghe start i-0f1e5ab2dde71cf13
```

### stop

Stops the instance. Accepts multiple arguments.

```
$ ghe stop i-0f1e5ab2dde71cf13
```

Or

```
$ ghe stop all
```

### terminate

Terminates the instance. Accepts multiple arguments.

```
$ ghe terminate i-0f1e5ab2dde71cf13
```

Or

```
$ ghe terminate all
```

### dns

Updates your */usr/local/etc/dnsmasq.d/ghe* file and restart dnsmasq.

```
$ ghe dns 3.112.56.237
```

The */usr/local/etc/dnsmasq.d/ghe* file example:

```
$ cat /usr/local/etc/dnsmasq.d/ghe
address=/ghe.example.com/3.112.56.237
```

### list

List your instances

```
$ ghe list
+-------------------------------+---------------------+-----------------+-------------------------------------------------------+--------------+---------+---------+
|             DATE              |     INSTANCE ID     |       NAME      |                      PUBLIC DNS                       |  PUBLIC IP   |  STATE  | VERSION |
+-------------------------------+---------------------+-----------------+-------------------------------------------------------+--------------+---------+---------+
| 2019-10-16 05:32:09 +0000 UTC | i-0ba7a56b7312aa9fb | ghe.example.com | ec2-54-248-53-72.ap-northeast-1.compute.amazonaws.com | 54.248.53.72 | running | 2.18.4  |
+-------------------------------+---------------------+-----------------+-------------------------------------------------------+--------------+---------+---------+
```

### init

Runs the initial configuration based on the *~/.config/ghe.yaml*

```
$ ghe init i-0f1e5ab2dde71cf13
```

### help

Prints the help message.

```
NAME:
   ghe - Launch GitHub Enterprise instances on AWS

USAGE:
   ghe [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
   lowply <lowply@github.com>

COMMANDS:
   ami        Get the Amazon AMI ID of the given GitHub Enterprise Server version, based on your AWS region
   launch     Launch a GitHub Enterprise Server instance of the given version on EC2
   list       List all GitHub Enterprise Server instance with detailed information
   start      Start GitHub Enterprise Server instances. Multiple argument can be passed
   stop       Stop GitHub Enterprise Server instances. Multiple argument, or 'all' can be passed
   terminate  Terminate GitHub Enterprise Server instances. Multiple argument, or 'all' can be passed
   dns        Update your dnsmasq config file to make your domain point to specific IP address
   configure  Initiate GitHub Enterprise Server configuration using the configuration file
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Install

```
go get github.com/lowply/ghe
```

## Config example

The config file path should be *~/.config/ghe.yaml*.

```
---
aws:
  profile: default
  region: ap-northeast-1
  type: r5.xlarge
  # type: r4.xlarge # Use r4 series for version 2.16 or less
  # type: m3.xlarge # Use m3 series for version 2.2 or less
  key_name: key
  sec_group: sg-b5eb8dd3
  subnet_id: subnet-96cae6fe
  excludes:
    - i-010a4fd3df2200508
    - i-0e4d91c8e6c4d3f29
basic:
  domain: "ghe.example.com"
  password: "password"
  license: "/path/to/license.ghl"
  pubkey: "/Users/you/.ssh/id_ed25519.pub"
  ssl_cert: "/path/to/cert.pem"
  ssl_key: "/path/to/key.pem"
ldap:
  host: "ldap.example.com"
  port: 636
  pass: "password"
  base: "ou=People,dc=example,dc=jp"
  dn: "cn=Manager,dc=example,dc=jp"
saml:
  url: "https://example.com/saml2"
  cert: "/path/to/saml.cer"
  issuer: "https://example.com/issuer/"
  format: "unspecified"
smtp:
  host: "smtp.example.com"
  support_email: "support@example.com"
  noreply_email: "noreply@example.com"
syslog:
  host: "syslog.example.com"
  cert: "/path/to/ca.pem"
```
