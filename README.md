unBALANCE
=========

*tl;dr* **unBALANCE** is an [unRAID](http://lime-technology.com) app to free up space from one of the disks in the array, by moving folders and files to other disks.

## Screenshot
![Screenshot](110-home.png)

## Upgrade Notes for v1.x.x
For those upgrading from previous versions (v0.7.4 and below), please take note of the following changes:

- The notifications system is based on unRAID's settings, so you need to set up unRAID's notifications first. This also means that you must be running 6.1 to receive emails and other unRAID alerts
- The configuration file format has changed, but the app will convert it upon first boot.

## Introduction
With unBALANCE, you select a disk that you want to have as much free space available as possible, and it will move the folders/files to other disks filling them up, as space allows.

The logic is quite simple:<br>
Get all elegible folders from the source disk<br>
Order the target disks by free space available<br>
For each disk, fill it up as much as possible with folders/files from the source disk, leaving some headroom (currently set at 450Mb).<br>

Internally, all move operations are handled by [diskmv](https://github.com/trinapicot/unraid-diskmv).

The array must be started for the app to work.

The first time you open the app, you are redirected to the settings page, where you can navigate your user shares, to select which folders you want to move.

You can Select an entire user share (/films in the screenshot below) or any folder(s) under the user shares (/films/bluray for example).

![Settings](110-settings.png)


## Install
There are 2 ways to install this application

- Docker app (preferred)<br>
Add the following repository in the Docker GUI<br>
https://github.com/jbrodriguez/docker-containers/tree/templates<br>
Add the container [jbrodriguez]/unbalance
The defaults are <br>
Port: 6237<br>
Volumes: <br>
"/mnt" (required)<br>
"/root" (required)<br>
"/usr/local/sbin" (required)<br>
"/path/to/config/dir" (required)<br>
"/path/to/log/dir" (not required)<br>
"/etc/localtime" (not required, to synchronize time with unRAID)<br>

- Manual
```Shell
# mkdir -p /boot/custom
# cd /boot/custom
# wget https://github.com/jbrodriguez/unbalance/releases/download/<enter latest version here>/unbalance-<enter latest version here>-linux-amd64.tar.gz -O - | tar -zxf - -C .
```
*NOTE*: If run manually, move operations will be performed as root user. Please take that into account.

## Running the app
Start the container or 

```Shell
# cd /boot/custom/unbalance
# ./unbalance
```
As mentioned previously, the app will show the Settings page the first time it's run. Choose the elegible folders now.

By default, the dry-run option is selected.

It means that operations are simulated, it only shows what it would actually do.

To perform the operations, uncheck the dry-run checkbox.

## Credits
This app uses the [diskmv](https://github.com/trinapicot/unraid-diskmv) script (check the [forum thread](http://lime-technology.com/forum/index.php?topic=36201.0) for additional information).

The icon was courteously created by [hernandito](http://lime-technology.com/forum/index.php?topic=39707.msg372508#msg372508) (fellow unRAID forums member)

It was built with:

- [Go](https://golang.org/) - Back End
- [echo](https://github.com/labstack/echo) - REST and websocket api
- [pubsub](https://github.com/tuxychandru/pubsub/) (slightly modified)
- [React](https://facebook.github.io/react/) - Front End
- [reactorx](https://github.com/jbrodriguez/reactorx) - Flux/Redux-like React 
- [flexboxgrid](http://flexboxgrid.com/) - CSS3 flex based grid system
framework
- [js-csp](https://github.com/ubolonton/js-csp) - Go-like concurrency for javascript
- [Webpack](https://webpack.github.io/) - Build toolchain

## License
[MIT license](http://jbrodriguez.mit-license.org)
