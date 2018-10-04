# Mantle: Gluing Container Linux together

This repository is a collection of utilities for developing Container Linux. Most of the
tools are for uploading, running, and interacting with Container Linux instances running
locally or in a cloud.

## Overview
Mantle is composed of many utilities:
 - `cork` for handling the Container Linux SDK
 - `gangue` for downloading from Google Storage
 - `kola` for launching instances and running tests
 - `kolet` an agent for kola that runs on instances
 - `ore` for interfacing with cloud providers
 - `plume` for releasing Container Linux
All of the utilities support the `help` command to get a full listing of their subcommands
and options.

## Tools

### cork
Cork is a tool that helps working with Container Linux images and the SDK.

#### cork create
Download and unpack the Container Linux SDK.

`cork create`

#### cork enter
Enter the SDK chroot, and optionally run a command. The command and its
arguments can be given after `--`.

`cork enter -- repo sync`

#### cork download-image
Download a Container Linux image into `$PWD/.cache/images`.

`cork download-image --platform=qemu`

#### Building Container Linux with cork
See [Modifying Container Linux](https://coreos.com/os/docs/latest/sdk-modifying-coreos.html) for
an example of using cork to build a Container Linux image.

### gangue
Gangue is a tool for downloading and verifying files from Google Storage with authenticated requests.
It is primarily used by the SDK.

#### gangue get
Get a file from Google Storage and verify it using GPG.

### kola
Kola is a framework for testing software integration in Container Linux instances
across multiple platforms. It is primarily designed to operate within
the Container Linux SDK for testing software that has landed in the OS image.
Ideally, all software needed for a test should be included by building
it into the image from the SDK.

Kola supports running tests on multiple platforms, currently QEMU, GCE,
AWS, VMware VSphere, and Packet. In the future systemd-nspawn and other
platforms may be added.
Local platforms do not rely on access to the Internet as a design
principle of kola, minimizing external dependencies. Any network
services required get built directly into kola itself. Machines on cloud
platforms do not have direct access to the kola so tests may depend on
Internet services such as discovery.etcd.io or quay.io instead.

Kola outputs assorted logs and test data to `_kola_temp` for later
inspection.

Kola is still under heavy development and it is expected that its
interface will continue to change.

By default, kola uses the `qemu` platform with the most recently built image
(assuming it is run from within the SDK).

#### kola run
The run command invokes the main kola test harness. It
runs any tests whose registered names matches a glob pattern.

`kola run <glob pattern>`

#### kola list
The list command lists all of the available tests.

#### kola spawn
The spawn command launches Container Linux instances.

#### kola mkimage
The mkimage command creates a copy of the input image with its primary console set
to the serial port (/dev/ttyS0). This causes more output to be logged on the console,
which is also logged in `_kola_temp`. This can only be used with QEMU images and must
be used with the `coreos_*_image.bin` image, *not* the `coreos_*_qemu_image.img`.

#### kola bootchart
The bootchart command launches an instance then generates an svg of the boot process
using `systemd-analyze`.

#### kola updatepayload
The updatepayload command launches a Container Linux instance then updates it by
sending an update to its update_engine. The update is the `coreos_*_update.gz` in the
latest build directory.

#### kola subtest parallelization
Subtests can be parallelized by adding `c.H.Parallel()` at the top of the inline function
given to `c.Run`. It is not recommended to utilize the `FailFast` flag in tests that utilize
this functionality as it can have unintended results.

#### kola test namespacing
The top-level namespace of tests should fit into one of the following categories:
1. Groups of tests targeting specific packages/binaries may use that
namespace (ex: `docker.*`)
2. Tests that target multiple supported distributions may use the
`coreos` namespace.
3. Tests that target singular distributions may use the distribution's
namespace.

#### kola test registration
Registering kola tests currently requires that the tests are registered
under the kola package and that the test function itself lives within
the mantle codebase.

Groups of similar tests are registered in an init() function inside the
kola package.  `Register(*Test)` is called per test. A kola `Test`
struct requires a unique name, and a single function that is the entry
point into the test. Additionally, userdata (such as a Container Linux
Config) can be be supplied. See the `Test` struct in
[kola/register/register.go](https://github.com/coreos/mantle/tree/master/kola/register/register.go)
for a complete list of options.

#### kola test writing
A kola test is a go function that is passed a `platform.TestCluster` to
run code against.  Its signature is `func(platform.TestCluster)`
and must be registered and built into the kola binary. 

A `TestCluster` implements the `platform.Cluster` interface and will
give you access to a running cluster of Container Linux machines. A test writer
can interact with these machines through this interface.

To see test examples look under
[kola/tests](https://github.com/coreos/mantle/tree/master/kola/tests) in the
mantle codebase.

#### kola native code
For some tests, the `Cluster` interface is limited and it is desirable to
run native go code directly on one of the Container Linux machines. This is
currently possible by using the `NativeFuncs` field of a kola `Test`
struct. This like a limited RPC interface.

`NativeFuncs` is used similar to the `Run` field of a registered kola
test. It registers and names functions in nearby packages.  These
functions, unlike the `Run` entry point, must be manually invoked inside
a kola test using a `TestCluster`'s `RunNative` method. The function
itself is then run natively on the specified running Container Linux instances.

For more examples, look at the
[coretest](https://github.com/coreos/mantle/tree/master/kola/tests/coretest)
suite of tests under kola. These tests were ported into kola and make
heavy use of the native code interface.

#### Manhole
The `platform.Manhole()` function creates an interactive SSH session which can
be used to inspect a machine during a test.

### kolet
kolet is run on kola instances to run native functions in tests. Generally kolet
is not invoked manually.

### ore
Ore provides a low level interface for each cloud provider. It has commands
related to launching instances on a variety of platforms (gcloud, aws,
azure, esx, and packet) within the latest SDK image. Ore mimics the underlying
api for each cloud provider closely, so the interface for each cloud provider
is different. See each providers `help` command for the available actions.

Note, when uploading to some cloud providers (e.g. gce) the image may need to be packaged
with a different --format (e.g. --format=gce) when running `image_to_vm.sh`

### plume
Plume is the Container Linux release utility. Releases are done in two stages,
each with their own command: pre-release and release. Both of these commands are idempotent.

#### plume pre-release
The pre-release command does as much of the release process as possible without making anything public.
This includes uploading images to cloud providers (except those like gce which don't allow us to upload
images without making them public).

### plume release
Publish a new Container Linux release. This makes the images uploaded by pre-release public and uploads
images that pre-release could not. It copies the release artifacts to public storage buckets and updates
the directory index.

#### plume index
Generate and upload index.html objects to turn a Google Cloud Storage
bucket into a publicly browsable file tree. Useful if you want something
like Apache's directory index for your software download repository.
Plume release handles this as well, so it does not need to be run as part of
the release process.

## Platform Credentials
Each platform reads the credentials it uses from different files. The `aws`, `do`, `esx` and `packet`
platforms support selecting from multiple configured credentials, call "profiles". The examples below
are for the "default" profile, but other profiles can be specified in the credentials files and selected
via the `--<platform-name>-profile` flag:
```
kola spawn -p aws --aws-profile other_profile
```

### aws
`aws` reads the `~/.aws/credentials` file used by Amazon's aws command-line tool.
It can be created using the `aws` command:
```
$ aws configure
```
To configure a different profile, use the `--profile` flag
```
$ aws configure --profile other_profile
```

The `~/.aws/credentials` file can also be populated manually:
```
[default]
aws_access_key_id = ACCESS_KEY_ID_HERE
aws_secret_access_key = SECRET_ACCESS_KEY_HERE
```

To install the `aws` command in the SDK, run:
```
sudo emerge --ask awscli
```

### azure
TBD (FIXME)

### do
`do` uses `~/.config/digitalocean.json`. This can be configured manually:
```
{
    "default": {
        "token": "token goes here"
    }
}
```

### esx
`esx` uses `~/.config/esx.json`. This can be configured manually:
```
{
    "default": {
        "server": "server.address.goes.here",
        "user": "user.goes.here",
        "password": "password.goes.here"
    }
}
```

### gce
`gce` uses the `~/.boto` file. When the `gce` platform is first used, it will print
a link that can be used to log into your account with gce and get a verification code
you can paste in. This will populate the `.boto` file.

See [Google Cloud Platform's Documentation](https://cloud.google.com/storage/docs/boto-gsutil)
for more information about the `.boto` file.

### packet
`packet` uses `~/.config/packet.json`. This can be configured manually:
```
{
	"default": {
		"api_key": "your api key here",
		"project": "project id here"
	}
}
```

### qemu
`qemu` is run locally and needs no credentials, but does need to be run as root.
