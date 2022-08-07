test:
	docker-compose build try-checkout-via-ssh
	docker-compose up --force-recreate --renew-anon-volumes --remove-orphans try-checkout-via-ssh
test-canon:
	docker-compose build canonical
	docker-compose up --force-recreate --renew-anon-volumes --remove-orphans canonical
