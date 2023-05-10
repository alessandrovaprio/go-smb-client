# Overview
Go Library for Samba2 Exported with C bindings.
It's aim is simplify the samba2 handling in Java.
It means it'll be callable from Java through JNI.

# Getting Started

## Prerequisites
 - First of all you need GO or Docker installed

## Compilation instruction
- If go installed
    ``` 
    $ make
    ```
- With Docker
    ``` 
    $ make dockerbuild
    ```

The binary will be saved in bin folder.

### File's explanation
- main.go contains all the function exposed with C bindings (special headers inside)
- client.go contains all the implementation in pure GO.

### Example
- Is present a java example usage with a test file.

# Libraries
Thanks to this library I simplified my workflow:
 - [go-smb2](https://github.com/hirochachacha/go-smb2)
# License
distributed under MIT.
