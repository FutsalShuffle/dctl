#!/bin/bash
docker-php-entrypoint
composer install
bin/console doc:mi:mi --no-interaction
php-fpm
