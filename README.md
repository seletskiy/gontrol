Gontrol
=======

Simple remote control via web-service written on Go.

What is it
==========

Gontrol brings possibility to use any device that can HTTP to control any other
device, where Go is supported.

One of the trivial usages is to use mobile phone as presentation clicker. No
app is required on the phone, only a Web browser.

How to use Gontrol as clicker
=============================

As an example, let's see how to use clicker in HTML-based presentations.

All you need is, step by step:

1. Launch `gontrol` in a directory with `web` subdirectory (see origin repo as
   an example);
2. In the HTML file of you presentation write something like this:
   ```js
   var controlSocket = new WebSocket("ws://YOUR_IP_ADDRESS_HERE:1444/s/")
   controlSocket.onmessage = function (event) {
       eval('var msg = ' + event.data); // or use JSON parser

       switch (msg['move'][0]) {
           case 'next':
               Flowtime.nextFragment(); // depends on framework you are using
               break;
           case 'prev':
               Flowtime.prevFragment();
               break
       }
   }
   ```
3. Open your presentation in Web browser on the controlled computer;
4. Open `http://YOUR_IP_ADDRESS:1444/` in the browser on your phone;
5. Control your presentation progress by clicking on buttons on your phone!

Note, that port can be changed to more pleasurable one, and `web` dir from the
origin repo can be used as starting point to control interface.
