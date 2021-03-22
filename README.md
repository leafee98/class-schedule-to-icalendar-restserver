# class-schedule-to-icalendar-restserver

the rest api server of class-schedule-to-icalendar, depend on [rest-api-grpc](github.com/leafee98/class-schedule-to-icalendar-rpcserver), require mariadb.

## install

### database configuration

enable event_schedular in maridb to allow auto delete expired token in database. see also on [here](https://mariadb.com/docs/reference/mdb/system-variables/event_scheduler/)

enable event_schedular in runtime.

```
https://mariadb.com/docs/reference/mdb/system-variables/event_scheduler/
```

enable event_schedular in configuraton file.

```
[mariadb]
event_scheduler=ON
```