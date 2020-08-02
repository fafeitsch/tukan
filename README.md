Tukan
===

*Tukan* is a REST client which supports to bulk-configuring of VoIP telephones.

Examples
---

1. Check which telephones are available in a certain network range and are able to be configured with
the provided credentials:
    ```shell script
    ?> tukan --login username -password securePassword scan 10.20.30.40:80+3
    http://10.20.30.40:80:
           Login successful
           Logout successful
    http://10.20.30.41:80:
           Login successful
           Logout successful
    http://10.20.30.42:80:
           Login successful
           Logout successful
    ```
2. Upload local telephone books:
   ```shell script
   ?> tukan --login pb-up -file /tmp/phonebook.xml 10.20.30.40:8080
   http://10.20.30.40:8080:
           Login successful
           Uploading Phone Book successful
           Logout successful
   ```
3. Download the whole phone configuration:
   ```shell script
   ?> go tukan backup --targetDir /tmp 127.0.0.1:8080
   http://10.20.30.40:8080:
           Login successful
           Downloading Parameters successful
           Logout successful
   ```
Simulation
---
In order to test Tukan, there is also a (very simple) VoIP endpoint simulator included
in this repository. This simulator responds to all endpoints which are needed by
Tukan and behaves (within limits) accordingly.
   
Settings
---
For the command line application, IP addresses can either be given as space separated list,
or in the notation w.x.y.z:port+N where N is the number of IP addresses to use. See first example for reference.

All other settings and commands are explained via the `--help` argument of Tukan.

Usage as library
---
Tukan consists of two parts: The command line application and the library. They are both
in the same package, but the library can be used easily for other projects. All exported
functions are documented.

Supported Hardware
---
For various reasons, I do not give a exhaustive list of compatible hardware. If you are
interested in using Tukan to configure your phones or using tukan as library in your own program,
please contact me.

Work in progress
---

The project is one of many hobbys of mine. I don't kown where it is headed, what features will be added,
or when it is finished.
 