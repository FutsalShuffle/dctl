#!/bin/bash
docker-php-entrypoint
composer install
php artisan migrate
php-fpm
