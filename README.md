Tukan
===

*Tukan* is a REST client which supports to bulk-configuring of VoIP telephones.

Please be careful with this tool, otherwise, the phone configurations
can be lost or the phones become insecure and unstable.
As declared in the license,
I give no warranty whatsoever to the usefulness or correctness of this tool.

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
   ?> tukan --login pb-up -sourceDir /tmp 10.20.30.40:8080
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

The main differences in the simulation are:
1. When posting parameters, the real phones will merge them with the existing ones, while
   the simulator simply overwrites them. For example, the request body `{"Phone Name": "Phone ABC"}` on *real* phones will
   only reset the phone name. The simulator will set all parameters to their empty values except
   the phone name.
2. When getting parameters, the real phones does not only send to actual field values,
   but also information about validation and possible values.
   The simulator responses with the same format as used for posting parameters.
   However, the UnmarshalJSON method of the parameters can deal with both variants.
   
Settings
---
For the command line application, IP addresses can either be given as space separated list,
or in the notation w.x.y.z:port+N where N is the number of IP addresses to use. See first example for reference.

All other settings and commands are explained via the `--help` argument of Tukan.

Usage as library
---
Tukan consists of two parts: The command line application and the library. They are both
in the same package, but the library can be used easily for other projects. All exported
functions are (or will be) documented.

Supported Hardware
---
For various reasons, I do not give a exhaustive list of compatible hardware. If you are
interested in using Tukan to configure your phones or using Tukan as library in your own program,
please contact me.

Work in progress
---

The project is one of many hobbys of mine. I don't know where it is headed, what features will be added,
or when it is finished.
 