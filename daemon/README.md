# SWE Run as Service

Not: For easy setup goto docker folder and follow docker-compose command.

### Manual installation

If you want to run the daemon as a service, you can use the following commands:

0. edit `.swed.config` file for your system
1. mkdir /opt/swed
2. cd /opt/swed
3. copy .swed-config /opt/swed
4. copy swed binary to /opt/swed
5. `chmod +x /opt/swed/swed`
6. copy swed.service to /etc/systemd/system
7. `systemctl enable swed.service`
8. `systemctl start swed.service`
9. watch for errors `journalctl -u swed.service -f`
