# go-UCSPMMetering
This application has been put together to assist with billing Cisco UCS platforms based on CPU utilisation.  As of the current status the Cisco UCS API does not expose the CPU percentage for each of the physical CPU's.  It is therefore not possible to have creative billing methods for how the systems are consumed.

It was required to create a way of monitoring, investigating and then summerising the utilisation on a weekly basis.

This tool brings together a couple of enterprise applications to allow the capture of this information, mainly it does the following;

* It will query Cisco UCS Performance Manager for all devices. (It will then exclude, network, compute and storage devices, to leave servers and hypervisors)
* Each device is then queried to get the associated hardware UUID. (If one is not in the Cisco UCS Performance Manager inventory, then it is ignored.)
* Each UCS Domain is then queried for all of the UUID's and returns information about the hardware it is associated with.
* The tool will then match each of the UUID's with the physical server UUID and produce an output of information. (Mainly ties a system to the Cisco UCS Server serial number)
* Once matched, each if the UUID's are then queried again in Cisco UCS Performance Manager for the report information over the specified period.


## Setting up your GO environment
Depending on your particular environment, there are a number of ways to setup and install GO.  This repo was developed on a MAC and was installed using Brew.  For instructions on installing HomeBrew, please check [here](https://brew.sh/); and then entering;
```fish
> brew install go
```

If you do not want to use HomeBrew or you are running on a different platform, you can install the GO language using a binary from here;
https://golang.org/dl/

Once this has completed, open a cmd or terminal window and check GO has been installed and configured correctly;

Enter <b>echo $GOPATH</b>, hopefully you will be presented with a path and should be ready to go.

```fish
> echo $GOPATH
/path/to/go/bin/src/pkg folders
```

## Testing your GO environment
Once you have completed the above, its time to create a very simple test script to ensure everything is ready.

Go to a path where you are happy to store the source code for your application, this could be anywhere, including your desktop, documents, root folder, etc.

Create a folder and enter the directory.  Create a new file called "main.go" and enter the following code into it;

```go
package main

import "fmt"

func main() {
    fmt.Println("GO is working!")
}
```

At the command line, change directory using cd to the directory where your main.go file is and execute the following;
```fish
> go run main.go
```

You should see as output, something similar to;

"GO is working!"

If you reached this point, everything is working and you are ready to run the included code!

## Getting the code
There are a couple of ways you can get the code, depending on how comfortable you are with the command line and development envrionments;

You could download the zip file, [here](https://github.com/robjporter/go-UCSPMMetering/archive/master.zip).

You could use the command line git command to clone the repository to your local machine;
1. At the command line, change directory using cd to the directory where the repository will be stored.
2. Enter, git clone https://github.com/robjporter/go-UCSPMMetering.git
3. You will see output similar to the following while it is copied.
```fish
Cloning into `go-UCSPMMetering`...
remote: Counting objects: 10, done.
remote: Compressing objects 100% (8/8), done.
remove: Total 10 (delta 1), reused 10 (delta 1)
unpacking objects: 100% (10/10), done.
```
4. Change into the new directory, cd go-UCSPMMetering.
5. Move onto setting up the application.

## Application dependencies
For the application to work correctly, we need to get one dependency and we can achieve that with the following, via the cmd line.
```fish
> go get -u github.com/robjporter/go-functions
```

## Setting up the application
You need to add the UCS and UCS Performance Manager systems to the application.  Your password will be encrypted before it is stored, however usernames will remain in plain text.  This should be a read only account on both systems, so should not cause too much of a security risk.

### Add UCS Domain
Repeat this process as many times as needed.
```go
> go run main.go add ucs --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Update UCS Domain
The update process will only succeed if the IP of the UCS Domain is already in the config file.
```go
> go run main.go update ucs --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Delete UCS Domain
The delete process will only succeed if the IP of the UCS Domain is already in the config file.
```go
> go run main.go delete ucs --ip=<IP>
```
### Show UCS Domains
To show the current configuration details for a UCS System;
```go
> go run main.go show ucs --ip=<IP>
```

### Add UCS Performance Manager
This action only needs to be done once, running it again will simply over write the config, as only a single UCS Performance Manager instance is permitted.
```go
> go run main.go add ucspm --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Update UCS Performance Manager
This will update the current UCS Performance Manager config.
```go
> go run main.go update ucspm --ip=<IP> --username=<USERNAME> --password=<PASSWORD>
```

### Delete UCS Performance Manager
As the application will only run against a single UCS Performance Manager, there is no need to specify any details, calling delete will remove the UCS Performance Manager config. 
```go
> go run main.go delete ucspm
```

### Show UCS Performance Manager
To show the current configuration details for the UCS Performance Manager;
```go
> go run main.go show ucspm
```

### Show All discoverable systems
To show all the currently entered system information;
```go
> go run main.go show all
```

## Running the application
Once the UCS and UCS Performance Manager systems have been added, the application is now ready to run.
```go
> go run main.go run
```

## Running the application for a specific month/year
You may wish to run the application and gather data for a specific month and/or year, you can achieve this by setting the correct flags;
### Current month and year
```go
> go run main.go run
```
### Current month in 2016
```go
> go run main.go run --year=2016
```
### Febraury of current year
```go
> go run main.go run --month=feb
```
### Specific month and year
```go
> go run main.go run --month=feb --year=2016
```

## Cleaning up after an application run
Once the application has been run, there will be several files generated (more if in debug mode) which you may wish to remove before running the application again.
```go
> go run main.go clean
```