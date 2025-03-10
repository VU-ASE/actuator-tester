# Usage

This service expect a JSON object sent over a UDP connection. Port is specified in `service.yaml` to the default value `:8080` and can be changed manually if desired.

Expected format of the incoming JSON:

```
{
    "channel": int,
    "value": float64
}
```
### Note 
`channel` expects possible values of 0 (for changing the steering angle), 1 (to engage the left motor) or 2 (to engage the right motor).

`value` expects a number between -1 and 1. 