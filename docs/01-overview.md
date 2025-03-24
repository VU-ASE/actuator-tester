# Overview

This is a service that can test the default [actuator](https://github.com/VU-ASE/actuator) through the the web interface or by accepting JSON objects over a UDP connection. It lets you test the motors and servo.

# Usage

This service expect a JSON object sent over a UDP connection. The default prot is `8080` and can be changed in the `service.yaml`.

The UDP payload should be a parsable JSON object in the following format:
```
{
    "channel": int,
    "value": float64
}
```

| Channel | Description |
|---------|-------------|
| 0       | Controls the steering servo and expects a value between -1.0 (full left) and 1.0 (full right) |
| 1       | Controls the **left** motor and expects a value between -1.0 (full reverse) and 1.0 (full forwards)|
| 2       | Controls the **right** motor and expects a value between -1.0 (full reverse) and 1.0 (full forwards)|

## Testing

To check if the motors work, you can run `make test` from the container - it will run `tester.py` in the `scripts` directory. Upon execution, you will be prompted to input the id of the rover you are trying to test. 