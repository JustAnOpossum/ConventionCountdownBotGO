#!/bin/sh

echo "Installing crontab"

killall crond

rm /var/spool/cron/crontabs/root
echo "$CRON MODE=send /bot" >> /var/spool/cron/crontabs/root

crond

echo "Starting bot..."

/bot