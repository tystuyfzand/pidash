PiDash
======

A simple, customizable dashboard page with draggable and resizable widgets on a grid. Easy implementation of modules/widgets using HTML, Javascript, and Lua.

![PiDash Preview](/screenshots/preview.png?raw=true)

Installing
----------

Currently, releases aren't available. The biggest issue is cross compiling the LESS compiler (which uses Duktape) for Go, and the same for SASS and others.

If you have experience with this and would like to help, please don't hesitate to contact me.

Setup
-----

Start the service (`service start pidash` or `systemctl start pidash`)

Configure the dashboard from any computer with a web browser and mouse/keyboard.

Configure the Raspberry Pi or other device to load http://127.0.0.1:8080

Note: If you modify any scripts/css, the application must be restarted. It compiles and minifies CSS/JS on start, even for modules.

Structure
---------

`html` - Contains all static files, images, css, js, less.

`modules` - Contains all modules to load, along with assets.

`src/meow.tf/dashboard` - Go source files. This isn't the normal way a Go application is setup, usually it's part of the main directory under a specific package, with a "cmd" package. This is what my projects are usually setup as, and may change.

Shortcuts
---------

__Adding widgets/modules - Ctrl+A__

Add a widget/module to the page

__Save layout - Ctrl+S__

Save the page/module layout and settings.

How does it work?
-----------------

Widgets are defined in modules. Modules can be anything from html like a basic clock, or even interactive html to do actions, like change a temperature on your thermostat!

Handlers are defined in lua scripts, with an easy interface to rendering views. The project uses an embedded Go lua implementation that provides an http and json interface for easy API implementation.

Frontend Dependencies/Libraries
-------------------------------

These are all bundled with known versions into the directory `html/assets`

* [jQuery](http://jquery.com/)
* [jQuery UI](https://jqueryui.com/)
* [Gridstack](http://gridstackjs.com)
* [OctoberCMS (Ajax Framework)](https://octobercms.com)
* [Bootstrap](https://getbootstrap.net)
* [Fontawesome](https://fontawesome.io)
* [Weather Icons](http://erikflowers.github.io/weather-icons/)
* [Bootbox](http://bootboxjs.com/)
* [Mousetrap](https://github.com/ccampbell/mousetrap)
* [Notify.js](https://notifyjs.com/)
* [Lodash](https://lodash.com/)