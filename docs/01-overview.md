# Overview

`actuator-tester` is a service that allows you to effortlessly debug any potential issues you might encounter when implementing your own actuator or simply to test the motors themselves. `actuator-tester` does not take any input from any other service and effectively replaces `controller` when you run it in the pipeline. You can directly interact with your `actuator` by sending your commands to this service via a specified UDP port.
